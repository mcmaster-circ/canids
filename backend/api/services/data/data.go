// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package data provides the data API service for the backend.
package data

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
	FieldNames    []string        `json:"fieldNames"`
	Fields        []string        `json:"fields"`
	Class         string          `json:"class"`
	Data          [][]interface{} `json:"data"`
	AvailableRows int             `json:"availableRows"`
}

const maxTableRows = 100

// dataHandler is "/api/data. It is responsible for populating a view with the data related to that view.
func dataHandler(ctx context.Context, s *state.State, a *auth.State, w http.ResponseWriter, r *http.Request) {
	// get user making current request + logging context
	_, l := jwtauth.FromContext(ctx), ctxlog.Log(ctx)
	w.Header().Set("Content-Type", "application/json")

	// get query parameters "start" and "end"
	v := r.URL.Query()
	visualizationUUID := v.Get("view")
	startStr := v.Get("start")
	endStr := v.Get("end")
	intervalStr := v.Get("interval")
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

	// get view from database
	view, _, err := elasticsearch.QueryViewByUUID(s, visualizationUUID)
	if err != nil {
		l.Errorf("error getting view with UUID '%s': %v", visualizationUUID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}

	// generate indexName to query
	indexName := "data-" + view.DataIndex + "-"

	// Get data in whatever way the given view class requires
	data := [][]interface{}{}
	availableRows := 0
	if (view.Class == elasticsearch.ViewBar) || (view.Class == elasticsearch.ViewPie) {
		// check that the view has the right amount of fields for this class
		if len(view.Fields) != 1 {
			l.Errorf("%s view: expected 1 fields, got %d", view.Class, len(view.Fields))
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}

		var keys []string
		var counts []int64

		if view.Name == fmt.Sprintf("%s %s", elasticsearch.DefaultViewName, view.Authorized) {
			keys, counts, err = elasticsearch.CountTotalDataInRange(s, view.Authorized, view.Fields[0], start, end)
		} else {
			keys, counts, err = elasticsearch.CountDataInRange(s, indexName, view.Authorized, view.Fields[0], start, end)
		}
		if err != nil {
			l.Error("error querying data conn: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}

		data = make([][]interface{}, 2)
		data[0] = []interface{}{}
		data[1] = []interface{}{}

		for _, key := range keys {
			data[0] = append(data[0], key)
		}
		for _, count := range counts {
			data[1] = append(data[1], count)
		}
	} else if view.Class == elasticsearch.ViewLine {
		// check that the view has the right amount of fields for this class
		if len(view.Fields) != 2 {
			l.Errorf("%s view: expected 2 fields, got %d", view.Class, len(view.Fields))
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}

		// parse interval from query parameter
		interval, err := strconv.ParseInt(intervalStr, 10, 64)
		if err != nil {
			l.Error("error parsing interval string: ", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GeneralResponse{
				Success: false,
				Message: "Could not parse interval string",
			})
			return
		}

		// make sure interval is greater than 0
		if interval <= 0 {
			l.Error("invalid interval: ", interval)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GeneralResponse{
				Success: false,
				Message: "Invalid interval, must be greater than 0",
			})
			return
		}

		// make sure is not over 10000 buckets
		buckets := (end.Unix() - start.Unix()) / interval
		if buckets > 10000 {
			l.Error("Too many buckets: ", buckets)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GeneralResponse{
				Success: false,
				Message: "Interval is too small for this time range",
			})
			return
		}

		// get data for the specified fields in the specified time range
		// TODO(Tanner)
		xdata, ydata, err := elasticsearch.QueryDataInRangeAggregated(s, indexName, view.Authorized, view.Fields[0], view.Fields[1], start, end, interval)
		if err != nil {
			l.Error("error querying data conn: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}

		data = [][]interface{}{
			xdata, ydata,
		}
	} else if view.Class == elasticsearch.ViewTable {
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
		if maxSize > maxTableRows {
			l.Error("invalid max size: ", maxSize)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(GeneralResponse{
				Success: false,
				Message: fmt.Sprintf("Invalid max size, must be under %d", maxTableRows),
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
		// TODO(Tanner)
		data, availableRows, err = elasticsearch.QueryDataInRange(s, indexName, view.Authorized, view.Fields, start, end, int(maxSize), int(from))
		if err != nil {
			l.Error("error querying data conn: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(InternalServerError)
			return
		}
	}

	// success
	l.Info("successfully queried data for asset")
	json.NewEncoder(w).Encode(dataResponse{
		FieldNames:    view.FieldNames,
		Fields:        view.Fields,
		Class:         string(view.Class),
		Data:          data,
		AvailableRows: availableRows,
	})
}
