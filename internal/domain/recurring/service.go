package recurring

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type RecurringService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
