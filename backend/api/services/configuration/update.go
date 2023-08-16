// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package configuration provides the configuration API service for the backend.
package configuration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// updateRequest is the format of the configuration update request.
type updateRequest struct {
	Configuration []state.DocumentSetting `json:"configuration"` // same structure as Elasticsearch document
}

// updateHandler is "/api/configuration/update". It is responsible for updating the
// existing configuration. Only an admin can request for the configuration to be updated.
func updateHandler(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// only admins can use this endpoint
	if current.Class != jwtauth.UserAdmin {
		l.Warn("non admin attempting to update settings")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Only an admin can update settings.",
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

	// print request
	fmt.Printf("received request to update settings: %+v", request)

	// update configuration
	err = state.UpdateSettings(s, request.Configuration)
	if err != nil {
		l.Error("cannot update settings ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully updated settings")
	out := GeneralResponse{
		Success: true,
		Message: "Successfully updated settings.",
	}
	json.NewEncoder(w).Encode(out)
}
