package main

import (
	"net/http"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	videos, err := app.videos.GetAll(0)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Videos = videos

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) search(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "search.tmpl.html", data)
}

func (app *application) videoAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "add-video.tmpl.html", data)
}

type addVideoForm struct {
	Title      string    `form:"title"`
	VideoURL   string    `form:"videoURL"`
	Rating     string    `form:"rating"`
	UploadDate time.Time `form:"uploadDate"`
	Creator    string    `form:"creator"`
	Actors     []string  `form:"actors"`
}

func (app *application) videoAddPost(w http.ResponseWriter, r *http.Request) {
	var form addVideoForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	// 1) Validate video information (url, rating, title)

	// 2) Validate creator exists by getting its id

	// 4) Validate actors by getting their ids

	// 4) Validate uploaded image, then save video thumbnail path and give it a name

	// 5) Insert video

	// 6) Insert video creator relations

	// 7) Insert video actor relations

	// file, _, err := r.FormFile("thumbnail")
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	// defer file.Close()

	// dst, err := os.Create(header.Filename)

	app.infoLog.Println(form)

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
