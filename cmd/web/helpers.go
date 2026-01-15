package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
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

// func (app *application) readIDParam(r *http.Request) (int64, error) {
// 	// When httprouter is parsing a request, any interpolated URL parameters will be
// 	// stored in the request context. We can use the ParamsFromContext() function to
// 	// retrieve a slice containing these parameter names and values.
// 	params := httprouter.ParamsFromContext(r.Context())
// 	// We can then use the ByName() method to get the value of the "id" parameter from
// 	// the slice. In our project all movies will have a unique positive integer ID, but
// 	// the value returned by ByName() is always a string. So we try to convert it to a
// 	// base 10 integer (with a bit size of 64). If the parameter couldn't be converted,
// 	// or is less than 1, we know the ID is invalid so we use the http.NotFound()
// 	// function to return a 404 Not Found response.
// 	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
// 	if err != nil || id < 1 {
// 		return 0, errors.New("invalid id parameter")
// 	}
// 	return id, nil
// }

// *** API HELPERS ***

func (app *application) readIDParam(r *http.Request) (int, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return int(id), nil
}

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Encode the data to JSON, returning the error if there was one.
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')
	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	// Note that it's OK if the provided header map is nil. Go doesn't throw an error
	// if you try to range over (or generally, read from) a nil map.
	maps.Insert(w.Header(), maps.All(headers))
	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the json.Decoder, and call the DisallowUnknownFields() method on it
	// before decoding. This means that if the JSON from the client now includes any
	// field which cannot be mapped to the target destination, the decoder will return
	// an error instead of just ignoring the field.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		// If there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		// A json.InvalidUnmarshalError error will be returned if we pass something
		// that is not a non-nil pointer to Decode(). We catch this and panic,
		// rather than returning an error to our handler. At the end of this chapter
		// we'll talk about panicking versus returning errors, and discuss why it's an
		// appropriate thing to do in this specific situation.
		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown
		// field "<name>"". We check for this, extract the field name from the error,
		// and interpolate it into our custom error message. Note that there's an open
		// issue at https://github.com/golang/go/issues/29035 regarding turning this
		// into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// Use the errors.As() function to check whether the error has the type
		// *http.MaxBytesError. If it does, then it means the request body exceeded our
		// size limit of 1MB and we return a clear error message.
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		// For anything else, return the error message as-is.
		default:
			return err
		}
	}
	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
