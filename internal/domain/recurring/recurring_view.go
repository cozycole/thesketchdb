package recurring

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type RecurringListResult struct {
	Recurring     []*models.RecurringRef
	TotalCount    int
	CharacterRefs []*models.CharacterRef
	CreatorRefs   []*models.CreatorRef
	PersonRefs    []*models.PersonRef
	ShowRefs      []*models.ShowRef
	TagRefs       []*models.TagRef
	Metadata      models.Metadata
	Filter        *models.Filter
}

func (s *RecurringService) ListRecurring(f *models.Filter, includeRefs bool) (RecurringListResult, error) {
	result := RecurringListResult{}
	recurring, metadata, err := s.Repos.Recurring.List(f)
	if err != nil {
		return result, fmt.Errorf("list recurring error: %w", err)
	}

	result.Metadata = metadata
	result.Filter = f
	result.Recurring = recurring
	result.TotalCount = metadata.TotalRecords
	return result, nil
}
