package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type TemplateData struct {
	Guides []*Guide
	Text   string
}

type Guide struct {
	Id      int
	Title   string
	Content string
	Author  string
	Created time.Time
	Updated time.Time
}

func HomeSiteHandler(w http.ResponseWriter, r *http.Request) {
	td := TemplateData{
		Text: "hello this text should be deleted until .. . . .  . ..  .. . .",
	}
	render(w, "./ui/templates/home.tmpl", td)
}

func CreateGuideHandler(w http.ResponseWriter, r *http.Request) {
	td := TemplateData{
		Text: "hello guide",
	}
	render(w, "./ui/templates/createguide.tmpl", td)
}

var td TemplateData

func ShowGuidesHandler(w http.ResponseWriter, r *http.Request) {

	if r.FormValue("title") != "" && r.FormValue("content") != "" {
		td.Guides = append(td.Guides, &Guide{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			Author:  "Anon",
		})

	}

	render(w, "./ui/templates/showguides.tmpl", td)
}

func render(w http.ResponseWriter, filename string, td TemplateData) {

	tp, err := template.ParseFiles(filename, "./ui/templates/base.layout.tmpl")
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
