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

// listResponse is the format of the add view request.
type listResponse struct {
	Success bool                         `json:"success"` // Success indicates if the request was successful
	Views   []elasticsearch.DocumentView `json:"views"`   // Views is a list of views for the requesting user
}

// listHandler is "/api/view/list". It is responsible for listing
// visualizations. For all users, it will return the list of saved views.
func listHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	_, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// output
	out := listResponse{}

	// fetch all views
	views, err := elasticsearch.AllView(s)
	if err != nil {
		l.Error("error fetching views ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	out.Views = views

	// if views array is nil, return empty list instead
	if out.Views == nil {
		out.Views = []elasticsearch.DocumentView{}
	}
	// return list of views
	out.Success = true
	json.NewEncoder(w).Encode(out)

}
