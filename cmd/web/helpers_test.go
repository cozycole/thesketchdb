package main

import (
	"net/http"
	"path"
	"testing"

	"sketchdb.cozycole.net/internal/assert"
	"sketchdb.cozycole.net/internal/utils"
)

// Checking that the decodePostForm function correctly
// marshals the request object into a addCreatorForm struct
func TestCreatorDecodePostForm(t *testing.T) {
	fields := map[string]string{
		"name":            "Test Name",
		"url":             "www.testsite.com",
		"establishedDate": "2024-09-10",
	}
	filepath := "./testdata/test-img.jpg"

	files := map[string]string{
		"profileImg": filepath,
	}

	buf, contentType, err := utils.CreateMultipartForm(fields, files)
	if err != nil {
		t.Fatal(err)
		return
	}

	r, err := http.NewRequest("POST", "/test/postform", buf)
	if err != nil {
		t.Fatal(err)
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

	buf, contentType, err = utils.CreateMultipartForm(fields, files)
	if err != nil {
		t.Fatal(err)
		return
	}

	r, err = http.NewRequest("POST", "/test/postform", buf)
	if err != nil {
		t.Fatal(err)
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

func TestVideoDecodePostForm(t *testing.T) {
	fields := map[string]string{
		"title":            "Test Name",
		"url":              "www.testsite.com",
		"uploadDate":       "2024-09-10",
		"rating":           "r",
		"creator":          "1",
		"peopleId[0]":      "1",
		"peopleId[1]":      "2",
		"peopleId[2]":      "3",
		"peopleText[0]":    "Tim",
		"peopleText[1]":    "Nathan",
		"peopleText[2]":    "James",
		"characterId[0]":   "4",
		"characterId[1]":   "5",
		"characterId[2]":   "6",
		"characterText[0]": "David",
		"characterText[1]": "Dave",
		"characterText[2]": "Davey",
	}

	filepath := "./testdata/test-img.jpg"
	filepath1 := "./testdata/test-thumbnail.jpg"
	filepath2 := "./testdata/test-img2.jpg"
	files := map[string]string{
		"thumbnail":             filepath,
		"characterThumbnail[0]": filepath1,
		"characterThumbnail[2]": filepath2,
	}

	buf, contentType, err := utils.CreateMultipartForm(fields, files)
	if err != nil {
		t.Fatal(err)
		return
	}

	r, err := http.NewRequest("POST", "/test/postform", buf)
	if err != nil {
		t.Fatal(err)
		return
	}
	r.Header.Add("content-type", contentType)

	app := newTestApplication(t)

	t.Run("Correct Form", func(t *testing.T) {
		var form addVideoForm

		app.decodePostForm(r, &form)

		assert.Equal(t, form.Title, fields["title"])
		assert.Equal(t, form.URL, fields["url"])
		assert.Equal(t, form.UploadDate, fields["uploadDate"])
		assert.Equal(t, form.Rating, fields["rating"])
		assert.Equal(t, form.CreatorID, 1)
		assert.DeepEqual(t, form.PersonIDs, []int{1, 2, 3})
		assert.DeepEqual(t, form.PersonInputs, []string{"Tim", "Nathan", "James"})
		assert.DeepEqual(t, form.CharacterIDs, []int{4, 5, 6})
		assert.DeepEqual(t, form.CharacterInputs, []string{"David", "Dave", "Davey"})
		assert.Equal(t, form.Thumbnail.Filename, path.Base(filepath))
		assert.Equal(t, len(form.CharacterThumbnails), 3)
		assert.Equal(t, form.CharacterThumbnails[0].Filename, path.Base(filepath1))
		assert.Equal(t, form.CharacterThumbnails[2].Filename, path.Base(filepath2))
	})
}
