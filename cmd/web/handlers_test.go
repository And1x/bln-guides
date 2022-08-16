package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/and1x/bln--h/testing_init"
)

func TestHomeSiteHandler(t *testing.T) {

	app := &app{
		infoLog:  log.New(ioutil.Discard, "", 0),
		errorLog: log.New(ioutil.Discard, "", 0),
		//infoLog:  log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime),
		//errorLog: log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile),
	} // mock app

	req := httptest.NewRequest("Get", "/", nil)
	res := httptest.NewRecorder()

	app.homeSiteHandler(res, req)

	want := http.StatusOK
	got := res.Code

	if got != want {
		t.Errorf("want %d but got %d", want, got)
	}
}
