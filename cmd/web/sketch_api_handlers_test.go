package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func TestCreateSketchAPI(t *testing.T) {
	app := newTestApplication(t)
	db := models.NewTestDb(t)
	app.services = NewServices(
		models.Repositories{
			Creators: &models.CreatorModel{DB: db},
			Sketches: &models.SketchModel{DB: db},
			Shows:    &models.ShowModel{DB: db},
		}, app.fileStorage)

	_ = insertTestCreator(t, &models.CreatorModel{DB: db})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.createSketchAPI)

	// test empty form
	form, formDataType, err := utils.CreateMultipartForm(
		map[string]string{},
		map[string]string{},
	)

	req := httptest.NewRequest(http.MethodPost, "/sketch", form)
	req.Header.Set("Content-Type", formDataType)
	handler.ServeHTTP(rr, req)

	// t.Log(rr.Body)
	assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)

	// test minimum viable form
	form, formDataType, err = utils.CreateMultipartForm(
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

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/sketch", form)
	req.Header.Set("Content-Type", formDataType)

	handler.ServeHTTP(rr, req)

	// t.Log(rr.Body)
	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestUpdateSketchAPI(t *testing.T) {
	app := newTestApplication(t)
	db := models.NewTestDb(t)
	app.services = NewServices(
		models.Repositories{
			Creators: &models.CreatorModel{DB: db},
			Sketches: &models.SketchModel{DB: db},
			Shows:    &models.ShowModel{DB: db},
		}, app.fileStorage)

	_ = insertTestCreator(t, &models.CreatorModel{DB: db})
	_ = insertTestCreator(t, &models.CreatorModel{DB: db})
	_ = insertTestSketch(t, &models.SketchModel{DB: db})

	handler := http.HandlerFunc(app.updateSketchAPI)

	// test empty form
	form, formDataType, err := utils.CreateMultipartForm(
		map[string]string{},
		map[string]string{},
	)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/sketch/1", form)
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", formDataType)

	handler.ServeHTTP(rr, req)

	t.Log(rr.Body)
	assert.Equal(t, rr.Code, http.StatusUnprocessableEntity)

	form, formDataType, err = utils.CreateMultipartForm(
		map[string]string{
			"title":      "Update Test",
			"uploadDate": "2000-06-20",
			"creatorId":  "2",
		},
		map[string]string{
			"thumbnail": "./testdata/test-profile-578x850.jpg",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/api/sketch/1", form)
	req.SetPathValue("id", "1")
	req.Header.Set("Content-Type", formDataType)

	handler.ServeHTTP(rr, req)

	t.Log(rr.Body)
	assert.Equal(t, rr.Code, http.StatusOK)
}
