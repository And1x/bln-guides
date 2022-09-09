package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *app) routes() http.Handler {

	// #1
	//defaultMiddleware := alice.New(app.recoverPanic, app.logging)
	//statefulMiddleware := alice.New(app.session.Enable)

	r := chi.NewRouter()
	//r.Get("/", statefulMiddleware.ThenFunc(app.homeSiteHandler).ServeHTTP)

	r.Use(app.recoverPanic, app.logging) // register our middleware to all routes

	r.Group(func(r chi.Router) { // group routes that should have subsequent(following) middleware
		r.Use(app.session.Enable, noSurf) // register session middleware // and csrf protection

		r.Get("/", http.HandlerFunc(app.homeSiteHandler))
		r.Get("/allguides", http.HandlerFunc(app.allGuidesHandler))
		r.Get("/guide/{id}", http.HandlerFunc(app.singleGuideHandler))

		r.Get("/user/register", http.HandlerFunc(app.registerUserFormHandler))
		r.Post("/user/register", http.HandlerFunc(app.registerUserHandler))
		r.Get("/user/login", http.HandlerFunc(app.loginUserFormHandler))
		r.Post("/user/login", http.HandlerFunc(app.loginUserHandler))

		// this group is only for authenticated users accessable
		r.Group(func(r chi.Router) {
			r.Use(app.requireAuth)
			r.Get("/createguide", http.HandlerFunc(app.createGuideFormHandler))
			r.Post("/createguide", http.HandlerFunc(app.createGuideHandler))
			r.Post("/deleteguide", http.HandlerFunc(app.deleteGuideHandler))
			r.Post("/editguide", http.HandlerFunc(app.editGuideHandler))
			r.Get("/editguide/{id}", http.HandlerFunc(app.editGuideFormHandler))
			r.Get("/user/settings", http.HandlerFunc(app.settingsUserFormHandler))
			r.Post("/user/settings", http.HandlerFunc(app.settingsUserHandler))
			r.Post("/user/logout", http.HandlerFunc(app.logoutUserHandler))
		})

	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// #1
	//return defaultMiddleware.Then(r)

	return r
}

/*
func (app *app) routes() http.Handler {

	defaultMiddleware := alice.New(app.recoverPanic, app.logging)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.homeSiteHandler))
	mux.Get("/allguides", http.HandlerFunc(app.allGuidesHandler))
	mux.Post("/deleteguide", http.HandlerFunc(app.deleteGuideHandler))
	mux.Get("/createguide", http.HandlerFunc(app.createGuideFormHandler))
	mux.Post("/createguide", http.HandlerFunc(app.createGuideHandler))
	mux.Post("/editguide", http.HandlerFunc(app.editGuideHandler))
	mux.Get("/editguide/:id", http.HandlerFunc(app.editGuideFormHandler))
	mux.Get("/guide/:id", http.HandlerFunc(app.singleGuideHandler))

	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fs))

	return defaultMiddleware.Then(mux)
}*/
