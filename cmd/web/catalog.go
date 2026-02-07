package main

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

func (app *application) getSketchCatalogResults(
	currentPage int,
	searchType string,
	filter *models.Filter,
) (*models.SearchResult, error) {
	sketches, metadata, err := app.sketches.Get(filter)
	if err != nil {
		return nil, fmt.Errorf("%s get error: %w", searchType, err)
	}

	// totalCount, err := app.sketches.GetCount(filter)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s search count error: %w", searchType, err)
	// }

	return &models.SearchResult{
		Type:             "sketch",
		SketchResults:    sketches,
		TotalSketchCount: metadata.TotalRecords,
	}, nil
}
