package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) seriesView(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	seriesId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	series, err := app.series.GetById(seriesId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	page, err := views.SeriesPageView(series, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Page = page
	app.render(r, w, http.StatusOK, "view-series.gohtml", "base", data)
}

type seriesFormPage struct {
	Title     string
	SeriesUrl string
	Form      seriesForm
}

func (app *application) seriesAddPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Page = seriesFormPage{
		Title: "Add Series",
		Form: seriesForm{
			Action: "/series/add",
		},
	}

	app.render(r, w, http.StatusOK, "series-form-page.gohtml", "base", data)
}

func (app *application) seriesAdd(w http.ResponseWriter, r *http.Request) {
	var form seriesForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateSeriesForm(&form)
	if !form.Valid() {
		form.Action = "/series/add"
		app.render(r, w, http.StatusUnprocessableEntity, "series-form-page.gohtml", "series-form", form)
		return
	}

	series := app.convertFormtoSeries(&form)

	thumbnailName, err := generateThumbnailName(form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	series.ThumbnailName = &thumbnailName

	slug := models.CreateSlugName(safeDeref(series.Title))
	series.Slug = &slug
	id, err := app.series.Insert(&series)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.saveLargeThumbnail(thumbnailName, "series", form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		app.series.Delete(id)
		return
	}

	w.Header().Add("Hx-Redirect", fmt.Sprintf("/series/%d/update", id))
}

func (app *application) seriesUpdatePage(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	seriesId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	series, err := app.series.GetById(seriesId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	form := app.convertSeriestoForm(series)
	form.Action = fmt.Sprintf("/series/%d/update", seriesId)
	form.ImageUrl = fmt.Sprintf(
		"%s/series/small/%s",
		app.baseImgUrl,
		safeDeref(series.ThumbnailName),
	)

	data := app.newTemplateData(r)
	data.Page = seriesFormPage{
		Title:     "Update Series",
		SeriesUrl: fmt.Sprintf("/series/%d/%s", seriesId, safeDeref(series.Slug)),
		Form:      form,
	}

	app.render(r, w, http.StatusOK, "series-form-page.gohtml", "base", data)
}

func (app *application) seriesUpdate(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	seriesId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	staleSeries, err := app.series.GetById(seriesId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}
	oldThumbnailName := safeDeref(staleSeries.ThumbnailName)

	var form seriesForm
	form.Action = fmt.Sprintf("/series/%d/update", seriesId)

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateSeriesForm(&form)
	if !form.Valid() {
		form.ImageUrl = fmt.Sprintf(
			"%s/series/small/%s",
			app.baseImgUrl,
			safeDeref(staleSeries.ThumbnailName),
		)

		app.render(r, w, http.StatusUnprocessableEntity, "series-form-page.gohtml", "series-form", form)
		return
	}

	thumbnailName := oldThumbnailName
	if form.Thumbnail != nil {
		var err error
		thumbnailName, err = generateThumbnailName(form.Thumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
		err = app.saveLargeThumbnail(thumbnailName, "series", form.Thumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}

	}

	updateSeries := app.convertFormtoSeries(&form)
	updateSeries.ThumbnailName = &thumbnailName

	slug := models.CreateSlugName(safeDeref(updateSeries.Title))
	updateSeries.Slug = &slug

	err = app.series.Update(&updateSeries)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if form.Thumbnail != nil && staleSeries.ThumbnailName != nil {
		err = app.deleteImage("series", *staleSeries.ThumbnailName)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	form.ImageUrl = fmt.Sprintf(
		"%s/series/small/%s",
		app.baseImgUrl,
		safeDeref(updateSeries.ThumbnailName),
	)

	app.render(r, w, http.StatusOK, "series-form-page.gohtml", "series-form", form)
}

func (app *application) deleteSeries(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	seriesId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	_, err = app.series.GetById(seriesId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
