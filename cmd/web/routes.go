package main

import (
	"net/http"

	"github.com/gorilla/mux"
	// "sketchdb.cozycole.net/ui"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	dir := "./ui/static/"
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	router.HandleFunc("/", app.home)
	router.HandleFunc("/snippet/view", app.snippetView)
	router.HandleFunc("/snippet/create", app.snippetCreate)
	router.HandleFunc("/search", app.searchPage)

	return router
}
