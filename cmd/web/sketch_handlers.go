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

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if ok && user.ID != nil {
		hasLike, _ := app.sketches.HasLike(*sketch.ID, *user.ID)
		sketch.Liked = &hasLike
	}

	data := app.newTemplateData(r)
	sketchPage, err := views.SketchPageView(sketch, tags, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = sketchPage
	app.infoLog.Printf("%+v/n", sketchPage)

	app.render(r, w, http.StatusOK, "view-sketch.gohtml", "base", data)
}

func (app *application) sketchAddPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Need to initialize form data since the template needs it to
	// render. It's a good place to put default values for the fields
	data.Forms.Sketch = &sketchForm{}
	data.Sketch = &models.Sketch{}
	app.render(r, w, http.StatusOK, "add-sketch.gohtml", "base", data)
}

func (app *application) sketchAdd(w http.ResponseWriter, r *http.Request) {
	var form sketchForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateAddSketchForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Sketch = &form
		data.Sketch = &models.Sketch{}
		app.render(r, w, http.StatusUnprocessableEntity, "add-sketch.gohtml", "base", data)
		return
	}

	sketch, err := convertFormToSketch(&form)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	slug := models.CreateSlugName(form.Title)
	slug = slug + "-" + models.GetTimeStampHash()

	sketch.Slug = &slug

	id, err := app.sketches.Insert(&sketch)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	*sketch.ID = id

	if sketch.Creator.ID != nil {
		err = app.sketches.InsertSketchCreatorRelation(id, *sketch.Creator.ID)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	thumbnailName, err := generateThumbnailName(form.Thumbnail)
	err = app.sketches.InsertThumbnailName(id, thumbnailName)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	sketch.ThumbnailName = &thumbnailName
	err = app.saveThumbnail(*sketch.ThumbnailName, "sketch", form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		// TODO: delete sketch entry here
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/sketch/%s", sketch.Slug), http.StatusSeeOther)
}

func (app *application) sketchUpdatePage(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
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

	cast, err := app.cast.GetCastMembers(*sketch.ID)
	if err != nil && errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	tags, err := app.tags.GetBySketch(*sketch.ID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	sketch.Cast = cast
	sketch.Tags = &tags
	app.infoLog.Printf("%+v", sketch.Tags)
	data := app.newTemplateData(r)
	data.Sketch = sketch
	data.Forms.Sketch = &sketchForm{}
	data.Forms.Cast = &castForm{}
	// data.Forms.SketchTags = &sketchTagsForm{}

	// need to instantiate empty struct to load
	// castUpdate form on the page
	data.CastMember = &models.CastMember{}
	app.render(r, w, http.StatusOK, "update-sketch.gohtml", "base", data)
}

func (app *application) sketchUpdate(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	var form sketchForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	oldSketch, err := app.sketches.GetById(sketchId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	app.validateUpdateSketchForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Sketch = oldSketch
		data.Forms.Sketch = &form
		app.render(r, w, http.StatusUnprocessableEntity, "sketch-form.gohtml", "sketch-form", data)
		return
	}

	sketch, err := convertFormToSketch(&form)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	var thumbnailName string
	if oldSketch.ThumbnailName != nil {
		thumbnailName = *oldSketch.ThumbnailName
	} else {
		thumbnailName = ""
	}

	if form.Thumbnail != nil {
		var err error
		thumbnailName, err = generateThumbnailName(sketch.ThumbnailFile)
		if err != nil {
			app.serverError(r, w, err)
			return
		}

		err = app.saveThumbnail(thumbnailName, "sketch", form.Thumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	*sketch.ID = sketchId
	sketch.ThumbnailName = &thumbnailName
	err = app.sketches.Update(&sketch)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.sketches.UpdateCreatorRelation(*sketch.ID, *sketch.Creator.ID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	if form.Thumbnail != nil && oldSketch.ThumbnailName != nil {
		err = app.deleteImage("sketch", *oldSketch.ThumbnailName)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	data := app.newTemplateData(r)
	data.Sketch = &sketch
	data.Forms.Sketch = &form
	app.render(r, w, http.StatusOK, "sketch-form.gohtml", "sketch-form", data)
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

func (app *application) sketchUpdateTags(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	var form sketchTagsForm
	r.ParseForm()
	tagStrIds, ok := r.Form["tagId[]"]
	tagNames, ok1 := r.Form["tagName[]"]
	app.infoLog.Print(tagStrIds)
	if !(ok && ok1) {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var tagIds []int
	for _, strId := range tagStrIds {
		id, err := strconv.Atoi(strId)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		tagIds = append(tagIds, id)
	}

	form.TagIds = tagIds
	form.TagInputs = tagNames

	app.validateSketchTagsForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.SketchTags = &form
		data.Tags = &[]*models.Tag{}
		app.render(r, w, http.StatusUnprocessableEntity, "tag-table.gohtml", "tag-table", data)
		return
	}

	tags, err := app.convertFormtoSketchTags(&form)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.sketches.BatchUpdateTags(sketchId, &tags)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Tags = &tags
	app.render(r, w, http.StatusOK, "tag-table.gohtml", "tag-table", data)
}
