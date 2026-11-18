package recurring

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type RecurringService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
