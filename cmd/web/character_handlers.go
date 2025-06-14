package main

import (
	"errors"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) characterView(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	characterdId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	character, err := app.characters.GetById(characterdId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	popularSketches, err := app.sketches.Get(
		&models.Filter{
			Limit:  16,
			Offset: 0,
			SortBy: "az",
			Characters: []*models.Character{
				&models.Character{ID: character.ID},
			},
		},
	)

	data := app.newTemplateData(r)
	page, err := views.CharacterPageView(
		character,
		popularSketches,
		app.baseImgUrl,
	)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = page
	app.render(r, w, http.StatusOK, "view-character.gohtml", "base", data)
}
