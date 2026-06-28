package views

import (
	"sketchdb.cozycole.net/internal/models"
)

type Grouping struct {
	ID       int
	Slug     string
	Title    string
	Sketches *SketchGallery
}

func GroupingView(g *models.Grouping, baseImgUrl string) Grouping {
	view := Grouping{}
	if g == nil || g.ID == nil {
		return view
	}

	view.ID = safeDeref(g.ID)
	view.Slug = safeDeref(g.Slug)
	view.Title = safeDeref(g.Title)
	view.Sketches, _ = SketchGalleryView(g.Sketches, baseImgUrl, "default", "sub")

	return view
}
