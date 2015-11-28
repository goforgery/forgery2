package main

import(
    "github.com/goforgery/forgery2"
)

func main() {
    app := f.CreateApp()
    app.Get("/", func(req *f.Request, res *f.Response, next func()) {
        res.Send("Hello world.")
    })
    app.Listen(3000)
}
