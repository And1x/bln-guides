package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *app) routes() http.Handler {

	r := chi.NewRouter()

	r.Use(app.recoverPanic, app.logging, app.setSecureHeader) // register our middleware to all routes

	r.Group(func(r chi.Router) { // group routes that should have subsequent(following) middleware

		r.Use(app.session.Enable, app.noSurf) // register session middleware // and csrf protection

		r.Get("/", http.HandlerFunc(app.homeSiteHandler))
		r.Get("/allguides", http.HandlerFunc(app.allGuidesHandler))
		r.Get("/guide/{id}", http.HandlerFunc(app.singleGuideHandler))

		r.Get("/user/register", http.HandlerFunc(app.registerUserFormHandler))
		r.Post("/user/register", http.HandlerFunc(app.registerUserHandler))
		r.Get("/user/login", http.HandlerFunc(app.loginUserFormHandler))
		r.Post("/user/login", http.HandlerFunc(app.loginUserHandler))

		// this group is only for authenticated users accessable
		r.Group(func(r chi.Router) {

			// Authentication active in Production, not while testing
			if app.inProduction {
				r.Use(app.requireAuth)
			}

			r.Get("/createguide", http.HandlerFunc(app.createGuideFormHandler))
			r.Post("/createguide", http.HandlerFunc(app.createGuideHandler))
			r.Post("/deleteguide", http.HandlerFunc(app.deleteGuideHandler))
			r.Get("/editguide/{id}", http.HandlerFunc(app.editGuideFormHandler))
			r.Post("/editguide", http.HandlerFunc(app.editGuideHandler))

			r.Post("/allguides", http.HandlerFunc(app.upvoteAllGuidesHandler))
			r.Post("/guide/{id}", http.HandlerFunc(app.upvoteSingleGuideHandler))

			r.Get("/user/profile", http.HandlerFunc(app.profileUserHandler))

			r.Get("/user/settings", http.HandlerFunc(app.settingsUserFormHandler))
			r.Post("/user/settings", http.HandlerFunc(app.settingsUserHandler))

			r.Get("/user/settings/password", http.HandlerFunc(app.settingsUserPwFormHandler))
			r.Post("/user/settings/password", http.HandlerFunc(app.settingsUserPwHandler))

			r.Get("/user/deposit", http.HandlerFunc(app.depositFormHandler))
			r.Post("/user/deposit", http.HandlerFunc(app.depositHandler))

			r.Post("/user/logout", http.HandlerFunc(app.logoutUserHandler))

		})

	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return r
}
