// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package dashboard provides the dashboard API service for the backend.
package dashboard

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

// getResponse is the format of the dashboard get request.
type getResponse struct {
	Success   bool                            `json:"success"`   // Success indicates if the request was successful
	Dashboard elasticsearch.DocumentDashboard `json:"dashboard"` // Dashboard is same format as Elasticsearch document
}

// getHandler is "/api/dashboard/get". It is responsible for fetching the
// dashboard. If the dashboard does not exist but has views, a new
// dashboard will be created and saved using the existing views. If no views
// exist, new views and a new dashboard will be created using the "defaultViews"
// list found in default.go.
func getHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	_, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// attempt to get dashboard
	dashboard, err := elasticsearch.GetDashboard(s)
	if err != nil {
		// could not find dashboard, need to provision new dashboard+
		dashboard, err := provisionDashboard(s)
		if err != nil {
			// could not provision new dashboard
			l.Error("error provisioning new dashboard ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}
		l.Info("successfully provisioned dashboard ", dashboard.UUID)
	}

	// successs
	json.NewEncoder(w).Encode(dashboard)
}
