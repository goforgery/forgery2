package f

import (
	"net/http"
	"net/url"
)

// A Request represents an HTTP request received by the server.
type Request struct {
	// The standard http.Request type
	*http.Request
	// The value of .URL.RequestURI() for easy access.
	// Note: this value may be changed by the Stackr.Handle() function.
	Url string
	// Set to the vlue of the matched portion of the .URL.RequestURI()
	MatchedUrl string
	// The value of .URL.RequestURI() for easy access.
	// Note: this value should NEVER be changed.
	OriginalUrl string
	// This property is a map containing the parsed request body.
	// This feature is provided by the bodyParser() middleware, though other body
	// parsing middleware may follow this convention as well.
	// This property defaults to {} when bodyParser() is used.
	Body map[string]string
	// This property is a map containing the parsed query-string, defaulting to {}.
	Query map[string]string
	// This property is a map of the files uploaded. This feature is provided
	// by the bodyParser() middleware, though other body parsing middleware may
	// follow this convention as well. This property defaults to {} when bodyParser() is used.
	Files map[string]interface{}
	// Holds general key/values over the lifetime of the request.
	Map map[string]interface{}
}

// Returns a new Request.
func CreateRequest(raw *http.Request) *Request {
	this := &Request{
		Request:     raw,
		Url:         raw.URL.RequestURI(),
		OriginalUrl: raw.URL.RequestURI(),
	}
	if this.Map == nil {
		this.Map = map[string]interface{}{}
	}
	return this
}

// Returns a Request that can be used for mocking in tests.
func CreateRequestMock() *Request {
	req := &http.Request{
		Header:     http.Header{},
		RequestURI: "/",
		URL:        new(url.URL),
	}
	return CreateRequest(req)
}
