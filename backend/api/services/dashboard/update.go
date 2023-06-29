// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package dashboard provides the dashboard API service for the backend.
package dashboard

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// updateRequest is the format of the dashboard update request.
type updateRequest struct {
	elasticsearch.DocumentDashboard // same format as Elasticsearch document
}

// updateHandler is "/api/dashboard/update". It is responsible for updating an
// existing dashboard. An admin can update the dashboard for everyone.
func updateHandler(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// reject request if standard user is making it
	if current.Class == jwtauth.UserStandard {
		l.Warn("standard user attempting to update dashboard")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Standard users can not update dashboard.",
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
	// ensure all fields are present
	if request.UUID == "" || request.Name == "" {
		l.Warn("not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "All fields must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return

	}
	// query current dashboard to update
	dashboard, esDocID, err := elasticsearch.QueryDashboardByUUID(s, request.UUID)
	if err != nil {
		l.Warn("invalid dashboard uuid ", request.UUID)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid dashboard UUID specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// ensure each view exists
	for _, view := range request.Views {
		// query the view
		_, _, err := elasticsearch.QueryViewByUUID(s, view)
		if err != nil {
			l.Warn("invalid view uuid provided ", view)
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Invalid view UUID specified.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	}
	// update dashboard with new parameters
	dashboard.Name = request.Name
	dashboard.Views = request.Views
	dashboard.Sizes = request.Sizes

	// index changes
	err = dashboard.Update(s, esDocID)
	if err != nil {
		l.Error("error updating dashboard ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully updated dashboard ", request.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "Successfully updated dashboard.",
	}
	json.NewEncoder(w).Encode(out)
}
