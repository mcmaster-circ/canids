// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package data provides the data API service for the backend.
package fields

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

type listField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type listReponse struct {
	Fields []listField `json:"fields"`
}

// listHandler is "/api/fields/list"
func listHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	_, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// get data conn fields from database
	fields, err := elasticsearch.GetAllDataMapping(s)
	if err != nil {
		l.Error("failed to get data conn mapping ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully queried for fields")
	json.NewEncoder(w).Encode(fields)
}
