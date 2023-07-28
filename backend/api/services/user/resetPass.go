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
	"github.com/mcmaster-circ/canids-v2/backend/libraries/email"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// resetPassRequest is the format of the reset password request.
type resetPassRequest struct {
	UUID string `json:"uuid"` // UUID of user to reset password for
}

// Pass is "/api/user/resetPass". It is responsible for sending a password reset
// link to the provided user. A standard user can not use this endpoint. An admin
// may request a password reset for any user. If the requesting user has
// permissions and the provided user email is valid, the user will be emailed a
// password reset link. If the requesting user does not have permissions or if
// the user email is not valid, an error will be returned.
func resetPassHandler(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// standard user can not reset password here
	if current.Class == jwtauth.UserStandard {
		l.Warn("standard user attempting to reset password")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Standard user can not reset password.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// attempt to parse request
	var request resetPassRequest
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
			Message: "Email field must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// query user from database
	existing, _, err := elasticsearch.QueryAuthByUUID(s, request.UUID)
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
	// generate reset token with provided expiry
	payload := &jwtauth.Payload{UUID: request.UUID}
	token, err := a.CreateToken(payload, auth.ResetDuration)
	if err != nil {
		l.Error("failed to generate password reset token ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	// send password reset email
	domain := s.Config.SendGridDomain
	resetRequest := "http://" + domain + "/reset?token=" + token
	err = email.SendPasswordReset(s, existing.Name, existing.UUID, resetRequest)
	if err != nil {
		l.Error("failed to send password reset ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully issued password reset for user account ", request.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "A password reset has been successfully issued for the user.",
	}
	json.NewEncoder(w).Encode(out)
}
