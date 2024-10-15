package mocks

import (
	"time"

	"sketchdb.cozycole.net/internal/models"
)

var mockCreator = &models.Creator{
	ID:              1,
	Name:            "Test Creator",
	URL:             "www.test-creator-page.com",
	Slug:            "test-creator-1",
	ProfileImage:    "test-creator-1.jpg",
	EstablishedDate: time.Now(),
}

type CreatorModel struct{}

func (m *CreatorModel) Insert(name, url, imgName, imgExt string, establishedDate time.Time) (int, string, string, error) {
	return 1, mockCreator.Slug, mockCreator.ProfileImage, nil
}

func (m *CreatorModel) Get(id int) (*models.Creator, error) {
	return mockCreator, nil
}

func (m *CreatorModel) GetBySlug(slug string) (*models.Creator, error) {
	return mockCreator, nil
}

func (m *CreatorModel) Exists(id int) (bool, error) {
	return true, nil
}
