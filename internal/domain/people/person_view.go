package people

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type PeopleListResult struct {
	People        []*models.PersonRef
	TotalCount    int
	CharacterRefs []*models.CharacterRef
	CreatorRefs   []*models.CreatorRef
	ShowRefs      []*models.ShowRef
	TagRefs       []*models.TagRef
	Metadata      models.Metadata
	Filter        *models.Filter
}

func (s *PersonService) ListPeople(f *models.Filter, includeRefs bool) (PeopleListResult, error) {
	result := PeopleListResult{}
	people, metadata, err := s.Repos.People.List(f)
	if err != nil {
		return result, fmt.Errorf("list people error: %w", err)
	}

	result.Metadata = metadata
	result.Filter = f
	result.People = people
	result.TotalCount = metadata.TotalRecords
	return result, nil
}
