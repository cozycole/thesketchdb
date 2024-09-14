package main

import (
	"net/http"
	"os"
	"path"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
)

// Checking that the decodePostForm function correctly
// marshals the request object into a addCreatorForm struct
func TestDecodePostForm(t *testing.T) {
	fields := map[string]string{
		"name":            "Test Name",
		"url":             "www.testsite.com",
		"establishedDate": "2024-09-10",
	}

	currDir, _ := os.Getwd()
	filepath := "./testdata/test-img.jpg"
	fullpath := path.Join(currDir, filepath)
	files := map[string]string{
		"profileImg": fullpath,
	}

	buf, contentType := createMultipartForm(t, fields, files)

	r, err := http.NewRequest("POST", "/test/postform", buf)
	if err != nil {
		t.Error(err)
		return
	}
	r.Header.Add("content-type", contentType)

	app := newTestApplication(t)

	t.Run("CorrectForm ExtraFields", func(t *testing.T) {
		var form addCreatorForm

		app.decodePostForm(r, &form)
		assert.Equal(t, form.Name, fields["name"])
		assert.Equal(t, form.URL, fields["url"])
		assert.Equal(t, form.EstablishedDate, fields["establishedDate"])
		assert.Equal(t, form.ProfileImage.Filename, path.Base(filepath))
	})
	// missing name as well
	fields = map[string]string{
		"url":             "www.testsite.com",
		"establishedDate": "2024-09-10",
	}
	files = map[string]string{}

	buf, contentType = createMultipartForm(t, fields, files)

	r, err = http.NewRequest("POST", "/test/postform", buf)
	if err != nil {
		t.Error(err)
		return
	}
	r.Header.Add("content-type", contentType)

	t.Run("No Image", func(t *testing.T) {
		var form addCreatorForm

		app.decodePostForm(r, &form)
		assert.Equal(t, form.ProfileImage, nil)
		assert.Equal(t, form.Name, "")
	})
}
