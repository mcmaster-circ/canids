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

// getResponse is the format of the group add request.
type getResponse struct {
	Success   bool                            `json:"success"`   // Success indicates if the request was successful
	Dashboard elasticsearch.DocumentDashboard `json:"dashboard"` // Dashboard is same format as Elasticsearch document
}

// getHandler is "/api/dashboard/get". It is responsible for fetching the group
// dashboard. If the dashboard does not exist but the group has views, a new
// dashboard will be created and saved using the existing views. If no views
// exist, new views and a new dashboard will be created using the "defaultViews"
// list found in default.go. A superuser can specify a "group" query parameter
// to fetch another group's dashboard.
func getHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// by default getting group dashboard of current user
	groupUUID := current.Group
	// superuser can request other dashboard
	if current.Class == jwtauth.UserSuperuser {
		groupParam := r.URL.Query().Get("group")
		if groupParam != "" {
			groupUUID = groupParam
		}
	}
	// attempt to get dashboard
	dashboard, _, err := elasticsearch.QueryDashboardByGroup(s, groupUUID)
	if err != nil {
		// could not find dashboard, need to provision new dashboard
		group, _, err := elasticsearch.QueryGroupByUUID(s, groupUUID)
		if err != nil {
			// invalid group provided
			l.Warn("invalid group ", groupUUID)
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Invalid group provided.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
		// ensure the group has authorized assets
		if len(group.Authorized) == 0 {
			l.Warn("no authorized assets available ", groupUUID)
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Group has no authorized assets to visualize.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
		// provision a new dashboard
		dashboard, err = provisionDashboard(s, group)
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
