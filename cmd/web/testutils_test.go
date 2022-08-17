package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/and1x/bln--h/pkg/models/mock"
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

	return &app{
		infoLog:  iL,
		errorLog: eL,
		guides:   &mock.GuidesModel{},
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

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

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
