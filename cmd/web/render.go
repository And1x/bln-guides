package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
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

// createTemplateCache loads all .tmpl in memory instead of reading these files on any page request
func createTemplateCache(tmplPath string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(tmplPath, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		// get just the filename without path
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(tmplPath, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

func (app *app) render(w http.ResponseWriter, r *http.Request, tmpl string, td *TemplateData) {

	ts, ok := app.templateCache[tmpl]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s is not in cache", tmpl))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}
