package main

import (
	"net/http"
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
	app.render(w, http.StatusOK, "admin-add.tmpl.html", data)
}

// type addVideoForm struct {
// 	Title string `form:"title"`
// 	VideoURL string `form:"videoURL"`
// 	Thumbnail

// }

func (app *application) videoAddPost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	// 1) Validate video information (valid date, url, rating)

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

	app.infoLog.Print(r.Form)

	w.WriteHeader(200)
	w.Write([]byte("OK"))

}
