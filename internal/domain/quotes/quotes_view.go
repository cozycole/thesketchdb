package quotes

import (
	"errors"
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type AdminQuoteData struct {
	Quotes          []*models.Quote
	TranscriptLines []*models.TranscriptLine
}

func (s *QuoteService) GetAdminQuotes(sketchId int) (AdminQuoteData, error) {
	data := AdminQuoteData{}
	transcript, err := s.Repos.Quotes.GetTranscriptBySketch(sketchId)
	if err != nil {
		return data, fmt.Errorf("get transcript error: %w", err)
	}

	quotes, err := s.Repos.Quotes.GetBySketch(sketchId)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return data, fmt.Errorf("get quotes error: %w", err)
	}

	data.Quotes = quotes
	data.TranscriptLines = transcript
	return data, nil
}
