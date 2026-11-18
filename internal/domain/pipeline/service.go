package pipeline

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type PipelineService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
