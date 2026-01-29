package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func TestCreateTestAPI(t *testing.T) {
	app := newTestApplication(t)
	db := models.NewTestDb(t)
	app.services = NewServices(
		models.Repositories{
			Creators: &models.CreatorModel{DB: db},
			Sketches: &models.SketchModel{DB: db},
			Shows:    &models.ShowModel{DB: db},
		}, app.fileStorage)

	_ = insertTestCreator(t, &models.CreatorModel{DB: db})
	form, formDataType, err := utils.CreateMultipartForm(
		map[string]string{
			"title":      "Test",
			"uploadDate": "2000-06-20",
			"creatorId":  "1",
		},
		map[string]string{
			"thumbnail": "./testdata/test-profile-578x850.jpg",
		},
	)

	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/sketch", form)
	req.Header.Set("Content-Type", formDataType)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.createSketchAPI)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected code: %d", rr.Code)
	}

	t.Log(rr.Body)
}
