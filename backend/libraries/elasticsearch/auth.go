// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package elasticsearch provides the simplified Elasticsearch interface.
package elasticsearch

import (
	"encoding/json"
	"errors"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
	"github.com/olivere/elastic"
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
	Group     string            `json:"group"`     // Group is the user's group
	Activated bool              `json:"activated"` // Activated if account is active
}

// Index will attempt to index the document to the "auth" index. It will return
// the newly created document ID or an error.
func (d *DocumentAuth) Index(s *state.State) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	result, err := client.Index().Index(indexAuth).BodyJson(d).Do(ctx)
	return result.Id, err
}

// Update will attempt to update the document in the "auth" with the provided
// Elasticsearch document ID. It will return an error if the transaction can not
// be performed.
func (d *DocumentAuth) Update(s *state.State, esDocID string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	_, err := client.Update().Index(indexAuth).Id(esDocID).
		Doc(map[string]interface{}{
			"uuid":      d.UUID,
			"password":  d.Password,
			"class":     d.Class,
			"name":      d.Name,
			"group":     d.Group,
			"activated": d.Activated,
		}).DetectNoop(true).Do(ctx)
	return err
}

// QueryAuthByUUID will attempt to query the "auth" index for a user, returning
// a DocumentAuth entry and document ID string. It may return an error if the
// query cannot be completed or if the user is not found.
func QueryAuthByUUID(s *state.State, uuid string) (DocumentAuth, string, error) {
	var d DocumentAuth
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for user with provided uuid
	termQuery := elastic.NewTermQuery("uuid.keyword", uuid)
	result, err := client.Search().Index(indexAuth).Query(termQuery).Do(ctx)
	if err != nil {
		return d, "", err
	}
	// ensure user was returned
	if result.Hits.TotalHits.Value == 0 {
		return d, "", errors.New("auth: no document with uuid found")
	}
	// select + parse user into DocumentAuth
	user := result.Hits.Hits[0]
	err = json.Unmarshal(user.Source, &d)
	if err != nil {
		return d, "", err
	}
	// successful query
	return d, user.Id, nil
}

// QueryAuthByGroup will attempt to query the "auth" index for all users
// belonging to a group. It may return an error if the query cannot be
// completed.
func QueryAuthByGroup(s *state.State, groupUUID string) ([]DocumentAuth, error) {
	var out []DocumentAuth
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for users with provided group uuid
	termQuery := elastic.NewTermQuery("group.keyword", groupUUID)
	result, err := client.Search().Index(indexAuth).Query(termQuery).Do(ctx)
	if err != nil {
		return out, err
	}
	// select + parse user into DocumentAuth, append to the output list
	for _, user := range result.Hits.Hits {
		var d DocumentAuth
		err = json.Unmarshal(user.Source, &d)
		if err != nil {
			return out, err
		}
		out = append(out, d)
	}
	// successful query
	return out, nil
}

// DeleteAuthByUUID will attempt to delete a document in the "auth" index with
// the specified UUID. It may return an error if the deletion cannot be completed.
func DeleteAuthByUUID(s *state.State, uuid string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	termQuery := elastic.NewTermQuery("uuid.keyword", uuid)
	_, err := client.DeleteByQuery(indexAuth).Query(termQuery).Do(ctx)
	return err
}

// UpdatePassword is for updating the password of an existing user in "auth". It
// accepts a state, the Elasticsearch document ID and a new password. It may
// return an error if the password cannot be updated.
func UpdatePassword(s *state.State, docID string, newPass string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	updates := map[string]interface{}{"password": newPass}
	_, err := client.Update().Index(indexAuth).Id(docID).Doc(updates).Do(ctx)
	return err
}
