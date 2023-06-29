// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package alarm provides the alarms API service for the backend.
package alarm

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

// RegisterRoutes registers routes to interact with the alarms.
func RegisterRoutes(s *state.State, a *jwtauth.Config, r *mux.Router) {
	// get alarms for dashboard /api/alarm/data
	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		dataHandler(r.Context(), s, a, w, r)
	})
}
