package f

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Application struct {
	Env      string
	handlers []handler
}

type handler struct {
	Route  string
	Handle func(*Request, *Response, func())
}

type Handle func(*Request, *Response, func())

type HandleError func(string, *Request, *Response, func())

// Create a new forgery application.
func CreateApp() *Application {
	this := &Application{}
	this.Env = os.Getenv("GO_ENV")
	if this.Env == "" {
		this.Env = "development"
	}
	return this
}

func (this *Application) Use(in ...interface{}) *Application {
	var route string
	var handle func(*Request, *Response, func())
	// Work out what we were given.
	for _, i := range in {
		switch i.(type) {
		case string:
			route = i.(string)
		case func(*Request, *Response, func()):
			handle = i.(func(*Request, *Response, func()))
		default:
			panic("stackr: Go home handler, you're drunk!")
		}
	}
	// If the route is empty make it "/".
	if len(route) == 0 {
		route = "/"
	}
	// Strip trailing slash
	if size := len(route); size > 1 && route[size-1] == '/' {
		route = route[:size-1]
	}
	this.handlers = append(this.handlers, handler{
		Route:  strings.ToLower(route),
		Handle: handle,
	})
	return this
}

func (this *Application) Handle(req *Request, res *Response, index int) {
	// For each call to Handle we want to catch anything that panics unless in development mode.
	defer func() {
		if this.Env != "development" {
			err := recover()
			if err != nil {
				res.Error = errors.New(fmt.Sprint(err))
			}
		}
	}()
	// If the response has been closed return.
	if res.Closed == true {
		return
	}
	var layer handler
	// Do we have another layer to use?
	if index >= len(this.handlers) {
		layer = handler{} // no
	} else {
		layer = this.handlers[index] // yes
		index++                   // increment the index by 1
	}
	// If there are no more layers and no headers have been sent return a 404.
	if layer.Handle == nil && res.HeaderSent == false {
		res.StatusCode = 404
		res.SetHeader("Content-Type", "text/plain")
		if req.Method == "HEAD" {
			res.End("")
			return
		}
		res.End("Cannot " + fmt.Sprint(req.Method) + " " + fmt.Sprint(req.OriginalUrl))
		return
	}
	// If there are no more layers and headers were sent then we are done so just return.
	if layer.Handle == nil {
		return
	}
	// Otherwise call the layer Handler.
	if strings.Contains(req.OriginalUrl, layer.Route) {
		// Set the value of Url to the portion after the matched layer.Route
		req.Url = strings.TrimPrefix(req.OriginalUrl, layer.Route)
		// Set the matched portion of the Url.
		req.MatchedUrl = layer.Route
		// Call the middleware function.
		layer.Handle(req, res, func() {
			// The value of next is a function that calls this function again, passing the index value.
			this.Handle(req, res, index)
		})
	} else {
		// Call this function, passing the index value.
		this.Handle(req, res, index)
	}
}

// ServeHTTP calls .Handle(req, res).
func (this *Application) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//...
}

// Listen for connections on HTTP.
func (this *Application) Listen(port int) {
	address := ":" + fmt.Sprint(port)
	log.Fatal(http.ListenAndServe(address, this))
}

// Listen for connections on HTTPS.
func (this *Application) ListenTLS(port int, certFile string, keyFile string) {
	address := ":" + fmt.Sprint(port)
	log.Fatal(http.ListenAndServeTLS(address, certFile, keyFile, this))
}
