package main

import (
	"errors"

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
	Filter              *models.Filter
	NoResults           bool
	PageURLParams       string
	CurrentPage         int
	Pages               []int
}

func (app *application) getSearchResults(filter *models.Filter) (*SearchResult, error) {

	videos, err := app.videos.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, err
	}

	videoCount, err := app.videos.GetCount(filter)
	if err != nil {
		return nil, err
	}

	people, err := app.people.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, err
	}

	peopleCount, err := app.people.GetCount(filter)
	if err != nil {
		return nil, err
	}

	creators, err := app.creators.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, err
	}

	creatorCount, err := app.creators.GetCount(filter)
	if err != nil {
		return nil, err
	}

	characters, err := app.characters.Get(filter)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return nil, err
	}

	characterCount, err := app.characters.GetCount(filter)
	if err != nil {
		return nil, err
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
