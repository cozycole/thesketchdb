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
		app.serverError(r, w, err)
		return
	} else if show.ID == nil {
		app.notFound(w)
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

	app.errorLog.Printf("ID:%s SNUM:%s ENUM:%s ERR:%s ERR2:%s", id, snum, enum, err, err2)
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

func (app *application) addShowPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Show = &models.Show{}
	data.Forms.Show = &showForm{}
	app.render(r, w, http.StatusOK, "add-show.gohtml", "base", data)
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
		data := app.newTemplateData(r)
		data.Forms.Show = &form
		data.Show = &models.Show{}
		app.render(r, w, http.StatusUnprocessableEntity, "add-show.gohtml", "base", data)
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
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Show = show
	for _, season := range data.Show.Seasons {
		for _, ep := range season.Episodes {
			app.infoLog.Printf("Episode %d\n", *ep.Number)
			for _, vid := range ep.Sketches {
				app.infoLog.Printf("Sketch %d\n", *vid.ID)
			}
		}
	}
	data.Episode = &models.Episode{}
	data.Forms.Show = &showForm{}
	data.Forms.Episode = &episodeForm{}
	for _, ep := range data.Show.Seasons[0].Episodes {
		if ep.Title != nil {
			app.infoLog.Println(*ep.Title)
		}

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

	app.validateUpdateShowForm(&form)
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

	data := app.newTemplateData(r)
	data.Season = season
	data.Episode = &models.Episode{}
	data.Forms.Episode = &episodeForm{}
	app.render(r, w, http.StatusOK, "season-form.gohtml", "season-form", data)
}

func (app *application) addEpisode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var form episodeForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.infoLog.Printf("EP FORM: %+v\n", form)

	app.validateEpisodeForm(&form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Episode = &form
		data.Episode = &models.Episode{}
		app.render(r, w, http.StatusUnprocessableEntity, "episode-form.gohtml", "episode-form", data)
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

	app.saveThumbnail(thumbnailName, "episode", form.Thumbnail)

	season, err := app.shows.GetSeason(*episode.SeasonId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.infoLog.Printf("Season: %+v\n", season)
	data := app.newTemplateData(r)
	data.Season = season
	data.Episode = &models.Episode{}
	data.Forms.Episode = &episodeForm{}
	data.Flash = flashMessage{Level: "success", Message: "Episode added!"}
	app.render(r, w, http.StatusOK, "season-form.gohtml", "season-form", data)
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

	app.validateUpdateEpisodeForm(&form)
	app.infoLog.Printf("%+v\n", form)
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Forms.Episode = &form
		data.Episode = oldEpisode
		app.render(r, w, http.StatusUnprocessableEntity, "episode-form.gohtml", "episode-form", data)
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

		err = app.saveThumbnail(thumbnailName, "episode", form.Thumbnail)
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

	data := app.newTemplateData(r)
	data.Forms.Episode = &form
	data.Episode = &episode
	data.Flash = flashMessage{
		Level:   "success",
		Message: "Episode updated successfully!",
	}

	app.render(r, w, http.StatusOK, "episode-form.gohtml", "episode-form", data)
}

func (app *application) deleteEpisode(w http.ResponseWriter, r *http.Request) {
	epIdParam := r.PathValue("id")
	epId, err := strconv.Atoi(epIdParam)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	data := app.newTemplateData(r)
	episode, err := app.shows.GetEpisode(epId)
	if err != nil {
		var status int
		if errors.Is(err, models.ErrNoRecord) {
			data.Flash = flashMessage{
				Level:   "error",
				Message: "404 Episode does not exist",
			}
			status = http.StatusNotFound
		} else {
			data.Flash = flashMessage{
				Level:   "error",
				Message: "500 Internal Server Error",
			}
			status = http.StatusInternalServerError
		}
		app.render(r, w, status, "flash-message.gohtml", "flash-message", data)
		return
	}

	if len(episode.Sketches) != 0 {
		data.Flash = flashMessage{
			Level:   "error",
			Message: "400 Cannot delete episode with sketches",
		}
		app.render(r, w, http.StatusBadRequest, "flash-message.gohtml", "flash-message", data)
		return
	}

	err = app.shows.DeleteEpisode(*episode.ID)
	if err != nil {
		data.Flash = flashMessage{
			Level:   "error",
			Message: "500 Internal Server Error",
		}
		app.render(r, w, http.StatusInternalServerError, "flash-message.gohtml", "flash-message", data)
		return
	}

	data.Flash = flashMessage{
		Level:   "success",
		Message: "Episode deleted successfully!",
	}
	app.render(r, w, http.StatusOK, "flash-message.gohtml", "flash-message", data)
}
