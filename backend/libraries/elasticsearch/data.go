// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package elasticsearch provides the simplified Elasticsearch interface.
package elasticsearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// index format : data-index-assetID-yyyy-mm-dd

// IndexDataField contains a list of fields per index.
type IndexDataField struct {
	Index  string      `json:"index"`
	Fields []DataField `json:"fields"`
}

// DataField contains a field name with field type.
type DataField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

var alarmFields = []string{"uid", "host", "timestamp", "id_orig_h", "id_orig_p", "id_orig_h_pos", "id_resp_h", "id_resp_p", "id_resp_h_pos"}

// Alarm contains the data for an alarm.
type Alarm struct {
	UID               string   `json:"uid"`
	Host              string   `json:"host"`
	Timestamp         string   `json:"timestamp"`
	SourceIP          string   `json:"id_orig_h"`
	SourcePort        int      `json:"id_orig_p"`
	SourceAlarms      []string `json:"id_orig_h_pos"`
	DestinationIP     string   `json:"id_resp_h"`
	DestinationPort   int      `json:"id_resp_p"`
	DestinationAlarms []string `json:"id_resp_h_pos"`
}

// IndexPayload attempts to index the provided payload under the index name. It
// will return the newly created document ID or an error.
func IndexPayload(s *state.State, indexName string, payload []byte) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	result, err := client.Index(indexName).Raw(bytes.NewReader(payload)).Do(ctx)
	if err != nil {
		return "", err
	}
	return result.Id_, nil
}

// GetAllDataMapping will fetch all data mappings. It returns a list of
// fields+type for each index type, or an error.
func GetAllDataMapping(s *state.State) ([]IndexDataField, error) {
	out := make([]IndexDataField, 0)

	// fetch all indexes
	indexes, err := GetIndexes(s)
	if err != nil {
		return nil, err
	}

	// generate list of prefixes
	prefixes := make(map[string]bool)
	for _, index := range indexes {
		parts := strings.Split(index, "-")
		// pattern data-fileName-assetID-yyyy-mm-dd
		if len(parts) != 4 {
			// ingore non-data files
			continue
		}
		prefix := fmt.Sprintf("data-%s", parts[1])
		// check if in prefix map
		_, ok := prefixes[prefix]
		if !ok {
			// add to prefix map
			prefixes[prefix] = true
		}
	}

	// iterate over all prefixes
	for index := range prefixes {
		field, err := GetDataMapping(s, index)
		if err != nil {
			// index does not have a document, skip
			continue
		}
		out = append(out, IndexDataField{
			Index:  strings.Split(index, "-")[1],
			Fields: field,
		})
	}
	return out, nil
}

// GetIndexes queries for a list of all indexes. Returns list of indicies or error.
func GetIndexes(s *state.State) ([]string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	indexes, err := client.Cat.Indices().Do(ctx)
	if err != nil {
		return nil, err
	}
	count := len(indexes)
	if count == 0 {
		return []string{}, nil
	}
	out := make([]string, count)
	for i, index := range indexes {
		out[i] = *index.Index
	}
	return out, nil
}

// GetDataMapping queries for the latest document and fetches the mapping for
// the document. It returns a list of field names and types in the mapping or an
// error.
func GetDataMapping(s *state.State, indexPrefix string) ([]DataField, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// Get latest doc
	latestDoc, err := client.Search().Index(indexPrefix + "*").
		Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		}).
		Sort(types.SortOptions{ // sort timestamp descending
			SortOptions: map[string]types.FieldSort{
				"timestamp": {
					Order: &sortorder.Desc,
				},
			},
		}).Size(1).Do(ctx)
	if err != nil {
		return []DataField{}, err
	}
	// ensure we got the 1 doc we requested
	if len(latestDoc.Hits.Hits) != 1 {
		return []DataField{}, errors.New(fmt.Sprintf("GetDataConnMapping: Expected 1 hit, got %d", len(latestDoc.Hits.Hits)))
	}
	// get the mapping for the index of the document found above
	mappingQuery, err := client.Indices.GetMapping().Index(latestDoc.Hits.Hits[0].Index_).Do(ctx)
	if err != nil {
		return []DataField{}, err
	}
	// Get the properties map from the query result above
	index, hasIndex := mappingQuery[latestDoc.Hits.Hits[0].Index_]
	if !hasIndex {
		return []DataField{}, errors.New("GetDataConnMapping: mapping query doesn't have 'index' field")
	}
	properties := index.Mappings.Properties

	// put the field names into an array
	fields := []DataField{}
	for propertyName, property := range properties {
		blob, _ := json.Marshal(property)
		var m map[string]interface{}
		_ = json.Unmarshal(blob, &m) // this is patently insane but they give me no better way to do it
		propertyType := m["type"].(string)
		fields = append(fields, DataField{
			Name: propertyName,
			Type: propertyType,
		})
	}
	return fields, nil
}

// ListDataAssets queries all indexes to fetch the asset names. It returns a
// list of assets or an error.
func ListDataAssets(s *state.State) ([]string, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// query for all index names
	indicesQuery, err := client.Cat.Indices().Do(ctx)
	if err != nil {
		return []string{}, err
	}
	// split index names to get the asset names and add them to a set
	assetNameSet := make(map[string]bool)
	for _, index := range indicesQuery {
		if strings.HasPrefix(*index.Index, "data-") {
			splitIndexName := strings.Split(*index.Index, "-")
			if len(splitIndexName) == 4 {
				assetName := splitIndexName[2]
				assetNameSet[assetName] = true
			} else {
				//TODO(Jon): error?
			}
		}
	}
	// get all names from the set and add them to an array
	result := []string{}
	for assetName := range assetNameSet {
		result = append(result, assetName)
	}
	return result, nil
}

// get alarms for a given asset in a given time range from a
func GetAlarms(s *state.State, indices []string, sources []string, destinations []string, start time.Time, end time.Time, size int, from int, sourceIP string, destIP string) ([]Alarm, int, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// return empty array if no sources or indices
	if len(sources) == 0 || len(indices) == 0 {
		return []Alarm{}, 0, nil
	}

	alarmSources := make([]interface{}, len(sources))
	for i, source := range sources {
		alarmSources[i] = source
	}

	alarmDestinations := make([]interface{}, len(destinations))
	for i, destination := range destinations {
		alarmDestinations[i] = destination
	}

	for i, index := range indices {
		indices[i] = fmt.Sprintf("data-%s-*", index)
	}

	r := types.Query{
		Range: map[string]types.RangeQuery{
			"timestamp": types.DateRangeQuery{
				From: start.Format(time.RFC3339),
				To:   end.Format(time.RFC3339),
			},
		},
	}

	origSources := types.Query{
		Terms: &types.TermsQuery{
			TermsQuery: map[string]types.TermsQueryField{
				"id_orig_h_pos": sources,
			},
		},
	}

	respSources := types.Query{
		Terms: &types.TermsQuery{
			TermsQuery: map[string]types.TermsQueryField{
				"id_resp_h_pos": destinations,
			},
		},
	}

	sourceIPQuery := types.Query{
		MatchPhrasePrefix: map[string]types.MatchPhrasePrefixQuery{
			"id_orig_h": {Query: sourceIP},
		},
	}

	destIPQuery := types.Query{
		MatchPhrasePrefix: map[string]types.MatchPhrasePrefixQuery{
			"id_resp_h": {Query: destIP},
		},
	}

	// query for all alarms in range and filter for either origSource or respSource being in alarmSources
	query := &types.Query{
		Bool: &types.BoolQuery{
			Must: []types.Query{
				r,
				origSources,
				respSources,
				sourceIPQuery,
				destIPQuery,
			},
		},
	}
	queryResult, err := client.Search().Index(strings.Join(indices, ",")).
		Query(query).Sort(types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			"timestamp": {
				Order: &sortorder.Desc,
			},
		},
	}).Size(size).From(from).Do(ctx)

	if err != nil {
		return []Alarm{}, 0, err
	}

	alarms := make([]Alarm, 0, len(queryResult.Hits.Hits))

	// loop through each alarm and unmarshal it into an Alarm struct
	for _, hit := range queryResult.Hits.Hits {
		var alarm Alarm
		err = json.Unmarshal(hit.Source_, &alarm)
		if err != nil {
			return alarms, 0, err
		}
		alarms = append(alarms, alarm)
	}

	return alarms, int(queryResult.Hits.Total.Value), nil
}

func QueryDataInRangeAggregated(s *state.State, indexPrefix string, xField string, yField string, start time.Time, end time.Time, interval int64) ([]interface{}, []interface{}, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// query for docs in the given time range
	query := &types.Query{
		Range: map[string]types.RangeQuery{
			"timestamp": types.DateRangeQuery{
				From: start.Format(time.RFC3339),
				To:   end.Format(time.RFC3339),
			},
		},
	}

	// aggregate time buckets given by interval (in seconds), average xfield and
	// yfield for each bucket

	ts := "timestamp"
	duration := fmt.Sprintf("%ds", interval)
	aggregation := types.DateHistogramAggregation{
		Field:         &ts,
		FixedInterval: duration,
	}

	aggX := types.AverageAggregation{
		Field: &xField,
	}
	aggY := types.AverageAggregation{
		Field: &yField,
	}

	// do query
	indexName := fmt.Sprintf("%s-*", indexPrefix)
	queryResult, err := client.Search().Index(indexName).Query(query).Size(0).Aggregations(map[string]types.Aggregations{
		"aggT": {
			DateHistogram: &aggregation,
			Aggregations: map[string]types.Aggregations{
				"aggX": {
					Avg: &aggX,
				},
				"aggY": {
					Avg: &aggY,
				},
			},
		},
	}).Do(ctx)

	if err != nil {
		return []interface{}{}, []interface{}{}, err
	}

	// get time histogram aggregation
	aggT, foundAggT := queryResult.Aggregations["aggT"].(*types.DateHistogramAggregate)
	if !foundAggT {
		// no aggT date histogram found, this probably mean the asset doesnt
		// have any indices yet
		return []interface{}{}, []interface{}{}, nil
	}

	// create arrays for the x and y data
	xresult := []interface{}{}
	yresult := []interface{}{}

	// process buckets from time aggregation
	for _, bucket := range aggT.Buckets.([]types.DateHistogramBucket) {
		// get x & y avg aggregations
		aggX, foundAggX := bucket.Aggregations["aggX"].(*types.AvgAggregate)
		aggY, foundAggY := bucket.Aggregations["aggY"].(*types.AvgAggregate)

		if foundAggX && foundAggY {
			// either get the averaged value or the date string from the bucket
			// if the field is 'timestamp'
			if xField == "timestamp" {
				xresult = append(xresult, bucket.KeyAsString)
			} else {
				xresult = append(xresult, aggX.Value)
			}

			if yField == "timestamp" {
				xresult = append(xresult, bucket.KeyAsString)
			} else {
				yresult = append(yresult, aggY.Value)
			}
		}
	}

	return xresult, yresult, nil
}

// QueryDataInRange queries the specified asset for all fields specified,
// returns an array of data for each field
func QueryDataInRange(s *state.State, indexPrefix string, fields []string, start time.Time, end time.Time, size int, from int) ([][]interface{}, int, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// query for all data conn documents for this asset in the given timerange,
	// sorted in descending time
	indexName := fmt.Sprintf("%s-*", indexPrefix)
	queryResult, err := client.Search().Index(indexName).
		Query(&types.Query{
			Range: map[string]types.RangeQuery{
				"timestamp": types.DateRangeQuery{
					From: start.Format(time.RFC3339),
					To:   end.Format(time.RFC3339),
				},
			},
		}).
		Sort(types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"timestamp": {
					Order: &sortorder.Desc,
				},
			},
		}).Size(size).From(from).Do(ctx)
	if err != nil {
		return [][]interface{}{}, 0, err
	}

	// create result array with an array for each field
	result := make([][]interface{}, len(fields))

	for i := range fields {
		result[i] = []interface{}{}
	}

	// unmarshal elasticsearch hits
	for _, hit := range queryResult.Hits.Hits {
		var d map[string]json.RawMessage
		err = json.Unmarshal(hit.Source_, &d)
		if err != nil {
			return result, 0, err
		}

		for i, field := range fields {
			result[i] = append(result[i], d[field])
		}
	}

	return result, int(queryResult.Hits.Total.Value), nil
}

func CountDataInRange(s *state.State, indexPrefix, field string, start time.Time, end time.Time) ([]string, []int64, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// Get the mapping
	mapping, err := GetDataMapping(s, indexPrefix)
	if err != nil {
		return []string{}, []int64{}, err
	}

	// find the type of this field
	fieldType := ""
	for _, fieldProp := range mapping {
		if fieldProp.Name == field {
			fieldType = fieldProp.Type
		}
	}

	if fieldType == "text" {
		field = fmt.Sprintf("%s.keyword", field)
	}

	agg := types.TermsAggregation{
		Field: &field,
	}

	// query for all data conn documents for this asset in the given timerange, sorted in ascending time
	indexName := fmt.Sprintf("%s-*", indexPrefix)
	queryResult, err := client.Search().Index(indexName).
		Query(&types.Query{
			Range: map[string]types.RangeQuery{
				"timestamp": types.DateRangeQuery{
					From: start.Format(time.RFC3339),
					To:   end.Format(time.RFC3339),
				},
			},
		}).Aggregations(map[string]types.Aggregations{
		"count": {
			Terms: &agg,
		},
	}).Do(ctx)
	if err != nil {
		return []string{}, []int64{}, err
	}

	keys := []string{}
	counts := []int64{}

	termsAgg, found := queryResult.Aggregations["count"].(*types.StringTermsAggregate)
	if !found {
		// no count terms aggregation found, this probably mean the asset doesnt have any indices yet
		return []string{}, []int64{}, nil
	}

	// unmarshal elasticsearch hits
	for _, bucket := range termsAgg.Buckets.([]types.StringTermsBucket) {
		key := fmt.Sprintf("%v", bucket.Key)
		keys = append(keys, key)
		counts = append(counts, bucket.DocCount)
	}

	return keys, counts, nil
}

func CountTotalDataInRange(s *state.State, field string, start time.Time, end time.Time) ([]string, []int64, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// Create Range Aggregation
	agg := types.RangeAggregation{
		Field: &field,
	}
	numOfBars := 10
	daysPerBar := int(end.Sub(start).Hours()/24) / numOfBars
	if daysPerBar == 0 { // If the amount of days is less than 1 default to 1 bar with the current day
		numOfBars, daysPerBar = 1, 1
	} else if daysPerBar < numOfBars { // If there are less days per bar than number of bars set, change the number of bars to amount of days per bar
		numOfBars = daysPerBar
	}
	for i := 0; i < numOfBars; i++ {
		rangeStart := start.AddDate(0, 0, i*daysPerBar)
		rangeEnd := start.AddDate(0, 0, (i+1)*daysPerBar)
		if rangeEnd.After(end) || (i+1) == numOfBars {
			rangeEnd = end
		}
		agg.Ranges = append(agg.Ranges, types.AggregationRange{
			From: rangeStart.Format(time.RFC3339),
			To:   rangeEnd.Format(time.RFC3339),
		})
	}

	// query for all data conn documents for this asset in the given timerange, sorted in ascending time
	indexName := "data-*"
	queryResult, err := client.Search().Index(indexName).
		Query(&types.Query{
			Range: map[string]types.RangeQuery{
				"timestamp": types.DateRangeQuery{
					From: start.Format(time.RFC3339),
					To:   end.Format(time.RFC3339),
				},
			},
		}).
		Aggregations(map[string]types.Aggregations{
			"count": {
				Range: &agg,
			},
		}).
		Do(ctx)
	if err != nil {
		return []string{}, []int64{}, err
	}

	keys := []string{}
	counts := []int64{}

	termsAgg, found := queryResult.Aggregations["count"].(*types.RangeAggregate)
	if !found {
		// no count terms aggregation found, this probably mean the asset doesnt have any indices yet
		return []string{}, []int64{}, nil
	}

	timeFormat := "02 Jan 2006"
	// unmarshal elasticsearch hits
	for i, bucket := range termsAgg.Buckets.([]types.RangeBucket) {
		// Create a condensed key name for ease of viewing
		startRange := start.AddDate(0, 0, i*daysPerBar)
		endRange := start.AddDate(0, 0, (i+1)*daysPerBar)
		if (i + 1) == numOfBars {
			endRange = end
		}
		key := fmt.Sprintf("%s - %s", startRange.Format(timeFormat), endRange.Format(timeFormat))
		keys = append(keys, key)
		counts = append(counts, bucket.DocCount)
	}

	return keys, counts, nil
}
