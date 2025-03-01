package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"sketchdb.cozycole.net/internal/models"
)

func (app *application) search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	q, _ := url.QueryUnescape(r.Form.Get("q"))
	htmxReq := r.Header.Get("HX-Request")
	page := r.Form.Get("page")
	currentPage, err := strconv.Atoi(page)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}

	assetType := r.Form.Get("type")
	if assetType == "" {
		assetType = "video"
	}

	results, err := app.getSearchResults(q, currentPage, assetType)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.SearchResults = results
	app.infoLog.Printf("%+v", results)

	w.Header().Add("HX-Push-Url", fmt.Sprintf("/search?q=%s&type=%s&page=%d", url.QueryEscape(q), assetType, currentPage))

	if htmxReq != "" {
		app.render(w, http.StatusOK, "search-result.tmpl.html", "search-result", data)
		return
	}

	app.render(w, http.StatusOK, "search.tmpl.html", "base", data)
}

type dropdownSearchResults struct {
	Results      []result
	Redirect     string // e.g. /person/add
	RedirectText string // e.g. "Add Person +"
}

type result struct {
	ID   int
	Text string
	Img  string
}

func (app *application) personSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	redirLink := "/person/add"
	redirText := "Add Person +"
	results := dropdownSearchResults{
		Redirect:     redirLink,
		RedirectText: redirText,
	}

	if q != "" {
		q = strings.Replace(q, " ", "", -1)
		dbResults, err := app.people.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(w, err)
			}
			return
		}

		if dbResults != nil {
			res := []result{}
			for _, row := range dbResults {
				r := result{}
				r.Text = *row.First + " " + *row.Last
				r.ID = *row.ID
				res = append(res, r)
			}

			results.Results = res
		}
	}

	data := app.newTemplateData(r)
	data.DropdownResults = results

	w.Header().Add("Hx-Trigger-After-Swap", "insertDropdownItem")

	app.render(w, http.StatusOK, "dropdown.tmpl.html", "", data)
}

func (app *application) characterSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	redirLink := "/character/add"
	redirText := "Add Character +"
	results := dropdownSearchResults{
		Redirect:     redirLink,
		RedirectText: redirText,
	}

	if q != "" {
		q = strings.Replace(q, " ", "", -1)
		dbResults, err := app.characters.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(w, err)
			}
			return
		}

		if dbResults != nil {
			res := []result{}
			for _, row := range dbResults {
				r := result{}
				r.Text = *row.Name
				r.ID = *row.ID
				res = append(res, r)
			}

			results.Results = res
		}
	}
	w.Header().Add("Hx-Trigger-After-Swap", "insertDropdownItem")

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(w, http.StatusOK, "dropdown.tmpl.html", "", data)
}

func (app *application) creatorSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	redirLink := "/creator/add"
	redirText := "Add Creator +"
	results := dropdownSearchResults{
		Redirect:     redirLink,
		RedirectText: redirText,
	}

	if q != "" {
		q = strings.Replace(q, " ", "", -1)
		creators, err := app.creators.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(w, err)
			}
			return
		}

		if creators != nil {
			res := []result{}
			for _, c := range creators {
				r := result{}
				r.Text = *c.Name
				r.ID = *c.ID
				res = append(res, r)
			}

			results.Results = res
		}
	}
	w.Header().Add("Hx-Trigger-After-Swap", "insertDropdownItem")

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(w, http.StatusOK, "dropdown.tmpl.html", "", data)
}

func (app *application) categorySearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	redirLink := "/category/add"
	redirText := "Add Category +"
	results := dropdownSearchResults{
		Redirect:     redirLink,
		RedirectText: redirText,
	}

	if q != "" {
		q = strings.Replace(q, " ", "", -1)
		categories, err := app.categories.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(w, err)
			}
			return
		}

		if categories != nil {
			res := []result{}
			for _, c := range *categories {
				r := result{}
				r.Text = *c.Name
				r.ID = *c.ID
				res = append(res, r)
			}

			results.Results = res
		}
	}
	w.Header().Add("Hx-Trigger-After-Swap", "insertDropdownItem")

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(w, http.StatusOK, "dropdown.tmpl.html", "", data)
}

func (app *application) tagSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	redirLink := "/tag/add"
	redirText := "Add Tag +"
	results := dropdownSearchResults{
		Redirect:     redirLink,
		RedirectText: redirText,
	}

	if q != "" {
		q = strings.Replace(q, " ", "", -1)
		tags, err := app.tags.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(w, err)
				return
			}
		}

		if tags != nil {
			res := []result{}
			for _, t := range *tags {
				r := result{}
				var text string
				if t.Category.Name != nil {
					text += *t.Category.Name + " / "
				}
				text += *t.Name
				r.Text = text
				r.ID = *t.ID
				res = append(res, r)
			}

			results.Results = res
		}
	}
	w.Header().Add("Hx-Trigger-After-Swap", "insertDropdownItem")

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(w, http.StatusOK, "dropdown.tmpl.html", "", data)
}
