// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package user provides the user API service for the backend.
package user

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

// RegisterRoutes registers routes to interact with the users.
func RegisterRoutes(s *state.State, a *jwtauth.Config, r *mux.Router) {
	// current user info /api/user/info
	r.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		infoHandler(r.Context(), s, a, w, r)
	})
	// list user /api/user/list
	r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		listHandler(r.Context(), s, a, w, r)
	})
	// add user /api/user/add
	r.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		addHandler(r.Context(), s, a, w, r)
	})
	// update user /api/user/update
	r.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		updateHandler(r.Context(), s, a, w, r)
	})
	// reset pass for other user /api/user/add
	r.HandleFunc("/resetPass", func(w http.ResponseWriter, r *http.Request) {
		resetPassHandler(r.Context(), s, a, w, r)
	})
	// delete other user /api/user/delete
	r.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		deleteHandler(r.Context(), s, a, w, r)
	})
}
