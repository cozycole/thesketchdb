package sketches

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type SketchService struct {
	Repos        models.Repositories
	ImgStore     fileStore.FileStorageInterface
	ArchiveStore fileStore.FileStorageInterface
}
