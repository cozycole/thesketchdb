package models

import (
	"testing"
	"time"

	"sketchdb.cozycole.net/internal/assert"
)

func TestPersonInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	m := PersonModel{db}

	tests := []struct {
		name        string
		first       string
		last        string
		imgName     string
		imgExt      string
		birthDate   time.Time
		wantImgName string
	}{
		{
			name:        "Valid Entry",
			first:       "Denis",
			last:        "O'Bell",
			imgName:     "denis-obell",
			imgExt:      ".jpg",
			birthDate:   time.Now(),
			wantImgName: "denis-obell-1.jpg",
		},
	}

	for _, tt := range tests {
		_, imgName, _, err := m.Insert(
			tt.first,
			tt.last,
			tt.imgName,
			tt.imgExt,
			tt.birthDate,
		)
		assert.NilError(t, err)
		assert.Equal(t, imgName, tt.wantImgName)
	}
}

func TestPersonSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	db := newTestDB(t)
	if err := restoreDbScript("../../sql/testdata/test1.sql"); err != nil {
		t.Fatal(err)
	}
	m := PersonModel{db}

	tests := []struct {
		name          string
		query         string
		wantSlug      string
		wantFirst     string
		wantLast      string
		wantImg       string
		wantBirthdate time.Time
	}{
		{
			name:          "Find Kyle",
			query:         "Ky",
			wantSlug:      "kyle-mooney-1",
			wantFirst:     "Kyle",
			wantLast:      "Mooney",
			wantImg:       "kyle-mooney-1.jpg",
			wantBirthdate: time.Date(1984, 9, 3, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		people, err := m.Search(tt.query)
		assert.NilError(t, err)
		if len(people) == 0 {
			t.Fatalf("No person found with query %s", tt.query)
		}

		person := people[0]
		assert.EqualPointer(t, person.First, &tt.wantFirst)
		assert.EqualPointer(t, person.Last, &tt.wantLast)
		assert.EqualPointer(t, person.Slug, &tt.wantSlug)
		assert.EqualPointer(t, person.ProfileImg, &tt.wantImg)
		assert.EqualPointer(t, person.BirthDate, &tt.wantBirthdate)
	}
}
