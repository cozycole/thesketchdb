package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) characterView(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	characterdId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	character, err := app.characters.GetById(characterdId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	popularSketches, err := app.sketches.Get(
		&models.Filter{
			Limit:  16,
			Offset: 0,
			SortBy: "az",
			Characters: []*models.Character{
				&models.Character{ID: character.ID},
			},
		},
	)

	data := app.newTemplateData(r)
	page, err := views.CharacterPageView(
		character,
		popularSketches,
		app.baseImgUrl,
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

func (app *application) addCharacterPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Page = characterFormPage{
		Title: "Add Character",
		Form: characterForm{
			Action: "/character/add",
		},
	}

	app.render(r, w, http.StatusOK, "character-form.gohtml", "base", data)
}

func (app *application) addCharacter(w http.ResponseWriter, r *http.Request) {
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
		exists, err := app.people.Exists(form.PersonID)
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

	id, err := app.characters.Insert(&character)
	if err != nil {
		app.serverError(r, w, err)
		app.characters.Delete(id)
		return
	}

	err = app.saveProfileImage(thumbName, "character", form.ProfileImage)
	if err != nil {
		app.serverError(r, w, err)
		app.characters.Delete(id)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/character/%d/%s", id, slug), http.StatusSeeOther)
}

func (app *application) updateCharacterPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	characterId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	character, err := app.characters.GetById(characterId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	characterForm := convertCharactertoForm(character)
	characterForm.ImageUrl = fmt.Sprintf("%s/character/%s", app.baseImgUrl, characterForm.ImageUrl)
	characterForm.Action = fmt.Sprintf("/character/%d/update", characterForm.ID)

	data := app.newTemplateData(r)
	data.Page = characterFormPage{
		Title: "Update Character",
		Form:  characterForm,
	}

	app.render(r, w, http.StatusOK, "character-form.gohtml", "base", data)
}

func (app *application) updateCharacter(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	characterId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	staleCharacter, err := app.characters.GetById(characterId)
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
		form.ImageUrl = fmt.Sprintf("%s/character/%s", app.baseImgUrl, oldProfileImgName)
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
		err = app.saveProfileImage(profileImgName, "character", form.ProfileImage)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	updatedCharacter := convertFormtoCharacter(&form)
	updatedCharacter.Image = &profileImgName
	slug := models.CreateSlugName(form.Name)
	updatedCharacter.Slug = &slug
	err = app.characters.Update(&updatedCharacter)
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
