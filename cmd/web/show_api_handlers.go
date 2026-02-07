package main

import (
	"net/http"
	"net/url"
	"strconv"

	"sketchdb.cozycole.net/internal/models"
)

func (app *application) listEpisodesAPI(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	page := r.Form.Get("page")
	selectedPage, err := strconv.Atoi(page)
	if err != nil || selectedPage < 1 {
		selectedPage = 1
	}

	pageSize := r.Form.Get("pageSize")
	selectedPageSize, err := strconv.Atoi(pageSize)
	if err != nil || selectedPageSize < 1 {
		selectedPageSize = 10
	}

	sort := r.Form.Get("sort")
	if sort == "" {
		sort = "popular"
	}

	query := r.Form.Get("query")
	if query == "" {
		query = r.Form.Get("q")
	}
	query, _ = url.QueryUnescape(query)

	app.infoLog.Printf("QUERY: '%s'", query)
	episodeList, err := app.services.Shows.ListEpisodes(
		&models.Filter{
			Query:    query,
			SortBy:   sort,
			PageSize: selectedPageSize,
			Page:     selectedPage,
		}, true)

	if err != nil {
		app.serverError(r, w, err)
		return
	}

	response := envelope{
		"episodes": episodeList.Episodes,
		"meta":     episodeList.Metadata,
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverError(r, w, err)
	}
}
