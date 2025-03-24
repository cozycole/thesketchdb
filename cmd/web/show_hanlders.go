package main

import (
	"errors"
	"net/http"
	// "path"
	// "time"
	"sketchdb.cozycole.net/internal/models"
	// "sketchdb.cozycole.net/internal/utils"
)

func (app *application) viewShow(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	show, err := app.shows.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	filter := &models.Filter{
		Limit:  8,
		SortBy: "latest",
		Shows:  []*models.Show{show},
	}

	videos, err := app.videos.Get(filter)
	if err != nil {
		app.serverError(w, err)
		return
	}

	cast, err := app.shows.GetShowCast(*show.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Show = show
	data.People = cast
	data.Videos = videos

	app.render(w, http.StatusOK, "view-show.tmpl.html", "base", data)
}
