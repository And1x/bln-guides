package main

import (
	"bytes"
	"net/http"
	"net/url"
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

func TestCreateGuideFormHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	resCode, _, _ := ts.get(t, "/createguide")

	if resCode != http.StatusOK {
		t.Errorf("want %d but got %d", http.StatusOK, resCode)
	}
}

// mock posted form with relevant values // title + content
func TestCreateGuideHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name        string
		title       string
		content     string
		wantResCode int
		wantResBody []byte
	}{
		{"Valid Form", "The Great Stick", "Long time ago, a great stick got...", http.StatusSeeOther, nil},
		{"Invalid Form - empty title", "", "why is it like that", http.StatusOK, []byte("This field cannot be blank!")},
		{"Invalid Form - empty content", "The empty c", "", http.StatusOK, []byte("This field cannot be blank!")},
		{"Invalid Form - title to long", "The ridicolously and unnecesary long title which will hopfully fail to be submitted but pass the test bc it's known.......", "bla", http.StatusOK, []byte("This field is too long")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("title", test.title)
			form.Add("content", test.content)

			resCode, _, body := ts.postForm(t, "/createguide", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", http.StatusOK, resCode)
			}

			if !bytes.Contains(body, test.wantResBody) {
				t.Errorf("want body %s to contain %q", body, test.wantResBody)
			}
		})
	}
}

// cases: valid id, invalid id(-1, dont existing, word )
func TestEditGuideFormHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name        string
		url         string
		wantResCode int
		wantResBody []byte
	}{
		{"Valid id", "/editguide/21", http.StatusOK, []byte("Cant stop, wont stop!")},
		{"Invalid id - negative nr", "/editguide/-1", http.StatusNotFound, nil},
		{"Invalid id - dont exist", "/editguide/2115", http.StatusNotFound, nil},
		{"Invalid id - word", "/editguide/twentyone", http.StatusNotFound, nil},
		{"Invalid id - empty id", "/editguide/", http.StatusNotFound, nil},
		{"Invalid id - float nr", "/editguide/4.321", http.StatusNotFound, nil},
		{"Invalid url - ending slash", "/editguide/21/", http.StatusNotFound, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			resCode, _, body := ts.get(t, test.url)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", test.wantResCode, resCode)
			}

			if !bytes.Contains(body, test.wantResBody) {
				t.Errorf("want body %s to contain %q", body, test.wantResBody)
			}
		})
	}
}
