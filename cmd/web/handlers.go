package main

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

var pageSize = 16

func (app *application) testing(w http.ResponseWriter, r *http.Request) {
	filter := models.Filter{
		Limit:  8,
		Offset: 0,
		SortBy: "latest",
	}

	latest, err := app.videos.Get(&filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Videos = latest

	app.render(r, w, http.StatusOK, "carousel-testing.tmpl.html", "base", data)
}

type HomePage struct {
	Featured        []*views.SketchThumbnail
	LatestSketches  []*views.SketchThumbnail
	PopularSketches []*views.SketchThumbnail
	Actors          []*models.Person
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// this will get replaced by a playlist at some point

	featured, err := app.videos.GetFeatured()
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	featuredSketchViews, err := views.FeaturedSketchesView(featured, app.baseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	filter := models.Filter{
		Limit:  8,
		Offset: 0,
		SortBy: "latest",
	}

	latest, err := app.videos.Get(&filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	latestViews, err := views.SketchThumbnailsView(latest, app.baseImgUrl, "")

	popularFilter := models.Filter{
		Limit:  20,
		Offset: 0,
	}

	popularSketches, err := app.videos.Get(&popularFilter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	popularSketchViews, err := views.SketchThumbnailsView(popularSketches, app.baseImgUrl, "")

	people, err := app.people.GetPeople([]int{1, 2, 3, 4, 5, 6, 7, 8})
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	homePageData := HomePage{
		Featured:        featuredSketchViews,
		LatestSketches:  latestViews,
		PopularSketches: popularSketchViews,
		Actors:          people,
	}

	data := app.newTemplateData(r)
	data.Page = homePageData

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

type Catalog struct {
	CatalogType string
	Catalog     any
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

	personIds := extractUrlParamIDs(r.URL.Query()["person"])
	var peopleFilter []*models.Person
	if len(personIds) > 0 {
		peopleFilter, err = app.people.GetPeople(personIds)
	}

	characterIds := extractUrlParamIDs(r.URL.Query()["character"])
	var characterFilter []*models.Character
	if len(characterIds) > 0 {
		characterFilter, err = app.characters.GetCharacters(characterIds)
	}

	creatorIds := extractUrlParamIDs(r.URL.Query()["creator"])
	var creatorFilter []*models.Creator
	if len(creatorIds) > 0 {
		creatorFilter, err = app.creators.GetCreators(&creatorIds)
	}

	tagIds := extractUrlParamIDs(r.URL.Query()["tag"])
	var tagFilter []*models.Tag
	if len(tagIds) > 0 {
		tagFilter, err = app.tags.GetTags(&tagIds)
	}

	limit := app.settings.pageSize
	offset := (currentPage - 1) * limit
	filter := &models.Filter{
		Query:      filterQuery,
		Characters: characterFilter,
		Creators:   creatorFilter,
		People:     peopleFilter,
		Tags:       tagFilter,
		SortBy:     sort,
		Limit:      limit,
		Offset:     offset,
	}

	results, err := app.getSketchCatalogResults(currentPage, "video", filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	sketchCount, err := app.videos.GetCount(filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	results.Filter = filter
	results.TotalVideoCount = sketchCount
	results.Query = url.QueryEscape(query)

	data := app.newTemplateData(r)

	totalPages := int(math.Ceil(float64(sketchCount) / float64(limit)))

	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	sketchCatalog, err := views.SketchCatalogView(
		results,
		currentPage,
		totalPages,
		isHxRequest && !isHistoryRestore,
		app.baseImgUrl,
	)

	if err != nil {
		app.serverError(r, w, err)
		return
	}

	url, err := views.BuildURL("/catalog/sketches", currentPage, filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	w.Header().Add("HX-Push-Url", url)
	app.infoLog.Println(url)

	if isHxRequest && !isHistoryRestore {
		app.infoLog.Println("TARGET: ", r.Header.Get("HX-Target"))
		if r.Header.Get("HX-Target") == "catalogSection" {
			app.render(r, w, http.StatusOK, "video-catalog.tmpl.html", "video-catalog", sketchCatalog)
		} else {
			app.render(r, w, http.StatusOK, "video-catalog-result.tmpl.html", "video-catalog-result", sketchCatalog.CatalogResult)
		}
		return
	}

	data.Page = Catalog{
		CatalogType: "Sketches",
		Catalog:     sketchCatalog,
	}

	app.render(r, w, http.StatusOK, "view-catalog.tmpl.html", "base", data)
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}
