package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) addCastPage(w http.ResponseWriter, r *http.Request) {
	// this is a sub template used by htmx to insert into update page. If you want
	// to make a separate page for it, check the headers and have different template loaded
	// based on whether htmx requested it or not
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	exists, err := app.sketches.Exists(sketchId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if !exists {
		app.notFound(w)
		return
	}

	// isHxRequest := r.Header.Get("HX-Request") == "true"
	// if isHxRequest {
	// 	return
	// }

	form := castForm{
		Action: fmt.Sprintf("/sketch/%d/cast", sketchId),
	}
	app.render(r, w, http.StatusOK, "cast-form.gohtml", "cast-form-modal", form)
}

func (app *application) addCast(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	var form castForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateCastForm(&form)
	if !form.Valid() {
		form.Action = fmt.Sprintf("/sketch/%d/cast", sketchId)
		app.render(r, w, http.StatusUnprocessableEntity, "cast-form.gohtml", "cast-form", form)
		return
	}

	castMember := convertFormtoCastMember(&form)

	if castMember.ThumbnailFile != nil {
		thumbName, err := generateThumbnailName(castMember.ThumbnailFile)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
		castMember.ThumbnailName = &thumbName
	}

	if castMember.ProfileFile != nil {
		profileName, err := generateThumbnailName(castMember.ProfileFile)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
		castMember.ProfileImg = &profileName
	}

	_, err = app.cast.Insert(sketchId, &castMember)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.saveCastImages(&castMember)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	cast, err := app.cast.GetCastMembers(sketchId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	castTableView := views.CastTableView(cast, sketchId, app.baseImgUrl)
	app.render(r, w, http.StatusOK, "cast-table.gohtml", "cast-table", castTableView)
}

func (app *application) updateCastPage(w http.ResponseWriter, r *http.Request) {
	castIdParam := r.PathValue("id")
	castId, err := strconv.Atoi(castIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	castMember, err := app.cast.GetById(castId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	form := convertCastMembertoForm(castMember)
	form.ID = safeDeref(castMember.ID)
	form.Action = fmt.Sprintf("/cast/%d/update", castId)
	form.ThumbnailName = fmt.Sprintf("%s/cast/thumbnail/%s",
		app.baseImgUrl, form.ThumbnailName)
	form.ProfileImage = fmt.Sprintf("%s/cast/profile/small/%s",
		app.baseImgUrl, form.ProfileImage)

	isHxRequest := r.Header.Get("HX-Request") == "true"

	if isHxRequest {
		// put the render below here and render full page if necessary
	}

	app.render(r, w, http.StatusOK, "cast-form.gohtml", "cast-form-modal", form)
}

func (app *application) updateCast(w http.ResponseWriter, r *http.Request) {
	castIdParam := r.PathValue("castId")
	castId, err := strconv.Atoi(castIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	var form castForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	staleMember, err := app.cast.GetById(castId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	currentThumbnail := safeDeref(staleMember.ThumbnailName)
	currentProfile := safeDeref(staleMember.ProfileImg)

	app.validateCastForm(&form)
	if !form.Valid() {
		form.ThumbnailName = fmt.Sprintf(
			"%s/cast/thumbnail/small/%s", app.baseImgUrl, currentThumbnail)
		form.ProfileImage = fmt.Sprintf(
			"%s/cast/profile/small/%s", app.baseImgUrl, currentProfile)
		form.Action = fmt.Sprintf("/cast/%d/update", castId)
		app.render(r, w, http.StatusUnprocessableEntity, "cast-form.gohtml", "cast-form", form)
		return
	}

	newMember := convertFormtoCastMember(&form)

	if form.CharacterThumbnail != nil {
		var err error
		currentThumbnail, err = generateThumbnailName(newMember.ThumbnailFile)
		if err != nil {
			app.serverError(r, w, err)
			return
		}

		err = app.saveMediumThumbnail(currentThumbnail, "/cast/thumbnail", form.CharacterThumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	if form.CharacterProfile != nil {
		var err error
		currentProfile, err = generateThumbnailName(newMember.ProfileFile)
		if err != nil {
			app.serverError(r, w, err)
			return
		}

		err = app.saveMediumProfile(
			currentProfile, "/cast/profile", form.CharacterProfile,
		)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	newMember.ThumbnailName = &currentThumbnail
	newMember.ProfileImg = &currentProfile

	err = app.cast.Update(&newMember)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	sketchId := safeDeref(staleMember.SketchID)
	cast, err := app.cast.GetCastMembers(sketchId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if form.CharacterThumbnail != nil && newMember.ThumbnailName != nil {
		app.deleteImage(fmt.Sprintf("cast/thumbnail"), safeDeref(staleMember.ThumbnailName))
	}

	if form.CharacterProfile != nil && newMember.ProfileImg != nil {
		app.deleteImage(fmt.Sprintf("cast/profile"), safeDeref(staleMember.ProfileImg))
	}

	castTable := views.CastTableView(cast, sketchId, app.baseImgUrl)
	app.render(r, w, http.StatusOK, "cast-table.gohtml", "cast-table", castTable)
}

func (app *application) orderCast(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	r.ParseForm()
	castPositionParams := r.PostForm["position"]
	castPositions, err := convertStringsToInts(castPositionParams)
	if err != nil {
		app.badRequest(w)
		return
	}

	err = app.cast.UpdatePositions(castPositions)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	cast, err := app.cast.GetCastMembers(sketchId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.CastTableView(cast, sketchId, app.baseImgUrl)
	app.render(r, w, http.StatusOK, "cast-table.gohtml", "cast-table", table)
}

func (app *application) deleteCast(w http.ResponseWriter, r *http.Request) {
	castIdParam := r.PathValue("castId")
	castId, err := strconv.Atoi(castIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	castMember, err := app.cast.GetById(castId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	err = app.cast.Delete(castId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.deleteImage(fmt.Sprintf("cast/thumbnail"), safeDeref(castMember.ThumbnailName))
	app.deleteImage(fmt.Sprintf("cast/profile"), safeDeref(castMember.ProfileImg))

	sketchId := safeDeref(castMember.SketchID)
	cast, err := app.cast.GetCastMembers(sketchId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.CastTableView(cast, sketchId, app.baseImgUrl)
	app.render(r, w, http.StatusOK, "cast-table.gohtml", "cast-table", table)
}
