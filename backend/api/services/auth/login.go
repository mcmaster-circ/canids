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
	UUID     string `json:"uuid"`
	Password string `json:"password"`
}

type GeneralResponse struct {
	Success bool   `json:"success"` // Success indicates if the request was successful
	Message string `json:"message"` // Message describes the request response
}

// Handles login requests
func loginHandler(s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {

	var request loginInfo
	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Decode request to json
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
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
		uuid := request.UUID
		password := request.Password

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

			return
		}

		l.Info("[login] incorrect username password combination")
	}

	// // Get user from elasticsearch
	// docID, _, err := elasticsearch.QueryAuthByUUID(s, request.UUID)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	out := GeneralResponse{
	// 		Success: false,
	// 		Message: "Invalid email/password",
	// 	}
	// 	json.NewEncoder(w).Encode(out)
	// 	return
	// }

	// // Check correct password
	// if !jwtauth.HashCompare(docID.Password, request.Password) {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	out := GeneralResponse{
	// 		Success: false,
	// 		Message: "Invalid email/password",
	// 	}
	// 	json.NewEncoder(w).Encode(out)
	// 	return
	// }

	// payload := jwtauth.Payload{
	// 	docID.UUID,
	// 	docID.Class,
	// 	docID.Name,
	// 	docID.Activated,
	// 	jwt.StandardClaims{},
	// }

	// // Generate JWT Token
	// token, err := a.CreateToken(&payload, auth.ExpireAge)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	out := GeneralResponse{
	// 		Success: false,
	// 		Message: "Failed to generate token",
	// 	}
	// 	json.NewEncoder(w).Encode(out)
	// 	return
	// }

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
	return false, payload
}
