package sketches

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type SketchService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
