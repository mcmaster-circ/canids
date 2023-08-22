package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/api/services/utils"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type approveIngestionRequest struct {
	UUID string `json:"uuid"` // Name of the ingestion engine
}

type GeneralResponse struct {
	Success bool   `json:"success"` // Success indicates if the request was successful
	Message string `json:"message"` // Message describes the request response
}

func approveIngestion(s *state.State, w http.ResponseWriter, r *http.Request) {

	var request approveIngestionRequest
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

	auth := waitList.getItem(request.UUID)

	err = utils.ValidateBasic(request.UUID)
	if err != nil {
		l.Warn("UUID name not specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "UUID " + err.Error(),
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// ensure ingestion client uuid does not already exist
	_, err = elasticsearch.QueryIngestionByUUID(s, request.UUID)
	if err == nil {
		// no error means we located a client
		l.Warn("uuid already exists ", request.UUID)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Ingestion with this name already defined.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	document := elasticsearch.DocumentIngestion{
		Key:     auth.Key,
		UUID:    request.UUID,
		Address: auth.Address,
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
	waitList.approve(request.UUID)

	resp := GeneralResponse{
		Success: true,
		Message: "Successfully created ingestion client",
	}

	l.Info("Created ingestion in elasticsearch")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
