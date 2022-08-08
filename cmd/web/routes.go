package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

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
}
