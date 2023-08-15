package elasticsearch

import (
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type DocumentIngestion struct {
	UUID string `json:"uuid"` // Represents the name of the ingestion engine
	Key  string `json:"key"`  // Represents the encryption key shared with the ingestion engine
}

const indexIngestion = "ingestion"

func (d *DocumentIngestion) Index(s *state.State) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	result, err := client.Index(indexIngestion).Document(d).Refresh(refresh.True).Do(ctx)
	return result.Id_, err
}

func QueryIngestionByUUID(s *state.State, uuid string) (DocumentIngestion, error) {
	var d DocumentIngestion
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for ingestion with provided uuid
	result, err := client.Search().Index(indexDashboard).Query(&types.Query{
		Term: map[string]types.TermQuery{
			"uuid.keyword": {Value: uuid},
		},
	}).Do(ctx)
	if err != nil {
		return d, err
	}
	// ensure ingestion was returned
	if result.Hits.Total.Value == 0 {
		return d, errors.New("ingestion: no document with uuid found")
	}
	// select + parse ingestion into DocumentIngestion
	ingestion := result.Hits.Hits[0]
	err = json.Unmarshal(ingestion.Source_, &d)
	if err != nil {
		return d, err
	}
	// successful query
	return d, nil
}

func DeleteIngestByUUID(s *state.State, uuid string) error {
	client, ctx := s.Elastic, s.ElasticCtx
	_, err := client.DeleteByQuery(indexBlacklist).Query(&types.Query{
		Term: map[string]types.TermQuery{
			"uuid.keyword": {Value: uuid},
		},
	}).Refresh(true).Do(ctx)
	return err
}
