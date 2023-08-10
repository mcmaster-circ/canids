package websocket

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type createIngestionRequest struct {
	UUID string `json:"uuid"` // Name of the ingestion engine
}

type createIngestionResponse struct {
	Key string `json:"key"` // Encryption key
}

type GeneralResponse struct {
	Success bool   `json:"success"` // Success indicates if the request was successful
	Message string `json:"message"` // Message describes the request response
}

func createIngestion(s *state.State, w http.ResponseWriter, r *http.Request) {

	var request createIngestionRequest
	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "application/json")

	// Decode request to json
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		l.Error("Failed to decode json", err)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Bad request format",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	//ensure that the blacklist name is not " " or "  "
	trimmed := strings.TrimSpace(request.UUID)

	// ensure name is not empty
	if request.UUID == "" || len(trimmed) == 0 {
		l.Warn("UUID name not specified specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "UUID field must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Generate key

	key := make([]byte, 32)

	_, err = rand.Read(key)
	if err != nil {
		l.Error("Failed to generate key", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact your system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	encodedKey := base64.StdEncoding.EncodeToString(key)

	document := elasticsearch.DocumentIngestion{
		Key:  encodedKey,
		UUID: request.UUID,
	}

	_, err = document.Index(s)
	if err != nil {
		l.Error("Failed to index ingestion", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact your system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Success

	resp := createIngestionResponse{
		Key: encodedKey,
	}

	l.Info("Created ingestion in elasticsearch")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
