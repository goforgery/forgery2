package f

/*
   The method provides the routing functionality for POST requests to the given "path".
*/
func (this *Application) Post(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("post", path, fn...)
}

/*
   The method provides the routing functionality for PUT requests to the given "path".
*/
func (this *Application) Put(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("put", path, fn...)
}

/*
   The method provides the routing functionality for HEAD requests to the given "path".
*/
func (this *Application) Head(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("head", path, fn...)
}

/*
   The method provides the routing functionality for DELETE requests to the given "path".
*/
func (this *Application) Delete(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("delete", path, fn...)
}

/*
   The method provides the routing functionality for OPTIONS requests to the given "path".
*/
func (this *Application) Options(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("options", path, fn...)
}

/*
   The method provides the routing functionality for TRACE requests to the given "path".
*/
func (this *Application) Trace(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("trace", path, fn...)
}

/*
   The method provides the routing functionality for COPY requests to the given "path".
*/
func (this *Application) Copy(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("copy", path, fn...)
}

/*
   The method provides the routing functionality for LOCK requests to the given "path".
*/
func (this *Application) Lock(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("lock", path, fn...)
}

/*
   The method provides the routing functionality for MKCOL requests to the given "path".
*/
func (this *Application) Mkcol(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("mkcol", path, fn...)
}

/*
   The method provides the routing functionality for MOVE requests to the given "path".
*/
func (this *Application) Move(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("move", path, fn...)
}

/*
   The method provides the routing functionality for PROPFIND requests to the given "path".
*/
func (this *Application) Propfind(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("propfind", path, fn...)
}

/*
   The method provides the routing functionality for PROPPATCH requests to the given "path".
*/
func (this *Application) Proppatch(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("proppatch", path, fn...)
}

/*
   The method provides the routing functionality for UNLOCK requests to the given "path".
*/
func (this *Application) Unlock(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("unlock", path, fn...)
}

/*
   The method provides the routing functionality for REPORT requests to the given "path".
*/
func (this *Application) Report(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("report", path, fn...)
}

/*
   The method provides the routing functionality for MKACTIVITY requests to the given "path".
*/
func (this *Application) Mkactivity(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("mkactivity", path, fn...)
}

/*
   The method provides the routing functionality for CHECKOUT requests to the given "path".
*/
func (this *Application) Checkout(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("checkout", path, fn...)
}

/*
   The method provides the routing functionality for MERGE requests to the given "path".
*/
func (this *Application) Merge(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("merge", path, fn...)
}

/*
   The method provides the routing functionality for M-SEARCH requests to the given "path".
*/
func (this *Application) Msearch(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("m-search", path, fn...)
}

/*
   The method provides the routing functionality for NOTIFY requests to the given "path".
*/
func (this *Application) Notify(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("notify", path, fn...)
}

/*
   The method provides the routing functionality for SUBSCRIBE requests to the given "path".
*/
func (this *Application) Subscribe(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("subscribe", path, fn...)
}

/*
   The method provides the routing functionality for UNSUBSCRIBE requests to the given "path".
*/
func (this *Application) Unsubscribe(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("unsubscribe", path, fn...)
}

/*
   The method provides the routing functionality for PATCH requests to the given "path".
*/
func (this *Application) Patch(path string, fn ...func(*Request, *Response, func())) {
    this.Verb("patch", path, fn...)
}
