package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes(staticRoute, imageStorageRoot string) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		fs := http.FileServer(http.Dir(staticRoute))
		r.Handle("/static/*", http.StripPrefix("/static/", fs))

		if app.fileStorage.Type() == "Local" {
			app.infoLog.Printf("Starting image file server rooted at %s\n", imageStorageRoot)
			imgFs := http.FileServer(http.Dir(imageStorageRoot))
			r.Handle("/images/*", http.StripPrefix(app.baseImgUrl, imgFs))
			app.infoLog.Printf("Image Url: %s\n", app.baseImgUrl)
		}
	})

	// public routes
	r.Group(func(r chi.Router) {
		r.Use(
			app.recoverPanic,
			app.sessionManager.LoadAndSave,
			app.logRequest,
			app.authenticate,
		)

		r.HandleFunc("/ping", ping)

		r.HandleFunc("/", app.home)
		r.HandleFunc("/browse", app.browse)
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

		r.Get("/person/{id}/{slug}", app.viewPerson)
		r.Get("/person/search", app.personSearch)

		r.Get("/character/{id}/{slug}", app.characterView)
		r.Get("/character/search", app.characterSearch)

		r.Get("/series/{id}/{slug}", app.seriesView)
		r.Get("/series/search", app.seriesSearch)

		r.Get("/recurring/{id}/{slug}", app.recurringView)
		r.Get("/recurring/search", app.recurringSearch)

		r.Get("/show/{id}/{slug}", app.viewShow)
		r.Get("/show/search", app.showSearch)
		r.Get("/show/{id}/{slug}/season", app.viewSeason)
		r.Get("/season/{id}/{slug}", app.viewSeason)

		r.Get("/episode/{id}/{slug}", app.viewEpisode)
		r.Get("/episode/search", app.episodeSearch)

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
			app.recoverPanic,
			app.sessionManager.LoadAndSave,
			app.logRequest,
			app.authenticate,
			app.requireAuthentication,
		)

		editorAdmin := []string{"editor", "admin"}
		admin := []string{"admin"}
		r.Get("/sketch/add", app.requireRoles(editorAdmin, app.sketchAddPage))
		r.Post("/sketch/add", app.requireRoles(editorAdmin, app.sketchAdd))
		r.Get("/sketch/{id}/update", app.requireRoles(editorAdmin, app.sketchUpdatePage))
		r.Post("/sketch/{id}/update", app.requireRoles(editorAdmin, app.sketchUpdate))
		r.Post("/sketch/{id}/tag", app.requireRoles(editorAdmin, app.sketchUpdateTags))

		r.Get("/sketch/{id}/cast", app.requireRoles(editorAdmin, app.addCastPage))
		r.Post("/sketch/{id}/cast", app.requireRoles(editorAdmin, app.addCast))
		r.Get("/cast/{id}/update", app.updateCastPage)
		r.Post("/cast/{castId}/update", app.requireRoles(editorAdmin, app.updateCast))
		r.Delete("/cast/{castId}", app.requireRoles(editorAdmin, app.deleteCast))
		r.Patch("/sketch/{id}/cast/order", app.requireRoles(editorAdmin, app.orderCast))

		r.Get("/show/add", app.requireRoles(editorAdmin, app.addShowPage))
		r.Post("/show/add", app.requireRoles(editorAdmin, app.addShow))
		r.Get("/show/{id}/update", app.requireRoles(editorAdmin, app.updateShowPage))
		r.Post("/show/{id}/update", app.requireRoles(editorAdmin, app.updateShow))

		r.Post("/show/{id}/season/add", app.requireRoles(editorAdmin, app.addSeason))
		r.Delete("/season/{id}", app.requireRoles(admin, app.deleteSeason))

		r.Get("/season/{id}/episode/add", app.requireRoles(editorAdmin, app.addEpisodeForm))
		r.Post("/season/{id}/episode/add", app.requireRoles(editorAdmin, app.addEpisode))
		r.Get("/episode/{id}/update", app.requireRoles(editorAdmin, app.updateEpisodeForm))
		r.Post("/episode/{id}/update", app.requireRoles(editorAdmin, app.updateEpisode))
		r.Delete("/episode/{id}", app.requireRoles(admin, app.deleteEpisode))

		r.Get("/creator/add", app.requireRoles(editorAdmin, app.addCreatorPage))
		r.Post("/creator/add", app.requireRoles(editorAdmin, app.addCreator))
		r.Get("/creator/{id}/update", app.requireRoles(editorAdmin, app.updateCreatorPage))
		r.Post("/creator/{id}/update", app.requireRoles(editorAdmin, app.updateCreator))

		r.Get("/person/add", app.requireRoles(editorAdmin, app.addPersonPage))
		r.Post("/person/add", app.requireRoles(editorAdmin, app.addPerson))
		r.Get("/person/{id}/update", app.requireRoles(editorAdmin, app.updatePersonPage))
		r.Post("/person/{id}/update", app.requireRoles(editorAdmin, app.updatePerson))

		r.Get("/character/add", app.requireRoles(editorAdmin, app.addCharacterPage))
		r.Post("/character/add", app.requireRoles(editorAdmin, app.addCharacter))
		r.Get("/character/{id}/update", app.requireRoles(editorAdmin, app.updateCharacterPage))
		r.Post("/character/{id}/update", app.requireRoles(editorAdmin, app.updateCharacter))

		r.Get("/category/add", app.requireRoles(editorAdmin, app.categoryAddPage))
		r.Post("/category/add", app.requireRoles(editorAdmin, app.categoryAdd))
		r.Get("/category/{id}/update", app.requireRoles(editorAdmin, app.categoryUpdatePage))
		r.Post("/category/{id}/update", app.requireRoles(editorAdmin, app.categoryUpdate))

		r.Get("/series/add", app.requireRoles(editorAdmin, app.seriesAddPage))
		r.Post("/series/add", app.requireRoles(editorAdmin, app.seriesAdd))
		r.Get("/series/{id}/update", app.requireRoles(editorAdmin, app.seriesUpdatePage))
		r.Post("/series/{id}/update", app.requireRoles(editorAdmin, app.seriesUpdate))

		r.Get("/recurring/add", app.requireRoles(editorAdmin, app.recurringAddPage))
		r.Post("/recurring/add", app.requireRoles(editorAdmin, app.recurringAdd))
		r.Get("/recurring/{id}/update", app.requireRoles(editorAdmin, app.recurringUpdatePage))
		r.Post("/recurring/{id}/update", app.requireRoles(editorAdmin, app.recurringUpdate))

		r.Get("/tag/add", app.requireRoles(editorAdmin, app.tagAddPage))
		r.Post("/tag/add", app.requireRoles(editorAdmin, app.tagAdd))
		r.Get("/tag/{id}/update", app.requireRoles(editorAdmin, app.tagUpdatePage))
		r.Post("/tag/{id}/update", app.requireRoles(editorAdmin, app.tagUpdate))
		r.Get("/tag/row", app.requireRoles(editorAdmin, app.tagRow))
	})

	return r
}
