package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *Application) recurringView(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	recurringId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	recurring, err := app.Recurring.GetById(recurringId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	page, err := views.RecurringPageView(recurring, app.Config.BaseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Page = page
	app.render(r, w, http.StatusOK, "view-recurring.gohtml", "base", data)
}

type recurringFormPage struct {
	Title        string
	RecurringUrl string
	Form         recurringForm
}

func (app *Application) recurringAddPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Page = recurringFormPage{
		Title: "Add recurring",
		Form: recurringForm{
			Action: "/recurring/add",
		},
	}

	app.render(r, w, http.StatusOK, "recurring-form-page.gohtml", "base", data)
}

func (app *Application) recurringAdd(w http.ResponseWriter, r *http.Request) {
	var form recurringForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.ErrorLog.Print(err)
		return
	}

	app.validateRecurringForm(&form)
	if !form.Valid() {
		form.Action = "/recurring/add"
		app.render(r, w, http.StatusUnprocessableEntity, "recurring-form-page.gohtml", "recurring-form", form)
		return
	}

	recurring := app.convertFormtoRecurring(&form)

	thumbnailName, err := generateThumbnailName(form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	recurring.ThumbnailName = &thumbnailName

	slug := models.CreateSlugName(safeDeref(recurring.Title))
	recurring.Slug = &slug
	id, err := app.Recurring.Insert(&recurring)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.saveLargeThumbnail(thumbnailName, "recurring", form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		app.Recurring.Delete(id)
		return
	}

	w.Header().Add("Hx-Redirect", fmt.Sprintf("/recurring/%d/update", id))
}

func (app *Application) recurringUpdatePage(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	recurringId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	recurring, err := app.Recurring.GetById(recurringId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	form := app.convertRecurringtoForm(recurring)
	form.Action = fmt.Sprintf("/recurring/%d/update", recurringId)
	form.ImageUrl = fmt.Sprintf(
		"%s/recurring/small/%s",
		app.Config.BaseImgUrl,
		safeDeref(recurring.ThumbnailName),
	)

	data := app.newTemplateData(r)
	data.Page = recurringFormPage{
		Title:        "Update Recurring",
		RecurringUrl: fmt.Sprintf("/recurring/%d/%s", recurringId, safeDeref(recurring.Slug)),
		Form:         form,
	}

	app.render(r, w, http.StatusOK, "recurring-form-page.gohtml", "base", data)
}

func (app *Application) recurringUpdate(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	recurringId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	stalerecurring, err := app.Recurring.GetById(recurringId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}
	oldThumbnailName := safeDeref(stalerecurring.ThumbnailName)

	var form recurringForm
	form.Action = fmt.Sprintf("/recurring/%d/update", recurringId)

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.ErrorLog.Print(err)
		return
	}

	app.validateRecurringForm(&form)
	if !form.Valid() {
		form.ImageUrl = fmt.Sprintf(
			"%s/recurring/small/%s",
			app.Config.BaseImgUrl,
			safeDeref(stalerecurring.ThumbnailName),
		)

		app.render(r, w, http.StatusUnprocessableEntity, "recurring-form-page.gohtml", "recurring-form", form)
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
		err = app.saveLargeThumbnail(thumbnailName, "recurring", form.Thumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}

	}

	updaterecurring := app.convertFormtoRecurring(&form)
	updaterecurring.ThumbnailName = &thumbnailName

	slug := models.CreateSlugName(safeDeref(updaterecurring.Title))
	updaterecurring.Slug = &slug

	err = app.Recurring.Update(&updaterecurring)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if form.Thumbnail != nil && stalerecurring.ThumbnailName != nil {
		err = app.deleteImage("recurring", *stalerecurring.ThumbnailName)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	form.ImageUrl = fmt.Sprintf(
		"%s/recurring/small/%s",
		app.Config.BaseImgUrl,
		safeDeref(updaterecurring.ThumbnailName),
	)

	app.render(r, w, http.StatusOK, "recurring-form-page.gohtml", "recurring-form", form)
}

func (app *Application) deleterecurring(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	recurringId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	_, err = app.Recurring.GetById(recurringId)
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
