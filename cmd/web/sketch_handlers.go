package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) sketchView(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(idParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	sketch, err := app.sketches.GetById(sketchId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	tags, err := app.tags.GetBySketch(sketchId)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	var userSketchInfo *models.UserSketchInfo
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if ok && user.ID != nil {
		hasLike, _ := app.sketches.HasLike(*sketch.ID, *user.ID)
		sketch.Liked = &hasLike
		userSketchInfo, err = app.users.GetUserSketchInfo(*user.ID, sketchId)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			app.serverError(r, w, err)
			return
		}
	}

	var userId *int
	if user == nil {
		userId = nil
	} else {
		userId = user.ID
	}
	quotes, err := app.quotes.GetBySketch(sketchId, userId)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	sketchPage, err := views.SketchPageView(sketch, quotes, tags, userSketchInfo, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = sketchPage

	app.render(r, w, http.StatusOK, "view-sketch.gohtml", "base", data)
}

func (app *application) sketchAddLike(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok || nil == user {
		app.infoLog.Println("User not logged in!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = app.users.AddLike(*user.ID, sketchId)
	if err != nil {
		// check if problem with primary key constraint
		app.badRequest(w)
		return
	}
}

func (app *application) sketchRemoveLike(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	err = app.users.RemoveLike(*user.ID, sketchId)
	if err != nil {
		app.badRequest(w)
		return
	}
}

func createSketchSlug(sketch *models.Sketch) string {
	var slugInput string
	if sketch.Episode != nil && sketch.Episode.GetShow() != nil {
		episode := sketch.Episode
		showString := safeDeref(episode.GetShow().Name)
		seasonNumber := safeDeref(episode.Season.Number)
		episodeNumber := safeDeref(episode.Number)
		slugInput += fmt.Sprintf("%s s%d e%d", showString, seasonNumber, episodeNumber)
	}

	if sketch.Creator != nil {
		slugInput += safeDeref(sketch.Creator.Name)
	}

	if slugInput == "" {
		return safeDeref(sketch.Title)
	}

	return models.CreateSlugName(slugInput + " " + safeDeref(sketch.Title))
}

func (app *application) sketchUpdateRating(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok || user.ID == nil {
		w.Header().Add("Hx-Redirect", "/login")
		http.Redirect(w, r, "/login", http.StatusOK)
		return
	}

	_, err = app.sketches.GetById(sketchId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	var form sketchRatingForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequest(w)
		return
	}

	app.validateSketchRatingForm(&form)
	if !form.Valid() {
		app.badRequest(w)
		return
	}

	// check if they already have a rating
	userSketchInfo, err := app.users.GetUserSketchInfo(*user.ID, sketchId)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	if safeDeref(userSketchInfo.Rating) == 0 {
		err = app.users.AddRating(*user.ID, sketchId, form.Rating)
	} else {
		err = app.users.UpdateRating(*user.ID, sketchId, form.Rating)
	}
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	sketch, err := app.sketches.GetById(sketchId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	userSketchInfo, err = app.users.GetUserSketchInfo(*user.ID, sketchId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	ratingView := views.SketchRatingView(userSketchInfo, sketch)
	app.render(r, w, http.StatusOK, "sketch-rating.gohtml", "sketch-rating", ratingView)
}

func (app *application) sketchDeleteRating(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok || user.ID == nil {
		w.Header().Add("Hx-Redirect", "/login")
		http.Redirect(w, r, "/login", http.StatusOK)
		return
	}

	_, err = app.sketches.GetById(sketchId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	// check if they already have a rating
	userSketchInfo, err := app.users.GetUserSketchInfo(*user.ID, sketchId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if safeDeref(userSketchInfo.Rating) != 0 {
		err = app.users.DeleteRating(*user.ID, sketchId)
	}
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	sketch, err := app.sketches.GetById(sketchId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	ratingView := views.SketchRatingView(&models.UserSketchInfo{}, sketch)
	app.render(r, w, http.StatusOK, "sketch-rating.gohtml", "sketch-rating", ratingView)
}
