# forgery2

[![Build Status](https://secure.travis-ci.org/ricallinson/forgery2.png?branch=master)](http://travis-ci.org/ricallinson/forgery2)

## Testing

    go test ./...

## Code Coverage

    go test -coverprofile=coverage.out; go tool cover -html=coverage.out -o=coverage.html
    open coverage.html
