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
	app.render(w, r, "createguide.page.tmpl", &TemplateData{Form: forms.New(nil)})
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

	app.render(w, r, "editguide.page.tmpl", &TemplateData{
		Guide: gid,
		Form:  forms.New(nil),
	})
}

// editGuideHandler updates a valid post request in the DB
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
		return
	}

	// throw client error for cases of invalid request like r.FormValue("delete") != "Delete"
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
// in Steps: // todo: better workflow??
// 1. Get User Reg. Form Data
// 2. Create User in DB with default(unusable) LNbits values
// 3. Create LNbits wallet/user
// 4. Update DB entry with newly created LNbits fields
func (app *app) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	// 1.Step
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "password")
	form.MinLength("password", 7)

	// add default LN- and Mail address if users leaves fields empty
	// format: name + @example.com
	if form.Get("lnaddr") == "" {
		form.Set("lnaddr", form.Get("name")+"@example.com")
	}
	if form.Get("email") == "" {
		form.Set("email", form.Get("name")+"@example.com")
	}

	form.ValidMail("lnaddr", "email")

	if !form.Valid() {
		app.render(w, r, "register.page.tmpl", &TemplateData{Form: form})
		return
	}

	name := form.Get("name")

	// 2.Step
	err = app.users.New(name, form.Get("password"), form.Get("lnaddr"), form.Get("email"))

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

	// 3.Step
	lnbuid, lnbadminkey, lnbinvoice, err := app.lnProvider.CreateUserWallet(name) //lnbits.CreateUserWallet(name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// 4.Step
	err = app.users.UpdateLNbByName(lnbuid, lnbadminkey, lnbinvoice, name)
	if err != nil {
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
	name := form.Get("name")
	id, err := app.users.Authenticate(name, form.Get("password"))
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
	app.session.Put(r, "userName", name) // if user successfully logged in we get his userName this way; no need to make a DB request...

	http.Redirect(w, r, "/createguide", http.StatusSeeOther)
}

// profile Handler shows Balance and Nav to Settings/Password change
func (app *app) profileHandler(w http.ResponseWriter, r *http.Request) {

	// Get invoiceKey from DB to call LnbitsAPI to get Balance
	loggedinUserId := app.session.GetInt(r, "userID") // todo: 1: see def:

	// if GetInvoiceKey or GetBalance fails -> show Balance currently not available but still render profile Page
	ikey, err := app.users.GetInvoiceKey(loggedinUserId)
	if err != nil {
		app.errorLog.Printf("couldn't receive incoiceKey from DB") // todo: errorlog fine here?
		app.render(w, r, "profile.page.tmpl", &TemplateData{
			StringMap: map[string]string{"Balance": fmt.Sprintln("Currently not available")},
		})
		return
	}

	balance, err := app.lnProvider.GetBalance(ikey)
	if err != nil {
		app.errorLog.Printf("couldn't get Balance from LNbits")
		app.render(w, r, "profile.page.tmpl", &TemplateData{
			StringMap: map[string]string{"Balance": fmt.Sprintln("Currently not available")},
		})
		return
	}

	app.render(w, r, "profile.page.tmpl", &TemplateData{
		StringMap: map[string]string{"Balance": fmt.Sprintf("%d sats", balance/1000)},
	})
}

// settingsUserFormHandler shows the Settings Page for the User
func (app *app) settingsUserFormHandler(w http.ResponseWriter, r *http.Request) {

	loggedinUserId := app.session.GetInt(r, "userID") // todo: 1:def: better with authentiction from dB like PW check instead take it from session
	if loggedinUserId <= 0 {
		app.clientError(w, http.StatusForbidden) // todo: StatusForbidden appropriate?
		return
	}

	userData, err := app.users.GetById(loggedinUserId)
	if err == models.ErrNoRows {
		app.clientError(w, http.StatusNotFound)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "usersetmail.page.tmpl", &TemplateData{
		User: userData, // todo: could this leak more Data than neccessary? maybe it's better just return needed User fields...
		Form: forms.New(nil),
	})
}

// settingsFormHandler saves changed user Settings in DB
func (app *app) settingsUserHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.ValidMail("lnaddr", "email")

	if !form.Valid() {
		app.render(w, r, "usersetmail.page.tmpl", &TemplateData{ // todo: instead to render just call settingsUserFormHandler per GET req and add session Flash MSG for Invalid Form?
			User: &models.User{
				LNaddr: form.Get("lnaddr"),
				Email:  form.Get("email"),
			},
			Form: form,
		})
		return
	}

	// update DB entry and check if mails already used / add then to form errors in case they are used
	err = app.users.UpdateByUid(app.authUserId(r), form.Get("lnaddr"), form.Get("email"))
	if err == models.ErrLnaddrAlreadyUsed || err == models.ErrEmailAlreadyUsed {
		switch {
		case err == models.ErrLnaddrAlreadyUsed:
			form.Errors.Add("lnaddr", "Lightning Address already exists")
		case err == models.ErrEmailAlreadyUsed:
			form.Errors.Add("email", "Email already exists")
		}
		app.render(w, r, "usersetmail.page.tmpl", &TemplateData{
			User: &models.User{
				LNaddr: form.Get("lnaddr"),
				Email:  form.Get("email"),
			},
			Form: form,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flashMsg", "Settings changed.")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// settingsUserPwFormHandler shows page to change User password
func (app *app) settingsUserPwFormHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "usersetpw.page.tmpl", &TemplateData{Form: forms.New(nil)})
}

func (app *app) settingsUserPwHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	// enough just to check newPassword hence repeat must equal it and old alredy got checked
	form.MinLength("newPassword", 7)
	if !form.Valid() {
		app.render(w, r, "usersetpw.page.tmpl", &TemplateData{Form: form})
		return
	}

	// Authenticate - check old password and name to authenticate and get UserID
	uid, err := app.users.Authenticate(app.getUserName(r), form.Get("oldPassword"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("oldPassword", "Password is incorrect") // todo: what if session name got tampered?
		app.render(w, r, "usersetpw.page.tmpl", &TemplateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// check if newPassword and Repeat Password are the same
	if form.Get("newPassword") != form.Get("confirmPassword") {
		form.Errors.Add("newPassword", "New Password is different to confirm new password")
		app.render(w, r, "usersetpw.page.tmpl", &TemplateData{Form: form})
		return
	}

	// Update DB
	err = app.users.UpdatePwByUid(uid, form.Get("newPassword"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// redirect with success message
	app.session.Put(r, "flashMsg", "Password changed.")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// logoutUserHandler removes the UserID from the session -> user isn't authenticated anymore
func (app *app) logoutUserHandler(w http.ResponseWriter, r *http.Request) {

	app.session.Remove(r, "userID") // todo remove userName aswell
	app.session.Put(r, "flashMsg", "Successfully logged out!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
