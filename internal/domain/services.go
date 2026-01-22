package services

import "sketchdb.cozycole.net/internal/models"

type Services struct {
	Sketches SketchService
}

func NewServices(repos models.Repositories) Services {
	return Services{
		Sketches: SketchService{Repos: repos},
	}
}

type SketchService struct {
	Repos models.Repositories
}
