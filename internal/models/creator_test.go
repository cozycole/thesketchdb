package models

import (
	"testing"
	"time"

	"sketchdb.cozycole.net/internal/assert"
)

func TestCreatorInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	m := CreatorModel{db}

	tests := []struct {
		name        string
		creatorName string
		url         string
		imgName     string
		imgExt      string
		birthDate   time.Time
		wantImgName string
	}{
		{
			name:        "Valid Entry",
			creatorName: "Travi$ Scott",
			imgName:     "travi-scott",
			imgExt:      ".jpg",
			birthDate:   time.Now(),
			wantImgName: "travi-scott-1.jpg",
		},
	}

	for _, tt := range tests {
		_, imgName, err := m.Insert(
			tt.creatorName,
			tt.url,
			tt.imgName,
			tt.imgExt,
			tt.birthDate,
		)
		assert.Equal(t, imgName, tt.wantImgName)
		assert.NilError(t, err)
	}
}
