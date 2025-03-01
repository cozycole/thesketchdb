package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"sketchdb.cozycole.net/internal/models"
)

var mimeToExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// 2 is inputed to look at the second frame for determining the Llongfile/Lshortfile and line number
	// for the logged output (since we don't want to log the line number here, but wherever it is called)
	app.errorLog.Output(2, trace)

	if app.debugMode {
		http.Error(w, trace, http.StatusInternalServerError)
	} else {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) badRequest(w http.ResponseWriter) {
	app.clientError(w, http.StatusBadRequest)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) unauthorized(w http.ResponseWriter) {
	app.clientError(w, http.StatusUnauthorized)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
// func (app *application) notFound(w http.ResponseWriter) {
// 	app.clientError(w, http.StatusNotFound)
// }

func (app *application) newTemplateData(r *http.Request) *templateData {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	var isEditor bool
	if ok {
		isEditor = user.Role == "admin" || user.Role == "editor"
	} else {
		isEditor = false
	}
	return &templateData{
		CurrentYear:     time.Now().Year(),
		ImageBaseUrl:    app.baseImgUrl,
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		Forms:           Forms{},
		IsEditor:        isEditor,
	}
}

func (app *application) render(w http.ResponseWriter, status int, page string, baseTemplate string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	// write template to buffer instead of straight to
	// the http.ResponseWriter
	var err error
	if baseTemplate != "" {
		err = ts.ExecuteTemplate(buf, baseTemplate, data)
	} else {
		err = ts.Execute(buf, data)
	}

	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
