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
	videos, err := app.videos.GetAll(8)
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

	limit := app.settings.pageSize
	offset := (currentPage - 1) * limit
	filter := &models.Filter{
		People: peopleFilter,
		Limit:  limit,
		Offset: offset,
	}

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

	params := url.Values{}

	params.Add("page", strconv.Itoa(result.CurrentPage))

	for _, p := range result.Filter.People {
		params.Add("person", strconv.Itoa(*p.ID))
	}

	u.RawQuery = params.Encode()

	return u.String(), nil
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}
