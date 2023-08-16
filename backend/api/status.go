// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package api provides the API service for the backend.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/uuid"
	"github.com/mcmaster-circ/canids-v2/backend/state"
	log "github.com/sirupsen/logrus"
)

// status is the return of "/status".
type status struct {
	Name        string `json:"name"`        // Name of application
	Build       string `json:"build"`       // Build is the Git build hash
	IsDocker    bool   `json:"isDocker"`    // IsDocker indicates if running in Docker
	ElasticPing bool   `json:"elasticPing"` // ElasticPing indicates of Elasticsearch is connected
	Time        string `json:"time"`        // Time is the current server time
	Uptime      string `json:"uptime"`      // Uptime is the uptime of the backend
}

// statusHandler is "/status". It returns the status of the backend.
func statusHandler(s *state.State, w http.ResponseWriter, r *http.Request) {
	// ping elasticsearch
	elasticPing := true
	esURI := fmt.Sprintf("http://%s:%s", s.Config.ElasticHost, s.Config.ElasticPort)
	_, _, err := s.Elastic.Ping(esURI).Do(s.ElasticCtx)
	if err != nil {
		// failed to ping elasticsearch
		elasticPing = false
	}
	// generate response
	now := time.Now().UTC()
	out := status{
		Build:       s.Hash,
		IsDocker:    s.IsDocker,
		ElasticPing: elasticPing,
		Time:        now.Format(time.RFC3339),
		Uptime:      now.Sub(s.Start).String(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// notFoundHandler is called when an unmatched route is requested. It returns a
// standard 404 message. Context is normally used for logging but a 404 bypass
// the initial middleware that generates the context.
func notFoundHandler(s *state.State, w http.ResponseWriter, r *http.Request) {
	// if real ip is available, use it
	realIP := r.Header.Get("X-Real-IP")
	addr := r.RemoteAddr
	if realIP != "" {
		addr = realIP
	}
	// log and return 404 error
	s.Log.WithFields(log.Fields{
		"request": uuid.Generate(),
		"addr":    addr,
		"method":  r.Method,
		"uri":     r.Host + r.RequestURI,
	}).Warn("[request] http 404 not found")

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 Not Found\n")
}
