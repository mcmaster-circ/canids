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

// updateRequest is the format of the update user request.
type updateRequest struct {
	Name      string `json:"name"`      // Name is desired user name
	UUID      string `json:"uuid"`      // UUID is desired user email
	Class     string `json:"class"`     // Class is desired user class
	Activated bool   `json:"activated"` // Activated is desired user activation status
}

// updateHandler is "/api/user/update". It allows for standard users and admins
// to update a user account. Standard user may change: name, email. Admin user
// may change: name, email, class of other users amd activation
// of other users. An admin cannot change their own class/activation to prevent
// accidental lockout. Another user is required for this. If a uuid query
// parameter is specified, it will attempt to modify another user. Else, it will
// attempt to modify the current user.
func updateHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// assume editing requesting user
	userToUpdate := current.UUID
	// check if updating other user
	uuid := r.URL.Query().Get("uuid")
	if uuid != "" {
		userToUpdate = uuid
	}
	// attempt to parse request
	var request updateRequest
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
	// ensure all fields are present
	if request.Name == "" || request.UUID == "" || request.Class == "" {
		l.Warn("not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "All fields must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// query user from database
	existing, esDocID, err := elasticsearch.QueryAuthByUUID(s, userToUpdate)
	if err != nil {
		l.Error("error getting user specified in request ", err)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid UUID provided.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// find if modifying self
	modifyingSelf := current.UUID == userToUpdate

	// find current user status
	isStandard := current.Class == jwtauth.UserStandard

	// check what parameters are being modified
	nameModified := existing.Name != request.Name
	emailModified := existing.UUID != request.UUID
	classModified := string(existing.Class) != request.Class
	activatedModified := existing.Activated != request.Activated

	// check if things have to change
	if !nameModified && !emailModified && !classModified && !activatedModified {
		l.Info("no changes in user update")
		out := GeneralResponse{
			Success: true,
			Message: "No changes were applied.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// standard user can only change self
	if isStandard && !modifyingSelf {
		l.Warn("standard user attempting to change other user")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Standard users can only modify their own account.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// standard can change: name, email
	if isStandard && (classModified || activatedModified) {
		l.Warn("standard user attempting to change class/activation")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Standard users can only modify name and email.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// prevent lockout: admins cannot modify own class/activation
	if modifyingSelf && (classModified || activatedModified) {
		l.Warn("user attempting to modify own class/activation")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Users can not modify their own class or activation.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// if email modified, make sure it doesn't already exist
	if emailModified {
		// query database for requested email
		_, _, err := elasticsearch.QueryAuthByUUID(s, request.UUID)
		if err == nil {
			// no error means user was found
			l.Warn("uuid already exists ", request.UUID)
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "UUID email address provided already has account.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	}
	// if class modified, ensure class is valid
	if classModified {
		_, ok := jwtauth.UserClassMap[request.Class]
		if !ok {
			l.Warn("invalid class ", request.Class)
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Invalid class provided.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	}

	// passed all the tests, update existing user as requested
	existing.Name = request.Name
	existing.UUID = request.UUID
	existing.Class = jwtauth.UserClassMap[request.Class]
	existing.Activated = request.Activated

	// attempt to commit changes to database
	err = existing.Update(s, esDocID)
	if err != nil {
		l.Error("failed to update account ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully updated user account ", existing.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "The user account has been successfully updated. Changes will be applied within 1 minute.",
	}
	json.NewEncoder(w).Encode(out)
}
