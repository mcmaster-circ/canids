// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package api provides the API service for the backend.
package api

import (
	"context"
	"embed"
	"fmt"
	"mime"
	"net/http"
	_ "net/http/pprof" // performance profiling
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/mcmaster-circ/canids-v2/backend/api/services/websocket"
	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/uuid"
	"github.com/mcmaster-circ/canids-v2/backend/state"
	log "github.com/sirupsen/logrus"
)

//go:embed all:frontend
var frontendContent embed.FS

// Start accepts the global state and the authentication state. It will register
// all routes and start the HTTP server. If the server fails to start an error
// will be returned.
func Start(s *state.State, a *jwtauth.Config, p *auth.State) error {
	// create main request router
	router := mux.NewRouter()
	router.StrictSlash(true)

	// register performance profiling
	router.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	// log all working requests on debug level
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// provide context to every request
			ctx := context.Background()

			// if real ip is available, use it
			realIP := r.Header.Get("X-Real-IP")
			addr := r.RemoteAddr
			if realIP != "" {
				addr = realIP
			}
			// update context with fields
			fields := s.Log.WithFields(log.Fields{
				"request": uuid.Generate(),
				"addr":    addr,
				"method":  r.Method,
				"uri":     r.Host + r.RequestURI,
			})
			ctx = ctxlog.WithFields(ctx, fields)

			// inject context into request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// create /api router with access to middleware
	secureRouter := router.PathPrefix("/api/").Subrouter()
	secureRouter.Use(func(next http.Handler) http.Handler {
		return auth.Middleware(s, a, next)
	})

	// register status, 404 handler
	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(s, w, r)
	})

	// subFilesystem, _ := fs.Sub(frontendContent, "frontend/out")
	// router.PathPrefix("/").Handler(http.FileServer(http.FS(subFilesystem)))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		fileExtension := ".html"

		// load index.html for paths ending in '/'
		if path[len(path)-1] == '/' {
			path = path + "index.html"
		}

		// add .html to paths that have no file extension
		if !strings.Contains(path, ".") {
			path = path + ".html"
		}

		contents, err := frontendContent.ReadFile(fmt.Sprintf("frontend/out%s", path))
		if err != nil {
			log.Printf("No file found at %s: %v", r.URL.Path, err)
			notFoundHandler(s, w, r)
			return
		}

		w.Header().Add("Content-Type", mime.TypeByExtension(fileExtension))
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(contents)))
		w.Write(contents)
	})

	// register all routes
	registerRoutes(s, a, p, router, secureRouter)

	// Start frame queue handler
	go websocket.HandleQueue(s)

	server := &http.Server{
		Addr:         ":6060",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	s.Log.Info("[main] backend now listening on :6060")
	return server.ListenAndServe()
}
