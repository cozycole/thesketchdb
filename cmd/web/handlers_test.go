package main

import (
	"fmt"
	"net/http"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
)

func TestCreatorAddPost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	const (
		validName      = "Test Name"
		validUrl       = "www.testurl.com"
		validDate      = "2024-09-10"
		validImgPath   = "./testdata/test-img.jpg"
		invalidImgPath = "./testdata/test-img.webp"
	)
	tests := []struct {
		testName string
		name     string
		url      string
		date     string
		imgPath  string
		wantCode int
	}{
		{
			testName: "Valid Submission",
			name:     validName,
			url:      validUrl,
			date:     validDate,
			imgPath:  validImgPath,
			wantCode: http.StatusSeeOther,
		},
		{
			testName: "Invalid Image",
			name:     validName,
			url:      validUrl,
			date:     validDate,
			imgPath:  invalidImgPath,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			testName: "Blank Submission",
			name:     "",
			url:      "",
			date:     "",
			imgPath:  invalidImgPath,
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			fields := map[string]string{
				"name":            tt.name,
				"url":             tt.url,
				"establishedDate": tt.date,
			}
			files := map[string]string{
				"profileImg": tt.imgPath,
			}
			code, _, body := ts.postMultipartForm(t, "/creator/add", fields, files)
			assert.Equal(t, code, tt.wantCode)

			// ensure inputs are returned in the form on 422
			if tt.wantCode == http.StatusUnprocessableEntity {
				assert.StringContains(t, body, fmt.Sprintf(`value="%s"`, tt.name))
				assert.StringContains(t, body, fmt.Sprintf(`value="%s"`, tt.url))
				assert.StringContains(t, body, fmt.Sprintf(`value="%s"`, tt.date))
			}
		})
	}
}

func TestPersonAddPost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	const (
		validFirst     = "Brad"
		validLast      = "Pitt"
		validDate      = "2024-09-10"
		validImgPath   = "./testdata/test-img.jpg"
		invalidImgPath = "./testdata/test-img.webp"
	)
	tests := []struct {
		testName string
		first    string
		last     string
		date     string
		imgPath  string
		wantCode int
	}{
		{
			testName: "Valid Submission",
			first:    validFirst,
			last:     validLast,
			date:     validDate,
			imgPath:  validImgPath,
			wantCode: http.StatusSeeOther,
		},
		{
			testName: "Invalid Image",
			first:    validFirst,
			last:     validLast,
			date:     validDate,
			imgPath:  invalidImgPath,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			testName: "Blank Submission",
			first:    "",
			last:     "",
			date:     "",
			imgPath:  invalidImgPath,
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			fields := map[string]string{
				"first":     tt.first,
				"last":      tt.last,
				"birthDate": tt.date,
			}
			files := map[string]string{
				"profileImg": tt.imgPath,
			}
			code, _, body := ts.postMultipartForm(t, "/person/add", fields, files)
			assert.Equal(t, code, tt.wantCode)

			// ensure inputs are returned in the form on 422
			if tt.wantCode == http.StatusUnprocessableEntity {
				assert.StringContains(t, body, fmt.Sprintf(`value="%s"`, tt.first))
				assert.StringContains(t, body, fmt.Sprintf(`value="%s"`, tt.last))
				assert.StringContains(t, body, fmt.Sprintf(`value="%s"`, tt.date))
			}
		})
	}
}

func TestVideoAddPost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()
	validPeople := []string{"1", "2", "3"}

	const (
		validTitle     = "Title"
		validUrl       = "www.url.com"
		validDate      = "2024-09-10"
		validRating    = "r"
		validImgPath   = "./testdata/test-img.jpg"
		invalidImgPath = "./testdata/test-img.webp"
		validCreator   = "1"
	)
	tests := []struct {
		testName string
		title    string
		url      string
		rating   string
		date     string
		imgPath  string
		creator  string
		people   []string
		wantCode int
	}{
		{
			testName: "Valid Submission",
			title:    validTitle,
			url:      validUrl,
			rating:   validRating,
			date:     validDate,
			imgPath:  validImgPath,
			creator:  validCreator,
			people:   validPeople,
			wantCode: http.StatusSeeOther,
		},
		{
			testName: "Invalid Image",
			title:    validTitle,
			url:      validUrl,
			rating:   validRating,
			date:     validDate,
			imgPath:  invalidImgPath,
			creator:  validCreator,
			people:   validPeople,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			testName: "Blank Submission",
			title:    "",
			url:      "",
			rating:   "",
			date:     "",
			imgPath:  "",
			creator:  "",
			people:   []string{"", "", ""},
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			fields := map[string]string{
				"title":      tt.title,
				"url":        tt.url,
				"rating":     tt.rating,
				"uploadDate": tt.date,
				"creator":    tt.creator,
				"people[0]":  tt.people[0],
				"people[1]":  tt.people[1],
				"people[2]":  tt.people[2],
			}
			files := map[string]string{
				"thumbnail": tt.imgPath,
			}
			code, _, body := ts.postMultipartForm(t, "/video/add", fields, files)
			fmt.Print(body)
			assert.Equal(t, code, tt.wantCode)

			// ensure inputs are returned in the form on 422
			if tt.wantCode == http.StatusUnprocessableEntity {
				assert.StringContains(t, body, fmt.Sprintf(`value="%s"`, tt.title))
				assert.StringContains(t, body, fmt.Sprintf(`value="%s"`, tt.url))
				assert.StringContains(t, body, fmt.Sprintf(`value="%s"`, tt.date))
				assert.StringContains(t, body, fmt.Sprintf(`value="%s"`, tt.rating))
			}
		})
	}
}
