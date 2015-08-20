package f

import (
	"errors"
	"fmt"
	"net/http"
)

// Response represents the response from an HTTP request.
type Response struct {
	// The standard http.ResponseWriter interface.
	Writer http.ResponseWriter
	// Ture if headers have been sent.
	HeaderSent bool
	// The HTTP status code to be return.
	StatusCode int
	// True if .End() has been called.
	Closed bool
	// Error. Populated by anything that wants to trigger an error.
	Error error
	// Events
	events map[string][]func()
}

/*
   Create a Mock http.ResponseWriter for testing.
*/

type MockResponseWriter struct {
	error   bool
	headers http.Header
	Written []byte
}

func (this *MockResponseWriter) Header() http.Header {
	return this.headers
}

func (this *MockResponseWriter) Write(data []byte) (int, error) {
	if this.error {
		return 0, errors.New("")
	}
	this.Written = data
	return len(data), nil
}

func (this *MockResponseWriter) WriteHeader(code int) {
	return
}

// Returns a new Response.
func CreateResponse(writer http.ResponseWriter) *Response {
	this := &Response{writer, false, 200, false, nil, map[string][]func(){}}
	return this
}

// Returns a Response that can be used for mocking in tests.
func CreateResponseMock(error bool) *Response {
	res := &MockResponseWriter{error, make(http.Header), []byte{}}
	return CreateResponse(res)
}

// Register a listener function for an event.
func (this *Response) On(event string, fn func()) {
	this.events[event] = append(this.events[event], fn)
}

// Emit an event calling all registered listeners.
func (this *Response) Emit(event string) {
	e, ok := this.events[event]
	if ok {
		for _, fn := range e {
			fn()
		}
	}
}

// Set a map of headers, calls SetHeader() for each item.
func (this *Response) SetHeaders(headers map[string]string) bool {
	for key, value := range headers {
		if this.SetHeader(key, value) == false {
			return false
		}
	}
	return true
}

// Set a single header.
// Note: all keys are converted to lower case.
func (this *Response) SetHeader(key string, value string) bool {
	//  If the headers have been sent nothing can be done so return false.
	if this.HeaderSent == true {
		return false
	}
	// http://www.w3.org/Protocols/rfc2616/rfc2616-sec4.html#sec4.2
	// Message headers are case-insensitive.
	if len(value) > 0 {
		this.Writer.Header().Set(key, value)
	}
	return true
}

// Remove the named header.
func (this *Response) RemoveHeader(key string) {
	this.Writer.Header().Del(key)
}

// Write any headers set to the client.
func (this *Response) writeHeaders() {
	this.Emit("header")
	this.HeaderSent = true
	this.Writer.WriteHeader(this.StatusCode)
}

/*
   Write bytes to the client.
*/
func (this *Response) WriteBytes(data []byte) bool {
	if this.HeaderSent == false {
		this.writeHeaders()
	}
	writen, err := this.Writer.Write(data)
	if err != nil {
		return false
	}
	return writen == len(data)
}

// Write data to the client.
func (this *Response) Write(data string) bool {
	if this.HeaderSent == false {
		this.writeHeaders()
	}
	if len(data) == 0 {
		return true
	}
	writen, err := fmt.Fprint(this.Writer, data)
	if err != nil {
		return false
	}
	return writen == len(data)
}

// Close the connection to the client.
func (this *Response) End(data string) bool {
	status := true
	status = this.Write(data)
	this.Closed = true
	return status
}
