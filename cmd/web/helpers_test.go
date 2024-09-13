package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
)

// Checking that the decodePostForm function correctly
// marshals the request object into a addCreatorForm struct
func TestDecodePostForm(t *testing.T) {
	currDir, _ := os.Getwd()
	filepath := "cmd/web/testdata/test-img.jpg"
	fullpath := path.Join(currDir, filepath)

	file, _ := os.Open(fullpath)
	defer file.Close()

	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	name := "Test Name"
	url := "www.testsite.com"
	establishedDate := "2024-09-10"

	x, _ := w.CreateFormField("name")
	x.Write([]byte(name))
	x, _ = w.CreateFormField("url")
	x.Write([]byte(url))
	x, _ = w.CreateFormField("establishedDate")
	x.Write([]byte(establishedDate))
	x, _ = w.CreateFormField("extraField")
	x.Write([]byte("randomExtra"))

	part, err := w.CreateFormFile("profileImg", path.Base(filepath))
	if err != nil {
		t.Error(err)
		return
	}

	io.Copy(part, file)

	r, err := http.NewRequest("POST", "/test/postform", buf)
	r.Header.Add("content-type", w.FormDataContentType())
	if err != nil {
		t.Error(err)
		return
	}
	w.Close()

	app := newTestApplication(t)

	t.Run("CorrectForm ExtraFields", func(t *testing.T) {
		var form addCreatorForm

		app.decodePostForm(r, &form)
		assert.Equal(t, form.Name, name)
		assert.Equal(t, form.URL, url)
		assert.Equal(t, form.EstablishedDate, establishedDate)
		assert.Equal(t, form.ProfileImage.Filename, path.Base(filepath))
	})

	buf = new(bytes.Buffer)
	w = multipart.NewWriter(buf)

	x, _ = w.CreateFormField("url")
	x.Write([]byte(url))
	x, _ = w.CreateFormField("establishedDate")
	x.Write([]byte(establishedDate))

	r, err = http.NewRequest("POST", "/test/postform", buf)
	r.Header.Add("content-type", w.FormDataContentType())
	if err != nil {
		t.Error(err)
		return
	}
	w.Close()

	t.Run("No Image", func(t *testing.T) {
		var form addCreatorForm

		app.decodePostForm(r, &form)
		assert.Equal(t, form.ProfileImage, nil)
		assert.Equal(t, form.Name, "")
	})
}
