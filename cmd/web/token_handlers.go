package main

import (
	"net/http"
)

func (app *application) createAdminToken(w http.ResponseWriter, r *http.Request) {
	// pull the authenticated user from context the same way your other handlers do
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	plaintext, err := app.tokens.Insert(userID)
	if err != nil {
		app.serverError(r, w, err)
		return
	}

	data := envelope{
		"token":   plaintext,
		"message": "Store this somewhere safe — it will not be shown again.\n",
	}

	app.writeJSON(w, http.StatusOK, data, nil)
}
