package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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
