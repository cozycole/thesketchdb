package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes(staticRoute, imageStorageRoot, imageUrl string) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		fs := http.FileServer(http.Dir(staticRoute))
		r.Handle("/static/*", http.StripPrefix("/static/", fs))

		app.infoLog.Printf("Starting image file server rooted at %s\n", imageStorageRoot)
		app.infoLog.Printf("Image Url: %s\n", imageUrl)

		imgFs := http.FileServer(http.Dir(imageStorageRoot))
		r.Handle("/images/*", http.StripPrefix(imageUrl, imgFs))
	})

	// public routes
	r.Group(func(r chi.Router) {
		r.Use(
			app.sessionManager.LoadAndSave,
			app.logRequest,
			app.authenticate,
		)

		r.HandleFunc("/ping", ping)

		r.HandleFunc("/", app.home)
		r.Get("/video/{slug}", app.videoView)
		r.Post("/video/like/{videoId}", app.videoAddLike)
		r.Delete("/video/like/{videoId}", app.videoRemoveLike)

		r.Get("/creator/{slug}", app.creatorView)

		r.Get("/person/{slug}", app.personView)
		r.Get("/person/search", app.personSearch)

		r.Get("/character/search", app.characterSearch)

		r.Get("/user/{username}", app.userView)

		r.Get("/search", app.search)

		r.Get("/signup", app.userSignup)
		r.Post("/signup", app.userSignupPost)
		r.Get("/login", app.userLogin)
		r.Post("/login", app.userLoginPost)
		r.Post("/logout", app.userLogoutPost)
	})

	// role routes
	r.Group(func(r chi.Router) {
		r.Use(
			app.sessionManager.LoadAndSave,
			app.logRequest,
			app.authenticate,
		)

		editorAdmin := []string{"editor", "admin"}
		// admin := []string{"admin"}
		r.Get("/video/add", app.requireRoles(editorAdmin, app.videoAdd))
		r.Post("/video/add", app.requireRoles(editorAdmin, app.videoAddPost))

		r.Get("/creator/add", app.requireRoles(editorAdmin, app.creatorAdd))
		r.Post("/creator/add", app.requireRoles(editorAdmin, app.creatorAddPost))

		r.Get("/person/add", app.requireRoles(editorAdmin, app.personAdd))
		r.Post("/person/add", app.requireRoles(editorAdmin, app.personAddPost))
	})

	return r
}
