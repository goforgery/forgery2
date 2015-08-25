package f

import (
	"github.com/ricallinson/httputils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"path/filepath"
)

const (
	TRUE  = "true"
	FALSE = "false"
)

type Application struct {
	Env      string
	handlers []handler
	// Application local variables are provided to all templates rendered within the application.
	// This is useful for providing helper functions to templates, as well as app-level data.
	Locals map[string]string
	// The Router middleware function.
	Router *Router
	// Has the Router been added to stackr.
	usedRouter bool
	// Stores the applications settings.
	settings map[string]string
	// The rendering engines assigned.
	engines map[string]Renderer
}

type handler struct {
	Route  string
	Handle func(*Request, *Response, func())
}

type Handle func(*Request, *Response, func())

// Create a new forgery application.
// * "env" Environment mode, defaults to $GO_ENV or "development"
// * "trust proxy" Enables reverse proxy support, disabled by default
// * "jsonp callback name" Changes the default callback name of "?callback="
// * "json spaces" JSON response spaces for formatting, defaults to 2 in development, 0 in production
// * "case sensitive routing" Enable case sensitivity, disabled by default, treating "/Foo" and "/foo" as the same
// * "strict routing" Enable strict routing, by default "/foo" and "/foo/" are treated the same by the router
// * X "view cache" Enables view template compilation caching, enabled in production by default
// * "view engine" The default engine extension to use when omitted
// * "views" The view directory path, defaulting to "./views"
func CreateApp() *Application {
	this := &Application{}
	this.Env = os.Getenv("GO_ENV")
	if this.Env == "" {
		this.Env = "development"
	}
	this.Locals = map[string]string{}
	this.Router = &Router{}
	this.settings = map[string]string{}
	this.engines = map[string]Renderer{}
	this.defaultConfiguration()
	return this
}

// Initialize application configuration.
func (this *Application) defaultConfiguration() {
	cwd, err := os.Getwd()
	if err != nil {
		panic("Cannot get current working directory!")
	}
	// default settings
	this.Enable("x-powered-by")
	this.Enable("etag")
	this.Set("env", os.Getenv("GO_ENV"))
	if this.Get("env") == "" {
		this.Set("env", "development")
	}
	// default configuration
	this.Configure(func() {
		this.Set("subdomain offset", "2")
		this.Set("views", filepath.Join(cwd, "views"))
		this.Set("jsonp callback name", "callback")
		this.Set("app path", "/")
	})
	this.Configure("development", func() {
		this.Set("json spaces", "  ")
	})
	this.Configure("production", func() {
		this.Enable("view cache")
	})
}

// Configure callback for zero or more envs,
// when no `env` is specified that callback will
// be invoked for all environments. Any combination
// can be used multiple times, in any order desired.
//
// Examples:
//
//    app.Configure(func() {
//      // executed for all envs
//    })
//
//    app.Configure("stage", func() {
//      // executed staging env
//    })
//
//    app.Configure("stage", "production", func() {
//      // executed for stage and production
//    })
//
// Note:
//
// These callbacks are invoked immediately, and
// are effectively sugar for the following:
//
// var env = os.Getenv("GO_ENV")
//
// switch (env) {
// case 'development':
// ...
// case 'stage':
// ...
// case 'production':
// ...
// }
func (this *Application) Configure(i ...interface{}) {
	var envs []string
	var fn func()
	// Look at the given inputs.
	for _, t := range i {
		switch t.(type) {
		case string:
			envs = append(envs, t.(string))
		case func():
			fn = t.(func())
		}
	}
	// If there are no envs call the func and return.
	if len(envs) == 0 {
		fn()
		return
	}
	// Loop over the envs until a match is found.
	// Then call the function.
	for _, e := range envs {
		if e == this.Get("env") {
			fn()
			return
		}
	}
}

// Returns the root of this app.
func (this *Application) Path() string {
	return this.Get("app path")
}

// Assigns setting "name" to "value".
func (this *Application) Set(n string, v ...string) string {
	if len(v) == 0 {
		return this.settings[n]
	}
	this.settings[n] = v[0]
	return v[0]
}

// Get setting "name" value.
// or;
// Provides the routing functionality for GET requests to the given "path".
func (this *Application) Get(path string, fn ...func(*Request, *Response, func())) string {
	// If there is no function then this is really a call to .Set()
	if len(fn) == 0 {
		return this.Set(path)
	}
	//Otherwise it's a call to .Verb()
	this.Verb("GET", path, fn...)
	return ""
}

// Set setting "name" to "true".
func (this *Application) Enable(n string) {
	this.Set(n, TRUE)
}

// Set setting "name" to "false".
func (this *Application) Disable(n string) {
	this.Set(n, FALSE)
}

// Check if setting "name" is enabled.
func (this *Application) Enabled(n string) bool {
	return this.Get(n) == TRUE
}

// Check if setting "name" is disabled.
func (this *Application) Disabled(n string) bool {
	return this.Get(n) == FALSE
}

// Register the given template engine callback as ext.
func (this *Application) Engine(ext string, fn Renderer) {
	this.engines[ext] = fn
}

// Render a "view" responding with the rendered string.
// This is the app-level variant of "res.render()", and otherwise behaves the same way.
func (this *Application) Render(view string, i ...interface{}) (string, error) {
	ext := filepath.Ext(view)
	if _, ok := this.engines[ext]; ok == false {
		return "", errors.New("Engine not found.")
	}
	file := filepath.Join(this.Get("views"), view)
	if _, err := os.Stat(file); err != nil || os.IsNotExist(err) {
		return "", errors.New("Failed to lookup view '" + file + "'")
	}
	i = append(i, this.Locals)
	t, err := this.engines[ext].Render(file, i...)
	if err != nil {
		return "", errors.New("Problem rendering view.")
	}
	return t, nil
}

// This method functions just like the app.Verb(verb, ...) method, however it matches all HTTP verbs.
func (this *Application) All(path string, fn ...func(*Request, *Response, func())) {
	for _, verb := range httphelp.Methods {
		this.Verb(verb, path, fn...)
	}
}

// The method provides the routing functionality in Forgery, where "verb" is one of the HTTP verbs,
// such as app.Verb("post", ...). Multiple callbacks may be given, all are treated equally,
// and behave just like middleware, with the one exception that these callbacks may invoke
// next('route') to bypass the remaining route callback(s). This mechanism can be used to perform
// pre-conditions on a route then pass control to subsequent routes when there is no reason to
// proceed with the route matched.
func (this *Application) Verb(verb string, path string, funcs ...func(*Request, *Response, func())) {
	if this.usedRouter == false {
		this.Router.CaseSensitive = this.Enabled("case sensitive routing")
		this.Router.Strict = this.Enabled("strict routing")
		this.Use(this.Router.Middleware(this))
		this.usedRouter = true
	}
	this.Router.AddRoute(verb, path, funcs...)
}

// Map logic to route parameters. For example when ":user" is
// present in a route path you may map user loading logic to
// automatically provide req.Map["user"] to the route, or perform
// validations on the parameter input.
func (this *Application) Param(p string, fn func(*Request, *Response, func())) {
	this.Router.AddParamFunc(p, fn)
}

// Utilize the given middleware `Handle` to the given `route`,
// defaulting to _/_. This "route" is the mount-point for the
// middleware, when given a value other than _/_ the middleware
// is only effective when that segment is present in the request's
// pathname.
//
// For example if we were to mount a function at _/admin_, it would
// be invoked on _/admin_, and _/admin/settings_, however it would
// not be invoked for _/_, or _/posts_.
//
// Examples:
//
//    var app = stackr.CreateServer();
//    app.Use(stackr.Favicon())
//    app.Use(stackr.Logger())
//    app.Use("/public", stackr.Static())
//
// If we wanted to prefix static files with _/public_, we could
// "mount" the `Static()` middleware:
//
//    app.Use("/public", stackr.Static(stackr.OptStatic{Root: "./static_files"}))
//
// This api is chainable, so the following is valid:
//
//    stackr.CreateServer().Use(stackr.Favicon()).Listen(3000);
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

// Handle server requests, punting them down the middleware stack.
// Note: this is a recursive function.
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
		index++                      // increment the index by 1
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
	this.Handle(CreateRequest(req), CreateResponse(res), 0)
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
