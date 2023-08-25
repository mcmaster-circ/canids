package auth

import (
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/api/services/utils"
	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/email"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type resetRequest struct {
	UUID string `json:"uuid"`
}

func requestResetHandler(s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "application/json")

	var request resetRequest

	// Decode json
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		l.Info("[request reset] failed to decode json")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Bad request format.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Ensure auth is activated
	if !s.AuthReady {
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Authentication not ready",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Ensure email is entered
	err = utils.ValidateBasic(request.UUID)
	if err != nil {
		l.Info("[request reset] password reset email not specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Email " + err.Error(),
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Retrieve user from elasticsearch
	user, _, err := elasticsearch.QueryAuthByUUID(s, request.UUID)
	if err != nil {
		l.Error("[request reset] cannot retreive user for password reset ", err)

		w.WriteHeader(http.StatusOK)
		out := GeneralResponse{
			Success: true,
			Message: "If this email address is registered, you will receive an email within the next few minutes.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Generate reset token
	payload := &jwtauth.Payload{UUID: user.UUID}
	token, err := a.CreateToken(payload, auth.ResetDuration)
	if err != nil {
		l.Error("[request reset] failed to generate password reset token ", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact the system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Send password reset email
	accessURL := s.Settings.AccessURL
	resetRequest := "http://" + accessURL + "/auth/reset-password?token=" + token
	err = email.SendPasswordReset(s, user.Name, user.UUID, resetRequest)
	if err != nil {
		l.Error("[request reset] failed to send password reset ", err)
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact the system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Success
	l.Info("[request reset] successfully issued password reset link for ", user.UUID)

	w.WriteHeader(http.StatusOK)
	out := GeneralResponse{
		Success: true,
		Message: "If this email address is registered, you will receive an email within the next few minutes.",
	}
	json.NewEncoder(w).Encode(out)
}
