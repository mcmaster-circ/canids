// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package api provides the API service for the backend.\

// SHOULD JUST BE BACK FOR EASE OF DEVELOPMENT
package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mcmaster-circ/canids-v2/backend/auth"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ctxlog"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/email"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/jwtauth"
	"github.com/mcmaster-circ/canids-v2/backend/state"
	"github.com/sirupsen/logrus"
)

// authPageType represents an authentication page type.
type authPageType string

const (
	// setupPage is the initial registration page
	setupPage authPageType = "setupPage"
	// registerPage is the user registration page
	registerPage authPageType = "registerPage"
	// loginPage is the user login page
	loginPage authPageType = "loginPage"
	// logoutPage is the user logout page
	logoutPage authPageType = "logoutPage"
	// reqeustPasswordPage is the forget password page
	requestPasswordPage authPageType = "requestPasswordReset"
	// resetPasswordPage is to reset a new password (from email link)
	resetPasswordPage authPageType = "resetPasswordPage"
)

// page represents a user authentication page.
type page struct {
	Page             authPageType // page is the type of authentication page
	SuccessMsg       string       // successMsg is an success message
	ErrorMsg         string       // errorMsg is an error message
	UserRegistration bool         // UserRegistration indicates of registration link is shown on login
	Version          string       // Version is the application version
}

// registerIndexAssets registers the index handlers.
//   - /: redirect to /setup, /login, /dashboard
//   - /login: login page
//   - /logout: logout page
//   - /setup: system setup page
//   - /requestReset: request password reset page
//   - /reset: password reset page (sent in email)
//   - /registration: user registration page (only if enabled)
func registerIndexAssets(s *state.State, a *jwtauth.Config, p *auth.State, r *mux.Router) {
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(s, a, p, w, r)
	})
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logoutHandler(s, a, p, w, r)
	})
	r.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
		setupHandler(s, a, p, w, r)
	})
	r.HandleFunc("/requestReset", func(w http.ResponseWriter, r *http.Request) {
		requestResetHandler(s, a, p, w, r)
	})
	r.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		resetHandler(s, a, p, w, r)
	})
	// only register "/register" route if user registration is enabled
	if s.Config.UserRegistration {
		r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
			registerHandler(s, a, p, w, r)
		})
	}
}

// loginHandler receives "/login" HTTP requests. It will redirect to the setup
// page if the system has not been initialized. If a GET request is peformed, it
// will return the login page. If a POST request is performed, it will capture
// "user" and "pass" form values, validating against the database. If the login
// is successful, the "X-State" and "X-Class" cookies will be set and the page
// will display a login successful, redirecting to the dashboard. If the login
// is not successful, an error is returned.
func loginHandler(s *state.State, a *jwtauth.Config, p *auth.State, w http.ResponseWriter, r *http.Request) {
	// get logger from request
	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// if authentication indexes are not provisioned, redirect to setup page
	if !s.AuthReady {
		l.Info("[login] authentication index not ready, redirecting to /setup")
		http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
		return
	}
	// if POST request then validate login, else display login page
	if r.Method == "POST" {
		// get form fields + validate login
		username := r.FormValue("user")
		password := r.FormValue("pass")
		success, user := validateLogin(s, l, username, password)
		if success {
			// generate + send cookie
			user.IssuedAt = time.Now().Unix()
			token, err := a.CreateToken(user, auth.ExpireAge)
			if err != nil {
				// can't issue new token, return login page with error
				l.Error("[login] failed to create authentication token ", err)
				p.AuthPage.Execute(w, page{
					Page:             loginPage,
					SuccessMsg:       "",
					ErrorMsg:         "Please contact the system administrator.",
					UserRegistration: s.Config.UserRegistration,
					Version:          s.Hash,
				})
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

			// login page with success message
			p.AuthPage.Execute(w, page{
				Page:             loginPage,
				SuccessMsg:       "Successful login. Now redirecting to dashboard.",
				ErrorMsg:         "",
				UserRegistration: s.Config.UserRegistration,
				Version:          s.Hash,
			})
			return
		}
		// login page with failure message
		l.Info("[login] incorrect username password combination")
		p.AuthPage.Execute(w, page{
			Page:             loginPage,
			SuccessMsg:       "",
			ErrorMsg:         "Incorrect login or inactive account. Please try again.",
			UserRegistration: s.Config.UserRegistration,
			Version:          s.Hash,
		})
		return
	}
	// no POST request, display login page
	p.AuthPage.Execute(w, page{
		Page:             loginPage,
		SuccessMsg:       "",
		ErrorMsg:         "",
		UserRegistration: s.Config.UserRegistration,
		Version:          s.Hash,
	})
}

// logoutHandler receives "/logout" HTTP requests. It will redirect to the setup
// page if the system has not been initialized. It will delete the "X-State" and
// "X-Class" cookies and return a logout success page.
func logoutHandler(s *state.State, a *jwtauth.Config, p *auth.State, w http.ResponseWriter, r *http.Request) {
	// get logger from request
	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// if authentication indexes are not provisioned, redirect to setup page
	if !s.AuthReady {
		http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
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

	// logout page with success message
	p.AuthPage.Execute(w, page{
		Page:       logoutPage,
		SuccessMsg: "You have been successfully signed out.",
		ErrorMsg:   "",
		Version:    s.Hash,
	})
}

// setupHandler receives "/setup" HTTP requests. It will redirect to the index
// page if the system has been initialized. If a GET request is peformed, it
// will return the setup page. If a POST request is performed, it will capture
// all form values, attempting to create an "auth" document in the database. If
// the creation is successful, the page will display a setup successful message,
// redirecting to the login. If the creation is not successful, an error is
// returned.
func setupHandler(s *state.State, a *jwtauth.Config, p *auth.State, w http.ResponseWriter, r *http.Request) {
	// get logger from request
	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// if authentication indexes are provisioned, redirect to index page
	if s.AuthReady {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "POST" {
		// POST fields
		formName := r.FormValue("name")
		formEmail := r.FormValue("email")
		formPass := r.FormValue("pass")
		formPassConfirm := r.FormValue("passConfirm")

		// validate all fields
		if formName == "" || formEmail == "" || formPass == "" || formPassConfirm == "" {
			// return setup page with fields error
			l.Info("[setup] not all fields specified")
			p.AuthPage.Execute(w, page{
				Page:       setupPage,
				SuccessMsg: "",
				ErrorMsg:   "All fields must be specified.",
				Version:    s.Hash,
			})
			return
		}
		// validate password and password confirmation
		if formPass != formPassConfirm {
			// return setup page with password error
			l.Info("[setup] password and confirmation password are not equal")
			p.AuthPage.Execute(w, page{
				Page:       setupPage,
				SuccessMsg: "",
				ErrorMsg:   "The passwords must be the same.",
				Version:    s.Hash,
			})
			return
		}
		// hash and salt password
		hashedPass, err := jwtauth.HashPassword(formPass)
		if err != nil {
			// return setup page with password error
			l.Error("[setup] cannot hash password ", err)
			p.AuthPage.Execute(w, page{
				Page:       setupPage,
				SuccessMsg: "",
				ErrorMsg:   "Please contact the system administrator.",
				Version:    s.Hash,
			})
			return
		}
		// create "auth" index
		err = elasticsearch.CreateIndex(s, "auth")
		if err != nil {
			// return setup page with general error
			l.Error("[setup] cannot create 'auth' index ", err)
			p.AuthPage.Execute(w, page{
				Page:       setupPage,
				SuccessMsg: "",
				ErrorMsg:   "Please contact the system administrator.",
				Version:    s.Hash,
			})
			return
		}
		// create user entry as admin
		user := elasticsearch.DocumentAuth{
			UUID:      formEmail,
			Password:  hashedPass,
			Class:     jwtauth.UserAdmin,
			Name:      formName,
			Activated: true,
		}
		// index user in "auth"
		docID, err := user.Index(s)
		if err != nil {
			// return setup page with general error
			l.Error("[setup] cannot index user ", err)
			p.AuthPage.Execute(w, page{
				Page:       setupPage,
				SuccessMsg: "",
				ErrorMsg:   "Please contact the system administrator.",
				Version:    s.Hash,
			})
			return
		}
		// update state that authentication is ready
		s.AuthReady = true
		// return setup page with success message
		l.Info("[setup] created new user auth/", docID)
		p.AuthPage.Execute(w, page{
			Page:       setupPage,
			SuccessMsg: "Successful setup. Now redirecting to login page.",
			ErrorMsg:   "",
			Version:    s.Hash,
		})
		return
	}
	// setup page with no success or error message
	p.AuthPage.Execute(w, page{
		Page:       setupPage,
		SuccessMsg: "",
		ErrorMsg:   "",
		Version:    s.Hash,
	})
}

// requestResetHandler receives "/requestReset" HTTP requests. It will redirect
// to the setup page if the system has not been initialized. If a GET request is
// performed, it will return the password reset page. If a POST request is
// performed, it will capture the email address and will display message that
// the email has been sent. For security purposes, it will not tell the user if
// the email was found in the system.
func requestResetHandler(s *state.State, a *jwtauth.Config, p *auth.State, w http.ResponseWriter, r *http.Request) {
	// get logger from request
	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// if authentication indexes are not provisioned, redirect to setup page
	if !s.AuthReady {
		http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "POST" {
		// ensure email is not null
		formEmail := r.FormValue("user")
		if formEmail == "" {
			l.Info("[request reset] password reset email not specified")
			p.AuthPage.Execute(w, page{
				Page:       requestPasswordPage,
				SuccessMsg: "",
				ErrorMsg:   "An email must be specified.",
				Version:    s.Hash,
			})
			return
		}
		// query user in database
		user, _, err := elasticsearch.QueryAuthByUUID(s, formEmail)
		if err != nil {
			l.Error("[request reset] cannot retreive user for password reset ", err)
			p.AuthPage.Execute(w, page{
				Page:       requestPasswordPage,
				SuccessMsg: "If this email address is registered, you will receive an email within the next few minutes.",
				ErrorMsg:   "",
				Version:    s.Hash,
			})
			return
		}
		// generate reset token with provided expiry
		payload := &jwtauth.Payload{UUID: user.UUID}
		token, err := a.CreateToken(payload, auth.ResetDuration)
		if err != nil {
			l.Error("[request reset] failed to generate password reset token ", err)
			p.AuthPage.Execute(w, page{
				Page:       requestPasswordPage,
				SuccessMsg: "",
				ErrorMsg:   "Please contact the system administrator.",
				Version:    s.Hash,
			})
			return
		}
		// send password reset email
		domain := s.Config.SendGridDomain
		resetRequest := "http://" + domain + "/reset?token=" + token
		err = email.SendPasswordReset(s, user.Name, user.UUID, resetRequest)
		if err != nil {
			l.Error("[request reset] failed to send password reset ", err)
			p.AuthPage.Execute(w, page{
				Page:       requestPasswordPage,
				SuccessMsg: "",
				ErrorMsg:   "Please contact the system administrator.",
				Version:    s.Hash,
			})
			return
		}
		l.Info("[request reset] successfully issued password reset link for ", user.UUID)
		p.AuthPage.Execute(w, page{
			Page:       requestPasswordPage,
			SuccessMsg: "If this email address is registered, you will receive an email within the next few minutes.",
			ErrorMsg:   "",
			Version:    s.Hash,
		})
		return
	}
	// password reset page with no success or error message
	p.AuthPage.Execute(w, page{
		Page:       requestPasswordPage,
		SuccessMsg: "",
		ErrorMsg:   "",
		Version:    s.Hash,
	})
}

// resetHandler receives "/reset" HTTP requests. It will redirect to the setup
// page if the system has not been initialized. If a GET request is performed,
// it will return the page to set a password. If a POST request is performed, it
// will attempt to update the hashed password in Elasticsearch.
func resetHandler(s *state.State, a *jwtauth.Config, p *auth.State, w http.ResponseWriter, r *http.Request) {
	// get logger from request
	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// if authentication indexes are not provisioned, redirect to setup page
	if !s.AuthReady {
		http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "POST" {
		// ensure token is present
		token := r.URL.Query().Get("token")
		if token == "" {
			l.Info("[reset] new password token not present")
			p.AuthPage.Execute(w, page{
				Page:       resetPasswordPage,
				SuccessMsg: "",
				ErrorMsg:   "Request not valid. Please try again.",
				Version:    s.Hash,
			})
			return
		}
		// parse + validate token
		payload, err := a.ParseToken(token)
		if err != nil {
			l.Error("[reset] new password token cannot be parsed ", err)
			p.AuthPage.Execute(w, page{
				Page:       resetPasswordPage,
				SuccessMsg: "",
				ErrorMsg:   "The link has expired. Please request a new email.",
				Version:    s.Hash,
			})
			return
		}
		// get new passwords from forms
		formPass := r.FormValue("pass")
		formPassConfirm := r.FormValue("passConfirm")

		// validate password fields
		if formPass == "" || formPassConfirm == "" {
			// return setup page with password error
			l.Info("[reset] not all fields specified")
			p.AuthPage.Execute(w, page{
				Page:       resetPasswordPage,
				SuccessMsg: "",
				ErrorMsg:   "All fields must be specified.",
				Version:    s.Hash,
			})
			return
		}
		// validate password and password confirmation
		if formPass != formPassConfirm {
			// return setup page with password error
			l.Info("[reset] password and confirmation password are not equal")
			p.AuthPage.Execute(w, page{
				Page:       resetPasswordPage,
				SuccessMsg: "",
				ErrorMsg:   "The passwords must be the same.",
				Version:    s.Hash,
			})
			return
		}
		// hash and salt new password
		hashedPass, err := jwtauth.HashPassword(formPass)
		if err != nil {
			// return setup page with password error
			l.Error("[reset] cannot hash password ", err)
			p.AuthPage.Execute(w, page{
				Page:       resetPasswordPage,
				SuccessMsg: "",
				ErrorMsg:   "Please contact the system administrator.",
				Version:    s.Hash,
			})
			return
		}
		// get existing Elasticsearch document ID
		_, id, err := elasticsearch.QueryAuthByUUID(s, payload.UUID)
		if err != nil {
			l.Error("[reset] cannot find user from email token ", err)
			p.AuthPage.Execute(w, page{
				Page:       resetPasswordPage,
				SuccessMsg: "",
				ErrorMsg:   "The link has expired. Please request a new email.",
				Version:    s.Hash,
			})
			return
		}
		// update password field for Elasticsearch document
		err = elasticsearch.UpdatePassword(s, id, hashedPass)
		if err != nil {
			l.Error("[reset] error updating password in auth document ", err)
			p.AuthPage.Execute(w, page{
				Page:       resetPasswordPage,
				SuccessMsg: "",
				ErrorMsg:   "Please contact the system administrator.",
				Version:    s.Hash,
			})
			return
		}
		// new password page with success message
		p.AuthPage.Execute(w, page{
			Page:       resetPasswordPage,
			SuccessMsg: "The password has been successfully updated.",
			ErrorMsg:   "",
			Version:    s.Hash,
		})
		return
	}
	// new password page with no success or error message
	p.AuthPage.Execute(w, page{
		Page:       resetPasswordPage,
		SuccessMsg: "",
		ErrorMsg:   "",
		Version:    s.Hash,
	})
}

// registerHandler receives "/register" HTTP requests. The handler is only
// registered if the "user_registration" flag is enabled in the state
// configuration. It will redirect to the setup page if the system has not been
// initialized. If a GET request is peformed, it will return the user
// registration page. If a POST request is performed, it will capture all form
// values, attempting to create a non-admin "auth" document in the database. If
// the creation is successful, the page will display a successful message,
// redirecting to a login. If the creation is is not successful, an error
// message is returned.
func registerHandler(s *state.State, a *jwtauth.Config, p *auth.State, w http.ResponseWriter, r *http.Request) {
	// get logger from request
	l := ctxlog.Log(r.Context())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// if authentication indexes are not provisioned, redirect to setup page
	if !s.AuthReady {
		http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
		return
	}
	if r.Method == "POST" {
		// POST fields
		formName := r.FormValue("name")
		formEmail := r.FormValue("email")
		formPass := r.FormValue("pass")
		formPassConfirm := r.FormValue("passConfirm")

		// validate all fields
		if formName == "" || formEmail == "" || formPass == "" || formPassConfirm == "" {
			// return register page with fields error
			l.Info("[register] not all fields specified")
			p.AuthPage.Execute(w, page{
				Page:       registerPage,
				SuccessMsg: "",
				ErrorMsg:   "All fields must be specified.",
				Version:    s.Hash,
			})
			return
		}
		// validate password and password confirmation
		if formPass != formPassConfirm {
			// return register page with password error
			l.Info("[register] password and confirmation password are not equal")
			p.AuthPage.Execute(w, page{
				Page:       registerPage,
				SuccessMsg: "",
				ErrorMsg:   "The passwords must be the same.",
				Version:    s.Hash,
			})
			return
		}
		// query elasticsearch by provided uuid
		_, _, err := elasticsearch.QueryAuthByUUID(s, formEmail)
		if err == nil {
			// return register page with user already exists error
			l.Error("[register] email already exists ", formEmail)
			p.AuthPage.Execute(w, page{
				Page:       registerPage,
				SuccessMsg: "",
				ErrorMsg:   "Email address already registered.",
				Version:    s.Hash,
			})
			return
		}
		// hash and salt password
		hashedPass, err := jwtauth.HashPassword(formPass)
		if err != nil {
			// return register page with password error
			l.Error("[register] cannot hash password ", err)
			p.AuthPage.Execute(w, page{
				Page:       registerPage,
				SuccessMsg: "",
				ErrorMsg:   "Please contact the system administrator.",
				Version:    s.Hash,
			})
			return
		}
		// create user entry as standard
		user := elasticsearch.DocumentAuth{
			UUID:      formEmail,
			Password:  hashedPass,
			Class:     jwtauth.UserStandard,
			Name:      formName,
			Activated: s.Config.UserActivated,
		}
		// index user in "auth"
		docID, err := user.Index(s)
		if err != nil {
			// return register page with general error
			l.Error("[register] cannot index user ", err)
			p.AuthPage.Execute(w, page{
				Page:       registerPage,
				SuccessMsg: "",
				ErrorMsg:   "Please contact the system administrator.",
				Version:    s.Hash,
			})
			return
		}
		// return register page with success message
		l.Info("[register] created new user auth/", docID)
		// if the account isn't activated yet, notify the user
		successMsg := "Successful registration. Now redirecting to login page."
		if !s.Config.UserActivated {
			successMsg = "Successful registration. An administrator must activate your account before you can sign in. Now redirecting to login page."
		}
		p.AuthPage.Execute(w, page{
			Page:       registerPage,
			SuccessMsg: successMsg,
			ErrorMsg:   "",
			Version:    s.Hash,
		})
		return
	}
	// register page with no success or error message
	p.AuthPage.Execute(w, page{
		Page:       registerPage,
		SuccessMsg: "",
		ErrorMsg:   "",
		Version:    s.Hash,
	})
}

// validateLogin takes the state, a uuid and a password string, returning a bool
// if the login is valid. If so, an auth Payload is returned containing the
// user's information.
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
