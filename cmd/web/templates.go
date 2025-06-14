package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"sketchdb.cozycole.net/internal/models"
	"sketchdb.cozycole.net/ui"
)

type flashMessage struct {
	Message string
	Level   string
}

// Define a templateData type to act as the holding
// struct for any dynamic data we want to pass to the
// HTML templates. Since the ExecuteTemplate only accepts one
// struct for inserting data and data can come from many sources,
// you need to combine it all into one
type templateData struct {
	BrowseSections  map[string][]*models.Sketch
	Categories      *[]*models.Category
	CSRFToken       string
	Cast            []*models.CastMember
	CastMember      *models.CastMember
	CatalogType     string
	Creator         *models.Creator
	CurrentYear     int
	DropdownResults dropdownSearchResults
	Episode         *models.Episode
	Episodes        []*models.Episode
	Featured        []*models.Sketch
	Flash           flashMessage
	Forms           Forms
	HtmxRequest     bool
	ImageBaseUrl    string
	IsAdmin         bool
	IsEditor        bool
	Person          *models.Person
	People          []*models.Person
	Season          *models.Season
	SectionType     string
	Show            *models.Show
	Tags            *[]*models.Tag
	ThumbnailType   string
	User            *models.User
	Sketch          *models.Sketch
	Sketches        []*models.Sketch
	Page            any
}

func formDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
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
	"dict":            dict,
	"derefString":     derefString,
	"formDate":        formDate,
	"printPersonName": printPersonName,
}

// Getting mapping of html page filename to template set for the page
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Add all pages and partials to the cache
	pages, err := fs.Glob(ui.Files, "html/**/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// fmt.Println(page)
		name := filepath.Base(page)

		patterns := []string{
			"html/base.gohtml",
			"html/partials/*.gohtml",
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
