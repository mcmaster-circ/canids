// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package auth provides the authentication state for the backend.
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
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

	// generate JWTState and AuthPage in State
	// if err = provisionAuth(s, a); err != nil {
	// 	s.Log.Error("[api] failed to provision Auth in api state")
	// 	return nil, err
	// }

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

// provisionAuth accepts the main program state and the API state. It generates
// the JWTState and AuthPage entries in the API State or returns an error.
// func provisionAuth(s *state.State, a *jwtauth.Config) (*jwtauth.Config, error) {
// 	// generate seed
// 	s.Log.Info("[api] initializing authentication secret")
// 	secret, err := jwtauth.GenerateSeed(SecretLength)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// generate JWTState
// 	s.Log.Info("[api] initializing JWT authentication state")
// 	a, err = jwtauth.Init(secret)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func DefaultUserSetup(s *state.State, a *jwtauth.Config) {

	err := elasticsearch.CreateIndex(s, "auth")
	if err != nil {
		// return setup page with general error
		s.Log.Error("[Default user setup] cannot create 'auth' index ", err)
		return
	}

	// Create default random password
	password, err := randomPass(32)
	if err != nil {
		s.Log.Error("[Default user setup] Failed to generate random password")
		return
	} else {
		s.Log.Info("[Default user setup] Default password: ", password)

	}

	// Hash and salt random password
	hashedPass, err := jwtauth.HashPassword(password)
	if err != nil {
		s.Log.Error("[Default user setup] Failed to hash password")
		return
	}

	user := elasticsearch.DocumentAuth{
		Name:      "Admin",
		UUID:      "admin@system.test",
		Class:     jwtauth.UserAdmin,
		Password:  hashedPass,
		Activated: true,
	}

	_, err = user.Index(s)
	if err != nil {
		s.Log.Error("[Default user setup] cannot index user", err)
		return
	}

	s.AuthReady = true

	s.Log.Info("[Default user setup] created new default admin user. Scroll up to find password. Email is admin@system.test")
}

func randomPass(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), err
}
