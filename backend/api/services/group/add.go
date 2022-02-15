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
	"github.com/mcmaster-circ/canids-v2/backend/libraries/uuid"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// addRequest is the format of the group add request.
type addRequest struct {
	Name       string   `json:"name"`       // Name is common group name
	Authorized []string `json:"authorized"` // Authorized is list of authorized assets
}

// addHandler is "/api/group/add". It is responsible for creating a new group.
// Only a superuser can request for a new group to be created. The group UUID
// must not exist and the group name must be unique.
func addHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// only superusers can use this endpoint
	if current.Class != jwtauth.UserSuperuser {
		l.Warn("non superuser attempting to create new group")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Only a superuser can create a new group.",
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
	// ensure name is not empty
	if request.Name == "" {
		l.Warn("group name not specified specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Name field must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// generate new group
	groupUUID := uuid.Generate()
	group := elasticsearch.DocumentGroup{
		UUID:       groupUUID,
		Name:       request.Name,
		Authorized: request.Authorized,
	}
	// retreive all groups
	groups, err := elasticsearch.AllGroup(s)
	if err != nil {
		l.Error("error fetching all groups ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	// ensure request name and UUID is unique
	for _, group := range groups {
		if request.Name == group.Name || groupUUID == group.UUID {
			l.Warn("group name in use")
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Group name already in use.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	}
	// index new group
	_, err = group.Index(s)
	if err != nil {
		l.Error("error indexing new group ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully created new group ", group.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "Group successfully created.",
	}
	json.NewEncoder(w).Encode(out)
}
