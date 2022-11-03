package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/and1x/bln-guides/pkg/forms"
	"github.com/and1x/bln-guides/pkg/models"
	"github.com/justinas/nosurf"
)

type TemplateData struct {
	Guide      *models.Guide
	Guides     []*models.Guide
	Form       *forms.Form
	User       *models.User
	StringMap  map[string]string // todo: add username/authuserid to this map??
	AuthUserId int               // todo: AuthUserId and UserName nedded or just User? we get em by session; no additional DB requests needed which seems like a advantage...
	UserName   string
	FlashMsg   string
	CSRFToken  string
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

// render executes the html template(tmpl)
func (app *app) render(w http.ResponseWriter, r *http.Request, tmpl string, td *TemplateData) {

	ts, ok := app.templateCache[tmpl]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s is not in cache", tmpl))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	buf.WriteTo(w)
}

// addDefaultData adds Data to TemplateData on any render call
func (app *app) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}

	td.AuthUserId = app.authUserId(r) // get User Id from session when User is logged in - if not -> id = 0
	td.UserName = app.getUserName(r)
	td.FlashMsg = app.session.PopString(r, "flashMsg")
	td.CSRFToken = nosurf.Token(r)
	return td
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
