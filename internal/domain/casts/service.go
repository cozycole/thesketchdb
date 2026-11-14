package casts

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type CastService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
