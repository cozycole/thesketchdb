package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)


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
	form.ThumbnailUrl = fmt.Sprintf("%s/episode/small/%s", app.baseImgUrl, form.ThumbnailName)
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
		form.ThumbnailUrl = fmt.Sprintf("%s/episode/small/%s", app.baseImgUrl, safeDeref(oldEpisode.Thumbnail))
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

func createEpisodeSlug(season *models.Season, epNumber int) string {
	var text string
	if season.Show != nil {
		text += safeDeref(season.Show.Name)
	}
	text += fmt.Sprintf(" s%d e%d", safeDeref(season.Number), epNumber)
	return models.CreateSlugName(text)
}
