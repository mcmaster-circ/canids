// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package view provides the view API service for the backend.
package view

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

// deleteRequest is the format of the view delete request.
type deleteRequest struct {
	UUID string `json:"uuid"` // UUID is unique view identifier
}

// deleteHandler is "/api/view/deleteHandler". It is responsible for deleting a
// view. A standard user cannot delete views. An admin can delete views for
// their group. A superuser can delete any view.
func deleteHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// only superusers can use this endpoint
	if current.Class != jwtauth.UserSuperuser {
		l.Warn("non superuser attempting to delete view")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Only a superuser can delete view.",
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

	// query elasticsearch to get view
	view, _, err := elasticsearch.QueryViewByUUID(s, request.UUID)
	if err != nil {
		l.Warn("invalid view uuid ", request.UUID)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid view UUID provided.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// if admin is deleting, ensure it is for same group
	if current.Class == jwtauth.UserAdmin {
		if view.Group != current.Group {
			l.Warn("admin user attempting to delete foreign views")
			w.WriteHeader(http.StatusForbidden)
			out := GeneralResponse{
				Success: false,
				Message: "Admin users can not delete foreign views.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	}

	// prevent deletion of a view that is currently in a dashboard
	dashboards, err := elasticsearch.AllDashboard(s)
	if err != nil {
		l.Error("failed to get dashboards ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	viewInDashboard := false
	for _, dashboard := range dashboards {
		for _, dashboardView := range dashboard.Views {
			if dashboardView == view.UUID {
				viewInDashboard = true
				break
			}
		}
	}

	if viewInDashboard {
		l.Warn("attempting to delete view that is in a dashboard")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Cannot delete views that are in a dashboard.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// delete view
	err = elasticsearch.DeleteViewByUUID(s, request.UUID)
	if err != nil {
		l.Error("failed to delete view ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully deleted view ", request.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "Successfully deleted view.",
	}
	json.NewEncoder(w).Encode(out)
}
