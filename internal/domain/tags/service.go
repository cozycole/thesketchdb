package tags

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type TagsService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
