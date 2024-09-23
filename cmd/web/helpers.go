package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
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

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
// func (app *application) notFound(w http.ResponseWriter) {
// 	app.clientError(w, http.StatusNotFound)
// }

func (app *application) newTemplateData(_ *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		// IsAuthenticated: app.isAutheticated(r),
		// CSRFToken:       nosurf.Token(r),
	}
}

// see ~/go/pkg/mod/github.com/go-playground/form/v4@v4.2.1/README.md
// for doc on the form decoder
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {

		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}

	// the following checks if struct tag with key "img" exists and sets
	// the given field with the file of name struct tag value
	// don't do any type checking since form decoder already did it
	v := reflect.ValueOf(dst).Elem()
	structType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := structType.Field(i)

		if tagValue := fieldType.Tag.Get("img"); tagValue != "" {

			fileHeaders, ok := r.MultipartForm.File[tagValue]
			if ok && len(fileHeaders) > 0 {
				newVal := reflect.ValueOf(fileHeaders[0])
				field.Set(newVal)
			}
		}
	}

	return nil
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	// write template to buffer instead of straight to
	// the http.ResponseWriter
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer
	w.WriteHeader(status)
	buf.WriteTo(w)
}
