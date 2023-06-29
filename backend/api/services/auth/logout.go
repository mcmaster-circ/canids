package auth

import (
	"net/http"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

func logoutHandler(s *state.State, a *jwtauth.Config, w http.ResponseWriter, r *http.Request) {

	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if !s.AuthReady {
		l.Info("[login] authentication not ready")

		// Create default user etc etc
		return
	}

	// delete the X-State cookie
	cookie := http.Cookie{
		Name:     "X-State",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true, // secure the cookie from JS attacks
	}
	// upgrade cookie security is site is accessible over SSL
	if s.Config.HTTPSEnabled {
		cookie.SameSite = http.SameSiteStrictMode
		cookie.Secure = true
	}
	http.SetCookie(w, &cookie)

	// delete the X-Class cookie
	cookie = http.Cookie{
		Name:     "X-Class",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true, // secure the cookie from JS attacks
	}
	// upgrade cookie security is site is accessible over SSL
	if s.Config.HTTPSEnabled {
		cookie.SameSite = http.SameSiteStrictMode
		cookie.Secure = true
	}
	http.SetCookie(w, &cookie)

	l.Info("[logout] token revoked, X-State and X-Class cookies cleared")

}
