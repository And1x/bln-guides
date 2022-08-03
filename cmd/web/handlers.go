package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/and1x/bln--h/pkg/models"
)

type TemplateData struct {
	Guide  *models.Guide
	Guides []*models.Guide
	Text   []string
}

func humandate(t time.Time) string {
	return t.Local().Format("02 Jan 2006 at 15:04")
	//return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humandate": humandate,
}

var td TemplateData

func (app *app) HomeSiteHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}

	if r.FormValue("title") != "" {

		id, err := app.guides.Insert(r.FormValue("title"), r.FormValue("content"), "anon")
		if err != nil {
			http.Error(w, "Couldnt Insert into DB", http.StatusInternalServerError) // todo: Dont expose intera - refactor msg
			log.Println(err)
			return // todo: revisit bc. it may be better wo. return to load home content - better ux
		}

		gg, err := app.guides.GetById(id, true)
		if err != nil {
			http.Error(w, "Cant get Guide by ID", http.StatusInternalServerError) // todo: Dont expose interna - refactor msg
			log.Println(err)                                                      // todo: err handling
		}
		td.Guide = gg
	} else { // show default home page
		td.Guide = &models.Guide{}
	}
	app.render(w, "./ui/templates/home.tmpl", td)
}

func (app *app) CreateGuideHandler(w http.ResponseWriter, r *http.Request) {
	td := TemplateData{}
	app.render(w, "./ui/templates/createguide.tmpl", td)
}

func (app *app) ShowGuidesHandler(w http.ResponseWriter, r *http.Request) {

	// dont't know if this is good design / used to be able to use 1html form for delete and edit. One form means 1action see html
	if r.FormValue("edit") == "Edit" {
		app.EditGuidesHandler(w, r)
		return
	}

	if r.FormValue("delete") == "Delete" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			fmt.Println(err) // todo: err handling
			return
		}
		err = app.guides.DeleteById(id)
		if err != nil {
			fmt.Println(err) // todo: err handling
			return
		}
		fmt.Println(r.FormValue("id"))
	}

	ga, err := app.guides.GetAll()
	if err != nil {
		http.Error(w, "couldn't get DB result", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	td.Guides = ga

	app.render(w, "./ui/templates/showguides.tmpl", td)
}

// EditGuidesHandler handles 2 kind of request
// 1. Shows Title and content by ID to edit in HTML Forms
// 2. If edit gots submitted it gets updated in DB and shown as 1.
func (app *app) EditGuidesHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println(err)
		return
	}

	if r.FormValue("submitEdit") == "Save" {
		err := app.guides.UpdateById(r.FormValue("title"), r.FormValue("content"), id)
		if err != nil {
			http.Error(w, "couldn't update DB query", http.StatusInternalServerError)
			log.Println(err) // todo: err handling
			return
		}
	}

	gid, err := app.guides.GetById(id, false)
	if err != nil {
		log.Println(err)
		return
	}

	td.Guide = gid
	fmt.Println(gid) // dont forget to delete

	app.render(w, "./ui/templates/editguide.tmpl", td)
}

func (app *app) render(w http.ResponseWriter, filename string, td TemplateData) {

	tp, err := template.New("base").Funcs(functions).ParseFiles(filename, "./ui/templates/base.layout.tmpl")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tp.Execute(w, td)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error while Executing tmpl", http.StatusInternalServerError)
	}
}
