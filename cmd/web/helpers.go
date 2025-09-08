package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"sketchdb.cozycole.net/internal/models"
)

var mimeToExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(r *http.Request, w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// 2 is inputed to look at the second frame for determining the Llongfile/Lshortfile and line number
	// for the logged output (since we don't want to log the line number here, but wherever it is called)
	app.errorLog.Output(2, trace)

	isHxRequest := r.Header.Get("HX-Request") == "true"
	if isHxRequest {
		data := app.newTemplateData(r)
		data.Flash = flashMessage{
			Level:   "error",
			Message: "500 Internal Server Error",
		}
		app.render(r, w,
			http.StatusInternalServerError,
			"flash-message.gohtml",
			"flash-message",
			data)
	} else if app.debugMode {
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

func extractUrlParamIDs(idParams []string) []int {
	var ids []int
	for _, idStr := range idParams {
		id, err := strconv.Atoi(idStr)
		if nil == err && id > 0 {
			ids = append(ids, id)
		}
	}

	return ids
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	var isEditor, isAdmin bool
	if ok {
		isEditor = safeDeref(user.Role) == "admin" || safeDeref(user.Role) == "editor"
		isAdmin = safeDeref(user.Role) == "admin"
	} else {
		isEditor = false
	}
	return &templateData{
		CurrentYear:  time.Now().Year(),
		Assets:       app.assets,
		ImageBaseUrl: app.baseImgUrl,
		Forms:        Forms{},
		User:         user,
		Origin:       app.settings.origin,
		IsEditor:     isEditor,
		IsAdmin:      isAdmin,
	}
}

func safeDeref[T any](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	var zero T
	return zero
}

func extractYouTubeVideoID(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	queryParams := parsedURL.Query()
	videoID := queryParams.Get("v")
	if videoID == "" {
		return "", fmt.Errorf("video ID not found in URL")
	}

	return videoID, nil
}

func parseTimestamp(input string) (int, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return 0, errors.New("invalid timestamp format, expected M:SS")
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %w", err)
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds: %w", err)
	}

	if seconds < 0 || seconds > 59 {
		return 0, errors.New("seconds must be between 0 and 59")
	}

	total := minutes*60 + seconds
	return total, nil
}

func secondsToMMSS(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

func ptr[T any](v T) *T {
	return &v
}

func convertStringsToInts(strs []string) ([]int, error) {
	ints := make([]int, len(strs))
	for i, s := range strs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at index %d: %v", i, err)
		}
		ints[i] = n
	}
	return ints, nil
}

func (app *application) isAutheticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

func (app *application) render(r *http.Request, w http.ResponseWriter, status int, page string, baseTemplate string, data any) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		r.Header.Set("HX-Request", "false")
		app.serverError(r, w, err)
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
		r.Header.Set("HX-Request", "false")
		app.serverError(r, w, err)
		return
	}

	// If the template is written to the buffer
	w.WriteHeader(status)
	buf.WriteTo(w)
}
