// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package blacklist provides the blacklist API service for the backend.
package blacklist

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/scheduler"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// deleteRequest is the format of the blacklist delete request.
type deleteRequest struct {
	UUID string `json:"uuid"` // UUID is a unique blacklist identifier
}

// deleteHandler is "/api/blacklist/delete". It is responsible ..
func deleteHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// only superusers can use this endpoint
	if current.Class != jwtauth.UserSuperuser {
		l.Warn("non superuser attempting to delete blacklist")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Only a superuser can delete blacklist.",
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

	// delete blacklist
	err = elasticsearch.DeleteBlacklistByUUID(s, request.UUID)
	if err != nil {
		l.Error("failed to delete blacklist ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// remove new IPs from blacklist
	go scheduler.Refresh(s)

	// success
	l.Info("successfully deleted blacklist ", request.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "Successfully deleted blacklist.",
	}
	json.NewEncoder(w).Encode(out)
}
