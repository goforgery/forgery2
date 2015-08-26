package httphelp

import (
	"net/http"
	"regexp"
	"strings"
)

func Fresh(req http.Header, res http.Header) bool {

	// Defaults.
	etagMatches := true
	notModified := true

	// Fields.
	modifiedSince := req.Get("if-modified-since")
	noneMatch := req.Get("if-none-match")
	lastModified := res.Get("last-modified")
	etag := res.Get("etag")
	cc := req.Get("cache-control")

	// Unconditional request.
	if modifiedSince == "" && noneMatch == "" {
		return false
	}

	// Check for no-cache cache request directive.
	if cc != "" && strings.Index(cc, "no-cache") != -1 {
		return false
	}

	var noneMatchSlice []string

	// Parse if-none-match.
	if noneMatch != "" {
		noneMatchSlice = regexp.MustCompile(" *, *").Split(noneMatch, -1)
	}

	// Search if-none-match
	if len(noneMatchSlice) > 0 {
		for _, t := range noneMatchSlice {
			if t == etag {
				etagMatches = true
			}
		}
		if etagMatches == false && noneMatchSlice[0] == "*" {
			etagMatches = true
		}
	}

	// if-modified-since
	if modifiedSince != "" {
		// modifiedSince = new Date(modifiedSince);
		// lastModified = new Date(lastModified);
		notModified = lastModified <= modifiedSince
	}

	return !!(etagMatches && notModified)
}
