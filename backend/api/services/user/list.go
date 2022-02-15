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

// listReponse is the format of the list user response.
type listResponse struct {
	Success bool    `json:"success"` // Success indicates if the request was successful
	Groups  []Group `json:"groups"`  // Groups is a list of groups
}

// Group represnts the list of groups in the system.
type Group struct {
	Group string `json:"group"` // Group is the group UUID
	Users []User `json:"users"` // Users is the list of users belonging to the group
}

// User represents the list of users in the system.
type User struct {
	Name             string `json:"name"`             // Name of user
	UUID             string `json:"uuid"`             // UUID is email of user
	Class            string `json:"class"`            // Class of user
	Activated        bool   `json:"activated"`        // Activated indicates if user is activated
	UpdatePermission bool   `json:"updatePermission"` // UpdatePermission indicates if current user can modify this user.
}

// listHandler is "/api/user/list". If the user requesting is standard or admin,
// it will return the list of users registered to the group. If the user
// requesting is superuser, it will return the list of users. It will also
// indicate if the user has permission to update the user, including password
// reset.
func listHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	current, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// output
	var out listResponse

	// list of groups to get users from
	var groupUUIDs []string

	// retreive users in current group if standard or admin
	if current.Class == jwtauth.UserAdmin || current.Class == jwtauth.UserStandard {
		groupUUIDs = append(groupUUIDs, current.Group)
	}
	// retreive users in all groups if superuser
	if current.Class == jwtauth.UserSuperuser {
		groups, err := elasticsearch.AllGroup(s)
		if err != nil {
			l.Error("error getting groups ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}
		// append all groups to list
		for _, group := range groups {
			groupUUIDs = append(groupUUIDs, group.UUID)
		}
	}
	// query for all users belonging to group
	for _, groupUUID := range groupUUIDs {
		// create a group response
		group := Group{
			Group: groupUUID,
		}
		users, err := elasticsearch.QueryAuthByGroup(s, groupUUID)
		if err != nil {
			l.Error("error getting users from group ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}
		// append each user to the group response
		for _, u := range users {
			// can update user if self or admin/superuser
			canUpdateUser := u.UUID == current.UUID
			if current.Class == jwtauth.UserAdmin || current.Class == jwtauth.UserSuperuser {
				canUpdateUser = true
			}
			group.Users = append(group.Users, User{
				Name:             u.Name,
				UUID:             u.UUID,
				Class:            string(u.Class),
				Activated:        u.Activated,
				UpdatePermission: canUpdateUser,
			})
		}
		// append group response to main output
		out.Groups = append(out.Groups, group)
	}

	// return list of users
	out.Success = true
	json.NewEncoder(w).Encode(out)
}
