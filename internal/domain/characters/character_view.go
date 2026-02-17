package characters

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type CharactersListResult struct {
	Characters  []*models.CharacterRef
	TotalCount  int
	CreatorRefs []*models.CreatorRef
	PersonRefs  []*models.PersonRef
	ShowRefs    []*models.ShowRef
	TagRefs     []*models.TagRef
	Metadata    models.Metadata
	Filter      *models.Filter
}

func (s *CharacterService) ListCharacters(f *models.Filter, includeRefs bool) (CharactersListResult, error) {
	result := CharactersListResult{}
	characters, metadata, err := s.Repos.Characters.List(f)
	if err != nil {
		return result, fmt.Errorf("list characters error: %w", err)
	}

	result.Metadata = metadata
	result.Filter = f
	result.Characters = characters
	result.TotalCount = metadata.TotalRecords
	return result, nil
}
