package characters

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type CharacterService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
