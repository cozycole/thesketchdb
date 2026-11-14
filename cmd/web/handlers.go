package main

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/cozycole/thesketchdb/cmd/web/views"
	"github.com/cozycole/thesketchdb/internal/models"
)

var pageSize = 16

type HomePage struct {
	Featured        []*views.SketchThumbnail
	RecentSketches  []*views.SketchThumbnail
	PopularSketches []*views.SketchThumbnail
	Actors          []*models.Person
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	featured, err := app.Sketches.GetFeatured()
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		app.serverError(r, w, err)
		return
	}

	featuredSketchViews, err := views.FeaturedSketchesView(featured, app.Config.BaseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	filter := models.Filter{
		Page:     1,
		PageSize: 8,
		SortBy:   "recent",
	}

	recentSketches, _, err := app.Sketches.Get(&filter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	recentSketchView, err := views.SketchThumbnailsView(recentSketches, app.Config.BaseImgUrl, "", true)

	popularFilter := models.Filter{
		SortBy:   "popular",
		Page:     1,
		PageSize: 20,
	}

	popularSketches, _, err := app.Sketches.Get(&popularFilter)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	popularSketchViews, err := views.SketchThumbnailsView(popularSketches, app.Config.BaseImgUrl, "", true)

	peopleIds := []int{61, 52, 42, 1, 2, 3, 39, 54, 56}
	// peopleIds := []int{5, 6, 7, 9, 2, 3, 1, 10, 4}
	people, err := app.People.GetPeople(peopleIds)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	idIndex := make(map[int]int)
	for i, id := range peopleIds {
		idIndex[id] = i
	}

	// Sort items based on the order slice
	sort.Slice(people, func(i, j int) bool {
		return idIndex[safeDeref(people[i].ID)] < idIndex[safeDeref(people[j].ID)]
	})

	homePageData := HomePage{
		Featured:        featuredSketchViews,
		RecentSketches:  recentSketchView,
		PopularSketches: popularSketchViews,
		Actors:          people,
	}

	data := app.newTemplateData(r)
	data.Page = homePageData

	app.render(r, w, http.StatusOK, "home.gohtml", "base", data)
}

func (app *Application) browse(w http.ResponseWriter, r *http.Request) {

	var sections []views.BrowseSectionDefinition
	for _, def := range views.BrowseSectionDefinitions {
		section := views.BrowseSectionDefinition{
			Title:  def.Title,
			Filter: def.Filter,
		}

		sketches, _, err := app.Sketches.Get(&section.Filter)
		if err != nil {
			app.serverError(r, w, err)
			return
		}

		section.Sketches = sketches

		sections = append(sections, section)
	}

	browsePage, err := views.BrowsePageView(sections, app.Config.BaseImgUrl)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Page = browsePage

	app.InfoLog.Printf("%+v\n", browsePage)

	app.render(r, w, http.StatusOK, "browse.gohtml", "base", data)
}

type Catalog struct {
	CatalogType string
	Catalog     any
}

func (app *Application) catalogView(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	page := r.Form.Get("page")
	currentPage, err := strconv.Atoi(page)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}

	sort := r.Form.Get("sort")
	if sort == "" {
		sort = "popular"
	}

	query, _ := url.QueryUnescape(r.Form.Get("query"))
	filterQuery := strings.Join(strings.Fields(query), " | ")

	personIds := extractUrlParamIDs(r.URL.Query()["person"])
	characterIds := extractUrlParamIDs(r.URL.Query()["character"])
	creatorIds := extractUrlParamIDs(r.URL.Query()["creator"])
	showIds := extractUrlParamIDs(r.URL.Query()["show"])
	tagIds := extractUrlParamIDs(r.URL.Query()["tag"])

	filter := &models.Filter{
		Query:        filterQuery,
		CharacterIDs: characterIds,
		CreatorIDs:   creatorIds,
		PersonIDs:    personIds,
		ShowIDs:      showIds,
		TagIDs:       tagIds,
		SortBy:       sort,
		Page:         currentPage,
		PageSize:     app.Config.PageSize,
	}

	results, err := app.Services.Sketches.ListSketches(filter, true)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := app.newTemplateData(r)

	isHxRequest := r.Header.Get("HX-Request") == "true"
	isHistoryRestore := r.Header.Get("HX-History-Restore-Request") == "true"
	sketchCatalog, err := views.SketchCatalogView(
		&results,
		isHxRequest && !isHistoryRestore,
		app.Config.BaseImgUrl,
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

	if isHxRequest && !isHistoryRestore {
		app.InfoLog.Println("TARGET: ", r.Header.Get("HX-Target"))
		if r.Header.Get("HX-Target") == "catalogSection" {
			app.render(r, w, http.StatusOK, "sketch-catalog.gohtml", "sketch-catalog", sketchCatalog)
		} else {
			app.render(r, w, http.StatusOK, "sketch-catalog-result.gohtml", "sketch-catalog-result", sketchCatalog.CatalogResult)
		}
		return
	}

	data.Page = Catalog{
		CatalogType: "Sketches",
		Catalog:     sketchCatalog,
	}

	app.render(r, w, http.StatusOK, "view-catalog.gohtml", "base", data)
}

func (app *Application) serveCMS(w http.ResponseWriter, r *http.Request) {
	cmsDistPath := "./cmsdist"
	if app.Config.DevEnv {
		cmsDistPath = "./ui/static/js/dist/cms"
	}

	p := filepath.Join(cmsDistPath, strings.TrimPrefix(r.URL.Path, "/admin"))

	if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
		http.ServeFile(w, r, p)
		return
	}

	// SPA router fallback
	http.ServeFile(w, r, filepath.Join(cmsDistPath, "index.html"))
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}
