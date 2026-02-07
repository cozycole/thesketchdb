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

func (app *application) viewSketchesAPI(w http.ResponseWriter, r *http.Request) {
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
	filterQuery := strings.Join(strings.Fields(query), " | ")

	personIds := extractUrlParamIDs(r.URL.Query()["person"])
	characterIds := extractUrlParamIDs(r.URL.Query()["character"])
	creatorIds := extractUrlParamIDs(r.URL.Query()["creator"])
	showIds := extractUrlParamIDs(r.URL.Query()["show"])
	tagIds := extractUrlParamIDs(r.URL.Query()["tag"])

	sketchList, err := app.services.Sketches.ListSketches(
		&models.Filter{
			Query:        filterQuery,
			CharacterIDs: characterIds,
			CreatorIDs:   creatorIds,
			PersonIDs:    personIds,
			ShowIDs:      showIds,
			TagIDs:       tagIds,
			SortBy:       sort,
			PageSize:     selectedPageSize,
			Page:         selectedPage,
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

	response := envelope{
		"filter_refs": filterRefs,
		"sketches":    sketchList.Sketches,
		"meta":        sketchList.Metadata,
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverError(r, w, err)
	}
}

func (app *application) adminGetSketchAPI(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("id param not defined"))
		return
	}

	sketch, err := app.services.Sketches.GetSketch(sketchId)
	if err != nil {
		if errors.Is(err, models.ErrNoSketch) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"sketch": sketch}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createSketchAPI(w http.ResponseWriter, r *http.Request) {
	var form sketchForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.validateSketchForm(&form)
	if !form.Valid() {
		app.failedValidationResponse(w, r, form.Validator.FieldErrors)
		return
	}

	formSketch := convertFormToSketch(&form)
	sketch, err := app.services.Sketches.CreateSketch(&formSketch, form.Thumbnail)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"sketch": sketch}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateSketchAPI(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("id param not defined"))
		return
	}

	var form sketchForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.validateSketchForm(&form)
	if !form.Valid() {
		app.failedValidationResponse(w, r, form.Validator.FieldErrors)
		return
	}

	sketch := convertFormToSketch(&form)

	sketch.ID = &sketchId
	file, _ := fileHeaderToBytes(form.Thumbnail)

	updatedSketch, err := app.services.Sketches.UpdateSketch(&sketch, file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"sketch": updatedSketch}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
