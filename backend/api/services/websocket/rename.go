package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/api/services/utils"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type renameIngestionRequest struct {
	UUID string `json:"uuid"` // UUID of the ingestion engine
	Name string `json:"name"` // Name of the ingestion engine
}

type GeneralResponse struct {
	Success bool   `json:"success"` // Success indicates if the request was successful
	Message string `json:"message"` // Message describes the request response
}

func renameIngestion(s *state.State, w http.ResponseWriter, r *http.Request) {

	var request renameIngestionRequest
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

	err = utils.ValidateBasic(request.Name)
	if err != nil {
		l.Warn("Name name not specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Name " + err.Error(),
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	s.Log.Println("Name: ", request.Name)

	document := elasticsearch.DocumentIngestion{
		Name: request.Name,
	}

	// query current ingestion to update
	existing, _, err := elasticsearch.QueryIngestionByUUID(s, request.UUID)
	if err != nil {
		l.Warn("invalid ingestion uuid ", request.UUID)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid ingestion UUID specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	document.Address = existing.Address
	document.Key = existing.Key
	document.UUID = existing.UUID

	s.Log.Println("Document: ", document)

	err = elasticsearch.DeleteIngestByUUID(s, existing.UUID)
	if err != nil {
		l.Error("Failed to update ingestion", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact your system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	_, err = document.Index(s)
	if err != nil {
		l.Error("Failed to update ingestion", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact your system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	resp := GeneralResponse{
		Success: true,
		Message: "Successfully updates ingestion client",
	}

	l.Info("Updated ingestion in elasticsearch")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
