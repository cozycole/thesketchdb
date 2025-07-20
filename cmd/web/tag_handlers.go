package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

type tagFormPage struct {
	Title string
	Form  tagForm
}

func (app *application) tagAddPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Page = tagFormPage{
		Title: "Add Tag",
		Form: tagForm{
			Action: "/tag/add",
		},
	}
	app.render(r, w, http.StatusOK, "add-tag.gohtml", "base", data)
}

func (app *application) tagAdd(w http.ResponseWriter, r *http.Request) {
	var form tagForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateTagForm(&form)
	if !form.Valid() {
		form.Action = "/tag/add"
		app.render(r, w, http.StatusUnprocessableEntity, "add-tag.gohtml", "tag-form", form)
		return
	}

	tag := convertFormtoTag(&form)
	tagSlug := safeDeref(tag.Name)
	if tag.Category != nil && tag.Category.ID != nil {
		category, err := app.categories.Get(safeDeref(tag.Category.ID))
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			app.serverError(r, w, err)
			return
		}
		tagSlug = safeDeref(category.Name) + " " + tagSlug
	}

	slug := models.CreateSlugName(tagSlug)
	tag.Slug = &slug
	id, err := app.tags.Insert(&tag)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	w.Header().Add("Hx-Redirect", fmt.Sprintf("/tag/%d/update", id))
}

func (app *application) tagUpdatePage(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	tagId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	tag, err := app.tags.Get(tagId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	form := convertTagtoForm(tag)
	form.Action = fmt.Sprintf("/tag/%d/update", tagId)

	data := app.newTemplateData(r)
	data.Page = tagFormPage{
		Title: "Update Tag",
		Form:  form,
	}

	app.render(r, w, http.StatusOK, "add-tag.gohtml", "base", data)
}

func (app *application) tagUpdate(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	tagId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	var form tagForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	_, err = app.tags.Get(tagId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	app.validateTagForm(&form)
	if !form.Valid() {
		form.Action = fmt.Sprintf("/tag/%d/update", tagId)
		app.render(r, w, http.StatusUnprocessableEntity, "add-tag.gohtml", "tag-form", form)
		return
	}

	tag := convertFormtoTag(&form)
	tagSlug := safeDeref(tag.Name)
	if tag.Category != nil && tag.Category.ID != nil {
		category, err := app.categories.Get(safeDeref(tag.Category.ID))
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			app.serverError(r, w, err)
			return
		}
		tagSlug = safeDeref(category.Name) + " " + tagSlug
	}

	slug := models.CreateSlugName(tagSlug)
	tag.Slug = &slug

	err = app.tags.Update(&tag)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.render(r, w, http.StatusOK, "add-tag.gohtml", "tag-form", form)
}

func (app *application) tagRow(w http.ResponseWriter, r *http.Request) {
	app.render(r, w, http.StatusOK, "tag-table.gohtml", "tag-row", views.TagRow{})
}
