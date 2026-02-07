package shows

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type ShowService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
