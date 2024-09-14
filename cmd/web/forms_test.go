package main

import (
	"testing"

	"sketchdb.cozycole.net/internal/assert"
)

// This is testing the validation of the struct that
// was unmarshalled by the decodePostForm function
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
