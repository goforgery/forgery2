package f

import (
	"encoding/json"
	"github.com/goforgery/httphelp"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// A Request represents an HTTP request received by the server.
type Request struct {
	// The standard http.Request type
	*http.Request
	// The value of .URL.RequestURI() for easy access.
	// Note: this value may be changed by the Stackr.Handle() function.
	Url string
	// Set to the vlue of the matched portion of the .URL.RequestURI()
	MatchedUrl string
	// The value of .URL.RequestURI() for easy access.
	// Note: this value should NEVER be changed.
	OriginalUrl string
	// This property is a map containing the parsed request body.
	// This feature is provided by the bodyParser() middleware, though other body
	// parsing middleware may follow this convention as well.
	// This property defaults to {} when bodyParser() is used.
	bodies map[string]string
	// This property is a map containing the parsed query-string, defaulting to {}.
	queries map[string]string
	// This property is a map of the files uploaded. This feature is provided
	// by the bodyParser() middleware, though other body parsing middleware may
	// follow this convention as well. This property defaults to {} when bodyParser() is used.
	files map[string]interface{}
	// Holds general key/values over the lifetime of the request.
	Map map[string]interface{}
	// The forgery.Response matched to this request.
	res *Response
	// The application server.
	app *Application
	// Return the remote address, or when "trust proxy" is enabled - the upstream address.
	ip string
	// When "trust proxy" is `true`, parse the "X-Forwarded-For" ip address list and return a slice,
	// otherwise an empty slice is returned. For example if the value were "client, proxy1, proxy2"
	// you would receive the slice {"client", "proxy1", "proxy2"} where "proxy2" is the furthest down-stream.
	ips []string
	// Returns the request URL pathname.
	Path string
	// Check if the request was issued with the "X-Requested-With" header field set to "XMLHttpRequest" (jQuery etc).
	Xhr bool
	// Return the protocol string "http" or "https" when requested with TLS.
	// When the "trust proxy" setting is enabled the "X-Forwarded-Proto" header field will be trusted.
	// If you're running behind a reverse proxy that supplies https for you this may be enabled.
	Protocol string
	// Check if a TLS connection is established. This is a short-hand for: "https" == req.Protocol
	Secure bool
	// This property is a map containing properties mapped to the named route "parameters".
	// For example if you have the route "/user/:name", then the "name" property is available
	// to you as req.params["name"]. This object defaults to {}.
	params map[string]string
	// The currently matched Route containing several properties such as the
	// route's original path string, the regexp generated, and so on.
	Route *Route
	//
	accepted []string
	//
	acceptedLanguages []string
	//
	acceptedCharsets []string
}

// Returns a new Request.
func CreateRequest(raw *http.Request, app *Application) *Request {
	this := &Request{}
	this.Request = raw
	this.Url = raw.URL.RequestURI()
	this.OriginalUrl = raw.URL.RequestURI()
	// Helpers for standard headers.
	this.Path = this.URL.Path
	this.Xhr = this.Header.Get("X-Requested-With") == "XMLHttpRequest"
	this.Protocol = this.URL.Scheme
	this.Secure = this.Protocol == "https"
	this.app = app
	// This could have been set by middleware so check if it's empty.
	if this.Map == nil {
		this.Map = map[string]interface{}{}
	}
	return this
}

// Set the Response this Response will use.
func (this *Request) SetResponse(res *Response) {
	this.res = res
}

// Returns the requesting IP.
func (this *Request) Ip() string {
	return this.Ips()[0]
}

// Returns the requesting IP chain.
func (this *Request) Ips() []string {
	if this.ips == nil {
		if t, v := this.app.Get("trust proxy"), this.Header.Get("X-Forwarded-For"); len(t) > 0 && len(v) > 0 {
			s := regexp.MustCompile(" *, *").Split(v, -1)
			this.ip = s[0]
			this.ips = s
		} else {
			this.ip = this.RemoteAddr
			this.ips = []string{}
		}
	}
	return this.ips
}

// Return the value for the given key if found in the request body.
func (this *Request) Body(key string) string {
	return this.Bodies()[key]
}

// Return a map of values passed as the request body.
func (this *Request) Bodies() map[string]string {
	if this.bodies == nil {
		this.bodies = map[string]string{}
		// Parse the form data if there was any.
		this.ParseForm()
		// Copy the first item for each key from this.Request.PostForm[] into this.body
		// The alternative is to access this.PostForm[] from this.Param() like so;
		//     v, ok = this.PostForm[n]
		for k, v := range this.PostForm {
			this.bodies[k] = v[0]
		}
	}
	return this.bodies
}

// Return the value for the given key if found in the request query.
func (this *Request) Query(key string) string {
	return this.Queries()[key]
}

// Return a map of values passed as the request query.
func (this *Request) Queries() map[string]string {
	if this.queries == nil {
		this.queries = map[string]string{}
		// Copy the first item for each key from this.Request.URL.Query() into this.query
		// This has a performance impact so is it worth it?
		// The alternative is to access this.URL.Query() from this.Param() like so;
		//     v, ok = this.URL.Query().Get(n)
		for k, v := range this.URL.Query() {
			this.queries[k] = v[0]
		}
	}
	return this.queries
}

// Return the value of param "name" when present. Lookup is performed in the following order:
//
// * Params
// * Body
// * Query
//
// Direct access to req.body, req.params, and req.query should be favored for clarity -
// unless you truly accept input from each object.
func (this *Request) Param(n string) string {
	var v string
	var ok bool
	v, ok = this.Params()[n]
	if ok {
		return v
	}
	v, ok = this.Bodies()[n]
	if ok {
		return v
	}
	v, ok = this.Queries()[n]
	if ok {
		return v
	}
	return ""
}

// Returns a map containing properties mapped to the named route "parameters".
// If a map is passed in it will replace the current map.
func (this *Request) Params(p ...map[string]string) map[string]string {
	if this.params == nil {
		this.params = map[string]string{}
	}
	if len(p) > 0 {
		this.params = p[0]
	}
	return this.params
}

// Return the value for the given key if found in the request files.
func (this *Request) File(key string) interface{} {
	return this.Files()[key]
}

// Return a map of values passed as the request files.
func (this *Request) Files() map[string]interface{} {
	if this.files == nil {
		this.files = map[string]interface{}{}
	}
	return this.files
}

// Contains the cookies sent by the user-agent.
func (this *Request) Cookie(n string, i ...interface{}) string {
	cookie, err := this.Request.Cookie(n)
	if err != nil {
		return ""
	}
	// Unescape the cookie.
	var v string
	var e error
	v, e = url.QueryUnescape(cookie.Value)
	v, e = Decode(v)
	if e != nil {
		return ""
	}
	// If no interface was given then just return the Value.
	if len(i) == 0 {
		return v
	}
	json.Unmarshal([]byte(v), i[0])
	return v
}

// Contains the signed cookies sent by the user-agent, unsigned and ready for use.
// Signed cookies are accessed by a different function to show developer intent, otherwise a
// malicious attack could be placed on `req.Cookie` values which are easy to spoof.
// Note that signing a cookie does not mean it is "hidden" nor encrypted, this simply
// prevents tampering as the secret used to sign is private.
func (this *Request) SignedCookie(n string, i ...interface{}) string {
	secret := this.app.Get("secret")
	if secret == "" {
		panic("f.Get(\"secret\") required for signed cookies")
	}
	v := this.Cookie(n)
	// Decrypt the cookie string value.
	v = Unsign(v, secret)
	// If no interface was given then just return the Value.
	if len(i) == 0 {
		return v
	}
	json.Unmarshal([]byte(v), i[0])
	return v
}

// Get the case-insensitive request header field.
func (this *Request) Get(f string) string {
	// Possible future bug.
	// http://www.w3.org/Protocols/rfc2616/rfc2616-sec4.html#sec4.2
	// Message headers are case-insensitive.
	return this.Header.Get(f)
}

// Check if the incoming request contains the "Content-Type" header field, and it matches the give mime "type".
func (this *Request) Is(t string) bool {
	h := this.Get("Content-Type")
	return strings.ToLower(h) == strings.ToLower(t)
}

// WARNING: Not complete!
// Check if the request is fresh - aka Last-Modified and/or the ETag still match,
// indicating that the resource is "fresh".
func (this *Request) Fresh() bool {
	if m := this.Method; m != "GET" && m != "HEAD" {
		return false
	}
	if s := this.res.StatusCode; (s >= 200 && s < 300) || 304 == s {
		return false
	}
	return httphelp.Fresh(this.Header, this.res.Writer.Header())
}

// Check if the request is stale - aka Last-Modified and/or the ETag do not match,
// indicating that the resource is "stale".
func (this *Request) Stale() bool {
	return this.Fresh() == false
}

// Return subdomains as a slice of strings.
func (this *Request) Subdomains() []string {
	offset, _ := strconv.Atoi(this.app.Get("subdomain offset"))
	subs := strings.Split(this.Host, ".")
	end := len(subs) - offset
	if end > 0 {
		return StringSliceReverse(subs[:end])
	}
	return []string{}
}

// Check if the given type is acceptable, returning true or false -
// in which case you should respond with 406 "Not Acceptable".
func (this *Request) Accepts(t string) bool {
	for _, v := range this.Accepted() {
		if strings.ToLower(t) == v {
			return true
		}
	}
	return false
}

// Check if the given "charset" is acceptable.
func (this *Request) AcceptsCharset(c string) bool {
	return strings.Index(strings.ToLower(this.Get("Accept-Charset")), strings.ToLower(c)) > -1
}

// Check if the given "lang" is acceptable.
func (this *Request) AcceptsLanguage(l string) bool {
	return strings.Index(strings.ToLower(this.Get("Accept-Language")), strings.ToLower(l)) > -1
}

// Return an slice of Accepted media types ordered from highest quality to lowest.
func (this *Request) Accepted() []string {
	if this.accepted == nil {
		a := this.Header.Get("Accept")
		this.accepted = this.processAccepted(a)
	}
	return this.accepted
}

// Return an slice of Accepted languages ordered from highest quality to lowest.
func (this *Request) AcceptedLanguages() []string {
	if this.acceptedLanguages == nil {
		a := this.Header.Get("Accept-Language")
		this.acceptedLanguages = this.processAccepted(a)
	}
	return this.acceptedLanguages
}

// Return an slice of Accepted charsets ordered from highest quality to lowest.
func (this *Request) AcceptedCharsets() []string {
	if this.acceptedCharsets == nil {
		a := this.Header.Get("Accept-Charset")
		this.acceptedCharsets = this.processAccepted(a)
	}
	return this.acceptedCharsets
}

// Return an slice of "accepted" ordered from highest quality to lowest.
func (this *Request) processAccepted(a string) (list []string) {
	for _, accept := range httphelp.ParseAccept(a) {
		if len(accept.SubType) > 0 {
			list = append(list, accept.Type+"/"+accept.SubType)
		} else {
			list = append(list, accept.Type)
		}
	}
	return
}
