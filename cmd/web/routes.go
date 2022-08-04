package main

import "net/http"

func (app *app) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.homeSiteHandler)
	mux.HandleFunc("/allguides", app.allGuidesHandler)
	mux.HandleFunc("/createguide", app.createGuideHandler)
	mux.HandleFunc("/editguide", app.editGuidesHandler)
	mux.HandleFunc("/guide", app.singleGuideHandler)

	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	return mux
}
