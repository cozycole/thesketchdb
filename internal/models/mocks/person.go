package mocks

import (
	"time"

	"sketchdb.cozycole.net/internal/models"
)

var mockPerson = &models.Person{
	ID:         1,
	First:      "Brad",
	Last:       "Pitt",
	ProfileImg: "brad-pitt-1.jpg",
	BirthDate:  time.Now(),
}

type PersonModel struct{}

func (m *PersonModel) Get(id int) (*models.Person, error) {
	return mockPerson, nil
}

func (m *PersonModel) Exists(id int) (bool, error) {
	switch id {
	case 1, 2, 3:
		return true, nil
	default:
		return false, nil
	}
}

func (m *PersonModel) Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, error) {
	return 1, "brad-pitt-1.jpg", nil
}
