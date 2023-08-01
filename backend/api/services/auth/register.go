package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

func RegisterRoutes(s *state.State, a *jwtauth.Config, r *mux.Router) {

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(s, a, w, r)
	})

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logoutHandler(s, a, w, r)
	})

	r.HandleFunc("/requestReset", func(w http.ResponseWriter, r *http.Request) {
		requestResetHandler(s, a, w, r)
	})

	r.HandleFunc("/resetPassword", func(w http.ResponseWriter, r *http.Request) {
		resetHandler(s, a, w, r)
	})

	r.HandleFunc("/registerUser", func(w http.ResponseWriter, r *http.Request) {
		registerUserHandler(s, a, w, r)
	})
}