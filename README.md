# forgery2

[![Build Status](https://secure.travis-ci.org/goforgery/forgery2.png?branch=master)](http://travis-ci.org/goforgery/forgery2)

__UNSTABLE VERSION 2.0__: The current stable version is here [forgery](https://github.com/ricallinson/forgery)

Forgery is a minimal and flexible golang web application framework, providing a robust set of features for building single and multi-page, web applications.

* Robust routing
* HTTP helpers (redirection, caching, etc)
* View system supporting 1 template engine (hopefully more will come)
* Content negotiation
* Focus on high performance
* Environment based configuration
* High test coverage

## Install

    go get github.com/goforgery/forgery2

## Use

Starts a web server.

```javascript
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
```

## Testing

    go test ./...

## Code Coverage

    go test -coverprofile=coverage.out; go tool cover -html=coverage.out -o=coverage.html
    open coverage.html

## Benchmark

    go test -bench=. -run=BENCH_ONLY

## Notes

This project started out as a clone of the superb Node.js library [Express](http://expressjs.com/).
