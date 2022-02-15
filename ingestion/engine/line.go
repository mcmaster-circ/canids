// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package engine

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"
)

// parseLine accepts a log entry. It parses the TSV log into a JSON byte string.
// For JSON logs, it validates the JSON byte string. It returns a JSON byte
// string or an error.
func parseLine(entry string, h *header) ([]byte, error) {
	if h == nil {
		// no header present, just validate JSON
		return validateJSON(entry)
	}
	// header present, validate JSON
	return processTSV(entry, h)
}

// validateJSON ensures the provided JSON line is valid. It return the byte representation or an error.
func validateJSON(entry string) ([]byte, error) {
	// validate JSON
	payload := []byte(entry)
	if json.Valid(payload) {
		return payload, nil
	}
	// bad JSON
	return nil, errBadJSON
}

// processTSV will parse a TSV file into JSON. It returns an error if the TSV headers does not match the available data.
func processTSV(entry string, h *header) ([]byte, error) {
	// get different columns
	columns := strings.Split(entry, h.separator)
	// data output map
	data := make(map[string]interface{})

	// iterate through all data column
	for i, column := range columns {
		field := h.fields[i]
		// if unset set, empty array
		if strings.Contains(h.types[i], "set") && (column == h.unsetField || column == h.emptyField) {
			data[field] = []string{}
			continue
		}
		// select on field type
		switch h.types[i] {
		case "time":
			// parse timestamp
			tsFloat, err := strconv.ParseFloat(column, 64)
			if err == nil {
				sec, dec := math.Modf(tsFloat)
				timestamp := time.Unix(int64(sec), int64(dec*(1e9)))
				data["timestamp"] = timestamp.Format(time.RFC3339)
			}
		case "port", "count", "int":
			// part integer
			val, err := strconv.ParseInt(column, 10, 64)
			if err == nil {
				data[field] = val
			} else {
				data[field] = column
			}
		case "interval", "double":
			// parse float
			val, err := strconv.ParseFloat(column, 64)
			if err == nil {
				data[field] = val
			} else {
				data[field] = column
			}
		case "bool":
			// parse boolean
			val, err := strconv.ParseBool(column)
			if err == nil {
				data[field] = val
			} else {
				data[field] = column
			}
		case "set[string]":
			// parse list of strings
			data[field] = strings.Split(column, h.setSeperator)
		default:
			if column != h.unsetField {
				data[field] = column
			} else {
				data[field] = ""
			}
		}
		// if field is still unset, set to null
		if data[field] == h.unsetField {
			data[field] = nil
		}
	}
	return json.Marshal(data)
}
