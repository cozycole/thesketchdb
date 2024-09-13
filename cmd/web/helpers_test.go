package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"testing"
)

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

	t.Run("CorrectFormWExtraFields", func(t *testing.T) {
		var form addCreatorForm

		app.decodePostForm(r, &form)
		if form.Name != name {
			t.Errorf("got %q; want %q", form.Name, name)
		}
		if form.URL != url {
			t.Errorf("got %q; want %q", form.URL, url)
		}
		if form.EstablishedDate != establishedDate {
			t.Errorf("got %q; want %q", form.EstablishedDate, establishedDate)
		}
		if form.ProfileImage.Filename != path.Base(filepath) {
			t.Errorf("got %q; want %q", form.ProfileImage.Filename, path.Base(filepath))
		}
	})
}
