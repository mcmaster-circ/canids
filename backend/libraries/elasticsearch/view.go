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
	indexView = "view"
)

// ViewClass indicates the type of visualization.
type ViewClass string

const (
	// ViewLine is a line chart
	ViewLine ViewClass = "line"
	// ViewBar is a bar chart
	ViewBar ViewClass = "bar"
	// ViewPie is a pie chart
	ViewPie ViewClass = "pie"
	// ViewTable is a data table
	ViewTable ViewClass = "table"
	// DefaultViewName is the name given to the default view created
	DefaultViewName string = "Data Ingested"
)

var (
	// ViewClassMap maps the string representation back to ViewClass
	ViewClassMap = map[string]ViewClass{
		"line":  ViewLine,
		"bar":   ViewBar,
		"pie":   ViewPie,
		"table": ViewTable,
	}
)

// DocumentView represents a document from the "view" index.
type DocumentView struct {
	UUID       string    `json:"uuid"`       // UUID is unique view identifier
	Group      string    `json:"group"`      // Group is the group UUID the visualization is accessible by
	Authorized string    `json:"authorized"` // Authorized are authorized assets (index) used for generating data
	Name       string    `json:"name"`       // Name is common visualization name
	Class      ViewClass `json:"class"`      // Class is the class of view
	DataIndex  string    `json:"index"`      // DataIndex is index fields are contained in
	Fields     []string  `json:"fields"`     // Fields is the array of fields from Authorized to be used in this view
	FieldNames []string  `json:"fieldNames"` // FieldNames is the array of common field names
}

// Index will attempt to index the document to the "view" index. It will return
// the newly created document ID or an error.
func (d *DocumentView) Index(s *state.State) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	result, err := client.Index().Index(indexView).BodyJson(d).Do(ctx)
	return result.Id, err
}

// Update will attempt to update the document in the "view" with the provided
// Elasticsearch document ID. It will return an error if the transaction can not
// be performed.
func (d *DocumentView) Update(s *state.State, esDocID string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	_, err := client.Update().Index(indexView).Id(esDocID).
		Doc(map[string]interface{}{
			"uuid":       d.UUID,
			"group":      d.Group,
			"authorized": d.Authorized,
			"name":       d.Name,
			"class":      d.Class,
			"index":      d.DataIndex,
			"fields":     d.Fields,
			"fieldNames": d.FieldNames,
		}).DetectNoop(true).Do(ctx)
	return err
}

// QueryViewByUUID will attempt to query the "view" index for a view, returning
// a DocumentView entry and document ID string. It may return an error if the
// query cannot be completed or if the view is not found.
func QueryViewByUUID(s *state.State, uuid string) (DocumentView, string, error) {
	var d DocumentView
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for view with provided uuid
	termQuery := elastic.NewTermQuery("uuid.keyword", uuid)
	result, err := client.Search().Index(indexView).Query(termQuery).Size(1000).Do(ctx)
	if err != nil {
		return d, "", err
	}
	// ensure view was returned
	if result.Hits.TotalHits.Value == 0 {
		return d, "", errors.New("view: no document with uuid found")
	}
	// select + parse view into DocumentView
	view := result.Hits.Hits[0]
	err = json.Unmarshal(view.Source, &d)
	if err != nil {
		return d, "", err
	}
	// successful query
	return d, view.Id, nil
}

// QueryViewByGroup will attempt to query the "view" index for all views
// belonging to a group. It may return an error if the query cannot be
// completed.
func QueryViewByGroup(s *state.State, groupUUID string) ([]DocumentView, error) {
	var out []DocumentView
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for views with provided group uuid
	termQuery := elastic.NewTermQuery("group.keyword", groupUUID)
	result, err := client.Search().Index(indexView).Query(termQuery).Size(1000).Do(ctx)
	if err != nil {
		return out, err
	}
	// select + parse views into DocumentView, append to the output list
	for _, view := range result.Hits.Hits {
		var d DocumentView
		err = json.Unmarshal(view.Source, &d)
		if err != nil {
			return out, err
		}
		out = append(out, d)
	}
	// successful query
	return out, nil
}

// DeleteViewByUUID will attempt to delete a document in the "view" index with
// the specified UUID. It may return an error if the deletion cannot be
// completed.
func DeleteViewByUUID(s *state.State, uuid string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	termQuery := elastic.NewTermQuery("uuid.keyword", uuid)
	_, err := client.DeleteByQuery(indexView).Query(termQuery).Size(1000).Do(ctx)
	return err
}

// AllView will attempt to query the "view" index and return all views in the
// system. It may return an error if the query cannot be completed.
func AllView(s *state.State) ([]DocumentView, error) {
	var out []DocumentView
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for all documents
	allQuery := elastic.NewMatchAllQuery()
	results, err := client.Search().Index(indexView).Query(allQuery).Size(1000).Do(ctx)
	if err != nil {
		return nil, err
	}
	// parse groups into DocumentView, append to out
	for _, group := range results.Hits.Hits {
		var d DocumentView
		err := json.Unmarshal(group.Source, &d)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}
