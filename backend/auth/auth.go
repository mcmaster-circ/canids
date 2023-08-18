// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package auth provides the authentication state for the backend.
package auth

import (
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

const (
	// SecretLength is length of authentication secret
	SecretLength = 128 // bytes (1024 bits), 64 is lowest acceptable value https://tools.ietf.org/html/rfc7518#section-3.2
	// ExpireAge is 3 days, how long a token is valid for
	ExpireAge = 72 * time.Hour
	// RenewAge is 1 minute, how old a valid token can be before it gets renewed
	RenewAge = 1 * time.Minute
	// ResetDuration is 24 hours, how long a reset token is valid for
	ResetDuration = 24 * time.Hour
)

// Provision initializes a State for the API microservice. It accepts the main
// program state. It will initialize the authentication state and authentication
// pages. It will return an initialized API  state or an error.
func Provision(s *state.State) (*jwtauth.Config, error) {
	s.Log.Info("[api] provisioning api state")

	// empty API state
	var err error
	var a *jwtauth.Config

	s.AuthReady = true

	s.Log.Info("[api] initializing authentication secret")
	secret, err := jwtauth.GenerateSeed(SecretLength)
	if err != nil {
		return nil, err
	}

	// Generate JWTState
	s.Log.Info("[api] initializing JWT authentication state")
	a, err = jwtauth.Init(secret)
	if err != nil {
		return nil, err
	}

	return a, nil
}
