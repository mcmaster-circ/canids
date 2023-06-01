// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package user provides the user API service for the backend.
package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/email"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// addRequest is the format of the add user request.
type addRequest struct {
	Name  string `json:"name"`  // Name of user to be created
	UUID  string `json:"uuid"`  // UUID is email of user to be created
	Class string `json:"class"` // Class of user to be created=
}

// addHandler is "/api/user/add". It is responsible for creating new users. An
// admin may create a user. An error will be returned if not all fields are
// specified or if the current user does not hold permissions to perform the action.
// If the account was successfully created, an email will be sent for the new user
// to set a password.
func addHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// reject request if standard user is making it
	if current.Class == jwtauth.UserStandard {
		l.Warn("standard user attempting to add user")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Standard users can not add users.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// attempt to parse request
	var request addRequest
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
	// ensure class is valid
	class, ok := jwtauth.UserClassMap[request.Class]
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
	// ensure user does not already exist
	_, _, err = elasticsearch.QueryAuthByUUID(s, request.UUID)
	if err == nil {
		// no error means we located a user
		l.Warn("uuid already exists ", request.UUID)
		out := GeneralResponse{
			Success: false,
			Message: "UUID email address provided already has account.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// create activated user for Elasticsearch
	user := elasticsearch.DocumentAuth{
		UUID:      request.UUID,
		Password:  "",
		Class:     class,
		Name:      request.Name,
		Activated: true,
	}
	// index user in database
	docID, err := user.Index(s)
	if err != nil {
		l.Error("cannot index user ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	l.Info("created new user auth/", docID)

	// generate account activation token (24 hour expiry)
	payload := &jwtauth.Payload{UUID: user.UUID}
	token, err := a.JWTState.CreateToken(payload, 24*time.Hour)
	if err != nil {
		l.Error("failed to generate acccount activation token ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	// send account activation email
	domain := s.Config.SendGridDomain
	resetRequest := "http://" + domain + "/reset?token=" + token
	err = email.SendNewReset(s, user.Name, user.UUID, resetRequest, current.Name, current.UUID)
	if err != nil {
		l.Error("failed to send account activation ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// everything was succesful, return success message
	l.Info("successfully created user account ", user.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "The user account has been successfully created. The user has been emailed to complete account activation.",
	}
	json.NewEncoder(w).Encode(out)
}
