package elasticsearch

import (
	"encoding/json"
	"errors"

	"github.com/mcmaster-circ/canids-v2/backend/state"
	"github.com/olivere/elastic"
)

type DocumentIngestion struct {
	UUID string `json:"uuid"` // Represents the name of the ingestion engine
	Key  string `json:"key"`  // Represents the encryption key shared with the ingestion engine
}

const indexIngestion = "ingestion"

func (d *DocumentIngestion) Index(s *state.State) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	result, err := client.Index().Index(indexIngestion).BodyJson(d).Do(ctx)
	return result.Id, err
}

func QueryIngestionByUUID(s *state.State, uuid string) (DocumentIngestion, error) {
	var d DocumentIngestion
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for ingestion with provided uuid
	termQuery := elastic.NewTermQuery("uuid.keyword", uuid)
	result, err := client.Search().Index(indexIngestion).Query(termQuery).Do(ctx)
	if err != nil {
		return d, err
	}
	// ensure ingestion was returned
	if result.Hits.TotalHits.Value == 0 {
		return d, errors.New("ingestion: no document with uuid found")
	}
	// select + parse ingestion into DocumentIngestion
	ingestion := result.Hits.Hits[0]
	err = json.Unmarshal(ingestion.Source, &d)
	if err != nil {
		return d, err
	}
	// successful query
	return d, nil
}
