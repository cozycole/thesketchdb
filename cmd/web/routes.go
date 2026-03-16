package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes(staticRoute string, serveStatic bool) http.Handler {
	r := chi.NewRouter()

	if serveStatic {
		r.Group(func(r chi.Router) {
			r.Use(app.recoverPanic)
			fs := http.FileServer(http.Dir(staticRoute))
			app.infoLog.Printf("Starting static file server rooted at %s\n", staticRoute)
			r.Handle("/static/*", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				fs.ServeHTTP(w, r)
			})))
		})
	}

	r.Group(func(r chi.Router) {
		r.Use(
			app.recoverPanic,
			app.secureHeaders,
			app.sessionManager.LoadAndSave,
			app.logRequest,
			app.authenticate,
		)

		// public site routes
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

		r.Post("/sketch/{id}/rating", app.sketchUpdateRating)
		r.Delete("/sketch/{id}/rating", app.sketchDeleteRating)

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

		r.HandleFunc("/ping", ping)

		// public site editor / admin routes
		r.Group(func(r chi.Router) {
			r.Use(
				app.requireAuthentication,
				app.requireRoles(editorAdmin),
			)

			r.Get("/admin*", app.serveCMS)

			r.Get("/sketch/add", app.sketchAddPage)
			r.Post("/sketch/add", app.sketchAdd)
			r.Get("/sketch/{id}/update", app.sketchUpdatePage)
			r.Post("/sketch/{id}/update", app.sketchUpdate)
			r.Post("/sketch/{id}/tag", app.sketchUpdateTags)

			r.Get("/sketch/{id}/cast", app.addCastPage)
			r.Post("/sketch/{id}/cast", app.addCast)

			r.Get("/cast/{id}/update", app.updateCastPage)
			r.Post("/cast/{castId}/update", app.updateCast)
			r.Delete("/cast/{castId}", app.deleteCast)
			r.Patch("/sketch/{id}/cast/order", app.orderCast)

			r.Get("/cast", app.castDropdown)
			r.Get("/cast/{id}/tags", app.castTagUpdateForm)
			r.Post("/cast/{id}/tags", app.castTagUpdate)

			r.Get("/show/add", app.addShowPage)
			r.Post("/show/add", app.addShow)
			r.Get("/show/{id}/update", app.updateShowPage)
			r.Post("/show/{id}/update", app.updateShow)

			r.Post("/show/{id}/season/add", app.addSeason)

			r.Get("/season/{id}/episode/add", app.addEpisodeForm)
			r.Post("/season/{id}/episode/add", app.addEpisode)
			r.Get("/episode/{id}/update", app.updateEpisodeForm)
			r.Post("/episode/{id}/update", app.updateEpisode)

			r.Get("/creator/add", app.addCreatorPage)
			r.Post("/creator/add", app.addCreator)
			r.Get("/creator/{id}/update", app.updateCreatorPage)
			r.Post("/creator/{id}/update", app.updateCreator)

			r.Get("/person/add", app.addPersonPage)
			r.Post("/person/add", app.addPerson)
			r.Get("/person/{id}/update", app.updatePersonPage)
			r.Post("/person/{id}/update", app.updatePerson)

			r.Get("/character/add", app.addCharacterPage)
			r.Post("/character/add", app.addCharacter)
			r.Get("/character/{id}/update", app.updateCharacterPage)
			r.Post("/character/{id}/update", app.updateCharacter)

			r.Get("/category/add", app.categoryAddPage)
			r.Post("/category/add", app.categoryAdd)
			r.Get("/category/{id}/update", app.categoryUpdatePage)
			r.Post("/category/{id}/update", app.categoryUpdate)

			r.Get("/series/add", app.seriesAddPage)
			r.Post("/series/add", app.seriesAdd)
			r.Get("/series/{id}/update", app.seriesUpdatePage)
			r.Post("/series/{id}/update", app.seriesUpdate)

			r.Get("/recurring/add", app.recurringAddPage)
			r.Post("/recurring/add", app.recurringAdd)
			r.Get("/recurring/{id}/update", app.recurringUpdatePage)
			r.Post("/recurring/{id}/update", app.recurringUpdate)

			r.Get("/tag/add", app.tagAddPage)
			r.Post("/tag/add", app.tagAdd)
			r.Get("/tag/{id}/update", app.tagUpdatePage)
			r.Post("/tag/{id}/update", app.tagUpdate)
			r.Get("/tag/row", app.tagRow)
		})

		// admin only public site routes
		r.Group(func(r chi.Router) {
			r.Use(app.requireRoles(adminOnly))

			r.Delete("/season/{id}", app.deleteSeason)
			r.Delete("/episode/{id}", app.deleteEpisode)
		})

		// api routes
		r.Route("/api/v1", func(r chi.Router) {
			// public api routes
			r.Get("/cast", app.listCastAPI)
			r.Get("/characters", app.listCharactersAPI)
			r.Get("/creators", app.listCreatorsAPI)
			r.Get("/episodes", app.listEpisodesAPI)
			r.Get("/people", app.listPeopleAPI)
			r.Get("/recurring-sketches", app.listRecurringAPI)
			r.Get("/sketch-series", app.listSeriesAPI)
			r.Get("/sketches", app.viewSketchesAPI)
			r.Get("/tags", app.listTagsAPI)

			// editor / admin API routes
			r.Group(func(r chi.Router) {
				r.Use(app.requireRoles(editorAdmin))
				r.Get("/admin/sketch/{id}", app.adminGetSketchAPI)
				r.Post("/admin/sketch", app.createSketchAPI)
				r.Put("/admin/sketch/{id}", app.updateSketchAPI)

				r.Get("/admin/sketch/{id}/cast", app.adminGetCastAPI)
				r.Post("/admin/sketch/{id}/cast", app.createCastAPI)
				r.Put("/admin/sketch/{id}/cast/{castId}", app.updateCastAPI)
				r.Delete("/admin/sketch/{id}/cast/{castId}", app.deleteCastAPI)
				r.Put("/admin/sketch/{id}/cast/order", app.updateCastOrderAPI)

				r.Get("/admin/sketch/{id}/quotes", app.adminGetQuotesAPI)
				r.Put("/admin/sketch/{id}/quotes", app.updateQuotesAPI)

				r.Get("/admin/sketch/{id}/videos", app.getSketchVideos)
				r.Post("/admin/sketch/{id}/upload-url", app.generateSketchVideoS3PutUrl)
				r.Post("/admin/sketch/{id}/video-uploaded", app.sketchVideoUploaded)
			})

			// admin only api routes
			r.Group(func(r chi.Router) {
				r.Use(app.requireRoles(adminOnly))
				r.Get("/admin/get-token", app.createAdminToken)
			})
		})
	})
	return r
}
