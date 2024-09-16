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
			code, _, body := ts.postMultipartForm(t, "/add/creator", fields, files)
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
