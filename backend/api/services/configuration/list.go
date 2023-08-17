// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package configuration provides the configuration API service for the backend.
package configuration

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type listResponse struct {
	Success       bool                    `json:"success"`       // Success indicates if the request was successful
	Configuration []state.DocumentSetting `json:"configuration"` // Configuation is the list of editable settings
}

// listHandler is "/api/configuration/list"
func listHandler(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	_, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// fetch all settings if admin
	settings, err := state.AllSettings(s)
	if err != nil {
		l.Error("error fetching all settings ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	out := listResponse{
		Success:       true,
		Configuration: settings,
	}

	// success
	l.Info("successfully queried for all settings")
	json.NewEncoder(w).Encode(out)
}
