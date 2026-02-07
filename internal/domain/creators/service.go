package creators

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type CreatorService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
