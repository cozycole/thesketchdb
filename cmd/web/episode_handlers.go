package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cozycole/thesketchdb/cmd/web/views"
	"github.com/cozycole/thesketchdb/internal/models"
)

func (app *Application) viewEpisode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	episodeId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		return
	}

	episode, err := app.Shows.GetEpisode(episodeId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	page, err := views.EpisodePageView(episode, app.Config.BaseImgUrl)
	app.InfoLog.Printf("%+v\n", page)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = page
	app.render(r, w, http.StatusOK, "view-episode.gohtml", "base", data)
}

func (app *Application) addEpisodeForm(w http.ResponseWriter, r *http.Request) {
	seasonIdParam := r.PathValue("id")
	seasonId, err := strconv.Atoi(seasonIdParam)
	if err != nil {
		app.badRequest(w)
		app.ErrorLog.Print(err)
		return
	}

	season, _ := app.Shows.GetSeason(seasonId)
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

func (app *Application) addEpisode(w http.ResponseWriter, r *http.Request) {
	seasonIdParam := r.PathValue("id")
	seasonId, err := strconv.Atoi(seasonIdParam)
	if err != nil {
		app.badRequest(w)
		app.ErrorLog.Print(err)
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

	season, err := app.Shows.GetSeason(seasonId)
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

	id, err := app.Shows.InsertEpisode(&episode)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	episode.ID = &id

	err = app.saveLargeThumbnail(thumbnailName, "episode", form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		app.Shows.DeleteEpisode(id)
		return
	}

	newSeason, err := app.Shows.GetSeason(*episode.Season.ID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.EpisodeTableView(
		newSeason, app.Config.BaseImgUrl,
		fmt.Sprintf("/show/%d/%s", safeDeref(season.Show.ID), safeDeref(season.Show.Slug)),
	)
	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "episode-table", table)
}

func (app *Application) updateEpisodeForm(w http.ResponseWriter, r *http.Request) {
	epIdParam := r.PathValue("id")
	epId, err := strconv.Atoi(epIdParam)
	if err != nil {
		app.badRequest(w)
		app.ErrorLog.Print(err)
		return
	}

	episode, err := app.Shows.GetEpisode(epId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(r, w, err)
		}
		return
	}

	form := convertEpisodeToForm(episode)
	form.ThumbnailUrl = fmt.Sprintf("%s/episode/small/%s", app.Config.BaseImgUrl, form.ThumbnailName)
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

func (app *Application) updateEpisode(w http.ResponseWriter, r *http.Request) {
	epIdParam := r.PathValue("id")
	epId, err := strconv.Atoi(epIdParam)
	if err != nil {
		app.badRequest(w)
		app.ErrorLog.Print(err)
		return
	}

	var form episodeForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	oldEpisode, err := app.Shows.GetEpisode(epId)
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
		form.ThumbnailUrl = fmt.Sprintf("%s/episode/small/%s", app.Config.BaseImgUrl, safeDeref(oldEpisode.Thumbnail))
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
	app.InfoLog.Printf("%+v\n", episode)
	err = app.Shows.UpdateEpisode(&episode)
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

	season, err := app.Shows.GetSeason(safeDeref(episode.Season.ID))
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.EpisodeTableView(season, app.Config.BaseImgUrl,
		fmt.Sprintf("/show/%d/%s", safeDeref(season.Show.ID), safeDeref(season.Show.Slug)),
	)

	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "episode-table", table)
}

func (app *Application) deleteEpisode(w http.ResponseWriter, r *http.Request) {
	epIdParam := r.PathValue("id")
	epId, err := strconv.Atoi(epIdParam)
	if err != nil {
		app.badRequest(w)
		app.ErrorLog.Print(err)
		return
	}

	episode, err := app.Shows.GetEpisode(epId)
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

	err = app.Shows.DeleteEpisode(*episode.ID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	season, err := app.Shows.GetSeason(safeDeref(episode.Season.ID))
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	table := views.EpisodeTableView(season, app.Config.BaseImgUrl,
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
