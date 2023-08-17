package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// listReponse is the format of the list user response.
type listResponse struct {
	Success bool        `json:"success"` // Success indicates if the request was successful
	Clients []Ingestion `json:"clients"` // Users is a list of users
}

// User represents the list of users in the system.
type Ingestion struct {
	UUID string `json:"uuid"` // Represents the name of the ingestion client
	Key  string `json:"key"`  // Represents the encryption key shared with the ingestion client
}

// listHandler is "/api/ingestion/list". It will return the list of clients
func listHandler(s *state.State, w http.ResponseWriter, r *http.Request) {

	// get user making current request + logging context
	w.Header().Set("Content-Type", "application/json")
	l := ctxlog.Log(r.Context())

	// output
	var out listResponse

	clients, err := elasticsearch.AllIngest(s)
	if err != nil {
		out := GeneralResponse{
			Success: false,
			Message: "Failed to retrieved ingestion clients.",
		}
		l.Error("error getting clients ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(out)
		return
	}

	for _, c := range clients {
		out.Clients = append(out.Clients, Ingestion{
			UUID: c.UUID,
			Key:  c.Key,
		})
	}

	// return list of users
	out.Success = true
	json.NewEncoder(w).Encode(out)
}
