package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"sketchdb.cozycole.net/internal/models"
)

var pageSize = 16

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	featured, err := app.videos.GetById(1)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	featured2, err := app.videos.GetById(2)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	featured3, err := app.videos.GetById(3)
	if err != nil {
		app.serverError(r, w, err)
		return
	}
	featured4, err := app.videos.GetById(4)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	members, err := app.cast.GetCastMembers(1)
	if err != nil && err != models.ErrNoRecord {
		app.serverError(r, w, err)
		return
	}

	featured.Cast = members

	filter := models.Filter{
		Limit:  8,
		Offset: 0,
	}

	videos, err := app.videos.Get(&filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Videos = videos
	data.Featured = []*models.Video{featured, featured2, featured3, featured4}

	for _, f := range data.Featured {
		app.infoLog.Printf("%s\n", *f.Title)
	}

	app.render(r, w, http.StatusOK, "home.tmpl.html", "base", data)
}

func (app *application) browse(w http.ResponseWriter, r *http.Request) {
	browseSections := make(map[string][]*models.Video)
	limit := 8
	offset := 0

	// First add "custom" sections (ex: latest, trending, recommended/because you liked X)
	latest, err := app.videos.Get(
		&models.Filter{
			Limit:  limit,
			Offset: offset,
			SortBy: "latest"})
	if err != nil {
		app.errorLog.Println(err)
	}
	browseSections["Latest"] = latest

	kyleId := 1
	actorVideos, err := app.videos.Get(
		&models.Filter{
			Limit:  limit,
			Offset: offset,
			SortBy: "az",
			People: []*models.Person{
				&models.Person{ID: &kyleId},
			},
		},
	)
	if err != nil {
		app.errorLog.Println(err)
	}
	browseSections["Sketches Featuring Kyle Mooney"] = actorVideos

	data := app.newTemplateData(r)
	data.BrowseSections = browseSections
	app.render(r, w, http.StatusOK, "browse.tmpl.html", "base", data)
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
	query, _ := url.QueryUnescape(r.Form.Get("query"))
	filterQuery := strings.Join(strings.Fields(query), " | ")

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

	tagIdParams := r.URL.Query()["tag"]

	var tagIds []int
	for _, idStr := range tagIdParams {
		id, err := strconv.Atoi(idStr)
		if nil == err && id > 0 {
			tagIds = append(tagIds, id)
		}
	}

	var tagFilter []*models.Tag
	if len(tagIds) > 0 {
		tagFilter, err = app.tags.GetTags(&tagIds)
	}

	limit := app.settings.pageSize
	offset := (currentPage - 1) * limit
	filter := &models.Filter{
		Query:    filterQuery,
		Creators: creatorFilter,
		People:   peopleFilter,
		Tags:     tagFilter,
		SortBy:   sort,
		Limit:    limit,
		Offset:   offset,
	}

	results, err := app.getCatalogResults(currentPage, "video", filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)

	results.Filter = filter
	data.SearchResults = results
	data.SearchResults.Query = url.QueryEscape(query)
	data.Query = query

	url, err := buildURL("/catalog/sketches", results)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	w.Header().Add("HX-Push-Url", url)
	app.infoLog.Println(url)

	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	if isHxRequest && !isHistoryRestore {
		app.render(r, w, http.StatusOK, "catalog-result.tmpl.html", "catalog-result", data)
		return
	}

	app.render(r, w, http.StatusOK, "view-catalog.tmpl.html", "base", data)
}

func buildURL(baseURL string, result *SearchResult) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	params := result.Filter.Params()
	params.Add("page", strconv.Itoa(result.CurrentPage))
	params.Set("query", result.Query)

	u.RawQuery = params.Encode()

	return u.String(), nil
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}
