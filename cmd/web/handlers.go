package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/internal/utils"
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

	data := app.newTemplateData(r)
	data.Video = video

	app.render(w, http.StatusOK, "view-video.tmpl.html", data)
}

func (app *application) creatorView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	creator, err := app.creators.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	videos, err := app.videos.GetByCreator(creator.ID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Creator = creator
	data.Videos = videos

	app.render(w, http.StatusOK, "view-creator.tmpl.html", data)
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

	app.validateAddCreatorForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add-creator.tmpl.html", data)
		return
	}

	date, _ := time.Parse(time.DateOnly, form.EstablishedDate)
	imgName := models.CreateSlugName(form.Name, maxFileNameLength)

	file, err := form.ProfileImage.Open()
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer file.Close()

	mimeType, err := utils.GetMultipartFileMime(file)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// the insert returns the fullImgName which is {fileName}-{id}.{ext}
	_, slug, fullImgName, err := app.creators.
		Insert(
			form.Name, form.URL, imgName,
			mimeToExt[mimeType], date,
		)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.fileStorage.SaveFile(path.Join("creator", fullImgName), file)
	if err != nil {
		// TODO: We gotta remove the db record on this error
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/creator/%s", slug), http.StatusSeeOther)
}

func (app *application) personAdd(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = addCreatorForm{}
	app.render(w, http.StatusOK, "add-person.tmpl.html", data)
}

func (app *application) personAddPost(w http.ResponseWriter, r *http.Request) {
	var form addPersonForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateAddPersonForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "add-person.tmpl.html", data)
		return
	}

	date, _ := time.Parse(time.DateOnly, form.BirthDate)
	imgName := models.CreateSlugName(form.First+" "+form.Last, maxFileNameLength)

	file, err := form.ProfileImage.Open()
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer file.Close()

	mimeType, err := utils.GetMultipartFileMime(file)
	if err != nil {
		app.serverError(w, err)
		return
	}

	_, fullImgName, err := app.people.
		Insert(
			form.First, form.Last, imgName,
			mimeToExt[mimeType], date,
		)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.fileStorage.SaveFile(path.Join("person", fullImgName), file)
	if err != nil {
		// TODO: We gotta remove the db record on this error
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/person/add", http.StatusSeeOther)
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

	date, _ := time.Parse(time.DateOnly, form.UploadDate)
	imgName := models.CreateSlugName(form.Title, maxFileNameLength)

	file, err := form.Thumbnail.Open()
	if err != nil {
		app.serverError(w, err)
		return
	}
	defer file.Close()

	mimeType, err := utils.GetMultipartFileMime(file)
	if err != nil {
		app.serverError(w, err)
		return
	}

	vidID, slug, thumbnailName, err := app.videos.Insert(
		form.Title, form.URL, strings.ToUpper(form.Rating),
		imgName, mimeToExt[mimeType], date,
	)
	if err != nil {
		app.serverError(w, err)
		return
	}

	width, height, err := utils.GetImageDimensions(file)
	if err != nil {
		// TODO: gotta remove the db record on this error and any after
		app.serverError(w, err)
		return
	}

	var dstFile io.Reader
	dstFile = file
	// This is stock youtube thumbnail dimensions but has black
	// top/bottom borders that need to be removed
	if width == 480 && height == 360 {
		rect := image.Rect(0, 45, 480, 315)
		dstFile, err = utils.CropImg(file, mimeToExt[mimeType], rect)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	err = app.fileStorage.SaveFile(path.Join("video", thumbnailName), dstFile)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.videos.InsertVideoCreatorRelation(vidID, form.CreatorID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, id := range form.PersonIDs {
		err = app.videos.InsertVideoPersonRelation(vidID, id)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/video/%s", slug), http.StatusSeeOther)
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}
