package main

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"sketchdb.cozycole.net/internal/models"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *Application) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scriptSrc := "script-src 'self' https://www.youtube.com"
		scriptSrcElem := "script-src https://static.cloudflareinsights.com"
		if app.Config.DebugMode {
			// browser sync injects script
			scriptSrc += " 'unsafe-inline'"
			scriptSrcElem += " 'unsafe-inline'"
		}

		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				scriptSrc+";"+
				scriptSrcElem+";"+
				"connect-src 'self' https://www.youtube.com https://*.digitaloceanspaces.com data:;"+
				"frame-src 'self' https://www.youtube.com;"+
				"style-src 'self' fonts.googleapis.com;"+
				"img-src 'self' data: blob: https://*.digitaloceanspaces.com https://i.ytimg.com https://*.ytimg.com https://www.gstatic.com;"+
				"font-src 'self' fonts.gstatic.com;",
		)

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
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

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Config.DebugMode {
			app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.SessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id != 0 {
			user, err := app.Users.GetById(id)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			user, err := app.Users.GetByToken(token, models.ScopeAPIKey)
			if err != nil {
				app.ErrorLog.Print(err)
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var (
	editorAdmin = []string{"editor", "admin"}
	adminOnly   = []string{"admin"}
)

func (app *Application) requireRoles(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(userContextKey).(*models.User)
			if !ok || !slices.Contains(roles, derefString(user.Role)) {
				app.unauthorized(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *Application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAutheticated(r) {
			app.SessionManager.Put(r.Context(), "postLoginRedirectURL", r.URL.Path)
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
