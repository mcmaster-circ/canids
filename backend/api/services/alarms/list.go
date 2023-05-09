// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package alarms provides the alarms API service for the backend.
package alarms

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type listField struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type listReponse struct {
	Fields []listField `json:"fields"`
}

// listHandler is "/api/alarms/list"
func listHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	_, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")
	l.Info("TESTS")

	// TODO(Russell): get alarm source files from elasticsearch
	fields := []listField{
		{
			Name: "firehol_abusers_1d",
			URL:  "https://iplists.firehol.org/files/firehol_abusers_1d.netset",
		},
		{
			Name: "firehol_abusers_30d",
			URL:  "https://iplists.firehol.org/files/firehol_abusers_30d.netset",
		},
		{
			Name: "firehol_anonymous",
			URL:  "https://iplists.firehol.org/files/firehol_anonymous.netset",
		},
		{
			Name: "firehol_level1",
			URL:  "https://iplists.firehol.org/files/firehol_level1.netset",
		},
		{
			Name: "firehol_level2",
			URL:  "https://iplists.firehol.org/files/firehol_level2.netset",
		},
		{
			Name: "firehol_level3",
			URL:  "https://iplists.firehol.org/files/firehol_level3.netset",
		},
	}

	// success
	l.Info("successfully queried for fields")
	json.NewEncoder(w).Encode(fields)
}
