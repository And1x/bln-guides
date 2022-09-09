package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError logs a stack trace + error msg and sends a http Status Error to the User
func (app *app) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n,%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends http Status Error to the User -
// (used for consistency eg. app.serverError and app.clientError)
func (app *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// authUserId returns userID form the user session
func (app *app) authUserId(r *http.Request) int {
	return app.session.GetInt(r, "userID")
}

// getUserName return userName from the user session
func (app *app) getUserName(r *http.Request) string {
	return app.session.GetString(r, "userName")
}

// isAuthorized checks if the users session ID fits to the guide.UserId he wants to edit/delete
// todo: probably better as middleware: however this needs to pull guide.Id from GET request(url)
// and POST request(form) - seems more cumbersome than just use it in handlers
func (app *app) isAuthorized(guideId int, w http.ResponseWriter, r *http.Request) bool {

	guide, err := app.guides.GetById(guideId, false)
	if err != nil {
		//app.clientError(w, http.StatusNotFound)
		return false
	}

	if app.session.GetInt(r, "userID") == guide.UserID {
		return true
	} else {
		return false
	}
}
