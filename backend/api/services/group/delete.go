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

// deleteRequest is the format of the group delete request.
type deleteRequest struct {
	UUID string `json:"uuid"` // UUID is unique group identifier
}

// deleteHandler is "/api/group/delete". It is responsible ..
func deleteHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// only superusers can use this endpoint
	if current.Class != jwtauth.UserSuperuser {
		l.Warn("non superuser attempting to delete group")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Only a superuser can delete group.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// attempt to parse request
	var request deleteRequest
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

	// ensure field is specified
	if request.UUID == "" {
		l.Warn("uuid field not specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Must specify UUID field.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// get group from database
	group, _, err := elasticsearch.QueryGroupByUUID(s, request.UUID)
	if err != nil {
		l.Error("failed to query for group ", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GeneralResponse{
			Success: false,
			Message: "Invalid group UUID.",
		})
		return
	}

	// query for users belonging to group
	users, err := elasticsearch.QueryAuthByGroup(s, request.UUID)
	if err != nil {
		// if error nil, users were found
		l.Error("failed to retreive users in group group ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// can only delete a group with no users
	if len(users) > 0 {
		l.Warn("attempting to delete a non-empty group")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "The group must be empty.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// prevent deletion of group that is currently in a view
	views, err := elasticsearch.AllView(s)
	if err != nil {
		l.Error("failed to get views ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	groupInView := false
	for _, view := range views {
		if view.Group == group.UUID {
			groupInView = true
			break
		}
	}

	if groupInView {
		l.Warn("attempting to delete group that is in a view")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Cannot delete groups that are in a view.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// prevent deletion of group assets that are currently in a view
	groupAssetInView := false
	for _, view := range views {
		for _, groupAsset := range group.Authorized {
			if view.Authorized == groupAsset {
				groupAssetInView = true
				break
			}
		}
	}

	if groupAssetInView {
		l.Warn("attempting to delete group whos asset is in a view")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Cannot delete groups whos assets are used in a view.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// delete group
	err = elasticsearch.DeleteGroupByUUID(s, request.UUID)
	if err != nil {
		l.Error("failed to delete group ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully deleted group ", request.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "Successfully deleted group.",
	}
	json.NewEncoder(w).Encode(out)
}
