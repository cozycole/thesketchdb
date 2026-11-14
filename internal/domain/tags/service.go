package tags

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type TagsService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
