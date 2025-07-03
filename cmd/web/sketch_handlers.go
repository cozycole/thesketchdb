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

	app.render(r, w, http.StatusOK, "view-sketch.gohtml", "base", data)
}

type sketchFormPage struct {
	SketchID    int
	Title       string
	SketchForm  sketchForm
	CastSection castSection
	TagTable    views.TagTable
}

func (app *application) sketchAddPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Need to initialize form data since the template needs it to
	// render. It's a good place to put default values for the fields
	data.Page = sketchFormPage{
		Title: "Add Sketch",
		SketchForm: sketchForm{
			Action: "/sketch/add",
		},
	}

	app.render(r, w, http.StatusOK, "sketch-form-page.gohtml", "base", data)
}

func (app *application) sketchAdd(w http.ResponseWriter, r *http.Request) {
	var form sketchForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateSketchForm(&form)
	if !form.Valid() {
		form.Action = "/sketch/add"
		app.render(r, w, http.StatusUnprocessableEntity, "sketch-form-page.gohtml", "sketch-form", form)
		return
	}

	sketch := convertFormToSketch(&form)
	slug := models.CreateSlugName(form.Title)
	sketch.Slug = &slug

	thumbName, err := generateThumbnailName(form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	sketch.ThumbnailName = &thumbName

	youtubeID, _ := extractYouTubeVideoID(*sketch.URL)
	if youtubeID != "" {
		sketch.YoutubeID = &youtubeID
	}

	id, err := app.sketches.Insert(&sketch)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if sketch.Creator != nil {
		err = app.sketches.InsertSketchCreatorRelation(id, *sketch.Creator.ID)
		if err != nil {
			app.serverError(r, w, err)
			app.sketches.Delete(id)
			return
		}
	}

	err = app.saveLargeThumbnail(thumbName, "sketch", form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		app.sketches.Delete(id)
		return
	}

	w.Header().Add("Hx-Redirect", fmt.Sprintf("/sketch/%d/update", id))
}

type castSection struct {
	SketchID  int
	CastTable views.CastTable
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

	form := convertSketchToForm(sketch)
	form.ImageUrl = fmt.Sprintf("%s/sketch/%s", app.baseImgUrl, safeDeref(sketch.ThumbnailName))
	form.Action = fmt.Sprintf("/sketch/%d/update", safeDeref(sketch.ID))

	cast, err := app.cast.GetCastMembers(*sketch.ID)
	if err != nil && errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	castTable := views.CastTableView(cast, sketchId, app.baseImgUrl)

	tags, err := app.tags.GetBySketch(*sketch.ID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Page = sketchFormPage{
		SketchID:   sketchId,
		Title:      "Update Sketch",
		SketchForm: form,
		CastSection: castSection{
			SketchID:  sketchId,
			CastTable: castTable,
		},
		TagTable: views.TagTableView(tags, sketchId),
	}

	app.render(r, w, http.StatusOK, "sketch-form-page.gohtml", "base", data)
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

	app.validateSketchForm(&form)
	if !form.Valid() {
		form.ImageUrl = fmt.Sprintf("%s/sketch/%s", app.baseImgUrl, safeDeref(oldSketch.ThumbnailName))
		app.render(r, w, http.StatusUnprocessableEntity, "sketch-form-page.gohtml", "sketch-form", form)
		return
	}

	sketch := convertFormToSketch(&form)

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

		err = app.saveLargeThumbnail(thumbnailName, "sketch", form.Thumbnail)
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

	if sketch.Creator != nil && sketch.Creator.ID != nil {
		err = app.sketches.UpdateCreatorRelation(*sketch.ID, *sketch.Creator.ID)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	if form.Thumbnail != nil && oldSketch.ThumbnailName != nil {
		err = app.deleteImage("sketch", *oldSketch.ThumbnailName)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	updatedSketch, _ := app.sketches.GetById(sketchId)
	newForm := convertSketchToForm(updatedSketch)
	newForm.ImageUrl = fmt.Sprintf("%s/sketch/%s", app.baseImgUrl, safeDeref(updatedSketch.ThumbnailName))
	isHxRequest := r.Header.Get("HX-Request") == "true"
	if isHxRequest {
		app.render(r, w, http.StatusOK, "sketch-form-page.gohtml", "sketch-form", newForm)
		return
	}

	cast, err := app.cast.GetCastMembers(*sketch.ID)
	if err != nil && errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	castTable := views.CastTableView(cast, sketchId, app.baseImgUrl)
	data := app.newTemplateData(r)
	data.Page = sketchFormPage{
		SketchID:   sketchId,
		Title:      "Update Sketch",
		SketchForm: newForm,
		CastSection: castSection{
			SketchID:  sketchId,
			CastTable: castTable,
		},
	}
	app.render(r, w, http.StatusOK, "sketch-form-page.gohtml", "base", newForm)
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
	if !ok {
		tagStrIds = []string{}
	}
	tagNames, ok := r.Form["tagName[]"]
	if !ok {
		tagNames = []string{}
	}

	// need to remove any duplicates first
	tagSet := make(map[int]struct{})

	for _, strId := range tagStrIds {
		id, err := strconv.Atoi(strId)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		if id > 0 {
			tagSet[id] = struct{}{}
		}
	}

	// Convert to slice if needed
	tagIds := make([]int, 0, len(tagSet))
	for id := range tagSet {
		tagIds = append(tagIds, id)
	}

	form.TagIds = tagIds
	form.TagInputs = tagNames

	tags, err := app.convertFormtoSketchTags(&form)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.validateSketchTagsForm(&form)
	if !form.Valid() {
		tagTable := views.TagTableView(tags, sketchId)
		tagTable.Error = form.MultiFieldErrors["tagId"][0]
		app.render(r, w, http.StatusUnprocessableEntity, "tag-table.gohtml", "tag-table", tagTable)
		return
	}

	err = app.sketches.BatchUpdateTags(sketchId, &tags)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.render(r, w, http.StatusOK, "tag-table.gohtml", "tag-table", views.TagTableView(tags, sketchId))
}
