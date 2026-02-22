package quotes

import (
	"sketchdb.cozycole.net/internal/fileStore"
	"sketchdb.cozycole.net/internal/models"
)

type QuoteService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
