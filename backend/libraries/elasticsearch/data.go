// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package elasticsearch provides the simplified Elasticsearch interface.
package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/state"
	"github.com/olivere/elastic"
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
	result, err := client.Index().Index(indexName).BodyString(string(payload)).Do(ctx)
	if err != nil {
		return "", err
	}
	return result.Id, nil
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
		if len(parts) != 6 {
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
	indexes, err := client.CatIndices().Do(ctx)
	if err != nil {
		return nil, err
	}
	count := len(indexes)
	if count == 0 {
		return []string{}, nil
	}
	out := make([]string, count)
	for i, index := range indexes {
		out[i] = index.Index
	}
	return out, nil
}

// GetDataMapping queries for the latest document and fetches the mapping for
// the document. It returns a list of field names and types in the mapping or an
// error.
func GetDataMapping(s *state.State, indexPrefix string) ([]DataField, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// Get latest doc
	latestDoc, err := client.Search().Index(indexPrefix+"*").
		Query(elastic.NewMatchAllQuery()).
		Sort("timestamp", false).Size(1).Do(ctx)
	if err != nil {
		return []DataField{}, err
	}
	// ensure we got the 1 doc we requested
	if len(latestDoc.Hits.Hits) != 1 {
		return []DataField{}, errors.New(fmt.Sprintf("GetDataConnMapping: Expected 1 hit, got %d", len(latestDoc.Hits.Hits)))
	}
	// get the mapping for the index of the document found above
	mappingQuery, err := client.GetMapping().Index(latestDoc.Hits.Hits[0].Index).Do(ctx)
	if err != nil {
		return []DataField{}, err
	}
	// Get the properties map from the query result above
	index, hasIndex := mappingQuery[latestDoc.Hits.Hits[0].Index].(map[string]interface{})
	if !hasIndex {
		return []DataField{}, errors.New("GetDataConnMapping: mapping query doesn't have 'index' field")
	}
	mappings, hasMapping := index["mappings"].(map[string]interface{})
	if !hasMapping {
		return []DataField{}, errors.New("GetDataConnMapping: mapping query doesn't have 'mappings' field")
	}
	properties, hasProperties := mappings["properties"].(map[string]interface{})
	if !hasProperties {
		return []DataField{}, errors.New("GetDataConnMapping: mapping query doesn't have 'properties' field")
	}
	// put the field names into an array
	fields := []DataField{}
	for propertyName, property := range properties {
		propertyType, hasType := (property.(map[string]interface{}))["type"].(string)
		if hasType {
			fields = append(fields, DataField{
				Name: propertyName,
				Type: propertyType,
			})
		}
	}
	return fields, nil
}

// ListDataAssets queries all indexes to fetch the asset names. It returns a
// list of assets or an error.
func ListDataAssets(s *state.State) ([]string, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// query for all index names
	indicesQuery, err := client.CatIndices().Do(ctx)
	if err != nil {
		return []string{}, err
	}
	// split index names to get the asset names and add them to a set
	assetNameSet := make(map[string]bool)
	for _, index := range indicesQuery {
		if strings.HasPrefix(index.Index, "data-") {
			splitIndexName := strings.Split(index.Index, "-")
			if len(splitIndexName) == 6 {
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
func GetAlarms(s *state.State, indices []string, sources []string, start time.Time, end time.Time, size int, from int) ([]Alarm, int, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// return empty array if no sources or indices
	if len(sources) == 0 || len(indices) == 0 {
		return []Alarm{}, 0, nil
	}

	alarmSources := make([]interface{}, len(sources))
	for i, source := range sources {
		alarmSources[i] = source
	}

	for i, index := range indices {
		indices[i] = fmt.Sprintf("data-%s-*", index)
	}

	r := elastic.NewRangeQuery("timestamp").
		From(start.Format(time.RFC3339)).
		To(end.Format(time.RFC3339))
	// query for all alarms in range and filter for either origSource or respSource being in alarmSources
	origSources := elastic.NewTermsQuery("id_orig_h_pos", alarmSources...)
	respSources := elastic.NewTermsQuery("id_resp_h_pos", alarmSources...)
	hasSource := elastic.NewBoolQuery().Should(origSources, respSources)
	query := elastic.NewBoolQuery().Must(r).Must(hasSource)
	queryResult, err := client.Search().Index(indices...).
		Query(query).Sort("timestamp", false).Size(size).From(from).Do(ctx)
	if err != nil {
		return []Alarm{}, 0, err
	}

	alarms := make([]Alarm, 0, len(queryResult.Hits.Hits))

	// loop through each alarm and unmarshal it into an Alarm struct
	for _, hit := range queryResult.Hits.Hits {
		var alarm Alarm
		err = json.Unmarshal(hit.Source, &alarm)
		if err != nil {
			return alarms, 0, err
		}
		alarms = append(alarms, alarm)
	}

	return alarms, int(queryResult.Hits.TotalHits.Value), nil
}

func QueryDataInRangeAggregated(s *state.State, indexPrefix string, xField string, yField string, start time.Time, end time.Time, interval int64) ([]interface{}, []interface{}, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// query for docs in the given time range
	query := elastic.NewRangeQuery("timestamp").From(start.Format(time.RFC3339)).To(end.Format(time.RFC3339))

	// aggregate time buckets given by interval (in seconds), average xfield and
	// yfield for each bucket
	aggregation := elastic.NewDateHistogramAggregation().Field("timestamp").Interval(fmt.Sprintf("%ds", interval)).
		SubAggregation("aggX", elastic.NewAvgAggregation().Field(xField)).
		SubAggregation("aggY", elastic.NewAvgAggregation().Field(yField))

	// do query
	indexName := fmt.Sprintf("%s-*", indexPrefix)
	queryResult, err := client.Search().Index(indexName).Query(query).Size(0).Aggregation("aggT", aggregation).Do(ctx)
	if err != nil {
		return []interface{}{}, []interface{}{}, err
	}

	// get time histogram aggregation
	aggT, foundAggT := queryResult.Aggregations.DateHistogram("aggT")
	if !foundAggT {
		// no aggT date histogram found, this probably mean the asset doesnt
		// have any indices yet
		return []interface{}{}, []interface{}{}, nil
	}

	// create arrays for the x and y data
	xresult := []interface{}{}
	yresult := []interface{}{}

	// process buckets from time aggregation
	for _, bucket := range aggT.Buckets {
		// get x & y avg aggregations
		aggX, foundAggX := bucket.Aggregations.Avg("aggX")
		aggY, foundAggY := bucket.Aggregations.Avg("aggY")

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
		Query(elastic.NewRangeQuery("timestamp").
			From(start.Format(time.RFC3339)).
			To(end.Format(time.RFC3339))).
		Sort("timestamp", false).Size(size).From(from).Do(ctx)
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
		err = json.Unmarshal(hit.Source, &d)
		if err != nil {
			return result, 0, err
		}

		for i, field := range fields {
			result[i] = append(result[i], d[field])
		}
	}

	return result, int(queryResult.Hits.TotalHits.Value), nil
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

	// query for all data conn documents for this asset in the given timerange, sorted in ascending time
	indexName := fmt.Sprintf("%s-*", indexPrefix)
	queryResult, err := client.Search().Index(indexName).
		Query(elastic.NewRangeQuery("timestamp").
			From(start.Format(time.RFC3339)).
			To(end.Format(time.RFC3339))).
		Aggregation("count", elastic.NewTermsAggregation().Field(field)).Do(ctx)
	if err != nil {
		return []string{}, []int64{}, err
	}

	keys := []string{}
	counts := []int64{}

	termsAgg, found := queryResult.Aggregations.Terms("count")
	if !found {
		// no count terms aggregation found, this probably mean the asset doesnt have any indices yet
		return []string{}, []int64{}, nil
	}

	// unmarshal elasticsearch hits
	for _, bucket := range termsAgg.Buckets {
		s.Log.Infof("Bucket %+v", bucket)
		key := fmt.Sprintf("%v", bucket.Key)
		keys = append(keys, key)
		counts = append(counts, bucket.DocCount)
	}

	return keys, counts, nil
}

func CountTotalDataInRange(s *state.State, field string, start time.Time, end time.Time) ([]string, []int64, error) {
	client, ctx := s.Elastic, s.ElasticCtx

	// Create Range Aggregation
	agg := elastic.NewRangeAggregation().Field(field)
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
		agg.AddRange(rangeStart.Format(time.RFC3339), rangeEnd.Format(time.RFC3339))
	}

	// query for all data conn documents for this asset in the given timerange, sorted in ascending time
	indexName := "data-*"
	queryResult, err := client.Search().Index(indexName).
		Query(elastic.NewRangeQuery("timestamp").
			From(start.Format(time.RFC3339)).
			To(end.Format(time.RFC3339))).
		Aggregation("count", agg).
		Do(ctx)
	if err != nil {
		return []string{}, []int64{}, err
	}

	keys := []string{}
	counts := []int64{}

	termsAgg, found := queryResult.Aggregations.Terms("count")
	if !found {
		// no count terms aggregation found, this probably mean the asset doesnt have any indices yet
		return []string{}, []int64{}, nil
	}

	timeFormat := "02 Jan 2006"
	// unmarshal elasticsearch hits
	for i, bucket := range termsAgg.Buckets {
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
