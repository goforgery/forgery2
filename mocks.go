package f

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
)

// Mock http.ResponseWriter for testing.
type MockResponseWriter struct {
	error   bool
	headers http.Header
	Buffer  *bytes.Buffer
}

// See http://golang.org/pkg/net/http/#ResponseWriter
func (this MockResponseWriter) Header() http.Header {
	return this.headers
}

// See http://golang.org/pkg/net/http/#ResponseWriter
func (this MockResponseWriter) Write(data []byte) (int, error) {
	if this.error {
		return 0, errors.New("Forced error.")
	}
	return this.Buffer.Write(data)
}

// See http://golang.org/pkg/net/http/#ResponseWriter
func (this MockResponseWriter) WriteHeader(code int) {
	//...
}

func CreateAppMock() (*Application, *Request, *Response, *bytes.Buffer) {
	app := CreateApp()
	req := CreateRequestMock(app)
	res, buf := CreateResponseMock(app, false)
	req.SetResponse(res)
	res.SetRequest(req)
	return app, req, res, buf
}

// Returns a Request that can be used for mocking in tests.
func CreateRequestMock(app *Application) *Request {
	req := &http.Request{
		Header:     http.Header{},
		RequestURI: "/",
		URL:        new(url.URL),
	}
	return CreateRequest(req, app)
}

// Returns a Response that can be used for mocking in tests.
func CreateResponseMock(app *Application, error bool) (*Response, *bytes.Buffer) {
	buf := bytes.NewBufferString("")
	res := MockResponseWriter{error, http.Header{}, buf}
	return CreateResponse(res, app), buf
}
