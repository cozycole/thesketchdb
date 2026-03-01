package tags

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type TagsListResult struct {
	Tags          []*models.Tag
	TotalCount    int
	Type          string
	CharacterRefs []*models.CharacterRef
	CreatorRefs   []*models.CreatorRef
	ShowRefs      []*models.ShowRef
	Metadata      models.Metadata
	Filter        *models.Filter
}

func (s *TagsService) ListTags(f *models.Filter, includeRefs bool) (TagsListResult, error) {
	result := TagsListResult{}
	tags, metadata, err := s.Repos.Tags.List(f)
	if err != nil {
		return result, fmt.Errorf("list tags error: %w", err)
	}

	result.Filter = f
	result.Metadata = metadata
	result.Tags = tags
	return result, nil
}
