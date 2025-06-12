package views

import (
	"fmt"
	"net/url"

	"sketchdb.cozycole.net/internal/models"
)

// Galleries become horizontally scrollable on small
// screens, but are grids on larger screens.
// Carousels are always horizontally scrollable
// Grids are never horizontally scrollable.

type SearchPage struct {
	Query                string
	EscapedQuery         string
	NoResults            bool
	PersonResults        *PersonGallery
	PersonResultCount    int
	CharacterResults     *CharacterGallery
	CharacterResultCount int
	CreatorResults       *CreatorGallery
	CreatorResultCount   int
	ShowResults          *ShowGallery
	ShowResultCount      int
	SketchResults        *SketchGallery
	SketchResultCount    int
}

func SearchPageView(results *models.SearchResult, query, baseImgUrl string, maxResults int) (*SearchPage, error) {
	page := SearchPage{}
	var err error

	page.PersonResultCount = results.TotalPersonCount
	if page.PersonResults, err = PersonGalleryView(results.PersonResults, baseImgUrl); err != nil {
		return nil, err
	}

	page.CreatorResultCount = results.TotalCreatorCount
	if page.CreatorResults, err = CreatorGalleryView(results.CreatorResults, baseImgUrl); err != nil {
		return nil, err
	}

	page.CharacterResultCount = results.TotalCharacterCount
	if page.CharacterResults, err = CharacterGalleryView(results.CharacterResults, baseImgUrl); err != nil {
		return nil, err
	}

	page.ShowResultCount = results.TotalShowCount
	if page.ShowResults, err = ShowGalleryView(results.ShowResults, baseImgUrl); err != nil {
		return nil, err
	}

	page.SketchResultCount = results.TotalVideoCount
	page.SketchResults, err = SketchGalleryView(results.VideoResults, baseImgUrl, "Default", "Full", maxResults)
	if err != nil {
		return nil, err
	}
	page.EscapedQuery = url.QueryEscape(query)
	page.SketchResults.SeeMoreUrl = fmt.Sprintf("/catalog/sketches?query=%s", page.EscapedQuery)
	page.SketchResults.SeeMore = page.SketchResultCount > maxResults

	page.NoResults = results.TotalVideoCount == 0 &&
		results.TotalPersonCount == 0 &&
		results.TotalCreatorCount == 0 &&
		results.TotalCharacterCount == 0 &&
		results.TotalShowCount == 0

	page.Query = query

	return &page, nil
}
