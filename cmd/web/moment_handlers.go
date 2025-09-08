package main

import (
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
)

type UpdateMoment struct {
	MomentID   int
	MomentForm momentForm
	QuoteForm  quoteForm
}

func (app *application) momentAdd(w http.ResponseWriter, r *http.Request) {
	var form momentForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateMomentForm(&form)
	if !form.Valid() {
		form.Action = "/moment/add"
		app.render(r, w, http.StatusUnprocessableEntity, "sketch-form-page.gohtml", "moment-form", form)
		return
	}

	moment := app.convertFormtoMoment(&form)

	_, err = app.moments.Insert(form.SketchID, &moment)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	moments, err := app.moments.GetBySketch(form.SketchID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	updateMoments := []UpdateMoment{}
	for _, m := range moments {
		mid := safeDeref(m.ID)

		momentForm := app.convertMomenttoForm(m)
		momentForm.Action = fmt.Sprintf("/moment/%d", mid)
		quoteForm := quoteForm{MomentID: mid, SketchID: form.SketchID}
		updateMoments = append(updateMoments, UpdateMoment{mid, momentForm, quoteForm})
	}

	app.render(r, w, http.StatusOK, "sketch-form-page.gohtml", "moments", updateMoments)
}

func (app *application) momentDelete(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	momentId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	err = app.moments.Delete(momentId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) momentUpdate(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	momentId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}
	var form momentForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateMomentForm(&form)
	if !form.Valid() {
		form.Action = fmt.Sprintf("/moment/%d", momentId)
		app.render(r, w, http.StatusUnprocessableEntity, "sketch-form-page.gohtml", "moment-form", form)
		return
	}

	moment := app.convertFormtoMoment(&form)

	err = app.moments.Update(&moment)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	moments, err := app.moments.GetBySketch(form.SketchID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	updateMoments := []UpdateMoment{}
	for _, m := range moments {
		mid := safeDeref(m.ID)
		momentForm := app.convertMomenttoForm(m)
		momentForm.Action = fmt.Sprintf("/moment/%d", mid)
		quoteForm := app.convertQuotestoForm(form.SketchID, mid, m.Quotes)
		updateMoments = append(updateMoments, UpdateMoment{mid, momentForm, quoteForm})
	}

	app.render(r, w, http.StatusOK, "sketch-form-page.gohtml", "moments", updateMoments)
}

func (app *application) quoteUpdate(w http.ResponseWriter, r *http.Request) {
	var form quoteForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	moment, err := app.moments.GetById(form.MomentID)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	sketchId := safeDeref(moment.Sketch.ID)

	app.validateQuoteForm(&form)
	if !form.Valid() {
		form.SketchID = sketchId
		app.padCastMemberInputs(&form)
		app.render(r, w, http.StatusUnprocessableEntity, "sketch-form-page.gohtml", "quote-table", form)
		return
	}

	quotes := app.convertFormtoQuotes(&form)
	err = app.moments.BatchUpdateQuotes(form.MomentID, quotes)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	moment, err = app.moments.GetById(form.MomentID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	updatedQuoteForm := app.convertQuotestoForm(form.SketchID, *moment.ID, moment.Quotes)

	updatedQuoteForm.Flash.Level = "success"
	updatedQuoteForm.Flash.Message = "Successfully updated quotes!"

	app.render(r, w, http.StatusOK, "sketch-form-page.gohtml", "quote-table", updatedQuoteForm)
}

// needed function for returning errored forms
func (app *application) padCastMemberInputs(form *quoteForm) {
	for i := range len(form.CastMemberID) {
		id := form.CastMemberID[i]
		if id == 0 {
			form.CastMemberName = append(form.CastMemberName, "")
			form.CastImageUrl = append(form.CastImageUrl, "")
		} else {
			cm, _ := app.cast.GetById(id)
			form.CastMemberName = append(form.CastMemberName, views.PrintCastBlurb(cm))
			imgUrl := fmt.Sprintf("%s/cast/profile/small/%s", app.baseImgUrl, safeDeref(cm.ProfileImg))
			form.CastImageUrl = append(form.CastImageUrl, imgUrl)
		}
	}
}
