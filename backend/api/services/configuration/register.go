// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package configuration provides the configuration API service for the backend.
package configuration

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

// RegisterRoutes registers routes to interact with the configuration.
func RegisterRoutes(s *state.State, a *jwtauth.Config, r *mux.Router) {
	// get configuration /api/configuration/list
	r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		listHandler(r.Context(), s, a, w, r)
	})
	// update configuration /api/configuration/update
	r.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		updateHandler(r.Context(), s, a, w, r)
	})
}
