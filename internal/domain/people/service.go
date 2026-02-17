package people

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type PersonService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
