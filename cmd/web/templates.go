package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strconv"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/ui"
)

// Define a templateData type to act as the holding
// struct for any dynamic data we want to pass to the
// HTML templates. Since the ExecuteTemplate only accepts one
// struct for inserting data and data can come from many sources,
// you need to combine it all into one
type templateData struct {
	BrowseSections  map[string][]*models.Video
	CSRFToken       string
	Cast            *[]*models.CastMember
	CastMember      *models.CastMember
	Creator         *models.Creator
	CurrentYear     int
	DropdownResults dropdownSearchResults
	Flash           string
	Forms           Forms
	ImageBaseUrl    string
	IsAuthenticated bool
	IsEditor        bool
	Person          *models.Person
	SearchResults   *SearchResult
	Tags            *[]*models.Tag
	User            *models.User
	Video           *models.Video
	Videos          []*models.Video
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("Jan 2, 2006")
}

func formDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func getYear(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return strconv.Itoa(t.Year())
}

func printPersonName(a *models.Person) string {
	if a == nil || a.First == nil {
		return ""
	}

	name := *a.First
	if a.Last == nil {
		return name
	}
	return name + " " + *a.Last
}

func dict(values ...any) map[string]any {
	if len(values)%2 != 0 {
		panic("invalid dict call")
	}
	m := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			panic("dict keys must be strings")
		}
		m[key] = values[i+1]
	}
	return m
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Init global variable which maps string func names to
// functions to be used within templates (since you can call
// functions from template). NOTE: The tempalte functions should only
// return a single value
var functions = template.FuncMap{
	"humanDate":       humanDate,
	"getYear":         getYear,
	"dict":            dict,
	"derefString":     derefString,
	"formDate":        formDate,
	"printPersonName": printPersonName,
}

// Getting mapping of html page filename to template set for the page
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Add all pages and partials to the cache
	pages, err := fs.Glob(ui.Files, "html/**/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		// Register the funcMap before parsing the files
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
