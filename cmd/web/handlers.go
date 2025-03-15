package main

import (
	"net/http"
	"net/url"
	"strconv"

	"sketchdb.cozycole.net/internal/models"
)

var maxFileNameLength = 50
var pageSize = 16

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	filter := models.Filter{
		Limit:  8,
		Offset: 0,
	}

	videos, err := app.videos.Get(&filter)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Videos = videos

	app.render(w, http.StatusOK, "home.tmpl.html", "base", data)
}

func (app *application) browse(w http.ResponseWriter, r *http.Request) {
	browseSections := make(map[string][]*models.Video)
	limit := 8
	offset := 0

	// First add "custom" sections (ex: latest, trending, recommended/because you liked X)
	latest, err := app.videos.GetLatest(limit, offset)
	if err != nil {
		app.errorLog.Println(err)
	}
	browseSections["Latest"] = latest

	jamesId := 3
	actorVideos, err := app.videos.GetByPerson(jamesId, limit, offset)
	if err != nil {
		app.errorLog.Println(err)
	}
	browseSections["Sketches Featuring James Hartnett"] = actorVideos

	data := app.newTemplateData(r)
	data.BrowseSections = browseSections
	app.render(w, http.StatusOK, "browse.tmpl.html", "base", data)
}

func (app *application) catalogView(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	page := r.Form.Get("page")
	currentPage, err := strconv.Atoi(page)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}

	sort := r.Form.Get("sort")
	if sort == "" {
		sort = "latest"
	}

	personIdParams := r.URL.Query()["person"]

	var personIds []int
	for _, idStr := range personIdParams {
		id, err := strconv.Atoi(idStr)
		if nil == err && id > 0 {
			personIds = append(personIds, id)
		}
	}

	var peopleFilter []*models.Person
	if len(personIds) > 0 {
		peopleFilter, err = app.people.GetPeople(&personIds)
	}

	creatorIdParams := r.URL.Query()["creator"]

	var creatorIds []int
	for _, idStr := range creatorIdParams {
		id, err := strconv.Atoi(idStr)
		if nil == err && id > 0 {
			creatorIds = append(creatorIds, id)
		}
	}

	var creatorFilter []*models.Creator
	if len(creatorIds) > 0 {
		creatorFilter, err = app.creators.GetCreators(&creatorIds)
		app.errorLog.Println(err)
	}

	limit := app.settings.pageSize
	offset := (currentPage - 1) * limit
	filter := &models.Filter{
		SortBy:   sort,
		Creators: creatorFilter,
		People:   peopleFilter,
		Limit:    limit,
		Offset:   offset,
	}
	app.infoLog.Printf("FILTER: %+v\n", filter)

	results, err := app.getCatalogResults(currentPage, "video", filter)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)

	results.Filter = filter
	data.SearchResults = results

	url, err := buildURL("/catalog", results)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Add("HX-Push-Url", url)

	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	if isHxRequest && !isHistoryRestore {
		app.render(w, http.StatusOK, "catalog-result.tmpl.html", "catalog-result", data)
		return
	}

	app.render(w, http.StatusOK, "view-catalog.tmpl.html", "base", data)
}

func buildURL(baseURL string, result *SearchResult) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	params := result.Filter.Params()
	params.Add("page", strconv.Itoa(result.CurrentPage))

	u.RawQuery = params.Encode()

	return u.String(), nil
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}
