package dashboard

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

type addViewRequest struct {
	View string `json:"view"`
}

func addViewHandler(ctx context.Context, s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	l := ctxlog.Log(r.Context())
	current := jwtauth.FromContext(ctx)

	// reject request if standard user is making it
	if current.Class == jwtauth.UserStandard {
		l.Warn("standard user attempting to update dashboard")
		w.WriteHeader(http.StatusForbidden)
		out := GeneralResponse{
			Success: false,
			Message: "Standard users can not update dashboard.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}
	// attempt to parse request
	var request addViewRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		l.Warn("invalid request format")
		w.WriteHeader(http.StatusBadRequest)
		out := GeneralResponse{
			Success: false,
			Message: "Bad request format.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	existingDashboard, err := elasticsearch.GetDashboard(s)
	if err != nil {
		l.Warn("Could not get dashboard")
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	existingDashboard.Views = append(existingDashboard.Views, request.View)
	_, docID, err := elasticsearch.QueryDashboardByUUID(s, existingDashboard.UUID)
	if err != nil {
		l.Warn("Could not get dashboard")
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	l.Println("Dashboard: ", existingDashboard)
	existingDashboard.Update(s, docID)
	if err != nil {
		l.Warn("Could not update dashboard")
		w.WriteHeader(http.StatusInternalServerError)
		out := GeneralResponse{
			Success: false,
			Message: "Please contact system administrator.",
		}
		json.NewEncoder(w).Encode(out)
		return
	}

	// Success
	out := GeneralResponse{
		Success: true,
		Message: "Successfully updated dashboard",
	}
	json.NewEncoder(w).Encode(out)
}
