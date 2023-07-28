// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package blacklist provides the blacklist API service for the backend.
package blacklist

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/scheduler"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// updateRequest is the format of the blacklist update request.
type updateRequest struct {
	elasticsearch.DocumentBlacklist // same structure as Elasticsearch document
}

// updateHandler is "/api/blacklist/update". It is responsible for updating an
// existing blacklist. Only an admin can request for a blacklist to be updated.
func updateHandler(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// only admins can use this endpoint
	if current.Class != jwtauth.UserAdmin {
		l.Warn("non admin attempting to update blacklist")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Only an admin can update blacklist.",
		}
		json.NewEncoder(w).Encode(out)
		return
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

	//ensure that the blacklist name is not " " or "  "
	trimmed := strings.TrimSpace(request.Name)

	// ensure all fields are specified
	if request.Name == "" || request.UUID == "" || len(trimmed) == 0 {
		l.Warn("not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "All fields must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// retreive all blacklists
	blacklists, err := elasticsearch.AllBlacklists(s)
	if err != nil {
		l.Error("error fetching all blacklists ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// ensure request name is unique
	for _, blacklist := range blacklists {
		if request.Name == blacklist.Name {
			l.Warn("blacklist name in use")
			w.WriteHeader(http.StatusBadRequest)
			out := GeneralResponse{
				Success: false,
				Message: "Blacklist name already in use.",
			}
			json.NewEncoder(w).Encode(out)
			return
		}
	}

	// query elasticsearch for document ID
	_, esDocID, err := elasticsearch.QueryBlacklistByUUID(s, request.UUID)
	if err != nil {
		l.Warn("invalid blacklist uuid ", request.UUID)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid blacklist UUID provided.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// update document
	err = request.Update(s, esDocID)
	if err != nil {
		l.Error("cannot update blacklist ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// remove new IPs from blacklist
	go scheduler.Refresh(s)

	// success
	l.Info("successfully updated blacklist ", request.UUID)
	out := GeneralResponse{
		Success: true,
		Message: "Successfully updated blacklist.",
	}
	json.NewEncoder(w).Encode(out)
}
