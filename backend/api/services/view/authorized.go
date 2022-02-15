// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package view provides the view API service for the backend.
package view

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

// authorizedResponse is the format of the authorized view.
type authorizedResponse struct {
	Success    bool     `json:"success"`    // Success indicates if the request was successful
	Authorized []string `json:"authorized"` // authorized is a list of available resources
}

// authorizedHandler is "/api/view/authorized". It is responsible for listing
// all assets accessible by the requesting user. A standard user cannot access
// any assets for new/updated visualizations. A standard can access authorized
// assets registered to their group and the aggregate. A superuser can access
// all authorized assets.
func authorizedHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// handler response
	var out authorizedResponse

	// standard user cannot access assets for creating visualizations
	if current.Class == jwtauth.UserStandard {
		json.NewEncoder(w).Encode(out)
		return
	}
	// retreive authorized in current group if admin
	if current.Class == jwtauth.UserAdmin {
		group, _, err := elasticsearch.QueryGroupByUUID(s, current.Group)
		if err != nil {
			l.Error("error getting groups ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}
		// append each resource in the authorized output list
		for _, resource := range group.Authorized {
			out.Authorized = append(out.Authorized, resource)
		}
	}
	// retreive authorized in all groups if superuser
	if current.Class == jwtauth.UserSuperuser {
		groups, err := elasticsearch.AllGroup(s)
		if err != nil {
			l.Error("error getting groups ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}
		// append all resources in the authorized output list
		for _, group := range groups {
			for _, resource := range group.Authorized {
				out.Authorized = append(out.Authorized, resource)
			}
		}
	}

	// can access aggregate information
	out.Authorized = append(out.Authorized, "Aggregate")

	// return list of authorized
	out.Success = true
	json.NewEncoder(w).Encode(out)
}
