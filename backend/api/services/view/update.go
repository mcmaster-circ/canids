// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package view provides the view API service for the backend.
package view

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// updateView is the format of the update view request.
type updateRequest struct {
	UUID       string   `json:"uuid"`       // UUID is unique view identifier
	Group      string   `json:"group"`      // Group is the group UUID the visualization is accessible by
	Authorized string   `json:"authorized"` // Authorized is authorized asset (index) used for generating data
	Name       string   `json:"name"`       // Name is common visualization name
	Class      string   `json:"class"`      // Class is the class of view
	DataIndex  string   `json:"index"`      // DataIndex is index fields are contained in
	Fields     []string `json:"fields"`     // Fields is the array of fields from Authorized to be used in this view
	FieldNames []string `json:"fieldNames"` // FieldNames is the array of common field names
}

// updateHandler is "/api/view/update". It is responsible for updating an
// existing view. A standard user cannot update views. An admin can update views
// for their group. A superuser can update all groups. The same restrictions
// regarding addHandler and authorized assets apply here.
func updateHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// reject request if standard user is making it
	if current.Class == jwtauth.UserStandard {
		l.Warn("standard user attempting to add view")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Standard users can not add views.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// attempt to parse request
	var request updateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		l.Warn("invalid request format")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Bad request format.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// ensure all fields are present
	if request.UUID == "" || request.Group == "" || request.Name == "" || request.Class == "" || request.DataIndex == "" {
		l.Warn("not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "All fields must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// ensure class is valid
	class, ok := elasticsearch.ViewClassMap[request.Class]
	if !ok {
		l.Warn("invalid class ", request.Class)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid class provided.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// ensure correct number of fields for each view class
	if (class == elasticsearch.ViewBar) || (class == elasticsearch.ViewPie) {
		if len(request.Fields) != 1 {
			l.Warnf("view class %s, expected 1 field, got %d", class, len(request.Fields))
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Bar/Pie views take 1 field.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	} else if class == elasticsearch.ViewLine {
		if len(request.Fields) != 2 {
			l.Warnf("view class %s, expected 2 field, got %d", class, len(request.Fields))
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Line views take 2 field.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	} else if class == elasticsearch.ViewTable {
		if len(request.Fields) == 0 {
			l.Warnf("view class %s, got no fields", class)
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Table view requires atleast one field.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	}

	// if admin is updating, ensure it is for same group
	if current.Class == jwtauth.UserAdmin && request.Group != current.Group {
		l.Warn("admin user attempting to add foreign views")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Admin users can not add foreign views.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// allowed admin or superuser here, ensure group actually exists
	group, _, err := elasticsearch.QueryGroupByUUID(s, request.Group)
	if err != nil {
		l.Warn("invalid group ", request.Group)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid group provided.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// admin can add authorized assets in group
	authorizedAssets := group.Authorized

	// a superuser can create visualizations with any authorized asset
	if current.Class == jwtauth.UserSuperuser {
		// need to fetch all authorized assets in all groups
		authorizedAssets = []string{}
		groups, err := elasticsearch.AllGroup(s)
		if err != nil {
			l.Error("cannot fetch all groups ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}
		// for each group, add each authorized asset to the list of authorized
		// assets
		for _, group := range groups {
			for _, authorized := range group.Authorized {
				authorizedAssets = append(authorizedAssets, authorized)
			}
		}
	}
	// ensure provided asset actually exists for requesting user
	validAsset := false
	for _, authorized := range authorizedAssets {
		if request.Authorized == authorized {
			// found valid asset
			validAsset = true
			break
		}
	}
	if request.Authorized != "Aggregate" && !validAsset {
		l.Warn("invalid authorized asset ", request.Authorized)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid asset provided.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// query elasticsearch for existing document ID
	_, esDocID, err := elasticsearch.QueryViewByUUID(s, request.UUID)
	if err != nil {
		l.Warn("invalid view uuid ", request.UUID)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid view UUID provided.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// create new document to update
	updatedDoc := elasticsearch.DocumentView{
		UUID:       request.UUID,
		Group:      request.Group,
		Authorized: request.Authorized,
		Name:       request.Name,
		Class:      class,
		DataIndex:  request.DataIndex,
		Fields:     request.Fields,
		FieldNames: request.FieldNames,
	}
	// update document
	err = updatedDoc.Update(s, esDocID)
	if err != nil {
		l.Error("cannot update view ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	// success
	l.Info("successfully updated view ", request.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "Successfully updated view.",
	}
	json.NewEncoder(w).Encode(out)
}
