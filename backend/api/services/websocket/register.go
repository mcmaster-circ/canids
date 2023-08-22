package websocket

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

func RegisterWS(s *state.State, r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleWebSocket(s, w, r)
	})
}

func RegisterUpdateFunctions(s *state.State, r *mux.Router) {
	r.HandleFunc("/getESMax", func(w http.ResponseWriter, r *http.Request) {
		getMaxHandler(s, w, r)
	})
	r.HandleFunc("/setESMax", func(w http.ResponseWriter, r *http.Request) {
		setMaxHandler(s, w, r)
	})
	r.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		deleteIngestion(s, w, r)
	})
	r.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		listHandler(s, w, r)
	})
	r.HandleFunc("/approve", func(w http.ResponseWriter, r *http.Request) {
		approveIngestion(s, w, r)
	})
	r.HandleFunc("/rename", func(w http.ResponseWriter, r *http.Request) {
		renameIngestion(s, w, r)
	})
}
