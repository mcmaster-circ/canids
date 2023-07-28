// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package blacklist provides the blacklist API service for the backend.
package blacklist

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
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

// RegisterRoutes registers routes to interact with the blacklist.
func RegisterRoutes(s *state.State, a *jwtauth.Config, r *mux.Router) {
	// get blacklist source lists for dashboard /api/blacklist/list
	r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		listHandler(r.Context(), s, a, w, r)
	})
	// add blacklist /api/blacklist/add
	r.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		addHandler(r.Context(), s, a, w, r)
	})
	// update blacklist /api/blacklist/update
	r.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		updateHandler(r.Context(), s, a, w, r)
	})
	// delete blacklist /api/blacklist/delete
	r.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		deleteHandler(r.Context(), s, a, w, r)
	})
}
