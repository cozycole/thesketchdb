package mocks

import (
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

var mockCharacter = &models.Character{
	ID:   utils.GetIntPtr(1),
	Slug: utils.GetStringPtr("david-s-pumpkins-1"),
	Name: utils.GetStringPtr("David S. Pumpkins"),
}

type CharacterModel struct{}

func (m *CharacterModel) GetBySlug(slug string) (*models.Person, error) {
	return mockPerson, nil
}

func (m *CharacterModel) Get(id int) (*models.Person, error) {
	return mockPerson, nil
}

func (m *CharacterModel) Exists(id int) (bool, error) {
	switch id {
	case 1, 2, 3:
		return true, nil
	default:
		return false, nil
	}
}

func (m *CharacterModel) Insert(first, last, imgName, imgExt string, birthDate time.Time) (int, string, string, error) {
	return 1, "brad-pitt-1.jpg", "brad-pitt-1", nil
}

func (m *CharacterModel) Search(query string) ([]*models.Character, error) {
	return []*models.Character{mockCharacter}, nil
}
