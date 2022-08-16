package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Testing the logging method (to learn, dont see much value in testing the logger)
// instead of writing to os.Stdout, a buffer is used
// check if buffer len is > 0
// missing parts: log.Ldate, log.Ltime and RemoteAddr
func TestLogging(t *testing.T) {

	// for testing purpose: use buffer to test instead of os.Stdout
	buf := &bytes.Buffer{}

	app := &app{
		infoLog: log.New(buf, "INFO:\t", log.Ldate|log.Ltime),
	}

	tests := []struct {
		name      string
		reqMethod string
		reqURI    string
	}{
		{
			name:      "basic GET request",
			reqMethod: "GET",
			reqURI:    "/test",
		},
		{
			name:      "basic POST request",
			reqMethod: "POST",
			reqURI:    "/test",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// record Response from a request
			rr := httptest.NewRecorder()
			r := httptest.NewRequest(test.reqMethod, test.reqURI, nil)

			// simple Handlerfunc to satisfiy Signature from logging method
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			})

			app.logging(next).ServeHTTP(rr, r)

			if buf.Len() <= 0 {
				t.Errorf("Nothing written to buffer")
			}

			if !bytes.Contains(buf.Bytes(), []byte(test.reqMethod)) || !bytes.Contains(buf.Bytes(), []byte(test.reqURI)) {
				t.Errorf("log expected to contain %s and %s but log contains %s", test.reqMethod, test.reqURI, buf.String())
			}

			buf.Reset()
		})
	}
}
