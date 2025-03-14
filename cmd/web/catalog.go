package main

import (
	"fmt"
	"math"

	"sketchdb.cozycole.net/internal/models"
)

func (app *application) getCatalogResults(currentPage int, searchType string, filter *models.Filter) (*SearchResult, error) {
	var result SearchResult
	var getErr, countErr error
	var totalCount, resultCount int
	app.infoLog.Printf("%+v\n", filter)
	switch searchType {
	default:
		var videos []*models.Video
		result.Type = "video"
		videos, getErr = app.videos.Get(filter)
		totalCount, countErr = app.videos.GetCount(filter)
		resultCount = len(videos)

		result.VideoResults = videos
	}

	if getErr != nil {
		return nil, fmt.Errorf("%s get error: %w", searchType, getErr)
	}

	if countErr != nil {
		return nil, fmt.Errorf("%s search count error: %w", searchType, countErr)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(app.settings.pageSize)))
	pageList := paginate(currentPage, totalPages)

	result.Pages = pageList
	result.NoResults = resultCount == 0
	result.CurrentPage = currentPage
	return &result, nil
}
