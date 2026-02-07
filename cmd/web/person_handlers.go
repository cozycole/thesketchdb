package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/external/wikipedia"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) viewPerson(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	personId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	person, err := app.people.GetById(personId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	popular, _, err := app.sketches.Get(
		&models.Filter{
			Page:      1,
			PageSize:  12,
			SortBy:    "popular",
			PersonIDs: []int{*person.ID},
		},
	)

	showCreatorCounts, err := app.people.GetCreatorShowCounts(personId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	stats, err := app.people.GetPersonStats(personId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	page, err := views.PersonPageView(
		person,
		stats,
		popular,
		showCreatorCounts,
		app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = page
	app.render(r, w, http.StatusOK, "view-person.gohtml", "base", data)
}

type personFormPage struct {
	Title string
	Form  personForm
}

func (app *application) addPersonPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Page = personFormPage{
		Title: "Add Person",
		Form: personForm{
			Action: "/person/add",
		},
	}
	app.render(r, w, http.StatusOK, "add-person.gohtml", "base", data)
}

func (app *application) addPerson(w http.ResponseWriter, r *http.Request) {
	var form personForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validatePersonForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Page = personFormPage{
			Title: "Add Person",
			Form:  form,
		}
		app.infoLog.Printf("%+v\n", form)
		app.render(r, w, http.StatusUnprocessableEntity, "add-person.gohtml", "base", data)
		return
	}

	person := convertFormtoPerson(&form)
	slug := models.CreateSlugName(views.PrintPersonName(&person))
	person.Slug = &slug

	thumbName, err := generateThumbnailName(form.ProfileImage)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	person.ProfileImg = &thumbName

	if person.WikiPage != nil {
		description, err := wikipedia.GetExtract(*person.WikiPage)
		if nil == err {
			person.Description = &description
		}
	}

	id, err := app.people.Insert(&person)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.saveLargeProfile(*person.ProfileImg, "person", form.ProfileImage)
	if err != nil {
		app.serverError(r, w, err)
		app.people.Delete(id)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/person/%d/%s", id, slug), http.StatusSeeOther)
}

func (app *application) updatePersonPage(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	personId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	person, err := app.people.GetById(personId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	personForm := convertPersontoForm(person)
	personForm.ImageUrl = fmt.Sprintf("%s/person/small/%s", app.baseImgUrl, personForm.ImageUrl)
	personForm.Action = fmt.Sprintf("/person/%d/update", personForm.ID)

	data := app.newTemplateData(r)
	data.Page = personFormPage{
		Title: "Update Person",
		Form:  personForm,
	}

	app.render(r, w, http.StatusOK, "add-person.gohtml", "base", data)
}

func (app *application) updatePerson(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	personId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	oldPerson, err := app.people.GetById(personId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	oldProfileImgName := safeDeref(oldPerson.ProfileImg)

	var form personForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validatePersonForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		form.ImageUrl = fmt.Sprintf("%s/person/small/%s", app.baseImgUrl, oldProfileImgName)
		form.Action = fmt.Sprintf("/person/%d/update", form.ID)
		data.Page = personFormPage{
			Title: "Update Person",
			Form:  form,
		}
		app.render(r, w, http.StatusUnprocessableEntity, "add-person.gohtml", "base", data)
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
		err = app.saveLargeProfile(profileImgName, "person", form.ProfileImage)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	newPerson := convertFormtoPerson(&form)
	newPerson.ProfileImg = &profileImgName
	slug := models.CreateSlugName(views.PrintPersonName(&newPerson))
	newPerson.Slug = &slug

	if newPerson.WikiPage != nil {
		description, err := wikipedia.GetExtract(*newPerson.WikiPage)
		if nil == err {
			newPerson.Description = &description
		}
	}

	err = app.people.Update(&newPerson)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if form.ProfileImage != nil && oldPerson.ProfileImg != nil {
		err = app.deleteImage("person", *oldPerson.ProfileImg)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/person/%d/%s", personId, slug), http.StatusSeeOther)
}
