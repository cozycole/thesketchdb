package people

import (
	"github.com/cozycole/thesketchdb/internal/fileStore"
	"github.com/cozycole/thesketchdb/internal/models"
)

type PersonService struct {
	Repos    models.Repositories
	ImgStore fileStore.FileStorageInterface
}
