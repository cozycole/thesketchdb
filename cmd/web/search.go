package main

import (
	"errors"
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

// NOTE: Query is defined on the Filter, SearchResult and templateData structs
// Given the search term: kenan snl
// - Filter.Query -> "kenan | snl"
// - SearchResult.Query -> "kenan+snl"
// - templateData.Query -> "kenan snl" (i.e. user facing)
type SearchResult struct {
	Type                string
	Query               string
	VideoResults        []*models.Video
	TotalVideoCount     int
	PersonResults       []*models.Person
	TotalPersonCount    int
	CreatorResults      []*models.Creator
	TotalCreatorCount   int
	CharacterResults    []*models.Character
	TotalCharacterCount int
	ShowResults         []*models.Show
	TotalShowCount      int
	Filter              *models.Filter
	NoResults           bool
	PageURLParams       string
	CurrentPage         int
	Pages               []int
}

func (app *application) getSearchResults(filter *models.Filter) (*SearchResult, error) {

	videos, err := app.videos.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("search video error: %s", err)
	}

	videoCount, err := app.videos.GetCount(filter)
	if err != nil {
		return nil, fmt.Errorf("search video count error: %s", err)
	}

	people, err := app.people.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("search person error: %s", err)
	}

	peopleCount, err := app.people.GetCount(filter)
	if err != nil {
		return nil, fmt.Errorf("search person count error: %s", err)
	}

	creators, err := app.creators.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("search creator error: %s", err)
	}

	creatorCount, err := app.creators.GetCount(filter)
	if err != nil {
		return nil, fmt.Errorf("search creator count error: %s", err)
	}

	characters, err := app.characters.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("search character error: %s", err)
	}

	characterCount, err := app.characters.GetCount(filter)
	if err != nil {
		return nil, fmt.Errorf("search character count error: %s", err)
	}

	shows, err := app.shows.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, fmt.Errorf("search show error: %s", err)
	}

	showCount, err := app.shows.GetCount(filter)
	if err != nil {
		return nil, fmt.Errorf("search show count error: %s", err)
	}

	noResults := peopleCount == 0 &&
		videoCount == 0 &&
		creatorCount == 0 &&
		characterCount == 0

	return &SearchResult{
		VideoResults:        videos,
		TotalVideoCount:     videoCount,
		PersonResults:       people,
		TotalPersonCount:    peopleCount,
		CreatorResults:      creators,
		TotalCreatorCount:   creatorCount,
		CharacterResults:    characters,
		TotalCharacterCount: characterCount,
		ShowResults:         shows,
		TotalShowCount:      showCount,
		NoResults:           noResults,
	}, nil
}

func paginate(currentPage, totalPages int) []int {
	var pages []int

	// Show the current page and two pages before and after
	start := currentPage - 1
	if start < 1 {
		start = 1
	}

	end := currentPage + 1
	if end > totalPages {
		end = totalPages
	}

	// Add the main range
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}

	// Add ellipsis and the last page if necessary
	if end < totalPages {
		if end+1 < totalPages {
			pages = append(pages, -1) // -1 represents "..."
		}
		pages = append(pages, totalPages)
	}

	// Add the first page and ellipsis if necessary
	if start > 1 {
		if start > 2 {
			pages = append([]int{-1}, pages...)
		}
		pages = append([]int{1}, pages...)
	}

	return pages
}
