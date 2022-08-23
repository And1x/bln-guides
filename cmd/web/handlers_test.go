package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

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
		author      string
		wantResCode int
		wantResBody []byte
	}{
		{"Valid Form", "The Great Stick", "Long time ago, a great stick got...", "blaa@getalby.com", http.StatusSeeOther, nil},
		{"Invalid Form - empty title", "", "why is it like that", "blaa@getalby.com", http.StatusOK, []byte("This field cannot be blank!")},
		{"Invalid Form - empty content", "The empty c", "", "blaa@getalby.com", http.StatusOK, []byte("This field cannot be blank!")},
		{"Invalid Form - title to long", "The ridicolously and unnecesary long title which will hopfully fail to be submitted but pass the test bc it's known.......", "bla", "blaa@getalby.com", http.StatusOK, []byte("This field is too long")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("title", test.title)
			form.Add("content", test.content)
			form.Add("author", test.author)

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

func TestEditGuideHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name        string
		id          string
		title       string
		content     string
		author      string
		wantResCode int
		wantResBody []byte
	}{
		{"Valid Form", "21", "The Great Stick", "Long time ago, a great stick got...", "blaa@getalby.com", http.StatusSeeOther, nil},
		{"Invalid Form - empty title", "21", "", "why is it like that", "blaa@getalby.com", http.StatusOK, []byte("This field cannot be blank!")},
		{"Invalid Form - empty content", "21", "The empty c", "", "blaa@getalby.com", http.StatusOK, []byte("This field cannot be blank!")},
		{"Invalid Form - title to long", "21", "The ridicolously and unnecesary long title which will hopfully fail to be submitted but pass the test bc it's known.......", "bla", "blaa@getalby.com", http.StatusOK, []byte("This field is too long")},
		{"Invalid Form - Id changed negative", "-45", "Hello", "how did u change the id?", "blaa@getalby.com", http.StatusNotFound, nil},
		//{"Invalid Form - Id changed no entry", "78", "Hello", "how did u change the id?", http.StatusNotFound, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("title", test.title)
			form.Add("content", test.content)
			form.Add("author", test.author)
			form.Add("id", test.id)

			resCode, _, body := ts.postForm(t, "/editguide", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", http.StatusOK, resCode)
			}

			if !bytes.Contains(body, test.wantResBody) {
				t.Errorf("want body %s to contain %q", body, test.wantResBody)
			}
		})
	}
}

func TestAllGuidesHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	resCode, _, body := ts.get(t, "/allguides")

	if resCode != http.StatusOK {
		t.Errorf("want %d but got %d", http.StatusOK, resCode)
	}

	if !bytes.Contains(body, []byte("Cant stop")) { // see mock.guide - title starts with Can't stop...
		t.Errorf("want body %s to contain %q", body, []byte("cant stop"))
	}
}

func TestDeleteGuideHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name        string
		id          string
		delete      string
		wantResCode int
	}{
		{"Valid", "21", "Delete", http.StatusSeeOther},
		{"Invalid Delete Form Value", "21", "Noi", http.StatusSeeOther},
		{"Invalid id", "-1", "Delete", http.StatusNotFound},
		{"Invalid id dont exist", "465", "Delete", http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("id", test.id)
			form.Add("delete", test.delete)

			resCode, _, _ := ts.postForm(t, "/deleteguide", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", http.StatusOK, resCode)
			}
		})
	}
}

func TestSingleGuideHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name        string
		url         string
		wantResCode int
		wantResBody []byte
	}{
		{"Valid id", "/guide/21", http.StatusOK, []byte("Cant stop, wont stop!")},
		{"Invalid id - negative nr", "/guide/-1", http.StatusNotFound, nil},
		{"Invalid id - dont exist", "/guide/2115", http.StatusNotFound, nil},
		{"Invalid id - word", "/guide/twentyone", http.StatusNotFound, nil},
		{"Invalid id - empty id", "/guide/", http.StatusNotFound, nil},
		{"Invalid id - float nr", "/guide/4.321", http.StatusNotFound, nil},
		{"Invalid url - ending slash", "/guide/21/", http.StatusNotFound, nil},
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
