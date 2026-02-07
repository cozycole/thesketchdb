package series

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type SeriesListResult struct {
	Series        []*models.SeriesRef
	TotalCount    int
	CharacterRefs []*models.CharacterRef
	CreatorRefs   []*models.CreatorRef
	PersonRefs    []*models.PersonRef
	ShowRefs      []*models.ShowRef
	TagRefs       []*models.TagRef
	Metadata      models.Metadata
	Filter        *models.Filter
}

func (s *SeriesService) ListSeries(f *models.Filter, includeRefs bool) (SeriesListResult, error) {
	result := SeriesListResult{}
	series, metadata, err := s.Repos.Series.List(f)
	if err != nil {
		return result, fmt.Errorf("list sketches error: %w", err)
	}

	result.Metadata = metadata
	result.Filter = f
	result.Series = series
	result.TotalCount = metadata.TotalRecords
	return result, nil
}
