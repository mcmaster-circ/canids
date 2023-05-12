// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package alarm provides the alarms API service for the backend.
package alarm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type dataResponse struct {
	Alarms        []elasticsearch.Alarm `json:"alarms"`
	AvailableRows int                   `json:"availableRows"`
}

const maxCards = 20

// dataHandler is "/api/data. It is responsible for populating a view with the data related to that view.
func getHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	_, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// get query parameters "start" and "end"
	v := r.URL.Query()
	indices := v["index"]
	sources := v["source"]
	startStr := v.Get("start")
	endStr := v.Get("end")
	maxSizeStr := v.Get("maxSize")
	fromStr := v.Get("from")

	// Parse "start" and "end" into time objects
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		l.Errorf("error parsing start time '%s': %v", startStr, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		l.Errorf("error parsing end time '%s': %v", endStr, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	availableRows := 0

	// parse maxSize from query parameter
	maxSize, err := strconv.ParseInt(maxSizeStr, 10, 32)
	if err != nil {
		l.Error("error parsing max size string: ", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GeneralResponse{
			Success: false,
			Message: "Could not parse max size string",
		})
		return
	}

	// make sure maxSize is greater than 0
	if maxSize <= 0 {
		l.Error("invalid max size: ", maxSize)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GeneralResponse{
			Success: false,
			Message: "Invalid max size, must be greater than 0",
		})
		return
	}

	// make sure maxSize is under or equal to max allowed value
	if maxSize > maxCards {
		l.Error("invalid max size: ", maxSize)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GeneralResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid max size, must be under %d", maxCards),
		})
		return
	}

	// parse 'from' from query parameter
	from, err := strconv.ParseInt(fromStr, 10, 32)
	if err != nil {
		l.Error("error parsing 'from' string: ", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(GeneralResponse{
			Success: false,
			Message: "Could not parse 'from' string",
		})
		return
	}

	// get data for the specified fields in the specified time range, sorted by timestamp
	data, availableRows, err := elasticsearch.GetAlarms(s, indices, sources, start, end, int(maxSize), int(from))
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
