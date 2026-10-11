package main

import (
	"net/http"

	"sketchdb.cozycole.net/internal/models"
)

const NO_EXPIRY = 0

func (app *application) createAPIKeyToken(w http.ResponseWriter, r *http.Request) {
	// pull the authenticated user from context the same way your other handlers do
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	token, err := app.tokens.New(userID, NO_EXPIRY, models.ScopeAPIKey)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := envelope{
		"token":   token.Plaintext,
		"message": "Store this somewhere safe — it will not be shown again.\n",
	}

	app.writeJSON(w, http.StatusOK, data, nil)
}
