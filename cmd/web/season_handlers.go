package main

import (
	"fmt"
	"errors"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) viewSeason(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	seasonId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		return
	}

	season, err := app.shows.GetSeason(seasonId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	show, err := app.shows.GetById(safeDeref(season.Show.ID))
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)

	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	if isHxRequest && !isHistoryRestore {
		format := r.URL.Query().Get("format")
		templateData := views.EpisodeGalleryView(season.Episodes, app.baseImgUrl, format, true)

		if isSeasonPath(r) {
			w.Header().Add("HX-Push-Url", fmt.Sprintf("/season/%d/%s", *season.ID, *season.Slug))
		}

		app.render(r, w, http.StatusOK, "episode-gallery.gohtml", "episode-gallery", templateData)
		return
	}

	page := views.SeasonPageView(show, season, app.baseImgUrl)
	data.Page = page
	app.render(r, w, http.StatusOK, "view-season.gohtml", "base", data)
}

func (app *application) addSeason(w http.ResponseWriter, r *http.Request) {

	var form seasonForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	app.validateSeasonForm(&form)
	if !form.Valid() {
		app.render(r, w, http.StatusUnprocessableEntity, "show-form-page.gohtml", "season-form", form)
		return
	}

	showIdParam := r.PathValue("id")
	showId, err := strconv.Atoi(showIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	show, err := app.shows.GetById(showId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	slug := models.CreateSlugName(safeDeref(show.Name) + fmt.Sprintf(" s%d", form.Number))
	newSeason := &models.Season{
		Number: &form.Number,
		Slug:   &slug,
		Show: &models.ShowRef{
			ID: show.ID,
		},
	}

	_, err = app.shows.AddSeason(newSeason)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	show, err = app.shows.GetById(showId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := views.SeasonDropdownsView(show, app.baseImgUrl)
	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "season-dropdowns", data)
}

func (app *application) deleteSeason(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	seasonIdParam := r.PathValue("id")
	seasonId, err := strconv.Atoi(seasonIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	err = app.shows.DeleteSeason(seasonId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

