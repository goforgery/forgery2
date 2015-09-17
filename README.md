# forgery2

[![Build Status](https://secure.travis-ci.org/goforgery/forgery2.png?branch=master)](http://travis-ci.org/goforgery/forgery2)

__CURRENTLY UNSTABLE VERSION__

Forgery is a minimal and flexible golang web application framework, providing a robust set of features for building single and multi-page, web applications.

    package main

    import(
        "github.com/goforgery/forgery2"
    )

    func init() {
        app := f.CreateApp()
        app.Get("/", func(req *f.Request, res *f.Response, next func()) {
            res.Send("Hello world.")
        })
        app.Listen(3000)
    }

* Robust routing
* HTTP helpers (redirection, caching, etc)
* View system supporting 1 template engine (hopefully more will come)
* Content negotiation
* Focus on high performance
* Environment based configuration
* High test coverage

## Testing

    go test ./...

## Code Coverage

    go test -coverprofile=coverage.out; go tool cover -html=coverage.out -o=coverage.html
    open coverage.html

## Benchmark

    go test -bench=. -run=BENCH_ONLY
