package f

import (
	"bytes"
	"github.com/goforgery/mustache"
	. "github.com/ricallinson/simplebdd"
	"testing"
)

func TestResponseRender(t *testing.T) {

	var app *Application
	var req *Request
	var res *Response
	var buf *bytes.Buffer

	BeforeEach(func() {
		app, req, res, buf = CreateAppMock()
		app.Set("views", "./fixtures/views")
	})

	Describe("Render()", func() {

		It("should return HTML", func() {
			app.Engine(".html", mustache.Create())
			app.Get("/foo", func(req *Request, res *Response, next func()) {
				res.Render("index.html", map[string]string{"title": "Mu"})
			})
			req.Method = "GET"
			req.OriginalUrl = "/foo"
			app.Handle(req, res, 0)
			html := buf.String()
			AssertEqual(html, "<h1>Mu</h1>")
		})

		It("should return error", func() {
			app.Engine(".md", mustache.Create())
			app.Get("/foo", func(req *Request, res *Response, next func()) {
				res.Render("index.html", map[string]string{"title": "Mu"})
			})
			req.Method = "GET"
			req.OriginalUrl = "/foo"
			app.Handle(req, res, 0)
			html := buf.String()
			AssertEqual(html, "View engine for '.html' not found.")
		})
	})
}
