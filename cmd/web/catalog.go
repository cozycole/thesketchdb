package main

import (
	"fmt"

	"github.com/cozycole/thesketchdb/internal/models"
)

func (app *Application) getSketchCatalogResults(
	currentPage int,
	searchType string,
	filter *models.Filter,
) (*models.SearchResult, error) {
	sketches, metadata, err := app.Sketches.Get(filter)
	if err != nil {
		return nil, fmt.Errorf("%s get error: %w", searchType, err)
	}

	// totalCount, err := app.Sketches.GetCount(filter)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s search count error: %w", searchType, err)
	// }

	return &models.SearchResult{
		Type:             "sketch",
		SketchResults:    sketches,
		TotalSketchCount: metadata.TotalRecords,
	}, nil
}
