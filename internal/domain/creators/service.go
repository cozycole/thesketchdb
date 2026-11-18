package creators

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type CreatorService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
