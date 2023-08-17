// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package auth provides the authentication state for the backend.
package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
	log "github.com/sirupsen/logrus"
)

// GeneralResponse is the structure of a general response.
type GeneralResponse struct {
	Success bool   `json:"success"` // Success indicates if the request was successful
	Message string `json:"message"` // Message describes the request response
}

var (
	// internalServerError is a server error.
	internalServerError = GeneralResponse{
		Success: false,
		Message: "500 Internal Server Error",
	}
	// unauthorizedError is a server error.
	unauthorizedError = GeneralResponse{
		Success: false,
		Message: "401 Unauthorized",
	}
)

// Middleware takes the global state, the auth state and the HTTP handler,
// conditionally forwarding the request to the next handler. The middleware will
// validate the X-State token. If the token is invalid or expired (older than
// ExpireAge), it will return a 401 Unauthorized error. If the token is valid
// and older than RenewAge, the middleware will query the database and update
// the token automatically. The request will then proceed normally. If the token
// is valid and not older than RenewAge, the request will proceed normally.
func Middleware(s *state.State, a *jwtauth.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get log context from request
		ctx := r.Context()
		l := ctxlog.Log(ctx)

		// if middleware disabled, just serve request
		if s.Settings.MiddlewareDisable {
			l.Info("[middleware]: middleware disabled")
			next.ServeHTTP(w, r)
			return
		}
		// middleware is enabled, get the state token in cookie
		cookie, err := r.Cookie("X-State")
		if err != nil || cookie.Value == "" {
			// cookie not present, return 401
			l.Warn("[middleware] missing X-State cookie")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(unauthorizedError)
			return
		}
		// validate the state token
		user, err := a.ParseToken(cookie.Value)
		if err != nil {
			// token is not valid, return 401
			l.Warn("[middleware] invalid X-State cookie: ", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(unauthorizedError)
			return
		}

		// user is authenticated at this point, inject authentication
		// information for logging
		l = l.WithFields(log.Fields{
			"uuid":  user.UUID,
			"class": string(user.Class),
		})
		ctx = ctxlog.WithFields(ctx, l)

		// if the token is valid and is older than the renew age, renew the
		// token
		renewTime := time.Unix(user.IssuedAt, 0).Add(RenewAge)
		if time.Now().After(renewTime) {
			// query user in database for state changes
			esUser, _, err := elasticsearch.QueryAuthByUUID(s, user.UUID)
			// if user can not be located or is no longer activated, revoke
			// access
			if err != nil || !esUser.Activated {
				// delete the X-State cookie
				cookie := http.Cookie{
					Name:     "X-State",
					MaxAge:   -1,
					Path:     "/",
					HttpOnly: true, // secure the cookie from JS attacks
				}
				// upgrade cookie security is site is accessible over SSL
				if s.Settings.HTTPSEnabled {
					cookie.SameSite = http.SameSiteStrictMode
					cookie.Secure = true
				}
				http.SetCookie(w, &cookie)
				l.Info("[middleware] token revoked, X-State cookie cleared")

				// delete the X-Class cookie
				cookie = http.Cookie{
					Name:     "X-Class",
					MaxAge:   -1,
					Path:     "/",
					HttpOnly: true, // secure the cookie from JS attacks
				}
				// upgrade cookie security is site is accessible over SSL
				if s.Settings.HTTPSEnabled {
					cookie.SameSite = http.SameSiteStrictMode
					cookie.Secure = true
				}
				http.SetCookie(w, &cookie)
				l.Info("[middleware] token revoked, X-Class cookie cleared")

				// return 401 message
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(unauthorizedError)
				return
			}
			// do not reovke token, synchronize token setings with database
			user.UUID = esUser.UUID
			user.Class = esUser.Class
			user.Name = esUser.Name
			user.Activated = esUser.Activated

			// update time + generate new token
			user.IssuedAt = time.Now().Unix()
			token, err := a.CreateToken(user, ExpireAge)
			if err != nil {
				// can't issue new token, return error
				l.Error("[middleware] failed to renew token", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(internalServerError)
				return
			}
			// generate new X-State cookie, send to browser
			cookie := http.Cookie{
				Name:     "X-State",
				Value:    token,
				Path:     "/",
				HttpOnly: true, // secure the cookie from JS attacks
			}
			// upgrade cookie security is site is accessible over SSL
			if s.Settings.HTTPSEnabled {
				cookie.SameSite = http.SameSiteStrictMode
				cookie.Secure = true
			}
			l.Debug("[middleware] renewed X-State cookie")
			http.SetCookie(w, &cookie)

			// generate new X-Class cookie, send to browser
			cookie = http.Cookie{
				Name:     "X-Class",
				Value:    string(user.Class),
				Path:     "/",
				HttpOnly: true, // secure the cookie from JS attacks
			}
			// upgrade cookie security is site is accessible over SSL
			if s.Settings.HTTPSEnabled {
				cookie.SameSite = http.SameSiteStrictMode
				cookie.Secure = true
			}
			l.Debug("[middleware] renewed X-Class cookie")
			http.SetCookie(w, &cookie)
		}

		// inject user authentication payload into context
		ctx = user.Context(ctx)
		l.Debug("[middleware] request")

		// cookie is present and state token is valid, process request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
