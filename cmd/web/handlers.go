package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/and1x/bln--h/pkg/forms"
	"github.com/and1x/bln--h/pkg/models"
	"github.com/go-chi/chi/v5"
)

func (app *app) homeSiteHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.tmpl", nil)
	//app.render(w, r, "home.page.tmpl", &TemplateData{})
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

	//get userID from session to know who created the guide
	loggedinUserId := app.session.GetInt(r, "userID")

	id, err := app.guides.Insert(form.Get("title"), form.Get("content"), loggedinUserId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flashMsg", "New Guide created!")

	http.Redirect(w, r, fmt.Sprintf("/guide/%d", id), http.StatusSeeOther)
}

func (app *app) editGuideFormHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	// not authorized user trys to edit other users guide/ or guide doesn't exists
	// todo: if guide doesn't exists, still a StatusForbidden response?
	if !app.isAuthorized(id, w, r) {
		app.clientError(w, http.StatusForbidden)
		return
	}

	gid, err := app.guides.GetById(id, false) // false bc edit in md not html
	if err == models.ErrNoRows {
		app.clientError(w, http.StatusNotFound)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	td := TemplateData{Guide: gid}

	app.render(w, r, "editguide.page.tmpl", &td)
}

// editGuideHandler updates a valid post request in the DB
// ### todo: if id got somehowe changed -
// ### integrate a way to that users cant change others guides or invalid ones -
// ### only negative numbers get checked so far
// ### authentication could solve this
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

	if !app.isAuthorized(id, w, r) {
		app.clientError(w, http.StatusForbidden)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content")
	form.MaxLength("title", 80)

	if !form.Valid() {
		app.render(w, r, "editguide.page.tmpl", &TemplateData{ // if invalid render with edited values not the ones before
			Guide: &models.Guide{
				Id:      id, // todo: ID useful here??
				Title:   form.Get("title"),
				Content: template.HTML(form.Get("content")),
			},
			Form: form,
		})
		return
	}

	err = app.guides.UpdateById(id, form.Get("title"), form.Get("content")) // form.Get returns validated values instead using r.PostFormValues
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/guide/%d", id), http.StatusSeeOther)
}

// allGuidesHandler lists all guides
func (app *app) allGuidesHandler(w http.ResponseWriter, r *http.Request) {

	ga, err := app.guides.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}
	td := TemplateData{Guides: ga}

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

		if !app.isAuthorized(id, w, r) {
			app.clientError(w, http.StatusForbidden)
			return
		}

		err = app.guides.DeleteById(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		app.session.Put(r, "flashMsg", "Your Guide got deleted!.")

		http.Redirect(w, r, "/allguides", http.StatusSeeOther)
	}

	// todo: Is this legit? - if handler gets called but can't delete respond with clientErr StatusBadRequest
	app.clientError(w, http.StatusBadRequest)
}

// singleGuideHandler handles via URL requested guide - in form like: .../guide/123
func (app *app) singleGuideHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	guide, err := app.guides.GetById(id, true)
	if err == models.ErrNoRows {
		app.clientError(w, http.StatusNotFound)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "singleguide.page.tmpl", &TemplateData{Guide: guide})
}

// registerUserFormHandler shows Form for Registration
func (app *app) registerUserFormHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.tmpl", &TemplateData{Form: forms.New(nil)})
}

// registerUserHandler creates a new DB entry with the Users Details
func (app *app) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "password")
	form.MinLength("password", 7)
	form.ValidMail("lnaddr", "email") // if empty no err bc then add name@blnguide.com

	if !form.Valid() {
		app.render(w, r, "register.page.tmpl", &TemplateData{Form: form})
		return
	}

	// Create new User in DB -- let user redo if email,lnaddr or name is already in DB
	username := form.Get("name")
	lnaddr := form.Get("lnaddr")
	email := form.Get("email")
	if lnaddr == "" {
		lnaddr = username + "@blnguide.lnd"
	}
	if email == "" {
		email = username + "@blnguide.com"
	}

	err = app.users.New(form.Get("name"), form.Get("password"), lnaddr, email)

	if err == models.ErrNameAlreadyUsed || err == models.ErrLnaddrAlreadyUsed || err == models.ErrEmailAlreadyUsed {
		switch {
		case err == models.ErrNameAlreadyUsed:
			form.Errors.Add("name", "Name already exists")
		case err == models.ErrLnaddrAlreadyUsed:
			form.Errors.Add("lnaddr", "Lightning Address already exists")
		case err == models.ErrEmailAlreadyUsed:
			form.Errors.Add("email", "Email already exists")
		}
		app.render(w, r, "register.page.tmpl", &TemplateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flashMsg", "Successfully registered. Please Login.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// loginUserFormHandler shows the Login Form
func (app *app) loginUserFormHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &TemplateData{Form: forms.New(nil)})
}

// loginUserHandler authenticates a user and creates a session for them
func (app *app) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("name"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Name or password is incorrect")
		app.render(w, r, "login.page.tmpl", &TemplateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	//app session
	app.session.Put(r, "flashMsg", "Successfully logged in")
	app.session.Put(r, "userID", id)

	http.Redirect(w, r, "/createguide", http.StatusSeeOther)
}

// logoutUserHandle removes the UserID from the session -> user insn't authenticated anymore
func (app *app) logoutUserHandler(w http.ResponseWriter, r *http.Request) {

	app.session.Remove(r, "userID")
	app.session.Put(r, "flashMsg", "Successfully logged out!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
