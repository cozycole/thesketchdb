package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"sketchdb.cozycole.net/cmd/web/views"
	"sketchdb.cozycole.net/internal/models"
)

func (app *application) search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	query, _ := url.QueryUnescape(r.Form.Get("query"))
	fmt.Println("QUERY", query)

	var results *models.SearchResult
	var err error
	if query != "" {
		filterQuery := strings.Join(strings.Fields(query), " | ")

		filter := &models.Filter{
			Query:    filterQuery,
			Page:     1,
			PageSize: app.settings.maxSearchResults,
			SortBy:   "popular",
		}

		results, err = app.getSearchResults(filter)
		if err != nil {
			app.serverError(r, w, err)
			return
		}
	}

	data := app.newTemplateData(r)
	data.Page, err = views.SearchPageView(
		results,
		query,
		app.baseImgUrl,
		app.settings.maxSearchResults,
	)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	app.render(r, w, http.StatusOK, "search.gohtml", "base", data)
}

type dropdownSearchResults struct {
	Results      []result
	Redirect     string // e.g. /person/add
	RedirectText string // e.g. "Add Person +"
}

type result struct {
	ImageUrl string
	ID       int
	Text     string
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
				app.serverError(r, w, err)
			}
			return
		}

		if dbResults != nil {
			res := []result{}
			for _, row := range dbResults {
				r := result{}
				r.ImageUrl = fmt.Sprintf("%s/person/small/%s", app.baseImgUrl, safeDeref(row.ProfileImg))
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

	app.render(r, w, http.StatusOK, "dropdown.gohtml", "", data)
}

func (app *application) episodeSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	results := dropdownSearchResults{}
	f := &models.Filter{
		Page:     1,
		PageSize: 10,
		Query:    q,
	}

	if q != "" {
		episodes, _, err := app.shows.ListEpisodes(f)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(r, w, err)
			}
			return
		}

		if episodes != nil {
			res := []result{}
			for _, e := range episodes {
				r := result{}
				r.ImageUrl = fmt.Sprintf("%s/episode/small/%s", app.baseImgUrl, safeDeref(e.Thumbnail))
				r.Text = views.PrintEpisodeName(e)
				r.ID = safeDeref(e.ID)
				res = append(res, r)
			}

			results.Results = res
		}
	}

	w.Header().Add("Hx-Trigger-After-Swap", "insertDropdownItem")

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(r, w, http.StatusOK, "dropdown.gohtml", "", data)

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
		characters, err := app.characters.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(r, w, err)
			}
			return
		}

		if characters != nil {
			res := []result{}
			for _, c := range characters {
				r := result{}
				r.ImageUrl = fmt.Sprintf("%s/character/small/%s", app.baseImgUrl, safeDeref(c.Image))
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

	app.render(r, w, http.StatusOK, "dropdown.gohtml", "", data)
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
				app.serverError(r, w, err)
			}
			return
		}

		if creators != nil {
			res := []result{}
			for _, c := range creators {
				r := result{}
				r.ImageUrl = fmt.Sprintf("%s/creator/small/%s", app.baseImgUrl, safeDeref(c.ProfileImage))
				r.ID = *c.ID
				r.Text = *c.Name
				res = append(res, r)
			}

			results.Results = res
		}
	}
	w.Header().Add("Hx-Trigger-After-Swap", "insertDropdownItem")

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(r, w, http.StatusOK, "dropdown.gohtml", "", data)
}

func (app *application) showSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	redirLink := "/show/add"
	redirText := "Add Show +"
	results := dropdownSearchResults{
		Redirect:     redirLink,
		RedirectText: redirText,
	}

	if q != "" {
		q = strings.Replace(q, " ", "", -1)
		shows, err := app.shows.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(r, w, err)
			}
			return
		}

		if shows != nil {
			res := []result{}
			for _, s := range shows {
				r := result{}
				r.ImageUrl = fmt.Sprintf("%s/show/small/%s", app.baseImgUrl, safeDeref(s.ProfileImg))
				r.ID = *s.ID
				r.Text = *s.Name
				res = append(res, r)
			}

			results.Results = res
		}
	}
	w.Header().Add("Hx-Trigger-After-Swap", "insertDropdownItem")

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(r, w, http.StatusOK, "dropdown.gohtml", "", data)
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
				app.serverError(r, w, err)
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

	app.render(r, w, http.StatusOK, "dropdown.gohtml", "", data)
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
				app.serverError(r, w, err)
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

	app.render(r, w, http.StatusOK, "dropdown.gohtml", "", data)
}

func (app *application) seriesSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	redirLink := "/series/add"
	redirText := "Add Series +"
	results := dropdownSearchResults{
		Redirect:     redirLink,
		RedirectText: redirText,
	}

	if q != "" {
		q = strings.Replace(q, " ", "", -1)
		series, err := app.series.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(r, w, err)
				return
			}
		}

		if series != nil {
			res := []result{}
			for _, s := range series {
				r := result{}
				r.ID = safeDeref(s.ID)
				r.Text = safeDeref(s.Title)
				r.ImageUrl = fmt.Sprintf("%s/series/small/%s", app.baseImgUrl, safeDeref(s.ThumbnailName))
				res = append(res, r)
			}

			results.Results = res
		}
	}
	w.Header().Add("Hx-Trigger-After-Swap", "insertDropdownItem")

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(r, w, http.StatusOK, "dropdown.gohtml", "", data)
}

func (app *application) recurringSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	redirLink := "/recurring/add"
	redirText := "Add Recurring +"
	results := dropdownSearchResults{
		Redirect:     redirLink,
		RedirectText: redirText,
	}

	if q != "" {
		q = strings.Replace(q, " ", "", -1)
		recurring, err := app.recurring.Search(q)
		if err != nil {
			if !errors.Is(err, models.ErrNoRecord) {
				app.serverError(r, w, err)
				return
			}
		}

		if recurring != nil {
			res := []result{}
			for _, s := range recurring {
				r := result{}
				r.ID = safeDeref(s.ID)
				r.Text = safeDeref(s.Title)
				r.ImageUrl = fmt.Sprintf("%s/recurring/small/%s", app.baseImgUrl, safeDeref(s.ThumbnailName))
				res = append(res, r)
			}

			results.Results = res
		}
	}
	w.Header().Add("Hx-Trigger-After-Swap", "insertDropdownItem")

	data := app.newTemplateData(r)
	data.DropdownResults = results

	app.render(r, w, http.StatusOK, "dropdown.gohtml", "", data)
}
