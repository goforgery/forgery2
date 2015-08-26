# Httphelp

HTTP Helper functions, Vars and Types for Go.

	go get github.com/ricallinson/httphelp

* `.Methods` A slice of all HTTP methods.
* `.StatusCodes` A map of all HTTP status codes and their textual description.
* `.ParseAccept()` A function to parse an accept header string returning a sorted slice of types.