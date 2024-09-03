package main

import (
	"net/http"
	"strings"

	"sketchdb.cozycole.net/internal/validator"
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

	// Need to initialize form data since the template needs it to
	// render.
	// It's a good place to put default values for the fields
	data.Form = addVideoForm{}
	app.render(w, http.StatusOK, "add-video.tmpl.html", data)
}

type addVideoForm struct {
	Title               string   `form:"title"`
	URL                 string   `form:"url"`
	Rating              string   `form:"rating"`
	UploadDate          string   `form:"uploadDate"`
	Creator             string   `form:"creator"`
	Actors              []string `form:"actors"`
	validator.Validator `form:"-"`
}

func (app *application) videoAddPost(w http.ResponseWriter, r *http.Request) {
	var form addVideoForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	// 1) Validate video information
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.URL), "url", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Rating), "rating", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.UploadDate), "uploadDate", "This field cannot be blank")
	form.CheckField(validator.ValidDate(form.UploadDate), "uploadDate", "Date must be of the format YYYY-MM-DD")
	form.CheckField(validator.PermittedValue(strings.ToLower(form.Rating), "pg", "pg-13", "r"),
		"rating",
		"Rating must be PG, PG-13 or R (case insensitive)",
	)

	if !form.Valid() {
		// sending a new html form with errors if it's not valid
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add-video.tmpl.html", data)
		return
	}

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
