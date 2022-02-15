// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package elasticsearch provides the simplified Elasticsearch interface.
package elasticsearch

import (
	"encoding/json"
	"errors"

	"github.com/mcmaster-circ/canids-v2/backend/state"
	"github.com/olivere/elastic"
)

const (
	indexDashboard = "dashboard"
)

// SizeClass indicates the size of each visualization.
type SizeClass string

const (
	// SizeFull takes up an entire dashboard row
	SizeFull SizeClass = "full"
	// SizeHalf takes up half a dashboard row
	SizeHalf SizeClass = "half"
)

var (
	// SizeClassMap maps the string representation back to SizeClasss
	SizeClassMap = map[string]SizeClass{
		"full": SizeFull,
		"half": SizeHalf,
	}
)

// DocumentDashboard represents a document from the "dashboard" index.
type DocumentDashboard struct {
	UUID  string      `json:"uuid"`  // UUID is unique dashboard identifier
	Group string      `json:"group"` // Group is group dashboard belongs to
	Name  string      `json:"name"`  // Name is dashboard name
	Views []string    `json:"views"` // Views is a list of views on the dashboard
	Sizes []SizeClass `json:"sizes"` // Sizes is a list of sizes corresponding to each view
}

// Index will attempt to index the document to the "dashboard" index. It will
// return the newly created document ID or an error.
func (d *DocumentDashboard) Index(s *state.State) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	result, err := client.Index().Index(indexDashboard).BodyJson(d).Do(ctx)
	return result.Id, err
}

// Update will attempt to update the document in the "dashboard" with the
// provided Elasticsearch document ID. It will return an error if the
// transaction can not be performed.
func (d *DocumentDashboard) Update(s *state.State, esDocID string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	_, err := client.Update().Index(indexDashboard).Id(esDocID).
		Doc(map[string]interface{}{
			"uuid":  d.UUID,
			"group": d.Group,
			"name":  d.Name,
			"views": d.Views,
			"sizes": d.Sizes,
		}).DetectNoop(true).Do(ctx)
	return err
}

// QueryDashboardByUUID will attempt to query the "dashboard" index for the
// dashboard with the specified UUID. It may return an error if the query cannot
// be completed.
func QueryDashboardByUUID(s *state.State, uuid string) (DocumentDashboard, string, error) {
	var d DocumentDashboard
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for dashboard with provided uuid
	termQuery := elastic.NewTermQuery("uuid.keyword", uuid)
	result, err := client.Search().Index(indexDashboard).Query(termQuery).Do(ctx)
	if err != nil {
		return d, "", err
	}
	// ensure dashboard was returned
	if result.Hits.TotalHits.Value == 0 {
		return d, "", errors.New("dashboard: no document with uuid found")
	}
	// select + parse dashboard into DocumentDashboard
	dashboard := result.Hits.Hits[0]
	err = json.Unmarshal(dashboard.Source, &d)
	if err != nil {
		return d, "", err
	}
	// successful query
	return d, dashboard.Id, nil
}

// QueryDashboardByGroup will attempt to query the "dashboard" index for the
// dashboard belonging to a group. It may return an error if the query cannot be
// completed.
func QueryDashboardByGroup(s *state.State, groupUUID string) (DocumentDashboard, string, error) {
	var d DocumentDashboard
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for dashboard with provided uuid
	termQuery := elastic.NewTermQuery("group.keyword", groupUUID)
	result, err := client.Search().Index(indexDashboard).Query(termQuery).Do(ctx)
	if err != nil {
		return d, "", err
	}
	// ensure dashboard was returned
	if result.Hits.TotalHits.Value == 0 {
		return d, "", errors.New("dashboard: no document for group uuid found")
	}
	// select + parse dashboard into DocumentDashboard
	dashboard := result.Hits.Hits[0]
	err = json.Unmarshal(dashboard.Source, &d)
	if err != nil {
		return d, "", err
	}
	// successful query
	return d, dashboard.Id, nil
}

// AllDashboard will attempt to query the "dashboard" index and return all dashboards in the
// system. It may return an error if the query cannot be completed.
func AllDashboard(s *state.State) ([]DocumentDashboard, error) {
	var out []DocumentDashboard
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for all documents
	allQuery := elastic.NewMatchAllQuery()
	results, err := client.Search().Index(indexDashboard).Query(allQuery).Size(1000).Do(ctx)
	if err != nil {
		return nil, err
	}
	// parse groups into DocumentDashboard, append to out
	for _, group := range results.Hits.Hits {
		var d DocumentDashboard
		err := json.Unmarshal(group.Source, &d)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}
