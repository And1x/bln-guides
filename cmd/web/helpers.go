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
