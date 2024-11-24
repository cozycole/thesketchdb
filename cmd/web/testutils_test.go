package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/form/v4"
	imgmock "sketchdb.cozycole.net/internal/img/mocks"
	"sketchdb.cozycole.net/internal/models/mocks"
	"sketchdb.cozycole.net/internal/utils"
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
		fileStorage:   &imgmock.FileStorage{},
		formDecoder:   formDecoder,
		templateCache: templateCache,
		videos:        &mocks.VideoModel{},
		creators:      &mocks.CreatorModel{},
		people:        &mocks.PersonModel{},
		characters:    &mocks.CharacterModel{},
		debugMode:     true,
	}
}

func resetMocks(app *application) {
	app.videos = &mocks.VideoModel{}
	app.creators = &mocks.CreatorModel{}
	app.people = &mocks.PersonModel{}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(_ *testing.T, h http.Handler) *testServer {
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

func (ts *testServer) postMultipartForm(t *testing.T, urlPath string, fields map[string]string, files map[string]string) (int, http.Header, string) {
	buf, contentType, err := utils.CreateMultipartForm(fields, files)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := ts.Client().Post(ts.URL+urlPath, contentType, buf)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
