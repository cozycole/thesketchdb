package characters

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type CharacterService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
