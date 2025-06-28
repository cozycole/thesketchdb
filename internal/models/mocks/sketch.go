package mocks

import (
	"slices"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

var mockSketch = &models.Sketch{
	ID:            1,
	Title:         "Test Title",
	URL:           "www.testvid.com",
	ThumbnailName: "a-test-thumbnail-1.jpg",
	Rating:        "g",
	Creator: &models.Creator{
		ID:              1,
		Name:            "Test Creator",
		URL:             "www.test-creator-page.com",
		ProfileImage:    "test-creator-1.jpg",
		EstablishedDate: time.Now(),
	},
	Cast: []*models.CastMember{
		{
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

type SketchModel struct {
	// sketchCreatorRel map[int][]int
	sketchPersonRel map[int][]int
}

func (m *SketchModel) Get(id int) (*models.Sketch, error) {
	return mockSketch, nil
}

func (m *SketchModel) GetAll(offset int) ([]*models.Sketch, error) {
	return []*models.Sketch{mockSketch}, nil
}

func (m *SketchModel) GetByCreator(id int) ([]*models.Sketch, error) {
	return []*models.Sketch{mockSketch}, nil
}

func (m *SketchModel) GetBySlug(slug string) (*models.Sketch, error) {
	return mockSketch, nil
}

func (m *SketchModel) Insert(sketch *models.Sketch) error {
	return nil
}

func (m *SketchModel) InsertSketchCreatorRelation(sketchId, creatorId int) error {
	return nil
}

func (m *SketchModel) InsertSketchPersonRelation(sketchId, personId, position int, characterId *int, imgName string) error {
	if m.sketchPersonRel == nil {
		m.sketchPersonRel = map[int][]int{}
	}

	val, ok := m.sketchPersonRel[sketchId]
	if ok {
		if slices.Contains(val, personId) {
			return models.ErrDuplicateVidPersonRel
		} else {
			m.sketchPersonRel[sketchId] = append(val, personId)
		}
	} else {
		m.sketchPersonRel[sketchId] = []int{personId}
	}
	return nil
}

func (m *SketchModel) Search(search string, offset int) ([]*models.Sketch, error) {
	t := time.Now()
	mockSketch.UploadDate = &t
	return []*models.Sketch{mockSketch}, nil
}
