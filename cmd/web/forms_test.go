package main

import (
	"mime/multipart"
	"testing"
	"time"

	"sketchdb.cozycole.net/internal/assert"
	"sketchdb.cozycole.net/internal/utils"
)

// This is testing the validation of the struct that
// was unmarshalled by the decodePostForm function
func TestValidateAddCreatorForm(t *testing.T) {
	// store in memory valid and invalid images
	var emptyMap map[string]string
	var emptySlice []string
	validImg, err := utils.CreateMultipartFileHeader("./testdata/test-img.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	invalidImg, err := utils.CreateMultipartFileHeader("./testdata/test-img.webp")
	if err != nil {
		t.Fatal(err)
		return
	}
	app := newTestApplication(t)

	tests := []struct {
		name           string
		form           *creatorForm
		fieldErrors    map[string]string
		nonFieldErrors []string
	}{
		{
			name: "Valid Submission",
			form: &creatorForm{
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
			form: &creatorForm{
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
			form: &creatorForm{
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
			app.validateCreatorForm(tt.form)
			assert.DeepEqual(t, tt.form.FieldErrors, tt.fieldErrors)
			assert.DeepEqual(t, tt.form.NonFieldErrors, tt.nonFieldErrors)
		})
	}
}

func TestValidateAddPersonForm(t *testing.T) {
	// store in memory valid and invalid images
	var emptyMap map[string]string
	var emptySlice []string
	validImg, err := utils.CreateMultipartFileHeader("./testdata/test-img.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	invalidImg, err := utils.CreateMultipartFileHeader("./testdata/test-img.webp")
	if err != nil {
		t.Fatal(err)
		return
	}
	app := newTestApplication(t)

	tests := []struct {
		name           string
		form           *personForm
		fieldErrors    map[string]string
		nonFieldErrors []string
	}{
		{
			name: "Valid Submission",
			form: &personForm{
				First:        "Brad",
				Last:         "Pitt",
				BirthDate:    "2024-09-07",
				ProfileImage: validImg,
			},
			fieldErrors:    emptyMap,
			nonFieldErrors: emptySlice,
		},
		{
			name: "Invalid Image",
			form: &personForm{
				First:        "Brad",
				Last:         "Pitt",
				BirthDate:    "2024-09-07",
				ProfileImage: invalidImg,
			},
			fieldErrors: map[string]string{
				"profileImg": "Uploaded file must be jpg or png",
			},
			nonFieldErrors: emptySlice,
		},
		{
			name: "Blank fields",
			form: &personForm{
				First:        "",
				Last:         "",
				BirthDate:    "",
				ProfileImage: nil,
			},
			fieldErrors: map[string]string{
				"first":      "This field cannot be blank",
				"last":       "This field cannot be blank",
				"birthDate":  "This field cannot be blank",
				"profileImg": "Please upload an image",
			},
			nonFieldErrors: emptySlice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.validatePersonForm(tt.form)
			assert.DeepEqual(t, tt.form.FieldErrors, tt.fieldErrors)
			assert.DeepEqual(t, tt.form.NonFieldErrors, tt.nonFieldErrors)
		})
	}
}

func TestValidateAddSketchForm(t *testing.T) {

	// store in memory valid and invalid images
	var emptyMap map[string]string
	var emptySlice []string
	validImg, err := utils.CreateMultipartFileHeader("./testdata/test-img.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	validImg1, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	validImg2, err := utils.CreateMultipartFileHeader("./testdata/test-img2.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}

	invalidImg, err := utils.CreateMultipartFileHeader("./testdata/test-img.webp")
	if err != nil {
		t.Fatal(err)
		return
	}

	app := newTestApplication(t)

	tests := []struct {
		name           string
		form           *addSketchForm
		fieldErrors    map[string]string
		nonFieldErrors []string
	}{
		{
			name: "Valid Submission",
			form: &addSketchForm{
				Title:               "Sketch Title",
				URL:                 "www.url.com",
				Rating:              "pg-13",
				UploadDate:          "2024-11-24",
				Thumbnail:           validImg,
				CreatorID:           1,
				PersonIDs:           []int{1, 2, 3},
				CharacterIDs:        []int{1, 2, 3},
				CharacterThumbnails: []*multipart.FileHeader{validImg, validImg1, validImg2},
			},
			fieldErrors:    emptyMap,
			nonFieldErrors: emptySlice,
		},
		{
			name: "Invalid Image",
			form: &addSketchForm{
				Title:               "Sketch Title",
				URL:                 "www.url.com",
				Rating:              "pg-13",
				UploadDate:          "2024-11-24",
				Thumbnail:           validImg,
				CreatorID:           1,
				PersonIDs:           []int{1, 2, 3},
				CharacterIDs:        []int{1, 2, 3},
				CharacterThumbnails: []*multipart.FileHeader{validImg, invalidImg, validImg2},
			},
			fieldErrors: map[string]string{
				"characterThumbnail[1]": "Uploaded file must be jpg or png",
			},
			nonFieldErrors: emptySlice,
		},
		{
			name: "Blank fields",
			form: &addSketchForm{
				Title:               "",
				URL:                 "",
				Rating:              "",
				UploadDate:          "",
				Thumbnail:           nil,
				CreatorID:           0,
				PersonIDs:           nil,
				CharacterIDs:        nil,
				CharacterThumbnails: nil,
			},
			fieldErrors: map[string]string{
				"title":      "This field cannot be blank",
				"url":        "This field cannot be blank",
				"rating":     "This field cannot be blank",
				"uploadDate": "This field cannot be blank",
				"thumbnail":  "Please upload an image",
				"creator":    "This field cannot be blank",
			},
			nonFieldErrors: emptySlice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.validateSketchForm(tt.form)
			assert.DeepEqual(t, tt.form.FieldErrors, tt.fieldErrors)
			assert.DeepEqual(t, tt.form.NonFieldErrors, tt.nonFieldErrors)
		})
	}
}

func TestConvertFormToSketch(t *testing.T) {
	validThumbnail, err := utils.CreateMultipartFileHeader("./testdata/test-thumbnail.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	validImg, err := utils.CreateMultipartFileHeader("./testdata/test-img.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	validImg2, err := utils.CreateMultipartFileHeader("./testdata/test-img2.jpg")
	if err != nil {
		t.Fatal(err)
		return
	}
	title := "Test Title"
	url := "www.test.com"
	rating := "pg"
	uploadDateStr := "2024-11-30"
	vidForm := addSketchForm{
		Title:               title,
		URL:                 url,
		Rating:              rating,
		UploadDate:          uploadDateStr,
		Thumbnail:           validThumbnail,
		PersonIDs:           []int{1, 2},
		PersonInputs:        []string{"Tim", "James"},
		CharacterIDs:        []int{1, 2},
		CharacterInputs:     []string{"Davey D", "Sammy S"},
		CharacterThumbnails: []*multipart.FileHeader{validImg, validImg2},
	}

	v, err := convertFormToSketch(&vidForm)
	if err != nil {
		t.Fatalf("unable to convert sketch: %s", err)
	}

	assert.Equal(t, v.Title, title)
	assert.Equal(t, v.URL, url)
	assert.Equal(t, v.Rating, rating)
	assert.Equal(t, *v.UploadDate, time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, *v.Cast[0].Actor.ID, 1)
	assert.Equal(t, *v.Cast[0].Character.ID, 1)
	assert.Equal(t, *v.Cast[1].Actor.ID, 2)
	assert.Equal(t, *v.Cast[1].Character.ID, 2)
	assert.Equal(t, v.Cast[0].ThumbnailFile.Filename, "test-img.jpg")
	assert.Equal(t, v.Cast[1].ThumbnailFile.Filename, "test-img2.jpg")
}
