package series

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type SeriesService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
