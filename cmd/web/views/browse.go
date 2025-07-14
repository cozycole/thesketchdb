package views

import (
	"sketchdb.cozycole.net/internal/models"
)

type BrowseSectionDefinition struct {
	Title    string
	Filter   models.Filter
	Sketches []*models.Sketch
}

var FEATURED_ID = 1
var BrowseSectionDefinitions = []BrowseSectionDefinition{
	{
		Title: "Featured Sketches",
		Filter: models.Filter{
			Limit:  10,
			Offset: 0,
			SortBy: "popular",
			Tags:   []*models.Tag{{ID: &FEATURED_ID}},
		},
	},
	{
		Title: "Popular",
		Filter: models.Filter{
			Limit:  10,
			Offset: 0,
			SortBy: "popular",
		},
	},
}

type BrowsePage struct {
	Sections []BrowseSection
}

type BrowseSection struct {
	Title         string
	SketchGallery *SketchGallery
	SeeAllUrl     string
}

func BrowsePageView(sections []BrowseSectionDefinition, baseImgUrl string) (BrowsePage, error) {
	page := BrowsePage{}

	var browseSections []BrowseSection
	for _, section := range sections {
		gallery, err := SketchGalleryView(
			section.Sketches,
			baseImgUrl,
			"base",
			"full",
			section.Filter.Limit,
		)
		if err != nil {
			return page, err
		}

		seeAllUrl, err := BuildURL(
			"/catalog/sketches",
			1,
			&section.Filter,
		)

		browseSections = append(browseSections, BrowseSection{section.Title, gallery, seeAllUrl})
	}

	page.Sections = browseSections
	return page, nil
}
