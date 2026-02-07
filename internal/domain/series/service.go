package series

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type SeriesService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
