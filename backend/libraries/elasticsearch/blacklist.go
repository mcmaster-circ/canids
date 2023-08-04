// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package elasticsearch provides the simplified Elasticsearch interface.
package elasticsearch

import (
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

const (
	indexBlacklist = "blacklist"
)

// DocumentBlacklist represents a document from the "blacklist" index.
type DocumentBlacklist struct {
	UUID string `json:"uuid"` // UUID is the unique blacklist identifier
	Name string `json:"name"` // Name is the blacklist display name
	URL  string `json:"url"`  // URL is the blacklist URL
}

// Index will attempt to index the document to the "blacklist" index. It will return
// the newly created document ID or an error.
func (d *DocumentBlacklist) Index(s *state.State) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	result, err := client.Index(indexBlacklist).Document(d).Refresh(refresh.True).Do(ctx)
	if err != nil {
		return "", err
	}
	return result.Id_, nil
}

// Update will attempt to update the document in the "blacklist" with the provided
// Elasticsearch document ID. It will return an error if the transaction can not
// be performed.
func (d *DocumentBlacklist) Update(s *state.State, esDocID string) error {
	client, ctx := s.Elastic, s.ElasticCtx

	_, err := client.Update(indexBlacklist, esDocID).
		Doc(map[string]interface{}{
			"uuid": d.UUID,
			"name": d.Name,
			"url":  d.URL,
		}).DetectNoop(true).Refresh(refresh.True).Do(ctx)
	return err
}

// QueryBlacklistByUUID will attempt to query the "blacklist" index for a blacklist,
// returning a DocumentBlacklist entry and document ID string. It may return an
// error if the query cannot be completed or if the blacklist is not found.
func QueryBlacklistByUUID(s *state.State, uuid string) (DocumentBlacklist, string, error) {
	var d DocumentBlacklist
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for blacklist with provided uuid
	result, err := client.Search().Index(indexBlacklist).Query(&types.Query{
		Term: map[string]types.TermQuery{
			"uuid.keyword": {Value: uuid},
		},
	}).Size(1000).Do(ctx)

	if err != nil {
		return d, "", err
	}


	// ensure blacklist was returned
	if result.Hits.Total.Value == 0 {
		return d, "", errors.New("blacklist: no document with uuid found")
	}

	// select + parse blacklist into DocumentBlacklist
	blacklist := result.Hits.Hits[0]
	err = json.Unmarshal(blacklist.Source_, &d)
	if err != nil {
		return d, "", err
	}

	// successful query
	return d, blacklist.Id_, nil
}

// AllBlacklists will attempt to query the "blacklist" index and return all blacklists in the
// system. It may return an error if the query cannot be completed.
func AllBlacklists(s *state.State) ([]DocumentBlacklist, error) {
	var out []DocumentBlacklist
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for all documents
	results, err := client.Search().Index(indexBlacklist).Query(&types.Query{
		MatchAll: &types.MatchAllQuery{},
	}).Size(1000).Do(ctx)

	if err != nil {
		return nil, err
	}
	// parse blacklists into DocumentBlacklist, append to out
	for _, blacklist := range results.Hits.Hits {
		var d DocumentBlacklist
		err := json.Unmarshal(blacklist.Source_, &d)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

// DeleteBlacklistByUUID will attempt to delete a document in the "blacklist" index with
// the specified UUID. It may return an error if the deletion cannot be
// completed.
func DeleteBlacklistByUUID(s *state.State, uuid string) error {
	client, ctx := s.Elastic, s.ElasticCtx

	_, err := client.DeleteByQuery(indexBlacklist).Query(&types.Query{
		Term: map[string]types.TermQuery{
			"uuid.keyword": {Value: uuid},
		},
	}).Refresh(true).Do(ctx)

	return err
}
