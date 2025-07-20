package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type SeriesPage struct {
	SeriesTitle     string
	Description     string
	Image           string
	UpdateSeriesUrl string
	PartCount       int
	Sketches        *SketchGallery
}

func SeriesPageView(series *models.Series, baseImgUrl string) (*SeriesPage, error) {
	page := SeriesPage{}
	page.SeriesTitle = safeDeref(series.Title)
	page.Description = safeDeref(series.Description)
	page.Image = fmt.Sprintf(
		"%s/series/%s",
		baseImgUrl,
		safeDeref(series.ThumbnailName),
	)

	page.UpdateSeriesUrl = fmt.Sprintf(
		"/series/%d/update",
		safeDeref(series.ID),
	)

	page.PartCount = len(series.Sketches)
	sketches, err := SketchGalleryView(
		series.Sketches,
		baseImgUrl,
		"default",
		"full",
		-1,
	)

	if err != nil {
		return nil, err
	}

	page.Sketches = sketches
	return &page, nil
}
