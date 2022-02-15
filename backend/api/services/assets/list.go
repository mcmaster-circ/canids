// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package data provides the data API service for the backend.
package assets

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

type listResponse struct {
	Assets []string `json:"assets"`
}

// listHandler is "/api/assets/list"
func listHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	_, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	assetNameSet := make(map[string]bool)

	// get all groups from database
	groups, err := elasticsearch.AllGroup(s)
	if err != nil {
		l.Error("error querying groups: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// get all authorized assets and add them to the set
	for _, group := range groups {
		for _, asset := range group.Authorized {
			if !assetNameSet[asset] {
				assetNameSet[asset] = true
			}
		}
	}

	// get all data-conn indices in the database and add them to the set
	dataConns, err := elasticsearch.ListDataAssets(s)
	if err != nil {
		l.Error("error querying asset names: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	for _, dataConn := range dataConns {
		if !assetNameSet[dataConn] {
			assetNameSet[dataConn] = true
		}
	}

	assetNames := []string{}
	for k := range assetNameSet {
		assetNames = append(assetNames, k)
	}

	l.Info("successfully queried for asset names")
	json.NewEncoder(w).Encode(listResponse{
		Assets: assetNames,
	})
}
