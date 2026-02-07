package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/form/v4"
	imgmock "sketchdb.cozycole.net/internal/fileStore/mocks"
	"sketchdb.cozycole.net/internal/models"
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
		debugMode:     true,
	}
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

func insertTestCreator(t *testing.T, model models.CreatorModelInterface, overrides ...func(*models.Creator)) *models.Creator {
	creator := &models.Creator{
		Name:         ptr("Test Creator"),
		Slug:         ptr("test-creator"),
		ProfileImage: ptr("test-img.jpg"),
	}

	for _, override := range overrides {
		override(creator)
	}

	_, err := model.Insert(creator)
	if err != nil {
		t.Fatal(err)
	}

	return creator
}

func insertTestSketch(t *testing.T, model models.SketchModelInterface, overrides ...func(*models.Sketch)) *models.Sketch {
	sketch := &models.Sketch{
		Title:         ptr("Test Sketch"),
		Slug:          ptr("test-sketch"),
		ThumbnailName: ptr("test-img.jpg"),
	}

	for _, override := range overrides {
		override(sketch)
	}

	_, err := model.Insert(sketch)
	if err != nil {
		t.Fatal(err)
	}

	return sketch
}
