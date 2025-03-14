package main

import (
	"fmt"
	"math"

	"sketchdb.cozycole.net/internal/models"
)

type SearchResult struct {
	Type           string
	Query          string
	ProfileResults []*models.ProfileResult
	VideoResults   []*models.Video
	Filter         *models.Filter
	NoResults      bool
	PageURLParams  string
	CurrentPage    int
	Pages          []int
}

func (app *application) getSearchResults(query string, currentPage int, searchType string) (*SearchResult, error) {
	limit := app.settings.pageSize
	offset := (currentPage - 1) * limit

	var videos []*models.Video
	var profiles []*models.ProfileResult
	var searchErr, countErr error
	var resultCount, totalCount int
	switch searchType {
	case "character":
		profiles, searchErr = app.characters.VectorSearch(query, limit, offset)
		totalCount, countErr = app.characters.SearchCount(query)
		resultCount = len(profiles)
	case "creator":
		profiles, searchErr = app.creators.VectorSearch(query)
		totalCount, countErr = app.creators.SearchCount(query)
		resultCount = len(profiles)
	case "person":
		app.infoLog.Println("searching person table...")
		profiles, searchErr = app.people.VectorSearch(query)
		totalCount, countErr = app.people.SearchCount(query)
		resultCount = len(profiles)
	// case "user":
	// 	app.infoLog.Println("searching user table...")
	default:
		// TODO: have the default be all and mix them together somehow
		app.infoLog.Println("searching video table...")
		videos, searchErr = app.videos.Search(query, limit, offset)
		totalCount, countErr = app.videos.SearchCount(query)
		resultCount = len(videos)
		searchType = "video"
	}

	if searchErr != nil {
		return nil, fmt.Errorf("%s search error: %w", searchType, searchErr)
	}

	if countErr != nil {
		return nil, fmt.Errorf("%s search count error: %w", searchType, countErr)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(app.settings.pageSize)))
	pageList := paginate(currentPage, totalPages)
	app.infoLog.Println(pageList)
	if app.debugMode {
		app.infoLog.Printf("Query: %s Page: %d Offset: %d Limit: %d | Total Count: %d\n", query, currentPage, offset, limit, totalCount)
	}

	return &SearchResult{
		Type:           searchType,
		Query:          query,
		ProfileResults: profiles,
		VideoResults:   videos,
		NoResults:      resultCount == 0,
		CurrentPage:    currentPage,
		Pages:          pageList,
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
