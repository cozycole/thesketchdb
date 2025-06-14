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
		r.HandleFunc("/testing", app.testing)

		r.HandleFunc("/browse", app.browse)
		r.HandleFunc("/", app.home)

		r.Get("/search", app.search)
		r.Get("/catalog/sketches", app.catalogView)
		// r.Get("/catalog/people", app.peopleCatalog)
		// r.Get("/catalog/characters", app.catalogView)
		// r.Get("/catalog/creators", app.catalogView)
		// r.Get("/catalog/shows", app.catalogView)

		r.Get("/categories", app.categoriesView)

		r.Get("/sketch/{id}/{slug}", app.sketchView)

		r.Post("/sketch/like/{id}", app.sketchAddLike)
		r.Delete("/sketch/like/{id}", app.sketchRemoveLike)

		r.Get("/creator/{id}/{slug}", app.creatorView)
		r.Get("/creator/search", app.creatorSearch)

		r.Get("/person/{id}/{slug}", app.personView)
		r.Get("/person/search", app.personSearch)

		r.Get("/character/{id}/{slug}", app.characterView)
		r.Get("/character/search", app.characterSearch)

		r.Get("/show/{id}/{slug}", app.viewShow)
		r.Get("/show/{id}/{slug}/season", app.viewSeason)
		r.Get("/show/{id}/{slug}/season/{snum}", app.viewSeason)
		r.Get("/show/{id}/{slug}/season/{snum}/episode/{enum}", app.viewEpisode)

		r.Get("/category/search", app.categorySearch)
		r.Get("/tag/search", app.tagSearch)

		r.Get("/user/{username}", app.userView)

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
		admin := []string{"admin"}
		r.Get("/sketch/add", app.requireRoles(editorAdmin, app.sketchAddPage))
		r.Post("/sketch/add", app.requireRoles(editorAdmin, app.sketchAdd))
		r.Get("/sketch/{id}/update", app.requireRoles(editorAdmin, app.sketchUpdatePage))
		r.Patch("/sketch/{id}", app.requireRoles(editorAdmin, app.sketchUpdate))
		r.Put("/sketch/{id}/tag", app.requireRoles(editorAdmin, app.sketchUpdateTags))

		r.Get("/cast/add", app.requireRoles(editorAdmin, app.addCastPage))
		r.Post("/sketch/{id}/cast", app.requireRoles(editorAdmin, app.addCast))
		r.Patch("/sketch/{id}/cast/{castId}", app.requireRoles(editorAdmin, app.updateCast))
		r.Patch("/sketch/{id}/tag", app.requireRoles(editorAdmin, app.sketchAdd))

		r.Get("/show/add", app.requireRoles(editorAdmin, app.addShowPage))
		r.Post("/show/add", app.requireRoles(editorAdmin, app.addShow))
		r.Get("/show/{id}/update", app.requireRoles(editorAdmin, app.updateShowPage))
		r.Patch("/show/{id}/update", app.requireRoles(editorAdmin, app.updateShow))

		r.Post("/show/{id}/season/add", app.requireRoles(editorAdmin, app.addSeason))
		r.Delete("/show/{id}/season/add", app.requireRoles(admin, app.updateShow))

		r.Post("/episode/add", app.requireRoles(editorAdmin, app.addEpisode))
		r.Patch("/episode/{id}", app.requireRoles(editorAdmin, app.updateEpisode))
		r.Delete("/episode/{id}", app.requireRoles(admin, app.deleteEpisode))

		r.Get("/creator/add", app.requireRoles(editorAdmin, app.creatorAdd))
		r.Post("/creator/add", app.requireRoles(editorAdmin, app.creatorAddPost))

		r.Get("/person/add", app.requireRoles(editorAdmin, app.personAdd))
		r.Post("/person/add", app.requireRoles(editorAdmin, app.personAddPost))

		r.Get("/category/add", app.requireRoles(editorAdmin, app.categoryAddPage))
		r.Post("/category/add", app.requireRoles(editorAdmin, app.categoryAdd))

		r.Get("/tag/add", app.requireRoles(editorAdmin, app.tagAddPage))
		r.Post("/tag/add", app.requireRoles(editorAdmin, app.tagAdd))
	})

	return r
}
