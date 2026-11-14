package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cozycole/thesketchdb/cmd/web/views"
	"github.com/cozycole/thesketchdb/internal/models"
)

func (app *Application) characterView(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	characterdId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	character, err := app.Characters.GetById(characterdId)
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
			Page:         1,
			PageSize:     12,
			SortBy:       "popular",
			CharacterIDs: []int{*character.ID},
		},
	)

	data := app.newTemplateData(r)
	page, err := views.CharacterPageView(
		character,
		popularSketches,
		app.Config.BaseImgUrl,
	)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = page
	app.render(r, w, http.StatusOK, "view-character.gohtml", "base", data)
}

type characterFormPage struct {
	Title string
	Form  characterForm
}

func (app *Application) addCharacterPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Page = characterFormPage{
		Title: "Add Character",
		Form: characterForm{
			Action: "/character/add",
		},
	}

	app.render(r, w, http.StatusOK, "character-form.gohtml", "base", data)
}

func (app *Application) addCharacter(w http.ResponseWriter, r *http.Request) {
	var form characterForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateCharacterForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Page = characterFormPage{
			Title: "Add Character",
			Form:  form,
		}
		app.render(r, w, http.StatusUnprocessableEntity, "character-form.gohtml", "base", data)
		return
	}

	if form.PersonID != 0 {
		exists, err := app.People.Exists(form.PersonID)
		if err != nil {
			app.serverError(r, w, err)
			return
		}

		if !exists {
			data := app.newTemplateData(r)
			form.AddFieldError("type", "Person does not exist, please select a dropdown")
			data.Page = characterFormPage{
				Title: "Add Character",
				Form:  form,
			}

			app.render(r, w, http.StatusUnprocessableEntity, "character-form.gohtml", "base", data)
			return
		}
	}

	character := convertFormtoCharacter(&form)
	slug := models.CreateSlugName(form.Name)
	character.Slug = &slug

	thumbName, err := generateThumbnailName(form.ProfileImage)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	character.Image = &thumbName

	id, err := app.Characters.Insert(&character)
	if err != nil {
		app.serverError(r, w, err)
		app.Characters.Delete(id)
		return
	}

	err = app.saveLargeProfile(thumbName, "character", form.ProfileImage)
	if err != nil {
		app.serverError(r, w, err)
		app.Characters.Delete(id)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/character/%d/%s", id, slug), http.StatusSeeOther)
}

func (app *Application) updateCharacterPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	characterId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.ErrorLog.Print(err)
		return
	}

	character, err := app.Characters.GetById(characterId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	characterForm := convertCharactertoForm(character)
	characterForm.ImageUrl = fmt.Sprintf("%s/character/small/%s", app.Config.BaseImgUrl, characterForm.ImageUrl)
	characterForm.Action = fmt.Sprintf("/character/%d/update", characterForm.ID)

	data := app.newTemplateData(r)
	data.Page = characterFormPage{
		Title: "Update Character",
		Form:  characterForm,
	}

	app.render(r, w, http.StatusOK, "character-form.gohtml", "base", data)
}

func (app *Application) updateCharacter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	characterId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.ErrorLog.Print(err)
		return
	}

	staleCharacter, err := app.Characters.GetById(characterId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	oldProfileImgName := safeDeref(staleCharacter.Image)

	var form characterForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateCharacterForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		form.ImageUrl = fmt.Sprintf("%s/character/small/%s", app.Config.BaseImgUrl, oldProfileImgName)
		form.Action = fmt.Sprintf("/character/%d/update", form.ID)
		data.Page = characterFormPage{
			Title: "Update Character",
			Form:  form,
		}
		app.render(r, w, http.StatusUnprocessableEntity, "character-form.gohtml", "base", data)
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
		err = app.saveLargeProfile(profileImgName, "character", form.ProfileImage)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	updatedCharacter := convertFormtoCharacter(&form)
	updatedCharacter.Image = &profileImgName
	slug := models.CreateSlugName(form.Name)
	updatedCharacter.Slug = &slug
	err = app.Characters.Update(&updatedCharacter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if form.ProfileImage != nil && staleCharacter.Image != nil {
		err = app.deleteImage("character", *staleCharacter.Image)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/character/%d/%s", characterId, slug), http.StatusSeeOther)
}
