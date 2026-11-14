package main

import (
	"errors"
	"fmt"

	"github.com/cozycole/thesketchdb/internal/models"
)

// NOTE: Query is defined on the Filter, SearchResult and templateData structs
// Given the search term: kenan snl
// - Filter.Query -> "kenan | snl"
// - SearchResult.Query -> "kenan+snl"
// - templateData.Query -> "kenan snl" (i.e. user facing)

func (app *Application) getSearchResults(filter *models.Filter) (*models.SearchResult, error) {

	sketches, _, err := app.Sketches.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("search sketch error: %s", err)
	}

	sketchCount, err := app.Sketches.GetCount(filter)
	if err != nil {
		return nil, fmt.Errorf("search sketch count error: %s", err)
	}

	people, err := app.People.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("search person error: %s", err)
	}

	peopleCount, err := app.People.GetCount(filter)
	if err != nil {
		return nil, fmt.Errorf("search person count error: %s", err)
	}

	creators, creatorMetadata, err := app.Creators.List(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("search creator error: %s", err)
	}

	characters, err := app.Characters.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("search character error: %s", err)
	}

	characterCount, err := app.Characters.GetCount(filter)
	if err != nil {
		return nil, fmt.Errorf("search character count error: %s", err)
	}

	shows, err := app.Shows.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("search show error: %s", err)
	}

	showCount, err := app.Shows.GetCount(filter)
	if err != nil {
		return nil, fmt.Errorf("search show count error: %s", err)
	}

	return &models.SearchResult{
		SketchResults:       sketches,
		TotalSketchCount:    sketchCount,
		PersonResults:       people,
		TotalPersonCount:    peopleCount,
		CreatorResults:      creators,
		TotalCreatorCount:   creatorMetadata.TotalRecords,
		CharacterResults:    characters,
		TotalCharacterCount: characterCount,
		ShowResults:         shows,
		TotalShowCount:      showCount,
	}, nil
}
