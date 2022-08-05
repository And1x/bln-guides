package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/and1x/bln--h/pkg/models"
)

type TemplateData struct {
	Guide   *models.Guide
	Guides  []*models.Guide
	InfoMsg map[string]string
}

func humandate(t time.Time) string {
	return t.Local().Format("02 Jan 2006 at 15:04")
	//return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humandate": humandate,
}

var td TemplateData // middlerware should make this unnecessary

func (app *app) homeSiteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(&td)

	if r.URL.Path != "/" && r.URL.Path != "/home" {
		app.clientError(w, http.StatusNotFound)
		return
	}

	app.render(w, "./ui/templates/home.tmpl", td)
}

// createGuideFormHandler gets called via "get" to show createguide Form
func (app *app) createGuideFormHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, "./ui/templates/createguide.tmpl", td)
}

func (app *app) createGuideHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		// return to create Guide Form in case no title and content
		if strings.TrimSpace(r.FormValue("title")) == "" || strings.TrimSpace(r.FormValue("content")) == "" { // todo: handle empty title ant content form directly in the handler?
			// todo: refactor hence after one empty post its saved in
			// emptyfields  so we get always this InfoMessage even when we don't post
			td.InfoMsg = map[string]string{}
			td.InfoMsg["emptyFields"] = "Please enter title and content"
			app.render(w, "./ui/templates/createguide.tmpl", td)
			return
		}

		id, err := app.guides.Insert(r.FormValue("title"), r.FormValue("content"), "anon")
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/guide?id=%d", id), http.StatusSeeOther)

	} else if r.Method == http.MethodGet { // todo: maybe better handled in routes.go
		app.createGuideFormHandler(w, r)
	} else {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
}

// EditGuidesHandler handles 2 kind of request
// 1. Shows Title and content by ID to edit in HTML Forms
// 2. If edit gots submitted - editGuidesHandler gets called again to save in DB and redirects to singleguide.tmpl
func (app *app) editGuidesHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound) // todo: StatusNotfound appropriate?
		return
	}

	if r.FormValue("submitEdit") == "Save" {
		err := app.guides.UpdateById(r.FormValue("title"), r.FormValue("content"), id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/guide?id=%d", id), http.StatusSeeOther)
	}

	gid, err := app.guides.GetById(id, false) // false bc edit in md not html
	if err != nil {
		app.serverError(w, err)
		return
	}

	td.Guide = gid
	fmt.Println(gid) // dont forget to delete

	app.render(w, "./ui/templates/editguide.tmpl", td)
}

func (app *app) allGuidesHandler(w http.ResponseWriter, r *http.Request) {

	// dont't know if this is good design / used to be able to use 1html form for delete and edit. One form means 1action see html
	if r.FormValue("edit") == "Edit" {
		app.editGuidesHandler(w, r)
		return
	}

	if r.FormValue("delete") == "Delete" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil || id < 1 {
			app.clientError(w, http.StatusNotFound)
			return
		}
		err = app.guides.DeleteById(id)
		if err != nil {
			app.serverError(w, err)
			return
		}
		fmt.Println(r.FormValue("id")) //todo: delete in production
	}

	ga, err := app.guides.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}
	td.Guides = ga

	app.render(w, "./ui/templates/allguides.tmpl", td)
}

// singleGuideHandler handles via URL requested guide - in form like: .../guide?id=123
func (app *app) singleGuideHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}
	guide, err := app.guides.GetById(id, true)
	if err != nil {
		app.serverError(w, err)
		return
	}
	td.Guide = guide

	app.render(w, "./ui/templates/singleguide.tmpl", td)
}

func (app *app) render(w http.ResponseWriter, filename string, td TemplateData) {

	tp, err := template.New("base").Funcs(functions).ParseFiles(filename, "./ui/templates/base.layout.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = tp.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
		return
	}
}
