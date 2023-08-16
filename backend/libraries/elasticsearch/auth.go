// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package elasticsearch provides the simplified Elasticsearch interface.
package elasticsearch

import (
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

const (
	indexAuth = "auth"
)

// DocumentAuth represents a document from the "auth" index.
type DocumentAuth struct {
	UUID      string            `json:"uuid"`      // UUID is unique user identifier
	Password  string            `json:"password"`  // Password is salted user password
	Class     jwtauth.UserClass `json:"class"`     // Class is user class
	Name      string            `json:"name"`      // Name is the user's name
	Activated bool              `json:"activated"` // Activated if account is active
}

// Index will attempt to index the document to the "auth" index. It will return
// the newly created document ID or an error.
func (d *DocumentAuth) Index(s *state.State) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	response, err := client.Index(indexAuth).Document(d).Do(ctx)
	if err != nil {
		return "", err
	}
	return response.Id_, err
}

// Update will attempt to update the document in the "auth" with the provided
// Elasticsearch document ID. It will return an error if the transaction can not
// be performed.
func (d *DocumentAuth) Update(s *state.State, esDocID string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	_, err := client.Update(indexAuth, esDocID).Doc(d).Do(ctx)
	return err
}

// QueryAuthByUUID will attempt to query the "auth" index for a user, returning
// a DocumentAuth entry and document ID string. It may return an error if the
// query cannot be completed or if the user is not found.
func QueryAuthByUUID(s *state.State, uuid string) (DocumentAuth, string, error) {
	var d DocumentAuth
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for user with provided uuid
	res, err := client.Search().
		Index(indexAuth).
		Request(&search.Request{
			Query: &types.Query{
				Term: map[string]types.TermQuery{
					"uuid.keyword": {Value: uuid},
				},
			},
		}).Do(ctx)
	if err != nil {
		return d, "", err
	}

	// ensure user was returned
	if res.Hits.Total.Value == 0 {
		return d, "", errors.New("auth: no document with uuid found")
	}
	// select + parse user into DocumentAuth
	user := res.Hits.Hits[0]
	err = json.Unmarshal(user.Source_, &d)

	// successful query
	return d, user.Id_, nil
}

// DeleteAuthByUUID will attempt to delete a document in the "auth" index with
// the specified UUID. It may return an error if the deletion cannot be completed.
func DeleteAuthByUUID(s *state.State, uuid string) error {
	client, ctx := s.Elastic, s.ElasticCtx

	_, err := client.DeleteByQuery(indexAuth).
		Query(&types.Query{
			Term: map[string]types.TermQuery{
				"uuid.keyword": {Value: uuid},
			},
		}).Do(ctx)
	return err
}

// UpdatePassword is for updating the password of an existing user in "auth". It
// accepts a state, the Elasticsearch document ID and a new password. It may
// return an error if the password cannot be updated.
func UpdatePassword(s *state.State, docID string, newPass string) error {
	client, ctx := s.Elastic, s.ElasticCtx

	updates := map[string]interface{}{"password": newPass}
	_, err := client.Update(indexAuth, docID).Doc(updates).Do(ctx)
	return err
}

// AllAuth will attempt to query the "auth" index and return all users in the
// system. It may return an error if the query cannot be completed.
func AllAuth(s *state.State) ([]DocumentAuth, error) {
	var out []DocumentAuth
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for all documents
	results, err := client.Search().Index(indexAuth).
		Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		}).Size(1000).Do(ctx)
	if err != nil {
		return nil, err
	}
	// parse document into DocumentAuth, append to out
	for _, document := range results.Hits.Hits {
		var d DocumentAuth
		err := json.Unmarshal(document.Source_, &d)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

func AuthIsActive(s *state.State) bool {
	client, ctx := s.Elastic, s.ElasticCtx

	isEmpty, err := client.Indices.Exists(indexAuth).Do(ctx)

	if err != nil {
		return false
	}

	return isEmpty

}
