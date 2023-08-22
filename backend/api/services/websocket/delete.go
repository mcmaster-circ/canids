package websocket

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type deleteIngestionRequest struct {
	UUID string `json:"uuid"`
}

func deleteIngestion(s *state.State, w http.ResponseWriter, r *http.Request) {

	var request deleteIngestionRequest
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

	//ensure that the uuid is not " " or "  "
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

	err = elasticsearch.DeleteIngestByUUID(s, request.UUID)
	if err != nil {
		l.Error("Failed to delete ingestion client", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	//Success

	del.update(request.UUID, false)

	w.WriteHeader(http.StatusOK)
	out := GeneralResponse{
		Success: true,
		Message: "Successfully deleted given uuid from ingestion index",
	}
	json.NewEncoder(w).Encode(out)
}
