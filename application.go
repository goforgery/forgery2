package f2

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Env   string
	stack []handler
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
	return this
}

func (this *Application) Handle(req *Request, res *Response, index int) {
	//...
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
