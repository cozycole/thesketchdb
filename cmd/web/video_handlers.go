package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/internal/models"
)

func (app *application) videoView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	video, err := app.videos.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if ok {
		hasLike, _ := app.videos.HasLike(video.ID, user.ID)
		video.Liked = hasLike
	}

	data := app.newTemplateData(r)
	if video.YoutubeID != nil && *video.YoutubeID != "" {
		videoUrl := fmt.Sprintf("https://www.youtube.com/watch?v=%s", *video.YoutubeID)
		video.URL = &videoUrl
	}
	data.Video = video

	app.render(w, http.StatusOK, "view-video.tmpl.html", "base", data)
}

func (app *application) videoAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Need to initialize form data since the template needs it to
	// render. It's a good place to put default values for the fields
	data.Form = addVideoForm{}
	app.render(w, http.StatusOK, "add-video.tmpl.html", "base", data)
}

func (app *application) videoAddPost(w http.ResponseWriter, r *http.Request) {
	var form addVideoForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateAddVideoForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add-video.tmpl.html", "base", data)
		return
	}

	video, err := convertFormToVideo(&form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = addVideoImageNames(&video)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// NOTE: This mutates the video struct by adding the newly created db serial id
	// to the id field
	err = app.videos.Insert(&video)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.saveVideoImages(&video)
	if err != nil {
		app.serverError(w, err)
		// TODO: delete video entry now
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/video/%d/%s", video.ID, video.Slug), http.StatusSeeOther)
}

func (app *application) videoAddLike(w http.ResponseWriter, r *http.Request) {
	videoIdParam := r.PathValue("videoId")
	videoId, err := strconv.Atoi(videoIdParam)
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

	err = app.users.AddLike(user.ID, videoId)
	if err != nil {
		// check if problem with primary key constraint
		app.badRequest(w)
		return
	}
}

func (app *application) videoRemoveLike(w http.ResponseWriter, r *http.Request) {
	videoIdParam := r.PathValue("videoId")
	videoId, err := strconv.Atoi(videoIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	err = app.users.RemoveLike(user.ID, videoId)
	if err != nil {
		app.badRequest(w)
		return
	}
}
