package elasticsearch

import (
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type DocumentIngestion struct {
	UUID    string `json:"uuid"`    // Represents the name of the ingestion client
	Key     string `json:"key"`     // Represents the encryption key shared with the ingestion client
	Address string `json:"address"` // Debug network address string
	Name    string `json:"string"`  // Set name for the ingestion client
}

const indexIngestion = "ingestion"

func (d *DocumentIngestion) Index(s *state.State) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	result, err := client.Index(indexIngestion).Document(d).Refresh(refresh.True).Do(ctx)
	return result.Id_, err
}

func QueryIngestionByUUID(s *state.State, uuid string) (DocumentIngestion, string, error) {
	var d DocumentIngestion
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for ingestion with provided uuid
	result, err := client.Search().Index(indexIngestion).Query(&types.Query{
		Term: map[string]types.TermQuery{
			"uuid.keyword": {Value: uuid},
		},
	}).Do(ctx)
	if err != nil {
		return d, "", err
	}
	// ensure ingestion was returned
	if result.Hits.Total.Value == 0 {
		return d, "", errors.New("ingestion: no document with uuid found")
	}
	// select + parse ingestion into DocumentIngestion
	ingestion := result.Hits.Hits[0]
	err = json.Unmarshal(ingestion.Source_, &d)
	if err != nil {
		return d, "", err
	}
	// successful query
	return d, ingestion.Id_, nil
}

func DeleteIngestByUUID(s *state.State, uuid string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	_, err := client.DeleteByQuery(indexIngestion).Query(&types.Query{
		Term: map[string]types.TermQuery{
			"uuid.keyword": {Value: uuid},
		},
	}).Refresh(true).Do(ctx)
	return err
}

// AllAuth will attempt to query the "auth" index and return all users in the
// system. It may return an error if the query cannot be completed.
func AllIngest(s *state.State) ([]DocumentIngestion, error) {
	var out []DocumentIngestion
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for all documents
	results, err := client.Search().Index(indexIngestion).
		Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		}).Size(1000).Do(ctx)
	if err != nil {
		return nil, err
	}
	// parse document into DocumentIngestion, append to out
	for _, document := range results.Hits.Hits {
		var d DocumentIngestion
		err := json.Unmarshal(document.Source_, &d)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

func (d *DocumentIngestion) Update(s *state.State, esDocID string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	_, err := client.Update(indexIngestion, esDocID).
		Doc(map[string]interface{}{
			"uuid":    d.UUID,
			"name":    d.Name,
			"address": d.Address,
			"key":     d.Key,
		}).DetectNoop(true).Do(ctx)
	return err
}
