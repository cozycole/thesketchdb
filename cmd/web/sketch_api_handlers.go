package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"sketchdb.cozycole.net/internal/models"
)

func (app *application) viewSketchesAPI(w http.ResponseWriter, r *http.Request) {
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

	limit := app.settings.pageSize
	offset := (currentPage - 1) * limit

	sketchList, err := app.services.Sketches.ListSketches(
		&models.Filter{
			Query:        filterQuery,
			CharacterIDs: characterIds,
			CreatorIDs:   creatorIds,
			PersonIDs:    personIds,
			ShowIDs:      showIds,
			TagIDs:       tagIds,
			SortBy:       sort,
			Limit:        limit,
			Offset:       offset,
		}, true)

	if err != nil {
		app.serverError(r, w, err)
		return
	}

	filterRefs := map[string]any{}
	if len(sketchList.CreatorRefs) > 0 {
		filterRefs["creators"] = sketchList.CreatorRefs
	}

	if len(sketchList.CharacterRefs) > 0 {
		filterRefs["characters"] = sketchList.CharacterRefs
	}

	if len(sketchList.PersonRefs) > 0 {
		filterRefs["people"] = sketchList.PersonRefs
	}

	if len(sketchList.ShowRefs) > 0 {
		filterRefs["shows"] = sketchList.ShowRefs
	}

	if len(sketchList.TagRefs) > 0 {
		filterRefs["tags"] = sketchList.TagRefs
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"total": sketchList.TotalCount, "filter_refs": filterRefs, "sketches": sketchList.Sketches}, nil)
	if err != nil {
		app.serverError(r, w, err)
	}
}

func (app *application) createSketchAPI(w http.ResponseWriter, r *http.Request) {
	var form sketchForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Print(err)
		return
	}

	app.validateSketchForm(&form)
	if !form.Valid() {
		form.Action = "/sketch/add"
		app.render(r, w, http.StatusUnprocessableEntity, "sketch-form-page.gohtml", "sketch-form", form)
		return
	}

	formSketch := convertFormToSketch(&form)
	sketch, err := app.services.Sketches.CreateSketch(&formSketch, form.Thumbnail)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"sketch": sketch}, nil)
	if err != nil {
		app.serverError(r, w, err)
	}
}
