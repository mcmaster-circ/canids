package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type maxSizeHolder struct {
	MaxIndexSize int `json:"maxIndexSize"`
}

// Handler for setting max elasticsearch index size
func setMaxHandler(s *state.State, w http.ResponseWriter, r *http.Request) {

	var request maxSizeHolder
	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "application/json")

	// Decode request to json
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		l.Error("Failed to decode json", err)
		w.WriteHeader(http.StatusBadRequest)
		out := auth.GeneralResponse{
			Success: false,
			Message: "Bad request format",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	if request.MaxIndexSize < 0 {
		out := auth.GeneralResponse{
			Success: false,
			Message: "Bad request format",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	l.Info("[ws] set max index size: ", request.MaxIndexSize)
	SetMaxElasticIndexSize(request.MaxIndexSize)
	w.WriteHeader(http.StatusOK)
	out := auth.GeneralResponse{
		Success: true,
		Message: "Successfully updated max index size",
	}
	json.NewEncoder(w).Encode(out)
}

// Handler for getting max elasticsearch index size
func getMaxHandler(s *state.State, w http.ResponseWriter, r *http.Request) {
	l := ctxlog.Log(r.Context())

	var response maxSizeHolder
	response.MaxIndexSize = GetMaxElasticIndexSize()

	l.Info("[ws] max index size: ", response.MaxIndexSize)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
