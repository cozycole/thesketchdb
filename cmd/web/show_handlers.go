package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/external/wikipedia"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) viewShowHome(w http.ResponseWriter, r *http.Request) {
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
		PageSize: 12,
		Page:     1,
		SortBy:   "popular",
		ShowIDs:  []int{*show.ID},
	}

	popular, _, err := app.sketches.Get(filter)
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
	pageData, err := views.ShowHomePageView(show, popular, cast, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = pageData
	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	if isHxRequest && !isHistoryRestore {
		app.render(r, w, http.StatusOK, "show-home.gohtml", "show-content", pageData)
		return
	}

	app.render(r, w, http.StatusOK, "show-home.gohtml", "base", data)
}

func (app *application) viewShowSketches(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id := r.PathValue("id")
	showId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
		app.errorLog.Print(err)
		return
	}

	page := r.Form.Get("page")
	currentPage, err := strconv.Atoi(page)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}

	sort := r.Form.Get("sort")
	if sort == "" {
		sort = "popular"
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
		PageSize: 12,
		Page:     currentPage,
		SortBy:   sort,
		ShowIDs:  []int{*show.ID},
	}

	results, err := app.services.Sketches.ListSketches(filter, true)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	pageData, err := views.ShowSketchesPageView(show, results, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = pageData

	// remove show id since we don't want it showing up in a show url
	filter.ShowIDs = nil
	url, err := views.BuildURL(
		fmt.Sprintf("/show/%d/%s/sketches", *show.ID, *show.Slug), currentPage, filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	w.Header().Add("HX-Push-Url", url)

	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	if isHxRequest && !isHistoryRestore {
		if r.Header.Get("HX-Target") == "showContent" {
			app.render(r, w, http.StatusOK, "show-sketches.gohtml", "show-content", pageData)

		} else {
			app.render(r, w, http.StatusOK, "sketches-result.gohtml", "sketches-result", pageData)
		}
		return
	}
	app.render(r, w, http.StatusOK, "show-sketches.gohtml", "base", data)
}

func (app *application) viewShowSeasons(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id := r.PathValue("id")
	showId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
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

	data := app.newTemplateData(r)
	pageData, err := views.ShowSeasonsPageView(show, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = pageData

	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	if isHxRequest && !isHistoryRestore {
		app.render(r, w, http.StatusOK, "show-seasons.gohtml", "show-content", pageData)
		return
	}
	app.render(r, w, http.StatusOK, "show-seasons.gohtml", "base", data)
}

func (app *application) viewShowGroupings(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	showId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
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

	groupings, err := app.shows.GetGroupings(showId)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	pageData, err := views.ShowExtrasPageView(show, groupings, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = pageData
	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	if isHxRequest && !isHistoryRestore {
		app.render(r, w, http.StatusOK, "show-extras.gohtml", "show-content", pageData)
		return
	}
	app.render(r, w, http.StatusOK, "show-extras.gohtml", "base", data)
}

func (app *application) viewShowCast(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	showId, err := strconv.Atoi(id)
	if err != nil {
		app.badRequest(w)
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

	cast, err := app.shows.GetShowCast(*show.ID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	pageData, err := views.ShowCastPageView(show, cast, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data.Page = pageData
	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	if isHxRequest && !isHistoryRestore {
		app.render(r, w, http.StatusOK, "show-cast.gohtml", "show-content", pageData)
		return
	}
	app.render(r, w, http.StatusOK, "show-cast.gohtml", "base", data)
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

	if show.WikiPage != nil {
		about, err := wikipedia.GetExtract(*show.WikiPage)
		if nil == err {
			show.About = &about
		}
	}

	id, err := app.shows.Insert(&show)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	show.ID = &id

	err = app.saveLargeProfile(*show.ProfileImg, "show", form.ProfileImg)
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
	form.ProfileImgUrl = fmt.Sprintf("%s/show/small/%s", app.baseImgUrl, safeDeref(show.ProfileImg))
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
	form.ProfileImgUrl = fmt.Sprintf("%s/show/small/%s", app.baseImgUrl, safeDeref(oldShow.ProfileImg))
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

	if newShow.WikiPage != nil {
		about, err := wikipedia.GetExtract(*newShow.WikiPage)
		app.infoLog.Printf(about)
		if nil == err {
			newShow.About = &about
		}
	}

	if form.ProfileImg != nil {
		var err error
		profileImg, err = generateThumbnailName(form.ProfileImg)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
		err = app.saveLargeProfile(profileImg, "show", form.ProfileImg)
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
	form.ProfileImgUrl = fmt.Sprintf("%s/show/small/%s", app.baseImgUrl, safeDeref(oldShow.ProfileImg))
	app.render(r, w, http.StatusOK, "show-form-page.gohtml", "show-form", form)
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
