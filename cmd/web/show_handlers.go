package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
		Limit:  12,
		SortBy: "popular",
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
	seasonId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		return
	}

	season, err := app.shows.GetSeason(seasonId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	show, err := app.shows.GetById(safeDeref(season.Show.ID))
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)

	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	if isHxRequest && !isHistoryRestore {
		format := r.URL.Query().Get("format")
		templateData := views.EpisodeGalleryView(season.Episodes, app.baseImgUrl, format, true)

		if isSeasonPath(r) {
			w.Header().Add("HX-Push-Url", fmt.Sprintf("/season/%d/%s", *season.ID, *season.Slug))
		}

		app.render(r, w, http.StatusOK, "episode-gallery.gohtml", "episode-gallery", templateData)
		return
	}

	page := views.SeasonPageView(show, season, app.baseImgUrl)
	data.Page = page
	app.render(r, w, http.StatusOK, "view-season.gohtml", "base", data)
}

func (app *application) viewEpisode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	episodeId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		return
	}

	episode, err := app.shows.GetEpisode(episodeId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	page, err := views.EpisodePageView(episode, app.baseImgUrl)
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
	ViewShowUrl     string
	ShowForm        showForm
	DisplaySeasons  bool
	SeasonDropdowns views.SeasonDropdowns
	SeasonForm      seasonForm
}

func (app *application) addShowPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Page = showFormPage{
		Title: "Add Show",
	}

	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "base", data)
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
		app.render(r, w, http.StatusUnprocessableEntity, "show-form-page.gohtml", "show-form", form)
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

	form := app.convertShowtoForm(show)
	form.ProfileImgUrl = fmt.Sprintf("%s/show/%s", app.baseImgUrl, safeDeref(show.ProfileImg))
	form.Action = fmt.Sprintf("/show/%d/update", showId)
	data := app.newTemplateData(r)
	data.Page = showFormPage{
		Title:           "Update Show",
		ShowID:          showId,
		ShowForm:        form,
		ViewShowUrl:     fmt.Sprintf("/show/%d/%s", showId, safeDeref(show.Slug)),
		DisplaySeasons:  true,
		SeasonDropdowns: views.SeasonDropdownsView(show, app.baseImgUrl),
		SeasonForm:      seasonForm{ShowID: showId},
	}
	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "base", data)
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
	form.Action = fmt.Sprintf("/show/%d/update", showId)
	form.ProfileImgUrl = fmt.Sprintf("%s/show/%s", app.baseImgUrl, safeDeref(oldShow.ProfileImg))
	if !form.Valid() {
		app.render(r, w, http.StatusUnprocessableEntity, "show-form.gohtml", "show-form", form)
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

	updatedShow, err := app.shows.GetById(showId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	form = app.convertShowtoForm(updatedShow)
	form.Action = fmt.Sprintf("/show/%d/update", showId)
	form.ProfileImgUrl = fmt.Sprintf("%s/show/%s", app.baseImgUrl, safeDeref(oldShow.ProfileImg))
	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "show-form", form)
}

func (app *application) addSeason(w http.ResponseWriter, r *http.Request) {

	var form seasonForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	app.validateSeasonForm(&form)
	if !form.Valid() {
		app.render(r, w, http.StatusUnprocessableEntity, "show-form-page.gohtml", "season-form", form)
		return
	}

	showIdParam := r.PathValue("id")
	showId, err := strconv.Atoi(showIdParam)
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

	slug := models.CreateSlugName(safeDeref(show.Name) + fmt.Sprintf(" s%d", form.Number))
	newSeason := &models.Season{
		Number: &form.Number,
		Slug:   &slug,
		Show: &models.Show{
			ID: show.ID,
		},
	}

	_, err = app.shows.AddSeason(newSeason)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	show, err = app.shows.GetById(showId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := views.SeasonDropdownsView(show, app.baseImgUrl)
	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "season-dropdowns", data)
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

	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "episode-form-modal", formModal)
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
		app.render(r, w, http.StatusUnprocessableEntity, "show-form-page.gohtml", "episode-form", form)
		return
	}

	episode := app.convertFormtoEpisode(&form)
	youtubeID, _ := extractYouTubeVideoID(safeDeref(episode.URL))
	if youtubeID != "" {
		episode.YoutubeID = &youtubeID
	}

	season, err := app.shows.GetSeason(seasonId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	slug := createEpisodeSlug(season, safeDeref(episode.Number))
	episode.Slug = &slug

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

	err = app.saveLargeThumbnail(thumbnailName, "episode", form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		app.shows.DeleteEpisode(id)
		return
	}

	newSeason, err := app.shows.GetSeason(*episode.Season.ID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.EpisodeTableView(
		newSeason, app.baseImgUrl,
		fmt.Sprintf("/show/%d/%s", safeDeref(season.Show.ID), safeDeref(season.Show.Slug)),
	)
	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "episode-table", table)
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
			safeDeref(episode.Season.Number),
			safeDeref(episode.Number),
		),
		Form: form,
	}

	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "episode-form-modal", formModal)
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
		app.render(r, w, http.StatusUnprocessableEntity, "show-form-page.gohtml", "episode-form", form)
		return
	}

	episode := app.convertFormtoEpisode(&form)
	youtubeID, _ := extractYouTubeVideoID(safeDeref(episode.URL))
	if youtubeID != "" {
		episode.YoutubeID = &youtubeID
	}

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

	season, err := app.shows.GetSeason(safeDeref(episode.Season.ID))
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.EpisodeTableView(season, app.baseImgUrl,
		fmt.Sprintf("/show/%d/%s", safeDeref(season.Show.ID), safeDeref(season.Show.Slug)),
	)

	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "episode-table", table)
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

	season, err := app.shows.GetSeason(safeDeref(episode.Season.ID))
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.EpisodeTableView(season, app.baseImgUrl,
		fmt.Sprintf("/show/%d/%s", safeDeref(season.Show.ID), safeDeref(season.Show.Slug)),
	)

	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "episode-table", table)
}

func getLatestSeasonNumber(show *models.Show) int {
	if show == nil {
		return 0
	}
	var max int
	for _, s := range show.Seasons {
		if safeDeref(s.Number) > max {
			max = safeDeref(s.Number)
		}
	}
	return max
}

func createEpisodeSlug(season *models.Season, epNumber int) string {
	var text string
	if season.Show != nil {
		text += safeDeref(season.Show.Name)
	}
	text += fmt.Sprintf(" s%d e%d", safeDeref(season.Number), epNumber)
	return models.CreateSlugName(text)
}

func isSeasonPath(r *http.Request) bool {
	rawURL := r.Header.Get("Hx-Current-URL")
	if rawURL == "" {
		return false
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return strings.HasPrefix(parsedURL.Path, "/season")
}
