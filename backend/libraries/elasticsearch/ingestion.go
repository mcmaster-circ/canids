package elasticsearch

import "github.com/mcmaster-circ/canids-v2/backend/state"

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
