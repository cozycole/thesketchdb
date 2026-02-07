package creators

import (
	"errors"

	"sketchdb.cozycole.net/internal/domain/shared"
	"sketchdb.cozycole.net/internal/models"
)

// you'll need to refactor the models into their respective
// domain to begin using service level filter otherwise it's
// a circular import
type CreatorFilter struct {
	shared.ListOptions
	PersonIDs []int
	TagIDs    []int
}

type CreatorListResult struct {
	Creators      []*models.CreatorRef
	CharacterRefs []*models.CharacterRef
	CreatorRefs   []*models.CreatorRef
	PersonRefs    []*models.PersonRef
	ShowRefs      []*models.ShowRef
	TagRefs       []*models.TagRef
	Metadata      *models.Metadata
	Filter        *CreatorFilter
}

func (s *CreatorService) ListCreators(f *models.Filter, includeRefs bool) (CreatorListResult, error) {
	result := CreatorListResult{}
	creators, metadata, err := s.Repos.Creators.List(f)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return result, err
	}

	result.Metadata = &metadata
	result.Creators = creators
	return result, nil
}
