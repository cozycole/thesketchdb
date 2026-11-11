package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *Application) creatorView(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	creatorId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	creator, err := app.Creators.GetById(creatorId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	popularSketches, _, err := app.Sketches.Get(
		&models.Filter{
			Page:       1,
			PageSize:   12,
			SortBy:     "popular",
			CreatorIDs: []int{*creator.ID},
		},
	)

	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	cast, err := app.Creators.GetCast(creatorId)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	page, err := views.CreatorPageView(
		creator,
		popularSketches,
		cast, app.Config.BaseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
	}

	data.Page = page
	app.render(r, w, http.StatusOK, "view-creator.gohtml", "base", data)
}

type creatorFormPage struct {
	Title string
	Form  creatorForm
}

func (app *Application) addCreatorPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Page = creatorFormPage{
		Title: "Add Creator",
		Form: creatorForm{
			Action: "/creator/add",
		},
	}

	app.render(r, w, http.StatusOK, "add-creator.gohtml", "base", data)
}

func (app *Application) addCreator(w http.ResponseWriter, r *http.Request) {
	var form creatorForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateCreatorForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Page = creatorFormPage{
			Title: "Add Creator",
			Form:  form,
		}
		app.render(r, w, http.StatusUnprocessableEntity, "add-creator.gohtml", "base", data)
		return
	}

	creator := convertFormtoCreator(&form)
	slug := models.CreateSlugName(form.Name)
	creator.Slug = &slug

	thumbName, err := generateThumbnailName(form.ProfileImage)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	creator.ProfileImage = &thumbName

	id, err := app.Creators.Insert(&creator)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.saveLargeProfile(thumbName, "creator", form.ProfileImage)
	if err != nil {
		app.serverError(r, w, err)
		app.Creators.Delete(id)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/creator/%d/%s", id, slug), http.StatusSeeOther)
}

func (app *Application) updateCreatorPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	creatorId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.ErrorLog.Print(err)
		return
	}

	creator, err := app.Creators.GetById(creatorId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	creatorForm := convertCreatortoForm(creator)
	creatorForm.ImageUrl = fmt.Sprintf("%s/creator/small/%s", app.Config.BaseImgUrl, creatorForm.ImageUrl)
	creatorForm.Action = fmt.Sprintf("/creator/%d/update", creatorForm.ID)

	data := app.newTemplateData(r)
	data.Page = creatorFormPage{
		Title: "Update Creator",
		Form:  creatorForm,
	}

	app.render(r, w, http.StatusOK, "add-creator.gohtml", "base", data)
}

func (app *Application) updateCreator(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	creatorId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.ErrorLog.Print(err)
		return
	}

	staleCreator, err := app.Creators.GetById(creatorId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	oldProfileImgName := safeDeref(staleCreator.ProfileImage)

	var form creatorForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateCreatorForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		form.ImageUrl = fmt.Sprintf("%s/creator/small/%s", app.Config.BaseImgUrl, oldProfileImgName)
		form.Action = fmt.Sprintf("/creator/%d/update", form.ID)
		data.Page = creatorFormPage{
			Title: "Update Creator",
			Form:  form,
		}
		app.render(r, w, http.StatusUnprocessableEntity, "add-creator.gohtml", "base", data)
		return
	}

	profileImgName := oldProfileImgName
	if form.ProfileImage != nil {
		var err error
		profileImgName, err = generateThumbnailName(form.ProfileImage)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
		err = app.saveLargeProfile(profileImgName, "creator", form.ProfileImage)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	updatedCreator := convertFormtoCreator(&form)
	updatedCreator.ProfileImage = &profileImgName
	slug := models.CreateSlugName(form.Name)
	updatedCreator.Slug = &slug

	err = app.Creators.Update(&updatedCreator)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if form.ProfileImage != nil && staleCreator.ProfileImage != nil {
		err = app.deleteImage("creator", *staleCreator.ProfileImage)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/creator/%d/%s", creatorId, slug), http.StatusSeeOther)

}
