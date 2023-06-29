// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package user provides the user API service for the backend.
package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// listReponse is the format of the list user response.
type listResponse struct {
	Success bool   `json:"success"` // Success indicates if the request was successful
	Users   []User `json:"users"`   // Users is a list of users
}

// User represents the list of users in the system.
type User struct {
	Name             string `json:"name"`             // Name of user
	UUID             string `json:"uuid"`             // UUID is email of user
	Class            string `json:"class"`            // Class of user
	Activated        bool   `json:"activated"`        // Activated indicates if user is activated
	UpdatePermission bool   `json:"updatePermission"` // UpdatePermission indicates if current user can modify this user.
}

// listHandler is "/api/user/list". It will return the list of users and
// also indicate if the user has permission to update the user, including
// password reset.
func listHandler(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// output
	var out listResponse

	// retreive users if admin
	if current.Class == jwtauth.UserAdmin {
		users, err := elasticsearch.AllAuth(s)
		if err != nil {
			l.Error("error getting users ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}

		for _, u := range users {
			// can update user if self or admin
			canUpdateUser := u.UUID == current.UUID
			if current.Class == jwtauth.UserAdmin {
				canUpdateUser = true
			}
			out.Users = append(out.Users, User{
				Name:             u.Name,
				UUID:             u.UUID,
				Class:            string(u.Class),
				Activated:        u.Activated,
				UpdatePermission: canUpdateUser,
			})
		}
	}

	// return list of users
	out.Success = true
	json.NewEncoder(w).Encode(out)
}
