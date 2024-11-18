package mocks

import (
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

var mockPerson = &models.Person{
	ID:         utils.GetIntPtr(1),
	First:      utils.GetStringPtr("Brad"),
	Last:       utils.GetStringPtr("Pitt"),
	ProfileImg: utils.GetStringPtr("brad-pitt-1.jpg"),
	BirthDate:  utils.GetTimePtr(time.Now()),
}

type PersonModel struct{}

func (m *PersonModel) GetBySlug(slug string) (*models.Person, error) {
	return mockPerson, nil
}

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

func (m *PersonModel) Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, string, error) {
	return 1, "brad-pitt-1.jpg", "brad-pitt-1", nil
}

func (m *PersonModel) Search(query string) ([]*models.Person, error) {
	return []*models.Person{mockPerson}, nil
}
