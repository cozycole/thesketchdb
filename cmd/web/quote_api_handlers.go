package main

import (
	"fmt"
	"net/http"
	"strconv"
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
		"transcript": quoteData.TranscriptLines,
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverError(r, w, err)
	}
}
