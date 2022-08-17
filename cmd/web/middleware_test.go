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
			res := httptest.NewRecorder()
			req := httptest.NewRequest(test.reqMethod, test.reqURI, nil)

			// simple Handlerfunc to satisfiy Signature from logging method
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			})

			app.logging(next).ServeHTTP(res, req)

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

// its working but the panic gets logged in terminal which makes it look messy
// func TestRecoverPanic(t *testing.T) {

// 	app := &app{
// 		errorLog: log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile), // errorLog is called by serverError hence its needed in app struct
// 	}

// 	// randomHandler initiates a panic - which should get recovered
// 	randomHandler := func(http.ResponseWriter, *http.Request) { panic("Help, a panic happend!") }
// 	panicHandler := app.recoverPanic(http.HandlerFunc(randomHandler))

// 	req := httptest.NewRequest("GET", "/", nil)
// 	res := httptest.NewRecorder()

// 	panicHandler.ServeHTTP(res, req)

// 	assert.Equal(t, http.StatusInternalServerError, res.Code)
// 	assert.Equal(t, res.Header().Values("Connection"), []string{"close"}) // check if in Headermap k="Connection" v="close" is set after panic recovery
// }
