package mocks

import (
	"slices"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

var mockVideo = &models.Video{
	ID:        1,
	Title:     "Test Title",
	URL:       "www.testvid.com",
	Thumbnail: "a-test-thumbnail-1.jpg",
	Rating:    "g",
	Creator: &models.Creator{
		ID:              1,
		Name:            "Test Creator",
		URL:             "www.test-creator-page.com",
		ProfileImage:    "test-creator-1.jpg",
		EstablishedDate: time.Now(),
	},
	Cast: []*models.CastMember{
		&models.CastMember{
			Position: utils.GetIntPtr(1),
			Actor: &models.Person{
				ID:    utils.GetIntPtr(1),
				Slug:  utils.GetStringPtr("first-last-1"),
				First: utils.GetStringPtr("First"),
				Last:  utils.GetStringPtr("Last"),
			},
		},
	},
}

type VideoModel struct {
	videoCreatorRel map[int][]int
	videoPersonRel  map[int][]int
}

func (m *VideoModel) Insert(title, video_url, rating, slug, imgExt string, upload_date time.Time) (int, string, string, error) {
	return 1, "test-img-1", "test-img-1.jpg", nil
}

func (m *VideoModel) Search(search string, offset int) ([]*models.Video, error) {
	t := time.Now()
	mockVideo.UploadDate = &t
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

func (m *VideoModel) GetByCreator(id int) ([]*models.Video, error) {
	return []*models.Video{mockVideo}, nil
}

func (m *VideoModel) InsertVideoCreatorRelation(vidId, creatorId int) error {
	return nil
}

func (m *VideoModel) InsertVideoPersonRelation(vidId, personId, position int) error {
	if m.videoPersonRel == nil {
		m.videoPersonRel = map[int][]int{}
	}

	val, ok := m.videoPersonRel[vidId]
	if ok {
		if slices.Contains(val, personId) {
			return models.ErrDuplicateVidPersonRel
		} else {
			m.videoPersonRel[vidId] = append(val, personId)
		}
	} else {
		m.videoPersonRel[vidId] = []int{personId}
	}
	return nil
}
