package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type setupInfo struct {
	UUID     string `json:"uuid"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func handleSetup(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r http.Request) {

	l := ctxlog.Log(ctx)
	var request setupInfo

	// Decode JSON response
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

	// ensure all fields are present
	if request.Name == "" || request.UUID == "" {
		l.Warn("not all fields specified")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "All fields must be specified.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Validate email
	validEmail := IsValidEmail(request.UUID)
	if !validEmail {
		l.Warn("invalid email")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Invalid email.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// ensure user does not already exist
	_, _, err = elasticsearch.QueryAuthByUUID(s, request.UUID)
	if err == nil {
		// no error means we located a user
		l.Warn("uuid already exists ", request.UUID)
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "UUID email address provided already has account.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	user := elasticsearch.DocumentAuth{
		UUID:      request.UUID,
		Password:  request.Password,
		Class:     jwtauth.UserAdmin,
		Name:      request.Name,
		Activated: true,
	}

	// Index new user in database
	docID, err := user.Index(s)
	if err != nil {
		l.Error("cannot index user ", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(InternalServerError)
		return
	}
	l.Info("created new user auth/", docID)

	// Send account verification email?

}

func IsValidEmail(email string) bool {
	// Regular expression pattern for email validation
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		return false
	}

	return match
}

var (
	// InternalServerError is the a JSON error message.
	InternalServerError = GeneralResponse{
		Success: false,
		Message: "500 Internal Server Error",
	}
)
