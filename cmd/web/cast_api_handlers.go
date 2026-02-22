package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) adminGetCastAPI(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	cast, err := app.services.Casts.GetAdminCast(sketchId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"cast":        cast.Cast,
		"screenshots": cast.Screenshots,
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverError(r, w, err)
	}
}

func (app *application) updateCastOrderAPI(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	var input struct {
		CastPositions []int `json:"castPositions"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if len(input.CastPositions) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = app.services.Casts.ReorderCast(sketchId, input.CastPositions)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) createCastAPI(w http.ResponseWriter, r *http.Request) {
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

	app.validateCastForm(&form, false)
	if !form.Valid() {
		app.failedValidationResponse(w, r, form.Validator.FieldErrors)
		return
	}

	castMember := convertFormtoCastMember(&form)
	castMember.SketchID = &sketchId

	thumbnail, _ := fileHeaderToBytes(form.CharacterThumbnail)
	profile, _ := fileHeaderToBytes(form.CharacterProfile)

	newCast, err := app.services.Casts.CreateCastMember(&castMember, thumbnail, profile)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cast": newCast}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCastAPI(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("sketch id is invalid"))
		return
	}

	var form castForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("Unable to decode multipart form"))
		return
	}

	app.validateCastForm(&form, true)
	if !form.Valid() {
		app.failedValidationResponse(w, r, form.Validator.FieldErrors)
		return
	}

	castMember := convertFormtoCastMember(&form)
	castMember.SketchID = &sketchId

	thumbnail, _ := fileHeaderToBytes(form.CharacterThumbnail)
	profile, _ := fileHeaderToBytes(form.CharacterProfile)

	updatedCast, err := app.services.Casts.UpdateCastMember(&castMember, thumbnail, profile)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cast": updatedCast}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCastAPI(w http.ResponseWriter, r *http.Request) {
	castIdParam := r.PathValue("castId")
	castId, err := strconv.Atoi(castIdParam)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.services.Casts.DeleteCastmember(castId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
