// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package alarm provides the alarms API service for the backend.
package alarm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type dataRequest struct {
	Index   []string `json:"index"`   // Index is the list of indices to search
	Source  []string `json:"source"`  // Source is the list of sources to search
	Dest    []string `json:"dest"`    // Dest is the list of destination alarms to search
	Start   string   `json:"start"`   // Start is the start time of the search
	End     string   `json:"end"`     // End is the end time of the search
	MaxSize int      `json:"maxSize"` // MaxSize is the maximum number of documents to return
	From    int      `json:"from"`    // From is the starting index of the search
}

type dataResponse struct {
	Alarms        []elasticsearch.Alarm `json:"alarms"`        // Alarms is the list of alarms
	AvailableRows int                   `json:"availableRows"` // AvailableRows is the number of rows available from elasticsearch
}

const maxCards = 20

// dataHandler is "/api/data. It is responsible for populating a view with the data related to that view.
func dataHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	_, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// attempt to parse request
	var request dataRequest
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

	// Parse "start" into time objects
	start, err := time.Parse(time.RFC3339, request.Start)
	if err != nil {
		l.Errorf("error parsing start time '%s': %v", request.Start, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// Parse "end" into time objects
	end, err := time.Parse(time.RFC3339, request.End)
	if err != nil {
		l.Errorf("error parsing end time '%s': %v", request.End, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	availableRows := 0

	// make sure maxSize is greater than 0
	if request.MaxSize <= 0 {
		l.Error("invalid max size: ", request.MaxSize)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GeneralResponse{
			Success: false,
			Message: "Invalid max size, must be greater than 0",
		})
		return
	}

	// make sure maxSize is under or equal to max allowed value
	if request.MaxSize > maxCards {
		l.Error("invalid max size: ", request.MaxSize)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GeneralResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid max size, must be under %d", maxCards),
		})
		return
	}

	// get data for the specified fields in the specified time range, sorted by timestamp
	data, availableRows, err := elasticsearch.GetAlarms(s, request.Index, request.Source, request.Dest, start, end, request.MaxSize, request.From)
	if err != nil {
		l.Error("error querying data conn: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// success
	l.Info("successfully queried data for asset")
	json.NewEncoder(w).Encode(dataResponse{
		Alarms:        data,
		AvailableRows: availableRows,
	})
}
