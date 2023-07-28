// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package blacklist provides the blacklist API service for the backend.
package blacklist

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type listResponse struct {
	Success    bool                              `json:"success"`    // Success indicates if the request was successful
	Blacklists []elasticsearch.DocumentBlacklist `json:"blacklists"` // Blacklists is the list of blacklisted IP sources
}

// listHandler is "/api/blacklist/list"
func listHandler(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	_, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// fetch all blacklists if admin
	blacklists, err := elasticsearch.AllBlacklists(s)
	if err != nil {
		l.Error("error fetching all blacklists ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	out := listResponse{
		Success:    true,
		Blacklists: blacklists,
	}

	// success
	l.Info("successfully queried for blacklists")
	json.NewEncoder(w).Encode(out)
}
