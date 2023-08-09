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
	r.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		getMaxHandler(s, w, r)
	})
	r.HandleFunc("/set/", func(w http.ResponseWriter, r *http.Request) {
		setMaxHandler(s, w, r)
	})
}
