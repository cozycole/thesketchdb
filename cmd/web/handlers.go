package main

import (
	"net/http"
	"path"
	"time"

	"sketchdb.cozycole.net/internal/models"
)

var maxFileNameLength = 40

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

func (app *application) creatorAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = addCreatorForm{}
	app.render(w, http.StatusOK, "add-creator.tmpl.html", data)
}

func (app *application) creatorAddPost(w http.ResponseWriter, r *http.Request) {
	var form addCreatorForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	validateAddCreatorForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add-creator.tmpl.html", data)
		return
	}

	date, _ := time.Parse(time.DateOnly, form.EstablishedDate)
	imgName := models.CreateImageName(form.Name, maxFileNameLength)

	file, err := form.ProfileImage.Open()
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer file.Close()
	// get image extension
	buf := make([]byte, 512)

	_, err = file.Read(buf)
	if err != nil {
		app.serverError(w, err)
		return
	}

	mimeType := http.DetectContentType(buf)

	// the insert returns the fullImgName which is {fileName}-{id}.{ext}
	_, fullImgName, err := app.creators.
		Insert(
			form.Name, form.URL, imgName,
			mimeToExt[mimeType], date,
		)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.fileStorage.SaveMultipartFile(path.Join("creator", fullImgName), file)
	if err != nil {
		// TODO: We gotta remove the db record on this error
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/add/creator", http.StatusSeeOther)
}

func (app *application) actorAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = addCreatorForm{}
	app.render(w, http.StatusOK, "add-actor.tmpl.html", data)
}

func (app *application) actorAddPost(w http.ResponseWriter, r *http.Request) {
	var form addActorForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	validateAddActorForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add-actor.tmpl.html", data)
		return
	}

	date, _ := time.Parse(time.DateOnly, form.BirthDate)
	imgName := models.CreateImageName(form.First+" "+form.Last, maxFileNameLength)

	file, err := form.ProfileImage.Open()
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer file.Close()
	// get image extension
	buf := make([]byte, 512)

	_, err = file.Read(buf)
	if err != nil {
		app.serverError(w, err)
		return
	}

	mimeType := http.DetectContentType(buf)

	_, fullImgName, err := app.actors.
		Insert(
			form.First, form.Last, imgName,
			mimeToExt[mimeType], date,
		)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.fileStorage.SaveMultipartFile(path.Join("actor", fullImgName), file)
	if err != nil {
		// TODO: We gotta remove the db record on this error
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/add/actor", http.StatusSeeOther)
}

func (app *application) videoAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Need to initialize form data since the template needs it to
	// render. It's a good place to put default values for the fields
	data.Form = addVideoForm{}
	app.render(w, http.StatusOK, "add-video.tmpl.html", data)
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
		app.render(w, http.StatusUnprocessableEntity, "add-video.tmpl.html", data)
		return
	}

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

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}
