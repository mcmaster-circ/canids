package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/api/services/utils"
	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type SetupRequest struct {
	Name            string `json:"name"`
	UUID            string `json:"uuid"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func setupUserHandler(s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {

	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "application/json")

	if elasticsearch.AuthIsActive(s) {
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "System already initialized.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	var request SetupRequest

	// Decode request json
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		l.Error("[setup] unable to decode json")

		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Bad request format.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Validate all fields
	err = utils.ValidateBasic(request.Name)
	if err != nil {
		l.Info("[setup] not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Fields " + err.Error(),
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	err = utils.ValidateBasic(request.UUID)
	if err != nil {
		l.Info("[setup] not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Fields " + err.Error(),
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	err = utils.ValidateBasic(request.Password)
	if err != nil {
		l.Info("[setup] not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Fields " + err.Error(),
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	err = utils.ValidateBasic(request.PasswordConfirm)
	if err != nil {
		l.Info("[setup] not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Fields " + err.Error(),
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Ensure passwords are the same
	if request.Password != request.PasswordConfirm {
		l.Info("[setup] password and confirmation password are not equal")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Passwords must match.",
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

	// Check if email already in elasticsearch
	_, _, err = elasticsearch.QueryAuthByUUID(s, request.UUID)
	if err == nil {
		l.Error("[setup] email already exists ", request.UUID)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Email already registered.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Hash and salt password
	hashedPass, err := jwtauth.HashPassword(request.Password)
	if err != nil {
		l.Error("[setup] cannot hash password ", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact the system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Create user for elastic search
	user := elasticsearch.DocumentAuth{
		UUID:      request.UUID,
		Name:      request.Name,
		Password:  hashedPass,
		Activated: true,
		Class:     jwtauth.UserAdmin,
	}

	// Index in "auth"
	_, err = user.Index(s)
	if err != nil {
		l.Error("[setup] cannot index user ", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact the system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	payload := &jwtauth.Payload{
		UUID:      request.UUID,
		Name:      request.Name,
		Activated: true,
		Class:     jwtauth.UserAdmin,
	}

	payload.IssuedAt = time.Now().Unix()
	token, err := a.CreateToken(payload, auth.ExpireAge)
	if err != nil {
		// can't issue new token, return login page with error
		l.Error("[setup] failed to create authentication token ", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact the system administrator.",
		}
		json.NewEncoder(w).Encode(out)
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
	if s.Settings.HTTPSEnabled {
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
	if s.Settings.HTTPSEnabled {
		cookie.SameSite = http.SameSiteStrictMode
		cookie.Secure = true
	}
	http.SetCookie(w, &cookie)

	l.Info("[login] token issued, X-State and X-Class cookies set")
	w.WriteHeader(http.StatusOK)
	out := GeneralResponse{
		Success: true,
		Message: "Successfully created your account and logged in",
	}
	json.NewEncoder(w).Encode(out)

}
