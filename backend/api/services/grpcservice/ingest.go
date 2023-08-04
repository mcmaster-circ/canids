// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package grpcservice provides gRPC streaming ingestion services.
package grpcservice

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/state"
	"github.com/sirupsen/logrus"
)

// ingest is triggered from the gRPC service. It indexes a chunk into the database.
func (s *IngestServer) ingest(frame *Frame) {
	// iterate over entries (lines)
	for _, entry := range frame.Payload {
		indexTime := time.Now().Format("2006-01-02")

		// dynamically fetch timestamp and inject GeoIP before indexing
		updated, ts, alarm := dynamicInjection(s.state, entry)
		if ts != nil {
			// index under packet capture time
			indexTime = ts.Format("2006-01-02")
		}

		// index the alarm if one was generated
		if len(alarm) > 0 {
			alarmIndex := fmt.Sprintf("data-%s.alarm-%s-%s", getElasticIndex(frame.FileName), frame.AssetID, indexTime)
			_, err := elasticsearch.IndexPayload(s.state, alarmIndex, alarm)
			if err != nil {
				s.state.Log.WithFields(logrus.Fields{
					"file_name": frame.FileName,
					"asset_id":  frame.AssetID,
					"index":     alarmIndex,
				}).Errorf("failed to index alarm: %s", err)
			}
		}

		// inject the possibly updated payload
		dataIndex := fmt.Sprintf("data-%s-%s-%s", getElasticIndex(frame.FileName), frame.AssetID, indexTime)
		_, err := elasticsearch.IndexPayload(s.state, dataIndex, updated)
		if err != nil {
			s.state.Log.WithFields(logrus.Fields{
				"file_name": frame.FileName,
				"asset_id":  frame.AssetID,
				"index":     dataIndex,
			}).Errorf("failed to index payload `%s`: %s", string(updated), err)
		}
	}
}

// TODO(Jon): update comment
// dynamicInjection will dynamically read the provided JSON record. Known fields
// (ip addresses, timestamps) will be parsed. If the JSON record can be parsed,
// a new JSON byte string will be returned containing GeoIP data, and the JSON
// timestamp will be returned. If the JSON record cannot be parsed, the existing
// record will be returned and the time will be nil.
func dynamicInjection(s *state.State, raw []byte) ([]byte, *time.Time, []byte) {
	// unmarshal data using general interface
	payload := make(map[string]interface{})
	err := json.Unmarshal(raw, &payload)
	if err != nil {
		return raw, nil, []byte{}
	}

	// fetch the RFC3339 string
	rawTime, ok := payload["timestamp"]
	if !ok {
		return raw, nil, []byte{}
	}
	timestring, ok := rawTime.(string)
	if !ok {
		return raw, nil, []byte{}
	}

	// parse timestamp
	t, err := time.Parse(time.RFC3339, timestring)
	if err != nil {
		return raw, nil, []byte{}
	}

	// remove "." from keys (Elasticsearch conflict)
	for key, val := range payload {
		if strings.Contains(key, ".") {
			newKey := strings.ReplaceAll(key, ".", "_")
			payload[newKey] = val
			delete(payload, key)
		}
	}

	// alarm test results
	var sourceIPPositive []string
	var sourceIPNegative []string
	var destIPPositive []string
	var destIPNegative []string

	// generate data from source IP address
	var sourceIP string
	hasSourceIP := false

	if rawSourceIP, hasRawSourceIP := payload["id_orig_h"]; hasRawSourceIP {
		sourceIP, hasSourceIP = rawSourceIP.(string)
	}

	if hasSourceIP {
		// inject GeoIP data from IP
		payload["id_orig_h_asn"] = geoIPASN(s, sourceIP)
		payload["id_orig_h_city"] = geoIPCity(s, sourceIP)
		payload["id_orig_h_country"] = geoIPCountry(s, sourceIP)

		// test against alarm ip sets
		sourceIPPositive, sourceIPNegative = s.AlarmManager.TestIP(sourceIP)
	}

	// generate data from destination IP address
	var destIP string
	hasDestIP := false

	if rawDestIP, hasRawDestIP := payload["id_resp_h"]; hasRawDestIP {
		destIP, hasDestIP = rawDestIP.(string)
	}

	if hasDestIP {
		// inject GeoIP data from IP
		payload["id_resp_h_asn"] = geoIPASN(s, destIP)
		payload["id_resp_h_city"] = geoIPCity(s, destIP)
		payload["id_resp_h_country"] = geoIPCountry(s, destIP)

		// test against alarm ip sets
		destIPPositive, destIPNegative = s.AlarmManager.TestIP(destIP)
	}

	// generate alarm payload if an alarm ip set matched
	var alarmPayload []byte
	if (len(sourceIPPositive) > 0) || (len(destIPPositive) > 0) {
		alarmFields := make(map[string]interface{})
		for key, val := range payload {
			alarmFields[key] = val
		}

		alarmFields["id_orig_h_pos"] = sourceIPPositive
		alarmFields["id_orig_h_neg"] = sourceIPNegative
		alarmFields["id_resp_h_pos"] = destIPPositive
		alarmFields["id_resp_h_neg"] = destIPNegative

		var err error
		alarmPayload, err = json.Marshal(alarmFields)
		if err != nil {
			alarmPayload = []byte{}
		}
	}

	/*
		// generate data from IP address fields
		for _, field := range []string{"id_orig_h", "id_resp_h"} {
			rawVal, ok := payload[field]
			if !ok {
				// field does not exist
				continue
			}
			rawIP, ok := rawVal.(string)
			if !ok {
				// field is not valid string
				continue
			}
			// inject GeoIP data from IP
			payload[field+"_asn"] = geoIPASN(s, rawIP)
			payload[field+"_city"] = geoIPCity(s, rawIP)
			payload[field+"_country"] = geoIPCountry(s, rawIP)
		}
	*/

	// regenerate JSON to capture injected fields
	data, err := json.Marshal(payload)
	if err != nil {
		return raw, &t, alarmPayload
	}

	return data, &t, alarmPayload
}

func getElasticIndex(fileName string) string {
	return strings.Replace(fileName, "-", "_", -1)
}
