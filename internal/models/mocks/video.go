package mocks

import (
	"time"

	"sketchdb.cozycole.net/internal/models"
)

var mockVideo = &models.Video{
	ID:         1,
	Title:      "Test Title",
	URL:        "www.testvid.com",
	Thumbnail:  "a-test-thumbnail-1.jpg",
	Rating:     "g",
	UploadDate: time.Now(),
	Creator: &models.Creator{
		ID:              1,
		Name:            "Test Creator",
		URL:             "www.test-creator-page.com",
		ProfileImage:    "test-creator-1.jpg",
		EstablishedDate: time.Now(),
	},
}

type VideoModel struct{}

func (m *VideoModel) Insert(title, video_url, rating, slug, imgExt string, upload_date time.Time) (int, string, string, error) {
	return 1, "test-img-1", "test-img-1.jpg", nil
}

func (m *VideoModel) Search(search string, offset int) ([]*models.Video, error) {
	return []*models.Video{mockVideo}, nil
}

func (m *VideoModel) GetAll(offset int) ([]*models.Video, error) {
	return []*models.Video{mockVideo}, nil
}

func (m *VideoModel) Get(id int) (*models.Video, error) {
	return mockVideo, nil
}

func (m *VideoModel) GetBySlug(slug string) (*models.Video, error) {
	return mockVideo, nil
}

func (m *VideoModel) InsertVideoCreatorRelation(vidId, creatorId int) error {
	return nil
}

func (m *VideoModel) InsertVideoActorRelation(vidId, actorId int) error {
	return nil
}
