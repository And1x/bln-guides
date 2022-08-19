package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/and1x/bln--h/pkg/forms"
	"github.com/and1x/bln--h/pkg/models"
)

//var td TemplateData // middlerware should make this unnecessary

func (app *app) homeSiteHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.tmpl", &TemplateData{})
}

// createGuideFormHandler gets called via "get" to show createguide Form
func (app *app) createGuideFormHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createguide.page.tmpl", &TemplateData{
		Form: forms.New(nil),
	})
}

func (app *app) createGuideHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content")
	form.MaxLength("title", 80)

	if !form.Valid() {
		app.render(w, r, "createguide.page.tmpl", &TemplateData{Form: form})
		return
	}

	id, err := app.guides.Insert(form.Get("title"), form.Get("content"), "anon")
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/guide/%d", id), http.StatusSeeOther)
}

func (app *app) editGuideFormHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	gid, err := app.guides.GetById(id, false) // false bc edit in md not html
	if err == models.ErrNoRows {
		app.clientError(w, http.StatusNotFound)
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	td := TemplateData{Guide: gid}

	app.render(w, r, "editguide.page.tmpl", &td)
}

func (app *app) editGuideHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content")
	form.MaxLength("title", 80)

	if !form.Valid() {
		app.render(w, r, "editguide.page.tmpl", &TemplateData{ // if invalid render with edited values not the ones before
			Guide: &models.Guide{
				Id:      id,
				Title:   form.Get("title"),
				Content: template.HTML(form.Get("content")),
			},
			Form: form,
		})
		return
	}

	err = app.guides.UpdateById(form.Get("title"), form.Get("content"), id) // form.Get returns validated values instead using r.PostFormValues
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/guide/%d", id), http.StatusSeeOther)
}

// allGuidesHandler lists all guides
func (app *app) allGuidesHandler(w http.ResponseWriter, r *http.Request) {
	td := TemplateData{}

	ga, err := app.guides.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}
	td.Guides = ga

	app.render(w, r, "allguides.page.tmpl", &td)
}

// deleteGuideHandler deletes selected Guide by id and redirects to allGuides
// accessed by allguides and singleguide pages
func (app *app) deleteGuideHandler(w http.ResponseWriter, r *http.Request) {

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
	http.Redirect(w, r, "/allguides", http.StatusSeeOther)
}

// singleGuideHandler handles via URL requested guide - in form like: .../guide?id=123
func (app *app) singleGuideHandler(w http.ResponseWriter, r *http.Request) {
	td := TemplateData{}

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

	app.render(w, r, "singleguide.page.tmpl", &td)
}
