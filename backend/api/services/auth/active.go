package auth

import (
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type IsActiveResponse struct {
	Active bool `json:"active"`
}

func IsActiveHandler(s *state.State, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	isActive := elasticsearch.AuthIsActive(s)
	w.WriteHeader(http.StatusOK)
	resp := IsActiveResponse{
		Active: isActive,
	}

	json.NewEncoder(w).Encode(resp)
}
