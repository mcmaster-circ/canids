// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package user provides the user API service for the backend.
package user

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

// deleteRequest is the format of the delete user request.
type deleteRequest struct {
	UUID string `json:"uuid"` // UUID is email of user to delete
}

// deleteHandler is "/api/user/delete". It is responsible for deleting users. A
// user can not delete their own account. A standard user can not delete an
// account. Admins can delete all users.
func deleteHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// standard user can not delete user
	if current.Class == jwtauth.UserStandard {
		l.Warn("standard user attempting to delete user")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Standard user can not delete users.",
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
	// ensure email field is present
	if request.UUID == "" {
		l.Warn("uuid field not specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "UUID field must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// user cannot delete self
	if current.UUID == request.UUID {
		l.Warn("user attempting to delete self")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "User can not delete own account.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// query user from database
	_, _, err = elasticsearch.QueryAuthByUUID(s, request.UUID)
	if err != nil {
		l.Error("error getting user specified in request ", err)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid email provided.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// delete user
	err = elasticsearch.DeleteAuthByUUID(s, request.UUID)
	if err != nil {
		l.Error("failed to delete user account ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully deleted user account ", request.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "The user account has been successfully deleted.",
	}
	json.NewEncoder(w).Encode(out)
}
