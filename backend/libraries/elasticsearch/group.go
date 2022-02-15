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
	indexGroup = "group"
)

// DocumentGroup represents a document from the "group" index.
type DocumentGroup struct {
	UUID       string   `json:"uuid"`       // UUID is unique group identifier
	Name       string   `json:"name"`       // Name is the group's display name
	Authorized []string `json:"authorized"` // Authorized is a list of assets the group has control of
}

// Index will attempt to index the document to the "group" index. It will return
// the newly created document ID or an error.
func (d *DocumentGroup) Index(s *state.State) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	result, err := client.Index().Index(indexGroup).BodyJson(d).Do(ctx)
	return result.Id, err
}

// Update will attempt to update the document in the "group" with the provided
// Elasticsearch document ID. It will return an error if the transaction can not
// be performed.
func (d *DocumentGroup) Update(s *state.State, esDocID string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	_, err := client.Update().Index(indexGroup).Id(esDocID).
		Doc(map[string]interface{}{
			"uuid":       d.UUID,
			"name":       d.Name,
			"authorized": d.Authorized,
		}).DetectNoop(true).Do(ctx)
	return err
}

// QueryGroupByUUID will attempt to query the "group" index for a group,
// returning a DocumentGroup entry and document ID string. It may return an
// error if the query cannot be completed or if the group is not found.
func QueryGroupByUUID(s *state.State, uuid string) (DocumentGroup, string, error) {
	var d DocumentGroup
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for group with provided uuid
	termQuery := elastic.NewTermQuery("uuid.keyword", uuid)
	result, err := client.Search().Index(indexGroup).Query(termQuery).Size(1000).Do(ctx)
	if err != nil {
		return d, "", err
	}
	// ensure group was returned
	if result.Hits.TotalHits.Value == 0 {
		return d, "", errors.New("group: no document with uuid found")
	}
	// select + parse group into DocumentGroup
	group := result.Hits.Hits[0]
	err = json.Unmarshal(group.Source, &d)
	if err != nil {
		return d, "", err
	}
	// successful query
	return d, group.Id, nil
}

// AllGroup will attempt to query the "group" index and return all groups in the
// system. It may return an error if the query cannot be completed.
func AllGroup(s *state.State) ([]DocumentGroup, error) {
	var out []DocumentGroup
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for all documents
	allQuery := elastic.NewMatchAllQuery()
	results, err := client.Search().Index(indexGroup).Query(allQuery).Size(1000).Do(ctx)
	if err != nil {
		return nil, err
	}
	// parse groups into DocumentGroup, append to out
	for _, group := range results.Hits.Hits {
		var d DocumentGroup
		err := json.Unmarshal(group.Source, &d)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

// DeleteGroupByUUID will attempt to delete a document in the "group" index with
// the specified UUID. It may return an error if the deletion cannot be
// completed.
func DeleteGroupByUUID(s *state.State, uuid string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	termQuery := elastic.NewTermQuery("uuid.keyword", uuid)
	_, err := client.DeleteByQuery(indexGroup).Query(termQuery).Size(1000).Do(ctx)
	return err
}
