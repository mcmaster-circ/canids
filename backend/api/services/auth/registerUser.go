package auth

import (
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/api/services/utils"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type registerUserRequest struct {
	Name            string `json:"name"`
	UUID            string `json:"uuid"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func registerUserHandler(s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {

	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "application/json")

	var request registerUserRequest

	// Decode request json
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		l.Error("[register] unable to decode json")

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
		l.Info("[register] not all fields specified")
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
		l.Info("[register] not all fields specified")
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
		l.Info("[register] not all fields specified")
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
		l.Info("[register] not all fields specified")
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
		l.Info("[register] password and confirmation password are not equal")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Passwords must match.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Check if email already in elasticsearch
	_, _, err = elasticsearch.QueryAuthByUUID(s, request.UUID)
	if err == nil {
		l.Error("[register] email already exists ", request.UUID)
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
		l.Error("[register] cannot hash password ", err)
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
		Activated: s.Config.UserActivated,
		Class:     jwtauth.UserStandard,
	}

	// Index in "auth"
	docID, err := user.Index(s)
	if err != nil {
		l.Error("[register] cannot index user ", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact the system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// return register page with success message
	l.Info("[register] created new user auth/", docID)
	// if the account isn't activated yet, notify the user
	successMsg := "Successful registration. Now redirecting to login page."
	if !s.Config.UserActivated {
		successMsg = "Successful registration. An administrator must activate your account before you can sign in. Now redirecting to login page."
	}
	w.WriteHeader(http.StatusOK)
	out := GeneralResponse{
		Success: true,
		Message: successMsg,
	}
	json.NewEncoder(w).Encode(out)
	return

}
