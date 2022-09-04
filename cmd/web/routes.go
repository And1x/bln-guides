package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *app) routes() http.Handler {

	defaultMiddleware := alice.New(app.recoverPanic, app.logging)

	r := chi.NewRouter()
	r.Get("/", http.HandlerFunc(app.homeSiteHandler))
	r.Get("/allguides", http.HandlerFunc(app.allGuidesHandler))
	r.Post("/deleteguide", http.HandlerFunc(app.deleteGuideHandler))
	r.Get("/createguide", http.HandlerFunc(app.createGuideFormHandler))
	r.Post("/createguide", http.HandlerFunc(app.createGuideHandler))
	r.Post("/editguide", http.HandlerFunc(app.editGuideHandler))
	r.Get("/editguide/{id}", http.HandlerFunc(app.editGuideFormHandler))
	r.Get("/guide/{id}", http.HandlerFunc(app.singleGuideHandler))

	r.Get("/user/register", http.HandlerFunc(app.registerUserFormHandler))
	r.Post("/user/register", http.HandlerFunc(app.registerUserHandler))

	r.Get("/user/login", http.HandlerFunc(app.loginUserFormHandler))
	r.Post("/user/login", http.HandlerFunc(app.loginUserHandler))
	r.Post("/user/logout", http.HandlerFunc(app.logoutUserHandler))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return defaultMiddleware.Then(r)
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
