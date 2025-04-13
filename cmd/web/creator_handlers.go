package main

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
)

func (app *application) creatorAddPost(w http.ResponseWriter, r *http.Request) {
	var form creatorForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateAddCreatorForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Creator = &form
		app.render(r, w, http.StatusUnprocessableEntity, "add-creator.tmpl.html", "base", data)
		return
	}

	date, _ := time.Parse(time.DateOnly, form.EstablishedDate)
	imgName := models.CreateSlugName(form.Name)

	file, err := form.ProfileImage.Open()
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	defer file.Close()

	mimeType, err := utils.GetMultipartFileMime(file)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	// the insert returns the fullImgName which is {fileName}-{id}.{ext}
	_, slug, fullImgName, err := app.creators.
		Insert(
			form.Name, form.URL, imgName,
			mimeToExt[mimeType], date,
		)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.fileStorage.SaveFile(path.Join("creator", fullImgName), file)
	if err != nil {
		// TODO: We gotta remove the db record on this error
		app.serverError(r, w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/creator/%s", slug), http.StatusSeeOther)
}

func (app *application) creatorView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	creator, err := app.creators.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	videos, err := app.videos.Get(
		&models.Filter{
			Limit:  16,
			Offset: 0,
			SortBy: "az",
			People: []*models.Person{
				&models.Person{ID: creator.ID},
			},
		},
	)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Creator = creator
	data.Videos = videos

	app.render(r, w, http.StatusOK, "view-creator.tmpl.html", "base", data)
}

func (app *application) creatorAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Forms.Creator = &creatorForm{}
	app.render(r, w, http.StatusOK, "add-creator.tmpl.html", "base", data)
}
