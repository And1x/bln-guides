/* @@@@@@@@@@@@@@@@@@@@

INFO:
Until I found a better solution, steps to do before testing:
- comment out middleware for user authentication in routes.go (r.Use(app.requireAuth))
- comment out .env path in main init function - comment old path

UNDO:
- handlers.go - createGuideHandler loggedinUserId is set to 1 -- @43 line

@@@@@@@@@@@@@@@@@@@@@ */

package main

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/and1x/bln--h/pkg/mock"
	"github.com/golangcollege/sessions"
)

// newTestApp returns a mocked app instance
// withLog to get more specific error messages
// without logs, outputs in terminal stay clean
func newTestApp(t *testing.T, withLog bool) *app {

	var iL, eL *log.Logger
	if withLog {
		iL = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
		eL = log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		iL = log.New(ioutil.Discard, "", 0)
		eL = log.New(ioutil.Discard, "", 0)
	}

	templateCache, err := createTemplateCache("./../../ui/templates/") // need to change path bc of 'go test'
	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New([]byte(os.Getenv("SESSION_SECRET")))
	session.Lifetime = 8 * time.Hour

	return &app{
		infoLog:       iL,
		errorLog:      eL,
		session:       session,
		templateCache: templateCache,
		guides:        &mock.GuidesModel{},
		users:         &mock.UserModel{},
		lnProvider:    &mock.LNbits{},
	}
}

// testServer type used for DI with methods like GET
type testServer struct {
	*httptest.Server
}

// newTestServer returns a new server
// testServer disabled to follow redirects for the client. Instead it returns the response for that specific req
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	// initalize new cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	// disable redirect to return the lass statuscode eg. 3xx
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// get sends a get request to the test server
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	res, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return res.StatusCode, res.Header, body
}

//postForm send a post request with form data to the test server
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {

	res, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return res.StatusCode, res.Header, body
}

var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

func extractCSRFToken(t *testing.T, body []byte) string {

	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}
