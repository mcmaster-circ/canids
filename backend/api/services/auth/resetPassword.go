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

type resetFields struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

// Resets the password of the user with email link
func resetHandler(s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {

	var request resetFields

	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		l.Info("[reset] failed to decode json")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Bad request format.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	if !s.AuthReady {
		l.Info("[reset] authentication not ready")
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Authentication not ready",
		}
		json.NewEncoder(w).Encode(out)

		return
	}

	// Check for token present
	if request.Token == "" {
		l.Info("[reset] new password token not present")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Request not valid. Please try again.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Parse token
	payload, err := a.ParseToken(request.Token)
	if err != nil {
		l.Error("[reset] new password token cannot be parsed ", err)

		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "The link has expired. Please request a new email link",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Validate passwords
	err = utils.ValidateBasic(request.Password)
	if err != nil {
		l.Info("[reset] not all fields specified")
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
		l.Info("[reset] not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Fields " + err.Error(),
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Check that passwords match
	if request.Password != request.PasswordConfirm {
		l.Info("[reset] password and confirmation password are not equal")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Passwords must match",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Hash and salt password
	hashedPass, err := jwtauth.HashPassword(request.Password)
	if err != nil {
		l.Error("[reset] cannot hash password ", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact the system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Find user coresponding to token
	_, id, err := elasticsearch.QueryAuthByUUID(s, payload.UUID)
	if err != nil {
		l.Error("[reset] cannot find user from email token ", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "The link has expired. Please request a new email link.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Update user password
	err = elasticsearch.UpdatePassword(s, id, hashedPass)
	if err != nil {
		l.Error("[reset] error updating password in auth document ", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact the system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	//Success
	w.WriteHeader(http.StatusOK)
	out := GeneralResponse{
		Success: true,
		Message: "Successfully updated password",
	}
	json.NewEncoder(w).Encode(out)
}
