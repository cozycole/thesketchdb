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
		name        string
		title       string
		url         string
		rating      string
		imgName     string
		imgExt      string
		birthDate   time.Time
		wantImgName string
	}{
		{
			name:        "Valid Entry",
			title:       "Test VIDEO",
			url:         "www.testurl.com/video",
			rating:      "R",
			imgName:     "test-video",
			imgExt:      ".jpg",
			birthDate:   time.Now(),
			wantImgName: "test-video-1.jpg",
		},
	}

	for _, tt := range tests {
		_, imgName, err := m.Insert(
			tt.title,
			tt.url,
			tt.rating,
			tt.imgName,
			tt.imgExt,
			tt.birthDate,
		)
		assert.Equal(t, imgName, tt.wantImgName)
		assert.NilError(t, err)
	}
}
