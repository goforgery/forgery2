package f

import (
    "testing"
)

func BenchmarkOneMiddleware(b *testing.B) {
    app := CreateApp()
    app.Use("/foo", func(req *Request, res *Response, next func()) {
        //...
    })
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := CreateRequestMock(app)
        res, _ := CreateResponseMock(app, false)
        req.SetResponse(res)
        res.SetRequest(req)
        req.OriginalUrl = "/foo"
        app.Handle(req, res, 0)
    }
}

func BenchmarkFourMiddlewares(b *testing.B) {
    app := CreateApp()
    app.Use("/foo", func(req *Request, res *Response, next func()) {
        //...
    })
    app.Use("/bar", func(req *Request, res *Response, next func()) {
        //...
    })
    app.Use("/baz", func(req *Request, res *Response, next func()) {
        //...
    })
    app.Use("/qux", func(req *Request, res *Response, next func()) {
        //...
    })
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := CreateRequestMock(app)
        res, _ := CreateResponseMock(app, false)
        req.SetResponse(res)
        res.SetRequest(req)
        req.OriginalUrl = "/qux"
        app.Handle(req, res, 0)
    }
}

func BenchmarkOneRoute(b *testing.B) {
    app := CreateApp()
    app.Get("/foo", func(req *Request, res *Response, next func()) {
        //...
    })
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := CreateRequestMock(app)
        res, _ := CreateResponseMock(app, false)
        req.SetResponse(res)
        res.SetRequest(req)
        req.OriginalUrl = "/foo"
        app.Handle(req, res, 0)
    }
}

func BenchmarkFourRoutes(b *testing.B) {
    app := CreateApp()
    app.Get("/foo", func(req *Request, res *Response, next func()) {
        //...
    })
    app.Get("/bar", func(req *Request, res *Response, next func()) {
        //...
    })
    app.Get("/baz", func(req *Request, res *Response, next func()) {
        //...
    })
    app.Get("/qux", func(req *Request, res *Response, next func()) {
        //...
    })
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := CreateRequestMock(app)
        res, _ := CreateResponseMock(app, false)
        req.SetResponse(res)
        res.SetRequest(req)
        req.OriginalUrl = "/qux"
        app.Handle(req, res, 0)
    }
}

func BenchmarkFourRoutesMatchFirst(b *testing.B) {
    app := CreateApp()
    app.Get("/foo", func(req *Request, res *Response, next func()) {
        //...
    })
    app.Get("/bar", func(req *Request, res *Response, next func()) {
        //...
    })
    app.Get("/baz", func(req *Request, res *Response, next func()) {
        //...
    })
    app.Get("/qux", func(req *Request, res *Response, next func()) {
        //...
    })
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := CreateRequestMock(app)
        res, _ := CreateResponseMock(app, false)
        req.SetResponse(res)
        res.SetRequest(req)
        req.OriginalUrl = "/foo"
        app.Handle(req, res, 0)
    }
}
