package models

import (
	"testing"
	"time"

	"sketchdb.cozycole.net/internal/assert"
	"sketchdb.cozycole.net/internal/utils"
)

func TestSketchInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	// We need to have existing characters for Insert to work
	if err := restoreDbScript("../../sql/testdata/test1.sql"); err != nil {
		t.Fatal(err)
	}

	m := SketchModel{db, 10}
	uploadDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		sketch Sketch
	}{
		{
			name: "Valid Entry",
			sketch: Sketch{
				Title:      "Test sketch",
				URL:        "www.testurl.com/sketch",
				Rating:     "R",
				Slug:       "test-sketch",
				UploadDate: &uploadDate,
				Creator:    &Creator{ID: 1},
				Cast: []*CastMember{
					{
						Position:  utils.GetIntPtr(0),
						Actor:     &Person{ID: utils.GetIntPtr(1)},
						Character: &Character{ID: utils.GetIntPtr(1)},
					},
					{
						Position:  utils.GetIntPtr(1),
						Actor:     &Person{ID: utils.GetIntPtr(2)},
						Character: &Character{ID: utils.GetIntPtr(2)},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		err := m.Insert(&tt.sketch)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSketchGet(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	if err := restoreDbScript("../../sql/testdata/test1.sql"); err != nil {
		t.Fatal(err)
	}

	m := SketchModel{db, 10}
	v, err := m.Get(1)
	if err != nil {
		t.Fatal(err)
	}

	uploadDate := time.Date(2008, 9, 8, 0, 0, 0, 0, time.UTC)
	emptyDescription := new(string)
	*emptyDescription = ""

	tests := []struct {
		name        string
		title       string
		url         string
		thumbnail   string
		uploadDate  time.Time
		pgRating    string
		description *string
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
			description: emptyDescription,
			creatorName: "nathanfielder",
			castSize:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, v.Title, tt.title)
			assert.Equal(t, v.URL, tt.url)
			assert.Equal(t, v.ThumbnailName, tt.thumbnail)
			assert.Equal(t, *v.UploadDate, tt.uploadDate)
			assert.Equal(t, v.Rating, tt.pgRating)
			assert.EqualPointer(t, v.Description, tt.description)
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

	m := SketchModel{db, 10}
	members, err := m.GetCastMembers(1)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(members); i++ {
		test := tests[i]
		member := *members[i]

		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, member.Actor.First, &test.first)
			assert.Equal(t, member.Actor.Last, &test.last)
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
