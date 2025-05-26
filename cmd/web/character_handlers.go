package main

import (
	"errors"
	"net/http"
	"strconv"

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

	popularSketches, err := app.videos.Get(
		&models.Filter{
			Limit:  16,
			Offset: 0,
			SortBy: "az",
			Characters: []*models.Character{
				&models.Character{ID: character.ID},
			},
		},
	)

	// stats, err := app.people.GetPersonStats(persondId)

	data := app.newTemplateData(r)
	data.CharacterPage.Character = character
	data.CharacterPage.Popular = popularSketches
	// data.PersonPage.Stats = stats

	app.render(r, w, http.StatusOK, "view-character.tmpl.html", "base", data)
}
