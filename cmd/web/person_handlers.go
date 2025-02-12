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

func (app *application) personView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	person, err := app.people.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	videos, err := app.videos.GetByPerson(*person.ID)

	data := app.newTemplateData(r)
	data.Person = person
	data.Videos = videos

	app.render(w, http.StatusOK, "view-person.tmpl.html", "base", data)
}

func (app *application) personAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = addPersonForm{}
	app.render(w, http.StatusOK, "add-person.tmpl.html", "base", data)
}

func (app *application) personAddPost(w http.ResponseWriter, r *http.Request) {
	var form addPersonForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateAddPersonForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add-person.tmpl.html", "base", data)
		return
	}

	date, _ := time.Parse(time.DateOnly, form.BirthDate)
	imgName := models.CreateSlugName(form.First+" "+form.Last, maxFileNameLength)

	file, err := form.ProfileImage.Open()
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer file.Close()

	mimeType, err := utils.GetMultipartFileMime(file)
	if err != nil {
		app.serverError(w, err)
		return
	}

	_, slug, fullImgName, err := app.people.
		Insert(
			form.First, form.Last, imgName,
			mimeToExt[mimeType], date,
		)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.fileStorage.SaveFile(path.Join("person", fullImgName), file)
	if err != nil {
		// TODO: We gotta remove the db record on this error
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/person/%s", slug), http.StatusSeeOther)
}
