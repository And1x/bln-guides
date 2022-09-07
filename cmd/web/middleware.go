package main

import (
	"fmt"
	"net/http"
)

// todo: Reasearch if necessary - couldn't find clear anwsers
// Im using template/html which already esacapse html injection...
/*
func setSecureHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}
*/

func (app *app) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// recover is a defered function to recover from a panic
func (app *app) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close") // user gets informed that the con is closed
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// requireAuth reidrects user if userId == 0 -> no user logged; else next hander gets called
func (app *app) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// if userId == 0 then User is not authenticated
		if app.authUserId(r) == 0 {
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
