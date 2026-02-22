package quotes

import (
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

	data.TranscriptLines = transcript
	return data, nil
}
