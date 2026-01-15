package services

import (
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
	Filter        *models.Filter
}

func (s *SketchService) ListSketches(f *models.Filter, includeRefs bool) (SketchListResult, error) {
	result := SketchListResult{}
	sketches, err := s.Repos.Sketches.Get(f)
	if err != nil {
		println("GET ERROR")
		return result, fmt.Errorf("list sketches error: %w", err)
	}

	totalCount, err := s.Repos.Sketches.GetCount(f)
	if err != nil {
		println("COUNT ERROR")
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

	result.Filter = f
	result.Sketches = sketches
	result.TotalCount = totalCount
	return result, nil
}
