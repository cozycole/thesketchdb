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

		r.HandleFunc("/browse", app.browse)

		r.HandleFunc("/", app.home)
		r.Get("/video/{slug}", app.videoView)

		r.Post("/video/like/{id}", app.videoAddLike)
		r.Delete("/video/like/{id}", app.videoRemoveLike)

		r.Get("/category/{slug}", app.categoryViewPage)

		r.Get("/creator/{slug}", app.creatorView)
		r.Get("/creator/search", app.creatorSearch)

		r.Get("/person/{slug}", app.personView)
		r.Get("/person/search", app.personSearch)

		r.Get("/character/search", app.characterSearch)

		r.Get("/category/search", app.categorySearch)
		r.Get("/tag/search", app.tagSearch)

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
		r.Get("/video/add", app.requireRoles(editorAdmin, app.videoAddPage))
		r.Post("/video/add", app.requireRoles(editorAdmin, app.videoAdd))
		r.Get("/video/{id}/update", app.requireRoles(editorAdmin, app.videoUpdatePage))
		r.Patch("/video/{id}", app.requireRoles(editorAdmin, app.videoUpdate))
		r.Put("/video/{id}/tag", app.requireRoles(editorAdmin, app.videoUpdateTags))

		r.Get("/cast/add", app.requireRoles(editorAdmin, app.addCastPage))
		r.Post("/video/{id}/cast", app.requireRoles(editorAdmin, app.addCast))
		r.Patch("/video/{id}/cast/{castId}", app.requireRoles(editorAdmin, app.updateCast))
		r.Patch("/video/{id}/tag", app.requireRoles(editorAdmin, app.videoAdd))

		r.Get("/creator/add", app.requireRoles(editorAdmin, app.creatorAdd))
		r.Post("/creator/add", app.requireRoles(editorAdmin, app.creatorAddPost))

		r.Get("/person/add", app.requireRoles(editorAdmin, app.personAdd))
		r.Post("/person/add", app.requireRoles(editorAdmin, app.personAddPost))

		r.Get("/category/add", app.requireRoles(editorAdmin, app.categoryAddPage))
		r.Post("/category/add", app.requireRoles(editorAdmin, app.categoryAdd))

		r.Get("/tag/add", app.requireRoles(editorAdmin, app.tagAddPage))
		r.Post("/tag/add", app.requireRoles(editorAdmin, app.tagAdd))

		// admin := []string{"admin"}
	})

	return r
}
