// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package user provides the user API service for the backend.
package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// infoHandler is "/api/user/info". It returns the authenticated user info.
func infoHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user from request
	current := jwtauth.FromContext(ctx)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(current)
}
