package main

import (
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/internal/models"
)

func (app *application) addCastPage(w http.ResponseWriter, r *http.Request) {
	var form castForm
	data := app.newTemplateData(r)
	data.CastMember = &models.CastMember{}
	data.Forms.Cast = &form
	// this is a sub template used by htmx to insert into update page. If you want
	// to make a separate page for it, check the headers and have different template loaded
	// based on whether htmx requested it or not
	app.render(r, w, http.StatusOK, "actor-input.tmpl.html", "actor-input", data)
}

func (app *application) addCast(w http.ResponseWriter, r *http.Request) {
	videoIdParam := r.PathValue("id")
	videoId, err := strconv.Atoi(videoIdParam)
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

	app.validateAddCast(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.CastMember = &models.CastMember{}
		data.Forms.Cast = &form
		app.render(r, w, http.StatusUnprocessableEntity, "actor-input.tmpl.html", "actor-input", data)
		return
	}

	castMember := convertFormtoCastMember(&form)

	castId, err := app.cast.Insert(videoId, &castMember)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	castMember.ID = &castId
	thumbName, err := generateThumbnailName(castMember.ThumbnailFile)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	castMember.ThumbnailName = &thumbName
	err = app.cast.InsertThumbnailName(*castMember.ID, *castMember.ThumbnailName)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.saveCastImages(&castMember)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	cast, _ := app.cast.GetCastMembers(videoId)
	data := app.newTemplateData(r)
	data.Cast = cast
	app.render(r, w, http.StatusOK, "cast-table.tmpl.html", "cast-table", data)
}

func (app *application) updateCast(w http.ResponseWriter, r *http.Request) {
	videoIdParam, castIdParam := r.PathValue("id"), r.PathValue("castId")
	videoId, err := strconv.Atoi(videoIdParam)
	castId, err2 := strconv.Atoi(castIdParam)
	if err != nil || err2 != nil {
		app.badRequest(w)
		return
	}
	app.infoLog.Printf("%d %d\n", videoId, castId)
}
