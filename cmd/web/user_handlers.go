package main

import (
	"errors"
	"net/http"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *Application) userView(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	user, err := app.Users.GetByUsername(username)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	favoriteSketches, err := app.Sketches.GetByUserLikes(*user.ID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	page, err := views.UserPageView(
		user,
		favoriteSketches,
		app.Config.BaseImgUrl,
	)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = page
	app.render(r, w, http.StatusOK, "view-user.gohtml", "base", data)
}
