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
