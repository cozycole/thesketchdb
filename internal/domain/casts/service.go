package casts

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type CastService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
