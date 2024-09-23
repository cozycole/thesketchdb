package main

import (
	"net/http"

	// "sketchdb.cozycole.net/ui"
	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	dir := "./ui/static/"
	fs := http.FileServer(http.Dir(dir))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	router.HandleFunc("/", app.home)

	router.HandleFunc("/ping", ping)

	router.Get("/search", app.search)
	router.Get("/video/add", app.videoAdd)
	router.Post("/video/add", app.videoAddPost)
	router.Get("/video/view/{slug}", app.videoView)

	router.Get("/creator/add", app.creatorAdd)
	router.Post("/creator/add", app.creatorAddPost)
	router.Get("/actor/add", app.actorAdd)
	router.Post("/actor/add", app.actorAddPost)

	return router
}
