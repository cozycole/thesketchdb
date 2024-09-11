package mocks

import (
	"time"

	"sketchdb.cozycole.net/internal/models"
)

var mockActor = &models.Actor{
	ID:         1,
	First:      "Brad",
	Last:       "Pitt",
	ProfileImg: "brad-pitt-1.jpg",
	BirthDate:  time.Now(),
}

type ActorModel struct{}

func (m *ActorModel) Get(id int) (*models.Actor, error) {
	return mockActor, nil
}

func (m *ActorModel) ExistsByName(fullname string) (int, error) {
	return 1, nil
}
