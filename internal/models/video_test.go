package models

import (
	"testing"
	"time"

	"sketchdb.cozycole.net/internal/assert"
)

func TestVideoInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	m := VideoModel{db, 10}

	tests := []struct {
		name      string
		title     string
		url       string
		rating    string
		slug      string
		imgExt    string
		birthDate time.Time
		wantSlug  string
		wantImg   string
	}{
		{
			name:      "Valid Entry",
			title:     "Test VIDEO",
			url:       "www.testurl.com/video",
			rating:    "R",
			slug:      "test-video",
			imgExt:    ".jpg",
			birthDate: time.Now(),
			wantSlug:  "test-video-1",
			wantImg:   "test-video-1.jpg",
		},
	}

	for _, tt := range tests {
		_, slug, imgName, err := m.Insert(
			tt.title,
			tt.url,
			tt.rating,
			tt.slug,
			tt.imgExt,
			tt.birthDate,
		)
		assert.Equal(t, slug, tt.wantSlug)
		assert.Equal(t, imgName, tt.wantImg)
		assert.NilError(t, err)
	}
}

func TestVideoGet(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	if err := restoreDbScript("../../sql/testdata/test1.sql"); err != nil {
		t.Fatal(err)
	}

	m := VideoModel{db, 10}
	v, err := m.Get(1)
	if err != nil {
		t.Fatal(err)
	}

	uploadDate := time.Date(2008, 9, 8, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		title       string
		url         string
		thumbnail   string
		uploadDate  time.Time
		pgRating    string
		description string
		creatorName string
		castSize    int
	}{
		{
			name:        "test1",
			title:       "Good Pals",
			url:         "https://www.youtube.com/watch?v=6aTqXkZHnQE",
			thumbnail:   "good-pals-1.jpg",
			uploadDate:  uploadDate,
			pgRating:    "PG",
			description: "",
			creatorName: "nathanfielder",
			castSize:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, v.Title, tt.title)
			assert.Equal(t, v.URL, tt.url)
			assert.Equal(t, v.Thumbnail, tt.thumbnail)
			assert.Equal(t, *v.UploadDate, tt.uploadDate)
			assert.Equal(t, v.Rating, tt.pgRating)
			assert.Equal(t, v.Description, tt.description)
			assert.Equal(t, v.Creator.Name, tt.creatorName)
			assert.Equal(t, len(v.Cast), tt.castSize)
		})
	}

}

func TestCastMembersGet(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)

	if err := restoreDbScript("../../sql/testdata/test1.sql"); err != nil {
		t.Fatal(err)
	}

	bday := time.Date(1983, 5, 13, 0, 0, 0, 0, time.UTC)
	desc := "this is the description"

	tests := []struct {
		name        string
		first       string
		last        string
		birthdate   *time.Time
		description *string
	}{
		{
			name:        "No Nils",
			first:       "Tim",
			last:        "Gilbert",
			birthdate:   &bday,
			description: &desc,
		},
		{
			name:        "Nils",
			first:       "James",
			last:        "Hartnett",
			birthdate:   nil,
			description: nil,
		},
	}

	m := VideoModel{db, 10}
	members, err := m.GetCastMembers(1)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(members); i++ {
		test := tests[i]
		member := *members[i]

		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, member.Actor.First, test.first)
			assert.Equal(t, member.Actor.Last, test.last)
			if member.Actor.BirthDate != nil {
				t.Log(bday)
				assert.Equal(t, *member.Actor.BirthDate, *test.birthdate)
			} else {
				assert.Equal(t, member.Actor.BirthDate, test.birthdate)
			}
			t.Log(member.Actor.Description)

			if member.Actor.Description != nil {
				assert.Equal(t, *member.Actor.Description, *test.description)
			} else {
				assert.Equal(t, member.Actor.Description, test.description)
			}

		})

	}

}
