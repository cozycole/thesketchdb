package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

type categoryPage struct {
	Title string
	Form  categoryForm
}

func (app *application) categoryAddPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Page = categoryPage{
		Title: "Add Category",
		Form: categoryForm{
			Action: "/category/add",
		},
	}

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
		form.Action = "/category/add"
		app.render(r, w, http.StatusUnprocessableEntity, "add-category.gohtml", "category-form", form)
		return
	}

	category := convertFormtoCategory(&form)
	slug := models.CreateSlugName(*category.Name)
	category.Slug = &slug
	id, err := app.categories.Insert(&category)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	w.Header().Add("Hx-Redirect", fmt.Sprintf("/category/%d/update", id))
}

func (app *application) categoryUpdatePage(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	categoryId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}
	category, err := app.categories.Get(categoryId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	form := convertCategoryToForm(category)
	form.Action = fmt.Sprintf("/category/%d/update", safeDeref(category.ID))

	isHxRequest := r.Header.Get("HX-Request") == "true"
	if isHxRequest {
		app.render(r, w, http.StatusOK, "add-category.gohtml", "category-form", form)
		return
	}

	data := app.newTemplateData(r)
	data.Page = categoryPage{
		Title: "Update Category",
		Form:  form,
	}

	app.render(r, w, http.StatusOK, "add-category.gohtml", "base", data)
}

func (app *application) categoryUpdate(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	categoryId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	staleCategory, err := app.categories.Get(categoryId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	var form categoryForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateCategoryForm(&form)
	if !form.Valid() {
		form.Action = fmt.Sprintf("/category/%d/update", safeDeref(staleCategory.ID))
		app.render(r, w, http.StatusUnprocessableEntity, "add-category.gohtml", "category-form", form)
		return
	}

	updatedCategory := convertFormtoCategory(&form)
	slug := models.CreateSlugName(safeDeref(updatedCategory.Name))
	updatedCategory.Slug = &slug
	err = app.categories.Update(&updatedCategory)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	form = convertCategoryToForm(&updatedCategory)
	app.render(r, w, http.StatusOK, "add-category.gohtml", "category-form", form)
}
