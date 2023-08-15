// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package blacklist provides the blacklist API service for the backend.
package blacklist

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/scheduler"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/uuid"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// addRequest is the format of the blacklist add request.
type addRequest struct {
	Name string `json:"name"` // Name is common blacklist name
	URL  string `json:"url"`  // URL is the URL of the blacklist source list
}

// addHandler is "/api/blacklist/add". It is responsible for creating a new blacklist.
// Only an admin can request for a new blacklist to be created. The blacklist UUID
// must not exist and the blacklist name must be unique.
func addHandler(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// only admins can use this endpoint
	if current.Class != jwtauth.UserAdmin {
		l.Warn("non admin attempting to create new blacklist")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Only an admin can create a new blacklist.",
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

	//ensure that the blacklist name is not " " or "  "
	trimmed := strings.TrimSpace(request.Name)

	// ensure name is not empty
	if request.Name == "" || len(trimmed) == 0 {
		l.Warn("blacklist name not specified specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Name field must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	//ensure url is not empty
	if request.URL == "" {
		l.Warn("blacklist url not specified specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Url field must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Ensure name of blacklist is not beginning or ending in whitespace
	for i, character := range request.Name {

		if unicode.IsSpace(character) {
			if i == 0 || i == (len(request.Name)-1) {
				l.Warn("Blacklist name cannot begin or end in whitespace")
				w.WriteHeader(http.StatusBadRequest)
				out := GeneralResponse{
					Success: false,
					Message: "Blacklist name cannot begin or end in whitespace",
				}
				json.NewEncoder(w).Encode(out)
				return
			}
		}
	}

	// Ensure url of blacklist is not beginning or ending in whitespace
	for i, character := range request.URL {

		if unicode.IsSpace(character) {
			if i == 0 || i == (len(request.Name)-1) {
				l.Warn("Blacklist url cannot begin or end in whitespace")
				w.WriteHeader(http.StatusBadRequest)
				out := GeneralResponse{
					Success: false,
					Message: "Blacklist url cannot begin or end in whitespace",
				}
				json.NewEncoder(w).Encode(out)
				return
			}
		}
	}

	// generate new blacklist
	blacklistUUID := uuid.Generate()
	blacklist := elasticsearch.DocumentBlacklist{
		UUID: blacklistUUID,
		Name: request.Name,
		URL:  request.URL,
	}

	// retreive all blacklists
	blacklists, err := elasticsearch.AllBlacklists(s)
	if err != nil {
		l.Error("error fetching all blacklists ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// ensure request name, url and UUID is unique
	for _, blacklist := range blacklists {
		if request.Name == blacklist.Name || blacklistUUID == blacklist.UUID {
			l.Warn("blacklist name in use")
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Blacklist name already in use.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}

		if request.URL == blacklist.URL {
			l.Warn("blacklist url in use")
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Blacklist url already in use.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	}

	// index new blacklist
	_, err = blacklist.Index(s)
	if err != nil {
		l.Error("error indexing new blacklist ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// get new IPs from blacklist
	go scheduler.Refresh(s)

	// success
	l.Info("successfully created new blacklist ", blacklist.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "Blacklist successfully created.",
	}
	json.NewEncoder(w).Encode(out)
}
