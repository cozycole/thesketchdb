package quotes

import (
	"sketchdb.cozycole.net/internal/models"
)

func (s *QuoteService) UpdateQuotes(sketchId int, quotes []*models.Quote, deleted []int) ([]*models.Quote, error) {
	err := s.Repos.Quotes.BatchUpdateQuotes(sketchId, quotes, deleted)
	if err != nil {
		return nil, err
	}

	for _, q := range quotes {
		if q.ID == nil {
			continue
		}

		castIds := []int{}
		tagIds := []int{}
		for _, t := range q.Tags {
			if t.ID != nil {
				tagIds = append(tagIds, *t.ID)
			}
		}

		for _, cm := range q.CastMembers {
			if cm.ID != nil {
				castIds = append(castIds, *cm.ID)
			}
		}

		err = s.Repos.Quotes.BatchUpdateQuoteCastMembers(*q.ID, castIds)
		if err != nil {
			return nil, err
		}

		err = s.Repos.Quotes.BatchUpdateQuoteTags(*q.ID, tagIds)
		if err != nil {
			return nil, err
		}
	}

	updatedQuotes, err := s.Repos.Quotes.GetBySketch(sketchId)
	if err != nil {
		return nil, err
	}

	return updatedQuotes, nil
}
