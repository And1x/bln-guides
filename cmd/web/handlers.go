package main

import (
	"html/template"
	"log"
	"net/http"
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

		gg, err := dbg.GetById(id)
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

// With goldmark - md lang to easily convert md to html
/*
func ShowGuidesHandler(w http.ResponseWriter, r *http.Request) {

	// specify goldmark extension
	md := goldmark.New(
		goldmark.WithExtensions(extension.TaskList),
		goldmark.WithExtensions(extension.Footnote),
	)

	var buf bytes.Buffer
	source := []byte(r.FormValue("content"))
	if err := md.Convert(source, &buf); err != nil {
		panic(err)
	}

	if r.FormValue("title") != "" && r.FormValue("content") != "" {
		td.Guides = append(td.Guides, &models.Guide{
			Title:   r.FormValue("title"),
			Content: template.HTML(buf.String()),
			Author:  "Anon",
		})

	}
	render(w, "./ui/templates/showguides.tmpl", td)
}
*/

func ShowGuidesHandler(w http.ResponseWriter, r *http.Request) {

	dbg := postgres.GuidesModel{DB: DB}

	ga, err := dbg.GetAll()
	if err != nil {
		log.Println(err) // handle this err better not just printing out
	}
	td.Guides = ga

	render(w, "./ui/templates/showguides.tmpl", td)
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
