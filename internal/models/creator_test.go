package models

// import (
// 	"testing"
// 	"time"
//
// 	"sketchdb.cozycole.net/internal/assert"
// )
//
// func TestCreatorInsert(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("models: skipping integration test")
// 	}
// 	db := newTestDB(t)
// 	m := CreatorModel{db}
//
// 	tests := []struct {
// 		name        string
// 		creatorName string
// 		url         string
// 		slug        string
// 		imgExt      string
// 		birthDate   time.Time
// 		wantImgName string
// 		wantSlug    string
// 	}{
// 		{
// 			name:        "Valid Entry",
// 			creatorName: "Travi$ Scott",
// 			slug:        "travi-scott",
// 			imgExt:      ".jpg",
// 			birthDate:   time.Now(),
// 			wantImgName: "travi-scott-1.jpg",
// 			wantSlug:    "travi-scott-1",
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		_, slug, imgName, err := m.Insert(
// 			tt.creatorName,
// 			tt.url,
// 			tt.slug,
// 			tt.imgExt,
// 			tt.birthDate,
// 		)
// 		assert.Equal(t, slug, tt.wantSlug)
// 		assert.Equal(t, imgName, tt.wantImgName)
// 		assert.NilError(t, err)
// 	}
// }
