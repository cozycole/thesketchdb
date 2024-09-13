package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
)

func TestValidateAddCreatorForm(t *testing.T) {
	// store in memory valid and invalid images
	var emptyMap map[string]string
	var emptySlice []string
	validImg := createMultipartFileHeader(t, "./testdata/test-img.jpg")
	invalidImg := createMultipartFileHeader(t, "./testdata/test-img.webp")

	tests := []struct {
		name           string
		form           *addCreatorForm
		fieldErrors    map[string]string
		nonFieldErrors []string
	}{
		{
			name: "Valid Submission",
			form: &addCreatorForm{
				Name:            "Valid Title",
				URL:             "https://validurl.com",
				EstablishedDate: "2024-09-07",
				ProfileImage:    validImg,
			},
			fieldErrors:    emptyMap,
			nonFieldErrors: emptySlice,
		},
		{
			name: "Invalid Image",
			form: &addCreatorForm{
				Name:            "Valid Title",
				URL:             "https://validurl.com",
				EstablishedDate: "2024-09-07",
				ProfileImage:    invalidImg,
			},
			fieldErrors: map[string]string{
				"profileImg": "Uploaded file must be jpg or png",
			},
			nonFieldErrors: emptySlice,
		},
		{
			name: "Blank fields",
			form: &addCreatorForm{
				Name:            "",
				URL:             "",
				EstablishedDate: "",
				ProfileImage:    nil,
			},
			fieldErrors: map[string]string{
				"name":            "This field cannot be blank",
				"url":             "This field cannot be blank",
				"establishedDate": "This field cannot be blank",
				"profileImg":      "Please upload an image",
			},
			nonFieldErrors: emptySlice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateAddCreatorForm(tt.form)
			assert.DeepEqual(t, tt.form.FieldErrors, tt.fieldErrors)
			assert.DeepEqual(t, tt.form.NonFieldErrors, tt.nonFieldErrors)
		})
	}
}

func createMultipartFileHeader(t *testing.T, filePath string) *multipart.FileHeader {
	t.Helper()

	file, err := os.Open(filePath)
	if err != nil {
		t.Error(err)
		return nil
	}
	defer file.Close()

	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)

	// create a new form and create a new file field
	formWriter := multipart.NewWriter(buffWriter)
	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Error(err)
		return nil
	}

	if _, err := io.Copy(formPart, file); err != nil {
		t.Error(err)
		return nil
	}

	formWriter.Close()

	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())

	// read the form components with max stored memory of 1MB
	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		t.Error(err)
		return nil
	}

	// return the multipart file header
	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		t.Error("multipart file not exists")
		return nil
	}

	return files[0]
}
