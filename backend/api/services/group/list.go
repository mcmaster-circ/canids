// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package group provides the group API service for the backend.
package group

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

// listResponse is the format of the list response.
type listResponse struct {
	Success bool                          `json:"success"` // Success indicates if the request was successful
	Current elasticsearch.DocumentGroup   `json:"current"` // Current is the group the requesting user is in
	Others  []elasticsearch.DocumentGroup `json:"others"`  // Others is a list of all groups only populated for superuser
}

// listHandler is "/api/group/list". It is responsible for fetching the group of
// the current standard or admin user. If a superuser is calling the API, it
// will also return the list of other groups.
func listHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// fetch all groups
	groups, err := elasticsearch.AllGroup(s)
	if err != nil {
		l.Error("error fetching all groups ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	// populate the response
	out := listResponse{Success: true}
	for _, group := range groups {
		if current.Group == group.UUID {
			// add the requesting user's group to current field
			out.Current = group
		} else if current.Class == jwtauth.UserSuperuser {
			// if superuser requesting, add all other fields
			out.Others = append(out.Others, group)
		}
	}

	// success
	json.NewEncoder(w).Encode(out)
}
