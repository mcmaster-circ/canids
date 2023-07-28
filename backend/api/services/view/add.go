// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package view provides the view API service for the backend.
package view

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/uuid"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// addRequest is the format of the add view request.
type addRequest struct {
	Name       string   `json:"name"`       // Name is common visualization name
	Class      string   `json:"class"`      // Class is the class of view
	DataIndex  string   `json:"index"`      // DataIndex is index fields are contained in
	Fields     []string `json:"fields"`     // Fields is the array of fields to be used in this view
	FieldNames []string `json:"fieldNames"` // FieldNames is the array of common field names
}

// addHandler is "/api/view/add". It is responsible for adding a new
// visualization. A standard user cannot create new visualizations.
// An admin may create visualizations for all assets.
func addHandler(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
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
	var request addRequest
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
	if request.Name == "" || request.Class == "" || request.DataIndex == "" {
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

	// ensure same number of fields and field names
	if len(request.Fields) != len(request.FieldNames) {
		l.Warnf("mismatched fields (%d) and field names (%d)", len(request.Fields), len(request.FieldNames))
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Must provide same number of fields and field names.",
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

	// create view for Elasticsearch
	viewUUID := uuid.Generate()
	view := elasticsearch.DocumentView{
		UUID:       viewUUID,
		Name:       request.Name,
		Class:      class,
		DataIndex:  request.DataIndex,
		Fields:     request.Fields,
		FieldNames: request.FieldNames,
	}
	// index view in database
	_, err = view.Index(s)
	if err != nil {
		l.Error("cannot index view ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	// everything was successful, return success message
	l.Info("successfully created new view ", view.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "View successfully created.",
	}
	json.NewEncoder(w).Encode(out)
}
