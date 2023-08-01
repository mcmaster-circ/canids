package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
	"github.com/sirupsen/logrus"
)

type loginInfo struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type GeneralResponse struct {
	Success bool   `json:"success"` // Success indicates if the request was successful
	Message string `json:"message"` // Message describes the request response
}

// Handles login requests
func loginHandler(s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {

	var request loginInfo
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

	if !s.AuthReady {
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Authentication not ready",
		}
		json.NewEncoder(w).Encode(out)

		return
	}

	if r.Method == "POST" {
		uuid := request.User
		password := request.Pass

		success, user := validateLogin(s, l, uuid, password)
		if success {
			// generate + send cookie
			user.IssuedAt = time.Now().Unix()
			token, err := a.CreateToken(user, auth.ExpireAge)
			if err != nil {
				// can't issue new token, return login page with error
				l.Error("[login] failed to create authentication token ", err)
				return
			}
			// generate new X-State cookie, send to browser
			cookie := http.Cookie{
				Name:     "X-State",
				Value:    token,
				Path:     "/",
				HttpOnly: true, // secure the cookie from JS attacks
			}
			// upgrade cookie security if site is accessible over SSL
			if s.Config.HTTPSEnabled {
				cookie.SameSite = http.SameSiteStrictMode
				cookie.Secure = true
			}
			http.SetCookie(w, &cookie)

			// generate new X-Class cookie, send to browser
			cookie = http.Cookie{
				Name:     "X-Class",
				Value:    string(user.Class),
				Path:     "/",
				HttpOnly: true, // secure the cookie from JS attacks
			}
			// upgrade cookie security if site is accessible over SSL
			if s.Config.HTTPSEnabled {
				cookie.SameSite = http.SameSiteStrictMode
				cookie.Secure = true
			}
			http.SetCookie(w, &cookie)

			l.Info("[login] token issued, X-State and X-Class cookies set")
			w.WriteHeader(http.StatusOK)
			out := GeneralResponse{
				Success: true,
				Message: "Successfully logged in",
			}
			json.NewEncoder(w).Encode(out)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Incorrect username and password combination",
		}
		json.NewEncoder(w).Encode(out)
		l.Info("[login] incorrect username password combination")
	}

}

func validateLogin(s *state.State, l *logrus.Entry, uuid string, pass string) (bool, *jwtauth.Payload) {

	payload := &jwtauth.Payload{}

	// query elasticsearch by provided uuid
	db, _, err := elasticsearch.QueryAuthByUUID(s, uuid)
	if err != nil {
		// error querying
		l.Error("[login] failed to find uuid in database ", err)
		return false, payload
	}
	// ensure user is activated
	if !db.Activated {
		l.Warn("[login] user is not activated ", db.UUID)
		return false, payload
	}
	// validate user
	if uuid == db.UUID && jwtauth.HashCompare(db.Password, pass) {
		// login successful, update and return payload
		payload.UUID = db.UUID
		payload.Class = db.Class
		payload.Name = db.Name
		payload.Activated = db.Activated
		return true, payload
	}
	// login not successful
	l.Info("Username and password do not match")
	return false, payload
}
