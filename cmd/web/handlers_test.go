package main

import (
	"net/http"
	"testing"
)

// func TestHomeSiteHandler(t *testing.T) {

// 	app := &app{
// 		infoLog:  log.New(ioutil.Discard, "", 0),
// 		errorLog: log.New(ioutil.Discard, "", 0),
// 		// infoLog:  log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime),
// 		// errorLog: log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile),
// 	} // mock app

// 	req := httptest.NewRequest("Get", "/", nil)
// 	res := httptest.NewRecorder()

// 	app.homeSiteHandler(res, req)

// 	want := http.StatusOK
// 	got := res.Code

// 	if got != want {
// 		t.Errorf("want %d but got %d", want, got)
// 	}
// }

func TestHomeSiteHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	resCode, _, _ := ts.get(t, "/")

	if resCode != http.StatusOK {
		t.Errorf("want %d but got %d", http.StatusOK, resCode)
	}
}
