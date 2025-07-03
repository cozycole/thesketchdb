package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) viewShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	showId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	show, err := app.shows.GetById(showId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	filter := &models.Filter{
		Limit:  8,
		SortBy: "az",
		Shows:  []*models.Show{show},
	}

	popular, err := app.sketches.Get(filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	cast, err := app.shows.GetShowCast(*show.ID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	page, err := views.ShowPageView(show, popular, cast, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = page
	app.render(r, w, http.StatusOK, "view-show.gohtml", "base", data)
}

func (app *application) viewSeason(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	snum := r.PathValue("snum")
	if snum == "" {
		snum = "1"
	}

	showId, err := strconv.Atoi(id)
	seasonNumber, err2 := strconv.Atoi(snum)
	if err != nil || err2 != nil {
		app.badRequest(w)
		app.errorLog.Printf("ID:%s SNUM:%s ERR:%s ERR2:%s", id, snum, err, err2)
		return
	}

	show, err := app.shows.GetById(showId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	var season *models.Season
	for _, s := range show.Seasons {
		if s.Number == nil {
			continue
		}

		if *s.Number == seasonNumber {
			season = s
		}
	}

	if season == nil {
		app.notFound(w)
		return
	}

	data := app.newTemplateData(r)

	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	data.SectionType = "sub"
	if isHxRequest && !isHistoryRestore {
		format := r.URL.Query().Get("format")
		app.infoLog.Println("FORMAT", format)
		templateData := views.SeasonSelectGalleryView(show.Seasons, season, app.baseImgUrl, format)
		app.infoLog.Printf("%+v\n", templateData)
		app.render(r, w, http.StatusOK, "season-select-gallery.gohtml", "season-select-gallery", templateData)
		return
	}

	page := views.SeasonPageView(show, season, app.baseImgUrl)
	data.Page = page
	app.infoLog.Printf("%+v\n", page.SeasonSelectGallery.EpisodeGallery.EpisodeThumbnails[0])

	app.render(r, w, http.StatusOK, "view-season.gohtml", "base", data)
}

func (app *application) viewEpisode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	snum := r.PathValue("snum")
	enum := r.PathValue("enum")

	showId, err := strconv.Atoi(id)
	seasonNumber, err2 := strconv.Atoi(snum)
	episodeNumber, err3 := strconv.Atoi(enum)
	if err != nil || err2 != nil || err3 != nil {
		app.badRequest(w)
		app.errorLog.Printf("ID:%s SNUM:%s ENUM:%s ERR:%s ERR2:%s", id, snum, enum, err, err2)
		return
	}

	show, err := app.shows.GetById(showId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	var season *models.Season
	for _, s := range show.Seasons {
		if s.Number == nil {
			continue
		}

		if *s.Number == seasonNumber {
			season = s
		}
	}

	if season == nil {
		app.notFound(w)
		return
	}

	var episode *models.Episode
	for _, e := range season.Episodes {
		if e.Number == nil {
			continue
		}

		if *e.Number == episodeNumber {
			episode = e
		}
	}

	if episode == nil {
		app.notFound(w)
		return
	}

	episode, err = app.shows.GetEpisode(*episode.ID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	page, err := views.EpisodePageView(show, episode, app.baseImgUrl)
	app.infoLog.Printf("%+v\n", page)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = page
	app.render(r, w, http.StatusOK, "view-episode.gohtml", "base", data)
}

type showFormPage struct {
	Title           string
	ShowID          int
	ShowForm        showForm
	SeasonDropdowns views.SeasonDropdowns
}

func (app *application) addShowPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Page = showFormPage{
		Title: "Add Show",
	}

	app.render(r, w, http.StatusOK, "update-show.gohtml", "base", data)
}

func (app *application) addShow(w http.ResponseWriter, r *http.Request) {
	var form showForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateShowForm(&form)
	if !form.Valid() {
		form.Action = "/show/add"
		app.render(r, w, http.StatusUnprocessableEntity, "update-show.gohtml", "show-form", form)
		return
	}

	show := app.convertFormtoShow(&form)
	slug := models.CreateSlugName(*show.Name)
	show.Slug = &slug

	thumbName, err := generateThumbnailName(form.ProfileImg)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	show.ProfileImg = &thumbName

	id, err := app.shows.Insert(&show)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	show.ID = &id

	err = app.saveProfileImage(*show.ProfileImg, "show", form.ProfileImg)
	if err != nil {
		app.shows.Delete(&show)
		app.serverError(r, w, err)
		return
	}

	isHxRequest := r.Header.Get("HX-Request") == "true"
	if isHxRequest {
		w.Header().Add("Hx-Redirect", fmt.Sprintf("/show/%d/update", id))
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/show/%d/%s", *show.ID, *show.Slug), http.StatusSeeOther)
}

func (app *application) updateShowPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	showId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	show, err := app.shows.GetById(showId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	// for _, season := range show.Seasons {
	// 	for _, ep := range season.Episodes {
	// 		app.infoLog.Printf("Episode %d\n", *ep.Number)
	// 		for _, vid := range ep.Sketches {
	// 			app.infoLog.Printf("Sketch %d\n", *vid.ID)
	// 		}
	// 	}
	// }

	form := app.convertShowtoForm(show)
	form.ProfileImgUrl = fmt.Sprintf("%s/show/%s", app.baseImgUrl, safeDeref(show.ProfileImg))
	form.Action = fmt.Sprintf("/show/%d/update", showId)
	data := app.newTemplateData(r)
	data.Page = showFormPage{
		Title:           "Update Show",
		ShowID:          showId,
		ShowForm:        form,
		SeasonDropdowns: views.SeasonDropdownsView(show, app.baseImgUrl),
	}
	app.render(r, w, http.StatusOK, "update-show.gohtml", "base", data)
}

func (app *application) updateShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	showId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	var form showForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	oldShow, err := app.shows.GetById(showId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	app.validateShowForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Show = &form
		data.Show = oldShow
		app.render(r, w, http.StatusUnprocessableEntity, "show-form.gohtml", "show-form", data)
		return
	}

	newShow := app.convertFormtoShow(&form)
	newShow.ID = &showId

	var profileImg string
	if oldShow.ProfileImg != nil {
		profileImg = *oldShow.ProfileImg
	}

	if form.ProfileImg != nil {
		var err error
		profileImg, err = generateThumbnailName(form.ProfileImg)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
		err = app.saveProfileImage(profileImg, "show", form.ProfileImg)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	newShow.ProfileImg = &profileImg
	err = app.shows.Update(&newShow)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if form.ProfileImg != nil && oldShow.ProfileImg != nil {
		err = app.deleteImage("show", *oldShow.ProfileImg)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	data := app.newTemplateData(r)
	data.Forms.Show = &form
	data.Show = &newShow
	data.Flash = flashMessage{
		Level:   "success",
		Message: "Show successfully updated!",
	}
	app.render(r, w, http.StatusOK, "show-form.gohtml", "show-form", data)
}

func (app *application) addSeason(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	showIdParam := r.PathValue("id")
	showId, err := strconv.Atoi(showIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	seasonId, err := app.shows.AddSeason(showId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	season, err := app.shows.GetSeason(seasonId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	show, err := app.shows.GetById(showId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	showUrl := fmt.Sprintf("/show/%d/%s", safeDeref(show.ID), safeDeref(show.Slug))
	data := views.SeasonDropdownView(season, app.baseImgUrl, showUrl)
	app.render(r, w, http.StatusOK, "update-show.gohtml", "season-dropdown", data)
}

func (app *application) deleteSeason(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	seasonIdParam := r.PathValue("id")
	seasonId, err := strconv.Atoi(seasonIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	err = app.shows.DeleteSeason(seasonId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *application) addEpisodeForm(w http.ResponseWriter, r *http.Request) {
	seasonIdParam := r.PathValue("id")
	seasonId, err := strconv.Atoi(seasonIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	season, _ := app.shows.GetSeason(seasonId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	formModal := views.FormModal{
		Title: fmt.Sprintf("Add Episode to Season %d", safeDeref(season.Number)),
		Form: episodeForm{
			Action:   fmt.Sprintf("/season/%d/episode/add", seasonId),
			SeasonId: seasonId,
		},
	}

	app.render(r, w, http.StatusOK, "update-show.gohtml", "episode-form-modal", formModal)
}

func (app *application) addEpisode(w http.ResponseWriter, r *http.Request) {
	seasonIdParam := r.PathValue("id")
	seasonId, err := strconv.Atoi(seasonIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	r.ParseForm()
	var form episodeForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateEpisodeForm(&form)
	if !form.Valid() {
		form.Action = fmt.Sprintf("/season/%d/episode/add", seasonId)
		app.render(r, w, http.StatusUnprocessableEntity, "update-show.gohtml", "episode-form", form)
		return
	}

	episode := app.convertFormtoEpisode(&form)

	thumbnailName, err := generateThumbnailName(form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	episode.Thumbnail = &thumbnailName

	id, err := app.shows.InsertEpisode(&episode)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	episode.ID = &id

	app.saveLargeThumbnail(thumbnailName, "episode", form.Thumbnail)

	season, err := app.shows.GetSeason(*episode.SeasonId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.EpisodeTableView(
		season, app.baseImgUrl,
		fmt.Sprintf("/show/%d/%s", safeDeref(season.ShowId), safeDeref(season.ShowSlug)),
	)
	app.render(r, w, http.StatusOK, "update-show.gohtml", "episode-table", table)
}

func (app *application) updateEpisodeForm(w http.ResponseWriter, r *http.Request) {
	epIdParam := r.PathValue("id")
	epId, err := strconv.Atoi(epIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	episode, err := app.shows.GetEpisode(epId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	form := convertEpisodeToForm(episode)
	form.ThumbnailUrl = fmt.Sprintf("%s/episode/%s", app.baseImgUrl, form.ThumbnailName)
	form.Action = fmt.Sprintf("/episode/%d/update", epId)

	formModal := views.FormModal{
		Title: fmt.Sprintf("Update Season %d Episode %d",
			safeDeref(episode.SeasonNumber),
			safeDeref(episode.Number),
		),
		Form: form,
	}

	app.render(r, w, http.StatusOK, "update-show.gohtml", "episode-form-modal", formModal)
}

func (app *application) updateEpisode(w http.ResponseWriter, r *http.Request) {
	epIdParam := r.PathValue("id")
	epId, err := strconv.Atoi(epIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	var form episodeForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	oldEpisode, err := app.shows.GetEpisode(epId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	app.validateEpisodeForm(&form)
	if !form.Valid() {
		form.Action = fmt.Sprintf("/episode/%d/update", epId)
		form.ThumbnailUrl = fmt.Sprintf("%s/episode/%s", app.baseImgUrl, safeDeref(oldEpisode.Thumbnail))
		app.render(r, w, http.StatusUnprocessableEntity, "update-show.gohtml", "episode-form", form)
		return
	}

	episode := app.convertFormtoEpisode(&form)

	var thumbnailName string
	if oldEpisode.Thumbnail != nil {
		thumbnailName = *oldEpisode.Thumbnail
	} else {
		thumbnailName = ""
	}

	if form.Thumbnail != nil {
		var err error
		thumbnailName, err = generateThumbnailName(form.Thumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}

		err = app.saveLargeThumbnail(thumbnailName, "episode", form.Thumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	episode.ID = &epId
	episode.Thumbnail = &thumbnailName
	app.infoLog.Printf("%+v\n", episode)
	err = app.shows.UpdateEpisode(&episode)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	if form.Thumbnail != nil && oldEpisode.Thumbnail != nil {
		err = app.deleteImage("episode", *oldEpisode.Thumbnail)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	season, err := app.shows.GetSeason(safeDeref(episode.SeasonId))
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.EpisodeTableView(season, app.baseImgUrl,
		fmt.Sprintf("/show/%d/%s", safeDeref(season.ShowId), safeDeref(season.ShowSlug)),
	)

	app.render(r, w, http.StatusOK, "update-show.gohtml", "episode-table", table)
}

func (app *application) deleteEpisode(w http.ResponseWriter, r *http.Request) {
	epIdParam := r.PathValue("id")
	epId, err := strconv.Atoi(epIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	episode, err := app.shows.GetEpisode(epId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	if len(episode.Sketches) != 0 {
		app.clientError(w, http.StatusConflict)
		return
	}

	err = app.shows.DeleteEpisode(*episode.ID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	season, err := app.shows.GetSeason(safeDeref(episode.SeasonId))
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.EpisodeTableView(season, app.baseImgUrl,
		fmt.Sprintf("/show/%d/%s", safeDeref(season.ShowId), safeDeref(season.ShowSlug)),
	)

	app.render(r, w, http.StatusOK, "update-show.gohtml", "episode-table", table)
}
