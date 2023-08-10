package websocket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/state"
	"github.com/sirupsen/logrus"
)

// Active index, number of entries
var activeIndices = map[string]int{}
var activeAlarmIndices = map[string]int{}

var maxSize = 1000000

// ingest is triggered from the gRPC service. It indexes a chunk into the database.
func ingest(frame *Frame, state *state.State, maxIndexSize int) {
	maxSize = maxIndexSize
	client, ctx := state.Elastic, state.ElasticCtx
	for _, entry := range frame.Payload {

		updated, _, alarm := dynamicInjection(state, entry)

		var selectedIndex string

		// Check if there are locally stored indices
		if len(activeIndices) > 0 {
			for index := range activeIndices {
				arr := strings.Split(index, "-")
				if len(arr) > 1 {
					//If found locally
					if arr[1] == frame.FileName {
						//Delete and increment index num if too many elements
						if activeIndices[index] >= maxSize {
							delete(activeIndices, index)
							indexNum, err := strconv.Atoi(arr[3])
							if err != nil {
								state.Log.WithFields(logrus.Fields{
									"file_name": frame.FileName,
									"index":     index,
								}).Errorf("strconv error: %s", err)
							} else {
								indexNum += 1
							}
							selectedIndex = fmt.Sprintf("data-%s-%s-%d", getElasticIndex(frame.FileName), frame.AssetID, indexNum)
							activeIndices[selectedIndex] = 0
						} else {
							selectedIndex = index
						}
					}
				} else {
					state.Log.WithFields(logrus.Fields{
						"file_name": frame.FileName,
						"index":     index,
					}).Errorf("Index in incorrect format")
				}
			}
		}

		// Not found locally
		if selectedIndex == "" {
			elasticIndices, err := client.IndexGet(fmt.Sprintf("data-%s-%s-*", getElasticIndex(frame.FileName), frame.AssetID)).Do(ctx)
			if err != nil {
				state.Log.WithFields(logrus.Fields{
					"file_name": frame.FileName,
					"asset_id":  frame.AssetID,
				}).Errorf("Error getting indices from elasticsearch: %s", err)
			} else {
				highestIndexNum := 0
				// Loop through indices from es that match
				for index, _ := range elasticIndices {
					arr := strings.Split(index, "-")
					if len(arr) > 1 {
						num, err := strconv.Atoi(arr[3])
						if err != nil {
							state.Log.WithFields(logrus.Fields{
								"file_name": frame.FileName,
								"index":     index,
							}).Errorf("strconv error: %s", err)
						} else {
							if num > highestIndexNum {
								highestIndexNum = num
							}
						}
					} else {
						state.Log.WithFields(logrus.Fields{
							"file_name": frame.FileName,
							"index":     index,
						}).Errorf("Index in incorrect format")
					}
				}
				//If found on es, set to highest number index
				if highestIndexNum > 0 {
					selectedIndex = fmt.Sprintf("data-%s-%s-%d", getElasticIndex(frame.FileName), frame.AssetID, highestIndexNum)

					currentSize, err := client.Count(selectedIndex).Do(ctx)
					if err != nil {
						state.Log.WithFields(logrus.Fields{
							"file_name": frame.FileName,
							"index":     selectedIndex,
						}).Errorf("Error getting current size of index: %s", err)
					} else {
						// Size check
						if currentSize < int64(maxSize) {
							activeIndices[selectedIndex] = int(currentSize)
						} else {
							highestIndexNum += 1
							selectedIndex = fmt.Sprintf("data-%s-%s-%d", getElasticIndex(frame.FileName), frame.AssetID, highestIndexNum)
							activeIndices[selectedIndex] = 0
						}
					}
				} else { //Doesnt exist anywhere
					selectedIndex = fmt.Sprintf("data-%s-%s-%d", getElasticIndex(frame.FileName), frame.AssetID, 1)
					activeIndices[selectedIndex] = 0
				}
			}
		}

		//Alarm
		if len(alarm) > 0 {
			var selectedAlarmIndex string

			// Check if there are locally stored indices
			if len(activeAlarmIndices) > 0 {
				for index := range activeAlarmIndices {
					arr := strings.Split(index, "-")
					if len(arr) > 1 {
						//If found locally
						if arr[1] == (frame.FileName + ".alarm") {
							//Delete and increment index num if too many elements
							if activeAlarmIndices[index] >= maxSize {
								delete(activeAlarmIndices, index)
								indexNum, err := strconv.Atoi(arr[3])
								if err != nil {
									state.Log.WithFields(logrus.Fields{
										"file_name": frame.FileName + ".alarm",
										"index":     index,
									}).Errorf("strconv error: %s", err)
								} else {
									indexNum += 1
								}
								selectedAlarmIndex = fmt.Sprintf("data-%s.alarm-%s-%d", getElasticIndex(frame.FileName), frame.AssetID, indexNum)
								activeAlarmIndices[selectedAlarmIndex] = 0
							} else {
								selectedAlarmIndex = index
							}
						}
					} else {
						state.Log.WithFields(logrus.Fields{
							"file_name": frame.FileName + ".alarm",
							"index":     index,
						}).Errorf("Index in incorrect format")
					}
				}
			}

			// Not found locally
			if selectedAlarmIndex == "" {
				elasticIndices, err := client.IndexGet(fmt.Sprintf("data-%s.alarm-%s-*", getElasticIndex(frame.FileName), frame.AssetID)).Do(ctx)
				if err != nil {
					state.Log.WithFields(logrus.Fields{
						"file_name": frame.FileName + ".alarm",
						"asset_id":  frame.AssetID,
					}).Errorf("Error getting indices from elasticsearch: %s", err)
				} else {
					highestIndexNum := 0
					// Loop through indices from es that match
					for index, _ := range elasticIndices {
						arr := strings.Split(index, "-")
						if len(arr) > 1 {
							num, err := strconv.Atoi(arr[3])
							if err != nil {
								state.Log.WithFields(logrus.Fields{
									"file_name": frame.FileName + ".alarm",
									"index":     index,
								}).Errorf("strconv error: %s", err)
							} else {
								if num > highestIndexNum {
									highestIndexNum = num
								}
							}
						} else {
							state.Log.WithFields(logrus.Fields{
								"file_name": frame.FileName + ".alarm",
								"index":     index,
							}).Errorf("Index in incorrect format")
						}
					}
					//If found on es, set to highest number index
					if highestIndexNum > 0 {
						selectedAlarmIndex = fmt.Sprintf("data-%s.alarm-%s-%d", getElasticIndex(frame.FileName), frame.AssetID, highestIndexNum)

						currentSize, err := client.Count(selectedAlarmIndex).Do(ctx)
						if err != nil {
							state.Log.WithFields(logrus.Fields{
								"file_name": frame.FileName,
								"index":     selectedAlarmIndex,
							}).Errorf("Error getting current size of index: %s", err)
						} else {
							// Size check
							if currentSize < int64(maxSize) {
								activeAlarmIndices[selectedAlarmIndex] = int(currentSize)
							} else {
								highestIndexNum += 1
								selectedAlarmIndex = fmt.Sprintf("data-%s.alarm-%s-%d", getElasticIndex(frame.FileName), frame.AssetID, highestIndexNum)
								activeAlarmIndices[selectedAlarmIndex] = 0
							}
						}
					} else { //Doesnt exist anywhere
						selectedAlarmIndex = fmt.Sprintf("data-%s.alarm-%s-%d", getElasticIndex(frame.FileName), frame.AssetID, 1)
						activeAlarmIndices[selectedAlarmIndex] = 0
					}
				}
			}

			// index the alarm
			alarmIndex := selectedAlarmIndex
			_, err := elasticsearch.IndexPayload(state, alarmIndex, alarm)
			if err != nil {
				state.Log.WithFields(logrus.Fields{
					"file_name": frame.FileName,
					"asset_id":  frame.AssetID,
					"index":     alarmIndex,
				}).Errorf("failed to index alarm: %s", err)
			}
		}

		// inject the possibly updated payload
		dataIndex := selectedIndex
		_, err := elasticsearch.IndexPayload(state, dataIndex, updated)
		if err != nil {
			state.Log.WithFields(logrus.Fields{
				"file_name": frame.FileName,
				"asset_id":  frame.AssetID,
				"index":     dataIndex,
			}).Errorf("failed to index payload: %s", err)
		} else {
			activeIndices[selectedIndex] = activeIndices[selectedIndex] + 1
		}
		state.Log.Printf("Indexing: %s\n", dataIndex)
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
