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

	_, _, body := ts.get(t, "/createguide")

	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name        string
		title       string
		content     string
		wantResCode int
		wantResBody []byte
	}{
		{"Valid Form", "The Great Stick", "Long time ago, a great stick got...", http.StatusSeeOther, nil},
		{"Invalid Form - empty title", "", "why is it like that", http.StatusOK, []byte("Required field!")},
		{"Invalid Form - empty content", "The empty c", "", http.StatusOK, []byte("Required field!")},
		{"Invalid Form - title to long", "The ridicolously and unnecesary long title which will hopfully fail to be submitted but pass the test bc it's known.......", "bla", http.StatusOK, []byte("Too long!")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("title", test.title)
			form.Add("content", test.content)
			form.Add("csrf_token", csrfToken)

			resCode, _, body := ts.postForm(t, "/createguide", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", test.wantResCode, resCode)
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
		{"Invalid id - word", "/editguide/twentyone", http.StatusBadRequest, nil},
		{"Invalid id - negative nr", "/editguide/-1", http.StatusBadRequest, nil},
		{"Invalid id - float nr", "/editguide/4.321", http.StatusBadRequest, nil},
		{"Invalid id - dont exist", "/editguide/2115", http.StatusNotFound, nil},
		{"Invalid id - empty id", "/editguide/", http.StatusNotFound, nil},
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

	_, _, body := ts.get(t, "/editguide/21")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name    string
		id      string
		title   string
		content string
		// UserID      int
		wantResCode int
		wantResBody []byte
	}{
		{"Valid Form", "21", "The Great Stick", "Long time ago, a great stick got...", http.StatusSeeOther, nil},
		{"Invalid Form - empty title", "21", "", "why is it like that", http.StatusOK, []byte("Required field!")},
		{"Invalid Form - empty content", "21", "The empty c", "", http.StatusOK, []byte("Required field!")},
		{"Invalid Form - title to long", "21", "The ridicolously and unnecesary long title which will hopfully fail to be submitted but pass the test bc it's known.......", "bla", http.StatusOK, []byte("Too long!")},
		{"Invalid Form - Id changed negative", "-45", "Hello", "how did u change the id?", http.StatusBadRequest, nil},
		{"Invalid Form - Id changed no entry", "78", "Hello", "how did u change the id?", http.StatusInternalServerError, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("title", test.title)
			form.Add("content", test.content)
			form.Add("id", test.id)
			form.Add("csrf_token", csrfToken)

			resCode, _, body := ts.postForm(t, "/editguide", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", test.wantResCode, resCode)
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

	_, _, body := ts.get(t, "/allguides")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name        string
		id          string
		delete      string
		wantResCode int
	}{
		{"Valid", "21", "Delete", http.StatusSeeOther},
		{"Invalid Delete Form Value", "21", "Noi", http.StatusBadRequest},
		{"Invalid id", "-1", "Delete", http.StatusBadRequest},
		{"Invalid id dont exist", "465", "Delete", http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("id", test.id)
			form.Add("delete", test.delete)
			form.Add("csrf_token", csrfToken)

			resCode, _, _ := ts.postForm(t, "/deleteguide", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", test.wantResCode, resCode)
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
		{"Invalid id - negative nr", "/guide/-1", http.StatusBadRequest, nil},
		{"Invalid id - word", "/guide/twentyone", http.StatusBadRequest, nil},
		{"Invalid id - float nr", "/guide/4.321", http.StatusBadRequest, nil},
		{"Invalid url - ending slash", "/guide/21/", http.StatusNotFound, nil},
		{"Invalid id - dont exist", "/guide/2115", http.StatusNotFound, nil},
		{"Invalid id - empty id", "/guide/", http.StatusNotFound, nil},
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

func TestRegisterUserFormHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	resCode, _, _ := ts.get(t, "/user/register")

	if resCode != http.StatusOK {
		t.Errorf("want %d but got %d", http.StatusOK, resCode)
	}
}

func TestRegisterUserHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/register")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name        string
		username    string
		password    string
		lnaddr      string
		email       string
		wantResCode int
	}{
		{"Valid without addresses", "anon", "password", "", "", http.StatusSeeOther},
		{"Valid complete", "anon", "password", "anon@ln.pay", "anon@mail.com", http.StatusSeeOther},

		// StatusOK on invalid Form Values -> user/register gets rendered again
		{"Invalid short password", "anon", "pas", "", "", http.StatusOK},
		{"Invalid empty name", "", "password", "", "", http.StatusOK},
		{"Invalid empty password", "", "password", "", "", http.StatusOK},
		{"Invalid  lnaddr", "anon", "password", "nolna", "", http.StatusOK},
		{"Invalid  email", "anon", "password", "", "nomail", http.StatusOK},

		// Already used values: (see mocked values in mock.users)
		{"Invalid used name", "Satu Naku", "password", "", "", http.StatusOK},
		{"Invalid used lnaddr", "anon", "password", "Satu@payme.com", "", http.StatusOK},
		{"Invalid used email", "anon", "password", "", "Satu@Naku.com", http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("name", test.username)
			form.Add("password", test.password)
			form.Add("lnaddr", test.lnaddr)
			form.Add("email", test.email)
			form.Add("csrf_token", csrfToken)

			resCode, _, _ := ts.postForm(t, "/user/register", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", test.wantResCode, resCode)
			}
		})
	}
}
func TestLoginUserFormHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	resCode, _, _ := ts.get(t, "/user/login")

	if resCode != http.StatusOK {
		t.Errorf("want %d but got %d", http.StatusOK, resCode)
	}
}

func TestLoginUserHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name        string
		username    string
		password    string
		wantResCode int
		wantResBody []byte
	}{
		{"Valid", "Satu Naku", "oldpassword", http.StatusSeeOther, nil},
		{"Invalid name", "anon", "oldpassword", http.StatusOK, []byte("Name or password is incorrect")},
		{"Invalid password", "Satu Naku", "wrongpw", http.StatusOK, []byte("Name or password is incorrect")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("name", test.username)
			form.Add("password", test.password)
			form.Add("csrf_token", csrfToken)

			resCode, _, body := ts.postForm(t, "/user/login", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", test.wantResCode, resCode)
			}

			if !bytes.Contains(body, test.wantResBody) {
				t.Errorf("want body %s to contain %q", body, test.wantResBody)
			}
		})
	}
}

func TestUpvoteAllGuidesHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/allguides")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name        string
		guideId     string
		wantResCode int
		wantResBody []byte
	}{
		{"Valid", "21", http.StatusOK, []byte("")},
		// todo: insuffiecient balance test -- have to change mocks to do that...
		{"Invalid guide don't exist", "1555", http.StatusInternalServerError, []byte("")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}

			form.Add("gid", test.guideId)
			form.Add("csrf_token", csrfToken)

			resCode, _, body := ts.postForm(t, "/allguides", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", test.wantResCode, resCode)
			}

			if !bytes.Contains(body, test.wantResBody) {
				t.Errorf("want body %s to contain %q", body, test.wantResBody)
			}
		})
	}
}

func TestUpvoteSingleGuideHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/guide/21")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name        string
		url         string
		wantResCode int
		wantResBody []byte
	}{
		{"Valid", "/guide/21", http.StatusOK, []byte("")},
		// todo: insuffiecient balance test -- have to change mocks to do that...
		{"Invalid guide don't exist", "/guide/1555", http.StatusInternalServerError, []byte("")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("csrf_token", csrfToken)

			resCode, _, body := ts.postForm(t, test.url, form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", test.wantResCode, resCode)
			}

			if !bytes.Contains(body, test.wantResBody) {
				t.Errorf("want body %s to contain %q", body, test.wantResBody)
			}
		})
	}
}

// todo: test failing routes?
func TestProfileUserHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	resCode, _, _ := ts.get(t, "/user/profile")

	if resCode != http.StatusOK {
		t.Errorf("want %d but got %d", http.StatusOK, resCode)
	}
}

func TestSettingsUserFormHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	resCode, _, _ := ts.get(t, "/user/settings")

	if resCode != http.StatusOK {
		t.Errorf("want %d but got %d", http.StatusOK, resCode)
	}
}

func TestSettingsUserHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/settings")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name        string
		lnaddr      string
		email       string
		upvote      string
		wantResCode int
		wantResBody []byte
	}{
		{"Valid", "ano@pay.me", "ano@mail.com", "25", http.StatusSeeOther, []byte("")},
		{"Invalid lnaddr", "noad", "ano@mail.com", "25", http.StatusOK, []byte("Invalid!")},
		{"Invalid email", "ano@pay.me", "ano", "25", http.StatusOK, []byte("Invalid!")},
		{"Invalid upvote amount", "ano@pay.me", "ano@mail.com", "ff5", http.StatusOK, []byte("Please enter a number &gt; 0!")},
		{"Invalid upvote negative amount", "ano@pay.me", "ano@mail.com", "-50", http.StatusOK, []byte("Please enter a number &gt; 0!")},
		{"Invalid lnaddr already exists", "Satu@payme.com", "ano@mail.com", "50", http.StatusOK, []byte("already exists")},
		{"Invalid email already exists", "ano@pay.me", "Satu@Naku.com", "50", http.StatusOK, []byte("already exists")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("lnaddr", test.lnaddr)
			form.Add("email", test.email)
			form.Add("upvote", test.upvote)
			form.Add("csrf_token", csrfToken)

			resCode, _, body := ts.postForm(t, "/user/settings", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", test.wantResCode, resCode)
			}

			if !bytes.Contains(body, test.wantResBody) {
				t.Errorf("want body %s to contain %q", body, test.wantResBody)
			}
		})
	}
}

func TestSettingsUserPwFormHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	resCode, _, _ := ts.get(t, "/user/settings/password")

	if resCode != http.StatusOK {
		t.Errorf("want %d but got %d", http.StatusOK, resCode)
	}
}

func TestSettingsUserPwHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/settings/password")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name            string
		oldPassword     string
		newPassword     string
		confirmPassword string
		wantResCode     int
		wantResBody     []byte
	}{
		{"Valid", "oldpassword", "newpassword", "newpassword", http.StatusSeeOther, []byte("")},
		{"Invalid old Password", "invalid", "newpassword", "newpassword", http.StatusOK, []byte("Password is incorrect")},
		{"Invalid confirm Password", "oldpassword", "newpassword", "invalid", http.StatusOK, []byte("New Password is different")},
		{"Invalid new Password to short", "oldpassword", "a1", "a1", http.StatusOK, []byte("Too short")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			form := url.Values{}
			form.Add("oldPassword", test.oldPassword)
			form.Add("newPassword", test.newPassword)
			form.Add("confirmPassword", test.confirmPassword)
			form.Add("csrf_token", csrfToken)

			resCode, _, body := ts.postForm(t, "/user/settings/password", form)

			if resCode != test.wantResCode {
				t.Errorf("want %d but got %d", test.wantResCode, resCode)
			}

			if !bytes.Contains(body, test.wantResBody) {
				t.Errorf("want body %s to contain %q", body, test.wantResBody)
			}
		})
	}
}

func TestLogoutUserHandler(t *testing.T) {
	app := newTestApp(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/")
	csrfToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", csrfToken)

	resCode, _, _ := ts.postForm(t, "/user/logout", form)

	if resCode != http.StatusSeeOther {
		t.Errorf("want %d but got %d", http.StatusSeeOther, resCode)
	}
}
