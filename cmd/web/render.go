package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/and1x/bln--h/pkg/forms"
	"github.com/and1x/bln--h/pkg/models"
)

type TemplateData struct {
	Guide  *models.Guide
	Guides []*models.Guide
	Form   *forms.Form
}

func humandate(t time.Time) string {

	if t.IsZero() {
		return ""
	}
	return t.Local().Format("02 Jan 2006 at 15:04")
	//return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humandate": humandate,
}

func (app *app) render(w http.ResponseWriter, filename string, td *TemplateData) {

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
