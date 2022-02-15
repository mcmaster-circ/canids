// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package api provides the API service for the backend.
package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// registerStaticAssets registers handlers to serve the static assets.
func registerStaticAssets(s *state.State, r *mux.Router) {
	// register authentication static assets
	r.HandleFunc("/static/main.css", func(w http.ResponseWriter, r *http.Request) {
		stylesheetHandler(s, w, r)
	})
	r.HandleFunc("/static/logo.png", func(w http.ResponseWriter, r *http.Request) {
		logoHandler(s, w, r)
	})
	r.HandleFunc("/static/load.js", func(w http.ResponseWriter, r *http.Request) {
		loadHandler(s, w, r)
	})
}

// stylesheetHandler receives "/static/main.css" HTTP requests. It returns a
// static stylesheet.
func stylesheetHandler(s *state.State, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "assets/main.css")
}

// logoHandler receives "/static/logo.png" HTTP requests. It returns a static
// logo.
func logoHandler(s *state.State, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, "assets/logo.png")
}

// loadHandler receives "/static/load.js" HTTP requests. It returns a static
// javascript field.
func loadHandler(s *state.State, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, "assets/load.js")
}
