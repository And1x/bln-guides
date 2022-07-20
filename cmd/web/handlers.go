package main

import (
	"html/template"
	"log"
	"net/http"
)

type TemplateData struct {
	Guide Guide
}

type Guide struct {
	Title   string
	Content string
	Author  string
}

func HomeSiteHandler(w http.ResponseWriter, r *http.Request) {

	tp, err := template.ParseFiles("./ui/templates/home.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tp.Execute(w, "Hello")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ShowGuidesHandler(w http.ResponseWriter, r *http.Request) {

	td := TemplateData{
		Guide: Guide{
			Title:   "The Great",
			Content: "Lirum ipsum trallalala",
			Author:  "Jk Rowling",
		},
	}

	render(w, "./ui/templates/showguides.tmpl", td)
}

func render(w http.ResponseWriter, filename string, td TemplateData) {

	tp, err := template.ParseFiles(filename)
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
