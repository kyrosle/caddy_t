# Main

Begin from the module `caddyhttp` , base of the thinking of web services are relevant with `http`. That's mean we can start from building request and response, and then enter a further step.



# Http Building

use `http.Request` and `http.Response` in go standard module

__use module__ :
* `net/http`

__details__ :


(modules/caddyhttp/caddyhttp.go)

`RequestMatcher` is to match a request
```go
type RequestMatcher interface {
	Match(*http.Request) bool
}
```

`Handler` is similar to `http.Handler`, but it may return `error`
```go
type Handler interface {
	ServerHttp(http.ResponseWriter, *http.Request) error
}
```

Similar with `http.HandlerFunc`
```go
type HandlerFunc func(http.ResponseWriter, *http.Request) error
```

`HandlerFunc` implements the `Handler` interface
```go
func (f HandlerFunc) ServerHttp(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}
```

`Middleware` chains one `Handler` to the next by passing the next `Handler` in the chain
```go
type Middleware func(Handler) Handler
```

`MiddlewareHandler` like a `Handler` with a third argument => `next Handler`
which never be nil, but may be no operation if this is the last handler in the chain.
__Handlers__ which act as middleware should call the next handler's `ServerHttp` method
so as to propagate the request down the chain properly.
__Headers__ which act as responders (content origins) need not invoke the next handler,
since the last handler in the chain should be the first to write the response.
```go
type MiddlewareHandler interface {
	ServeHttp(http.ResponseWriter, *http.Request, Handler) error
}
```

`emptyHandler` is used as a no-op (no operation) handler
```go
var emptyHandler Handler = HandlerFunc(func(http.ResponseWriter, *http.Request) error { return nil })
```

(modules/caddyhttp/error.go)

`HandlerError` is a serializable representation of the an error from within a HTTP handler
```go
type HandlerError struct {
	Err        error  // original error value and message
	StatusCode int    // HTTP status code to associate with the error
	ID         string // generated; for identifying this error in logs
	Trace      string // produced from call stack
}
```

// ErrorCtxKey is the context key to use when storing
// an error (for use with context.Context).
```go
const ErrorCtxKey = caddy.CtxKey("handler_chain_error")
```

(modules/caddyhttp.go)

An implicit suffix middleware that, if reached, sets the StatusCode to the
error stored in the ErrorCtxKey. This is to prevent situations where hte
Error chain does not actually handle the error(for instance, it matches only on some errors)
```go
var errorEmptyHandler Handler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	httpError := r.Context().Value(ErrorCtxKey)
	if handlerError, ok := httpError.(HandlerError); ok {
		w.WriteHeader(handlerError.StatusCode)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return nil
})
```

---

(modules/caddyhttp/responsematchers.go)

`ResponseMatcher` is a type which can determine if an HTTP response matches some criteria
```go
type ResponseMatcher struct {
    // If set, one of these status codes would be required.
    // A one-digit status can be used to represent all codes
    // in that class (e.g. 3 for all 3xx codes)
	StatusCode []int       `json:"status_code,omitempty"`

    // If set, each headers specified must be one of the
    // specified values, with the same logic used by the
    // [request header matcher]
	Headers    http.Header `json:"headers,omitempty"`
}
```

(modules/caddyhttp/routes.go)

