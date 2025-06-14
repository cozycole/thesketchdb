package main

import (
	"errors"
	"fmt"
	"net/http"

	"sketchdb.cozycole.net/internal/models"
)

func (app *application) categoriesView(w http.ResponseWriter, r *http.Request) {
	categories, err := app.categories.GetAll()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Categories = &categories
	app.render(r, w, http.StatusOK, "view-categories.gohtml", "base", data)
}

func (app *application) categoryAddPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Forms.Category = &categoryForm{}
	app.render(r, w, http.StatusOK, "add-category.gohtml", "base", data)
}

func (app *application) categoryAdd(w http.ResponseWriter, r *http.Request) {
	var form categoryForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateCategoryForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Category = &form
		app.render(r, w, http.StatusUnprocessableEntity, "add-category.gohtml", "base", data)
		return
	}

	category := convertFormtoCategory(&form)
	slug := models.CreateSlugName(*category.Name)
	category.Slug = &slug
	_, err = app.categories.Insert(&category)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Category added: %s", *category.Name))

	data := app.newTemplateData(r)
	data.Forms.Category = &categoryForm{}
	app.render(r, w, http.StatusOK, "add-category.gohtml", "base", data)
}
