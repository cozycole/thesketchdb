package models

import (
	"testing"
	"time"

	"sketchdb.cozycole.net/internal/assert"
)

func TestActorInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	db := newTestDB(t)
	m := ActorModel{db}

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
		_, imgName, err := m.Insert(
			tt.first,
			tt.last,
			tt.imgName,
			tt.imgExt,
			tt.birthDate,
		)
		assert.Equal(t, imgName, tt.wantImgName)
		assert.NilError(t, err)
	}
}
