// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package elasticsearch provides the simplified Elasticsearch interface.
package elasticsearch

import "github.com/mcmaster-circ/canids-v2/backend/state"

// CreateIndex will attempt to create an index with the specified name. It may
// return an error if the index cannot be created.
func CreateIndex(s *state.State, indexName string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	_, err := client.CreateIndex(indexName).Do(ctx)
	return err
}

// DeleteIndex will attempt to delete the index with the specified name. It may
// return an error if the index cannot be deleted.
func DeleteIndex(s *state.State, indexName string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	_, err := client.DeleteIndex(indexName).Do(ctx)
	return err
}
