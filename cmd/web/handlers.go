package main

import (
	"net/http"

	"sketchdb.cozycole.net/internal/models"
)

var maxFileNameLength = 50
var pageSize = 16

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	videos, err := app.videos.GetAll(8)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Videos = videos

	app.render(w, http.StatusOK, "home.tmpl.html", "base", data)
}

func (app *application) browse(w http.ResponseWriter, r *http.Request) {
	browseSections := make(map[string][]*models.Video)
	limit := 8
	offset := 0

	// First add "custom" sections (ex: latest, trending, recommended/because you liked X)
	latest, err := app.videos.GetLatest(limit, offset)
	if err != nil {
		app.errorLog.Println(err)
	}
	browseSections["Latest"] = latest

	jamesId := 3
	actorVideos, err := app.videos.GetByPerson(jamesId)
	if err != nil {
		app.errorLog.Println(err)
	}
	browseSections["Sketches Featuring James Hartnett"] = actorVideos

	data := app.newTemplateData(r)
	data.BrowseSections = browseSections
	app.render(w, http.StatusOK, "browse.tmpl.html", "base", data)
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}
