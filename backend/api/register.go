// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package api provides the API service for the backend.
package api

import (
	"github.com/gorilla/mux"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/alarm"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/assets"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/auth"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/blacklist"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/configuration"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/dashboard"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/data"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/fields"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/user"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/view"
	"github.com/mcmaster-circ/canids-v2/backend/api/services/websocket"
	authUtils "github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// registerRoutes registers all routes required for the HTTP service. Routes
// registered with the unsecure router will not require authentication for
// access. Routes registered with the secure router will require authentication.
func registerRoutes(s *state.State, a *jwtauth.Config, p *authUtils.State, unsecure *mux.Router, secure *mux.Router) {
	// register index assets
	registerIndexAssets(s, a, p, unsecure)

	// register static assets: /static
	registerStaticAssets(s, unsecure)

	// register auth service, do not require authentication: /api/auth
	auth.RegisterRoutes(s, a, unsecure.PathPrefix("/api/auth/").Subrouter())

	// register user service, require authentication: /api/user
	user.RegisterRoutes(s, a, secure.PathPrefix("/user/").Subrouter())

	// register view service, require authentication: /api/view
	view.RegisterRoutes(s, a, secure.PathPrefix("/view/").Subrouter())

	// register dashboard service, require authentication: /api/dashboard
	dashboard.RegisterRoutes(s, a, secure.PathPrefix("/dashboard/").Subrouter())

	// register data service, require authentication: /api/data
	data.RegisterRoutes(s, a, secure.PathPrefix("/data/").Subrouter())

	// register fields service, require authentication: /api/fields
	fields.RegisterRoutes(s, a, secure.PathPrefix("/fields/").Subrouter())

	// register alarm service, require authentication: /api/alarm
	alarm.RegisterRoutes(s, a, secure.PathPrefix("/alarm/").Subrouter())

	// register assets service, require authentication: /api/blacklist
	blacklist.RegisterRoutes(s, a, secure.PathPrefix("/blacklist/").Subrouter())

	// register assets service, require authentication: /api/configuration
	configuration.RegisterRoutes(s, a, secure.PathPrefix("/configuration/").Subrouter())

	// register assets service, require authentication: /api/assets
	assets.RegisterRoutes(s, a, secure.PathPrefix("/assets/").Subrouter())

	// Register ws for ingestion
	websocket.RegisterWS(s, unsecure.PathPrefix("/websocket/").Subrouter())

	// Register functions for updating ingestion
	websocket.RegisterUpdateFunctions(s, secure.PathPrefix("/ingestion/").Subrouter())
}
