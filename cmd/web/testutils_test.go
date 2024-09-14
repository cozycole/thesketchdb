package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/go-playground/form/v4"
	"sketchdb.cozycole.net/internal/models/mocks"
)

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	return &application{
		errorLog:      log.New(io.Discard, "", 0),
		infoLog:       log.New(io.Discard, "", 0),
		videos:        &mocks.VideoModel{},
		creators:      &mocks.CreatorModel{},
		actors:        &mocks.ActorModel{},
		formDecoder:   formDecoder,
		templateCache: templateCache,
		debugMode:     true,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	// REMEMBER: change this to NewTLSServer once https is enabled
	ts := httptest.NewServer(h)

	// disable redirect-following by executing a function
	// for all 3xx responses, and http.ErrUselastResponse is returned
	// which forces the client to immediately return
	// the received response
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

func (ts *testServer) postMultipartForm(t *testing.T, urlPath string, fields map[string]string, files map[string]string) {
	// create
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	// write text fields
	for name, val := range fields {
		x, err := w.CreateFormField(name)
		if err != nil {
			t.Fatal(err)
		}
		x.Write([]byte(val))
	}

	for name, filepath := range files {
		file, err := os.Open(filepath)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		part, err := w.CreateFormFile(name, path.Base(filepath))
		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatal(err)
		}
	}

	w.Close()

	// rs, err := ts.Client().Post(ts.URL+urlPath, form)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// defer rs.Body.Close()
	// body, err := io.ReadAll(rs.Body)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// bytes.TrimSpace(body)

	// return rs.StatusCode, rs.Header, string(body)
}

func createMultipartForm(t *testing.T, fields map[string]string, files map[string]string) (*bytes.Buffer, string) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	// write text fields
	for name, val := range fields {
		x, err := w.CreateFormField(name)
		if err != nil {
			t.Fatal(err)
		}
		x.Write([]byte(val))
	}

	for name, filepath := range files {
		file, err := os.Open(filepath)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		part, err := w.CreateFormFile(name, path.Base(filepath))
		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatal(err)
		}
	}

	w.Close()
	return buf, w.FormDataContentType()

}
