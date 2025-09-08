package main

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"sketchdb.cozycole.net/internal/models"
)

type contextKey string

const userContextKey = contextKey("user")

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Acts as a triger to make Go's HTTP server automatically
				// close the current connection
				w.Header().Set("Connection", "close")
				app.serverError(r, w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.debugMode {
			app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.users.GetById(id)
		if err != nil && err != models.ErrNoRecord {
			app.serverError(r, w, err)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireRoles(roles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(userContextKey).(*models.User)
		if ok && !slices.Contains(roles, derefString(user.Role)) {
			app.unauthorized(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAutheticated(r) {
			app.sessionManager.Put(r.Context(), "postLoginRedirectURL", r.URL.Path)
			w.Header().Add("Hx-Redirect", "/login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Also we want to make sure these pages that require
		// authentication aren't stored in browser cache
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)

	})
}
