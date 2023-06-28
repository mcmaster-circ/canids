package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type loginInfo struct {
	UUID     string `json:"uuid"`
	Password string `json:"password"`
}

func loginHandler(ctx context.Context, s *state.State, w http.ResponseWriter, r *http.Request) error {

	var request loginInfo

	// Decode request to json
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		out := GeneralResponse{
			Success: false,
			Message: "Bad request format",
		}
		json.NewEncoder(w).Encode(out)
		return nil
	}

	// Get user from elasticsearch
	docID, _, err := elasticsearch.QueryAuthByUUID(s, request.UUID)
	if err != nil {
		out := GeneralResponse{
			Success: false,
			Message: "Invalid email/password",
		}
		json.NewEncoder(w).Encode(out)
		return nil
	}

	// Check correct password
	if !jwtauth.HashCompare(docID.Password, request.Password) {
		out := GeneralResponse{
			Success: false,
			Message: "Invalid email/password",
		}
		json.NewEncoder(w).Encode(out)
		return nil
	}

	// Generate seed for JWT
	seed, err := jwtauth.GenerateSeed(SecretLength)
	if err != nil {
		return err
	}

}
