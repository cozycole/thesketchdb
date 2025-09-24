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

	moments, err := app.moments.GetBySketch(sketchId)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	sketchPage, err := views.SketchPageView(sketch, moments, tags, userSketchInfo, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = sketchPage

	app.render(r, w, http.StatusOK, "view-sketch.gohtml", "base", data)
}

type sketchFormPage struct {
	SketchID       int
	SketchUrl      string
	Title          string
	SketchForm     sketchForm
	CastSection    castSection
	TagTable       views.TagTable
	Moments        []UpdateMoment
	MomentForm     momentForm
	EmptyQuoteForm quoteForm
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
	// need to get these to create slug
	if sketch.Episode != nil {
		sketch.Episode, _ = app.shows.GetEpisode(safeDeref(sketch.Episode.ID))
	}
	if sketch.Creator != nil {
		sketch.Creator, _ = app.creators.GetById(safeDeref(sketch.Creator.ID))
	}

	slug := createSketchSlug(&sketch)
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

// REMEMBER TO UPDATE THIS FORM ON ADDITIONAL FIELDS
var emptyQuoteForm = quoteForm{
	QuoteID:        []int{0},
	CastMemberID:   []int{0},
	CastImageUrl:   []string{""},
	CastMemberName: []string{""},
	LineText:       []string{""},
	Funny:          []string{""},
	LineType:       []string{""},
	TagCount:       []int{0},
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
	form.ImageUrl = fmt.Sprintf("%s/sketch/small/%s", app.baseImgUrl, safeDeref(sketch.ThumbnailName))
	form.Action = fmt.Sprintf("/sketch/%d/update", safeDeref(sketch.ID))

	cast, err := app.cast.GetCastMembers(*sketch.ID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	castTable := views.CastTableView(cast, sketchId, app.baseImgUrl)

	tags, err := app.tags.GetBySketch(*sketch.ID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	moments, err := app.moments.GetBySketch(sketchId)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	// couple moment update form with their respective quote table form
	updateMoments := []UpdateMoment{}
	for _, m := range moments {
		mid := safeDeref(m.ID)
		momentForm := app.convertMomenttoForm(m)
		momentForm.Action = fmt.Sprintf("/moment/%d", mid)
		quoteForm := app.convertQuotestoForm(sketchId, mid, m.Quotes)
		updateMoments = append(updateMoments, UpdateMoment{mid, momentForm, quoteForm})
	}

	emptyQuoteForm.SketchID = sketchId
	data := app.newTemplateData(r)
	data.Page = sketchFormPage{
		SketchID:   sketchId,
		SketchUrl:  fmt.Sprintf("/sketch/%d/%s", sketchId, safeDeref(sketch.Slug)),
		Title:      "Update Sketch",
		SketchForm: form,
		CastSection: castSection{
			SketchID:  sketchId,
			CastTable: castTable,
		},
		Moments:        updateMoments,
		MomentForm:     momentForm{SketchID: sketchId, Action: "/moment/add"},
		EmptyQuoteForm: emptyQuoteForm, // for template
		TagTable:       views.TagTableView(tags, sketchId),
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
		form.ImageUrl = fmt.Sprintf("%s/sketch/small/%s", app.baseImgUrl, safeDeref(oldSketch.ThumbnailName))
		app.render(r, w, http.StatusUnprocessableEntity, "sketch-form-page.gohtml", "sketch-form", form)
		return
	}

	sketch := convertFormToSketch(&form)
	if sketch.Episode != nil {
		sketch.Episode, _ = app.shows.GetEpisode(safeDeref(sketch.Episode.ID))
	}
	if sketch.Creator != nil {
		sketch.Creator, _ = app.creators.GetById(safeDeref(sketch.Creator.ID))
	}

	slug := createSketchSlug(&sketch)
	sketch.Slug = &slug

	youtubeID, _ := extractYouTubeVideoID(*sketch.URL)
	if youtubeID != "" {
		sketch.YoutubeID = &youtubeID
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
	newForm.ImageUrl = fmt.Sprintf("%s/sketch/small/%s", app.baseImgUrl, safeDeref(updatedSketch.ThumbnailName))
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

func createSketchSlug(sketch *models.Sketch) string {
	var slugInput string
	if sketch.Episode != nil {
		episode := sketch.Episode
		showString := safeDeref(episode.Show.Name)
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
