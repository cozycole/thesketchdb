package pipeline

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type PipelineService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
