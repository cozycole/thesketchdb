package views

import (
	"fmt"

	"sketchdb.cozycole.net/internal/models"
)

type RecurringPage struct {
	RecurringTitle     string
	Description        string
	Image              string
	UpdateRecurringUrl string
	SketchCount        int
	Sketches           *SketchGallery
}

func RecurringPageView(recurring *models.Recurring, baseImgUrl string) (*RecurringPage, error) {
	page := RecurringPage{}
	page.RecurringTitle = safeDeref(recurring.Title)
	page.Description = safeDeref(recurring.Description)
	page.Image = fmt.Sprintf(
		"%s/recurring/%s",
		baseImgUrl,
		safeDeref(recurring.ThumbnailName),
	)

	page.UpdateRecurringUrl = fmt.Sprintf(
		"/recurring/%d/update",
		safeDeref(recurring.ID),
	)

	page.SketchCount = len(recurring.Sketches)
	sketches, err := SketchGalleryView(
		recurring.Sketches,
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
