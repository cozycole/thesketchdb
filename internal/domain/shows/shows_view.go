package shows

import (
	"errors"

	"sketchdb.cozycole.net/internal/domain/shared"
	"sketchdb.cozycole.net/internal/models"
)

// you'll need to refactor the models into their respective
// domain to begin using service level filter otherwise it's
// a circular import
type ShowFilter struct {
	shared.ListOptions
	PersonIDs []int
	TagIDs    []int
}

type EpisodeListResult struct {
	Episodes      []*models.EpisodeRef
	CharacterRefs []*models.CharacterRef
	PersonRefs    []*models.PersonRef
	ShowRefs      []*models.ShowRef
	TagRefs       []*models.TagRef
	Metadata      *models.Metadata
	Filter        *ShowFilter
}

func (s *ShowService) ListEpisodes(f *models.Filter, includeRefs bool) (EpisodeListResult, error) {
	result := EpisodeListResult{}
	episodes, metadata, err := s.Repos.Shows.ListEpisodes(f)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return result, err
	}

	result.Metadata = &metadata
	result.Episodes = episodes
	return result, nil
}
