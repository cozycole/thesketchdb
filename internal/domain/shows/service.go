package shows

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type ShowService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
