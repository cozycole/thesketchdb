package main

import (
	"fmt"
	"net/http"
	"strconv"

	"sketchdb.cozycole.net/internal/models"
)

func (app *application) adminGetQuotesAPI(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("id param not defined"))
		return
	}

	quoteData, err := app.services.Quotes.GetAdminQuotes(sketchId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"quotes":     quoteData.Quotes,
		"transcript": quoteData.TranscriptLines,
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverError(r, w, err)
	}
}

func (app *application) updateQuotesAPI(w http.ResponseWriter, r *http.Request) {
	sketchIdParam := r.PathValue("id")
	sketchId, err := strconv.Atoi(sketchIdParam)
	if err != nil {
		app.badRequest(w)
		return
	}

	type upsertQuoteDTO struct {
		ID          *int    `json:"id"`
		StartTimeMs *int    `json:"startTimeMs"`
		EndTimeMs   *int    `json:"endTimeMs"`
		Text        *string `json:"text"`
		CastIDs     []int   `json:"cast"`
		TagIDs      []int   `json:"tags"`
	}

	var input struct {
		Quotes    []upsertQuoteDTO `json:"upsert"`
		DeleteIds []int            `json:"delete"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	quotes := []*models.Quote{}
	for _, q := range input.Quotes {
		modelQ := models.Quote{
			ID:          q.ID,
			StartTimeMs: q.StartTimeMs,
			EndTimeMs:   q.EndTimeMs,
			Text:        q.Text,
		}
		cast := []*models.CastMember{}
		for _, id := range q.CastIDs {
			cast = append(cast, &models.CastMember{
				ID: &id,
			})
		}
		tags := []*models.Tag{}
		for _, id := range q.TagIDs {
			tags = append(tags, &models.Tag{
				ID: &id,
			})
		}

		modelQ.CastMembers = cast
		modelQ.Tags = tags
		quotes = append(quotes, &modelQ)
	}

	updatedQuotes, err := app.services.Quotes.UpdateQuotes(sketchId, quotes, input.DeleteIds)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"quotes": updatedQuotes,
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
