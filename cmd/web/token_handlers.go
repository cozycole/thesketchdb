package main

import (
	"net/http"

	"github.com/cozycole/thesketchdb/internal/models"
)

const NO_EXPIRY = 0

func (app *Application) createAPIKeyToken(w http.ResponseWriter, r *http.Request) {
	// pull the authenticated user from context the same way your other handlers do
	userID := app.SessionManager.GetInt(r.Context(), "authenticatedUserID")

	token, err := app.Tokens.New(userID, NO_EXPIRY, models.ScopeAPIKey)
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
