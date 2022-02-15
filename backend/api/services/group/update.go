// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package group provides the group API service for the backend.
package group

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

// updateRequest is the format of the group update request.
type updateRequest struct {
	elasticsearch.DocumentGroup // same structure as Elasticsearch document
}

// updateHandler is "/api/group/update". It is responsible for updating an
// existing group. Only a superuser can request for a group to be updated.
func updateHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// only superusers can use this endpoint
	if current.Class != jwtauth.UserSuperuser {
		l.Warn("non superuser attempting to update group")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Only a superuser can update group.",
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

	// ensure all fields are specified
	if request.Name == "" || request.UUID == "" {
		l.Warn("not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "All fields must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// query elasticsearch for document ID
	group, esDocID, err := elasticsearch.QueryGroupByUUID(s, request.UUID)
	if err != nil {
		l.Warn("invalid group uuid ", request.UUID)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid group UUID provided.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// find assets that would be deleted by this update
	assetsToBeDeleted := []string{}
	for _, existingAsset := range group.Authorized {
		existingAssetFound := false
		for _, newAsset := range request.Authorized {
			if existingAsset == newAsset {
				existingAssetFound = true
				break
			}
		}

		if !existingAssetFound {
			assetsToBeDeleted = append(assetsToBeDeleted, existingAsset)
		}
	}

	// prevent deletion of group assets that are currently in a view
	views, err := elasticsearch.AllView(s)
	if err != nil {
		l.Error("failed to get views ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	groupAssetInView := false
	for _, view := range views {
		for _, groupAsset := range assetsToBeDeleted {
			if view.Authorized == groupAsset {
				groupAssetInView = true
				break
			}
		}
	}

	if groupAssetInView {
		l.Warn("attempting to delete a group asset that is in a view")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Cannot delete groups assets that are used in a view.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// update document
	err = request.Update(s, esDocID)
	if err != nil {
		l.Error("cannot update group ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully updated group ", request.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "Successfully updated group.",
	}
	json.NewEncoder(w).Encode(out)
}
