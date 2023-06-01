// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package view provides the view API service for the backend.
package view

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// GeneralResponse is the structure of a general response.
type GeneralResponse struct {
	Success bool   `json:"success"` // Success indicates if the request was successful
	Message string `json:"message"` // Message describes the request response
}

var (
	// InternalServerError is the a JSON error message.
	InternalServerError = GeneralResponse{
		Success: false,
		Message: "500 Internal Server Error",
	}
)

// RegisterRoutes registers routes to interact with the views.
func RegisterRoutes(s *state.State, a *auth.State, r *mux.Router) {
	// list of saved visualizations /api/view/list
	r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		listHandler(r.Context(), s, a, w, r)
	})
	// add new visualization /api/view/add
	r.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		addHandler(r.Context(), s, a, w, r)
	})
	// update visualization /api/view/update
	r.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		updateHandler(r.Context(), s, a, w, r)
	})
	// delete visualization /api/view/delete
	r.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		deleteHandler(r.Context(), s, a, w, r)
	})
}
