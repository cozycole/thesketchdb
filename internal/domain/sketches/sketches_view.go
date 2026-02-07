package sketches

import (
	"errors"
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type SketchListResult struct {
	Sketches      []*models.SketchRef
	TotalCount    int
	CharacterRefs []*models.CharacterRef
	CreatorRefs   []*models.CreatorRef
	PersonRefs    []*models.PersonRef
	ShowRefs      []*models.ShowRef
	TagRefs       []*models.TagRef
	Metadata      models.Metadata
	Filter        *models.Filter
}

func (s *SketchService) ListSketches(f *models.Filter, includeRefs bool) (SketchListResult, error) {
	result := SketchListResult{}
	sketches, metadata, err := s.Repos.Sketches.Get(f)
	if err != nil {
		return result, fmt.Errorf("list sketches error: %w", err)
	}

	totalCount, err := s.Repos.Sketches.GetCount(f)
	if err != nil {
		return result, fmt.Errorf("list sketches total count error: %w", err)
	}

	// Fetched if needed for filter chips
	if includeRefs {
		if len(f.PersonIDs) > 0 {
			result.PersonRefs, _ = s.Repos.People.GetPersonRefs(f.PersonIDs)
		}

		if len(f.CharacterIDs) > 0 {
			result.CharacterRefs, _ = s.Repos.Characters.GetCharactersRefs(f.CharacterIDs)
		}

		if len(f.CreatorIDs) > 0 {
			result.CreatorRefs, _ = s.Repos.Creators.GetCreatorRefs(f.CreatorIDs)
		}

		if len(f.ShowIDs) > 0 {
			result.ShowRefs, _ = s.Repos.Shows.GetShowRefs(f.ShowIDs)
		}

		if len(f.TagIDs) > 0 {
			result.TagRefs, _ = s.Repos.Tags.GetTagRefs(f.TagIDs)
		}
	}

	result.Metadata = metadata
	result.Filter = f
	result.Sketches = sketches
	result.TotalCount = totalCount
	return result, nil
}

func (s *SketchService) GetSketch(id int) (*models.Sketch, error) {
	sketch, err := s.Repos.Sketches.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return nil, models.ErrNoSketch
		} else {
			return nil, err
		}

	}

	return sketch, nil
}
