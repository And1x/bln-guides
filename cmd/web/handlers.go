package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/and1x/bln--h/pkg/models"
	"github.com/and1x/bln--h/pkg/models/postgres"
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

func HomeSiteHandler(w http.ResponseWriter, r *http.Request) {

	if r.FormValue("title") != "" {

		dbg := &postgres.GuidesModel{DB: DB} // dbg = database guides

		id, err := dbg.Insert(r.FormValue("title"), r.FormValue("content"), "anon")
		if err != nil {
			log.Fatal(err)
		}

		gg, err := dbg.GetById(id, true)
		if err != nil {
			log.Println(err)
		}
		td.Guide = gg
	} else { // show default home page
		td.Guide = &models.Guide{}
	}
	render(w, "./ui/templates/home.tmpl", td)
}

func CreateGuideHandler(w http.ResponseWriter, r *http.Request) {
	td := TemplateData{}
	render(w, "./ui/templates/createguide.tmpl", td)
}

func ShowGuidesHandler(w http.ResponseWriter, r *http.Request) {

	// dont't know if this is good design / used to be able to use 1html form for delete and edit. One form means 1action see html
	if r.FormValue("edit") == "Edit" {
		EditGuidesHandler(w, r)
		return
	}

	dbg := postgres.GuidesModel{DB: DB}

	if r.FormValue("delete") == "Delete" {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			fmt.Println(err)
			return
		}
		err = dbg.DeleteById(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(r.FormValue("id"))
	}

	ga, err := dbg.GetAll()
	if err != nil {
		log.Println(err) // handle this err better not just printing out
	}
	td.Guides = ga

	render(w, "./ui/templates/showguides.tmpl", td)
}

// EditGuidesHandler handles 2 kind of request
// 1. Shows Title and content by ID to edit in HTML Forms
// 2. If edit gots submitted it gets updated in DB and shown as 1.
func EditGuidesHandler(w http.ResponseWriter, r *http.Request) {

	dbg := postgres.GuidesModel{DB: DB}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println(err)
		return
	}

	if r.FormValue("submitEdit") == "Save" {
		err := dbg.UpdateById(r.FormValue("title"), r.FormValue("content"), id)
		if err != nil {
			log.Println(err)
			return
		}
	}

	gid, err := dbg.GetById(id, false)
	if err != nil {
		log.Println(err)
		return
	}

	td.Guide = gid
	fmt.Println(gid)

	render(w, "./ui/templates/editguide.tmpl", td)
}

func render(w http.ResponseWriter, filename string, td TemplateData) {

	tp, err := template.New("base").Funcs(functions).ParseFiles(filename, "./ui/templates/base.layout.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tp.Execute(w, td)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error while Executing tmpl", http.StatusInternalServerError)
	}
}
