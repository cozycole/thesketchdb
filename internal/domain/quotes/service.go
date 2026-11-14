package quotes

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type QuoteService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
