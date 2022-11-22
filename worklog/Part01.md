# Main

Begin from the module `caddyhttp` , base of the thinking of web services are relevant with `http`. That's mean we can start from building request and response, and then enter a further step.

Having : 
* `RequestMatcher` (modules/caddyhttp/caddyhttp.go)
* `Handler`
* `MiddlewareHandler`
* Error handlers
* `ResponseMatcher` (modules/caddyhttp/responsematchers.go)
* `Route` (modules/caddyhttp/routes.go)
* `RouteList` (modules/caddyhttp/routes.go)

# Http Building

use `http.Request` and `http.Response` in go standard module

__use module__ :
* `net/http`
* `encoding/json`

__details__ :

* `RequestMatcher`
* `ResponseMatcher`
* `Handler`


## `RequestMatcher` (modules/caddyhttp/caddyhttp.go)
`RequestMatcher` Is to match a request
```go
type RequestMatcher interface {
	Match(*http.Request) bool
}
```

## `Handler` 
Is similar to `http.Handler`, but it may return `error`
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

## `MiddlewareHandler`
like a `Handler` with a third argument => `next Handler`
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

## Error handlers

---

`emptyHandler` is used as a no-op (no operation) handler
```go
var emptyHandler Handler = HandlerFunc(func(http.ResponseWriter, *http.Request) error { return nil })
```

---

Structure `HandlerError` (modules/caddyhttp/error.go)

`HandlerError` is a serializable representation of the an error from within a HTTP handler.

```go
type HandlerError struct {
	Err        error  // original error value and message
	StatusCode int    // HTTP status code to associate with the error
	ID         string // generated; for identifying this error in logs
	Trace      string // produced from call stack
}
```
`ErrorCtxKey` is the context key to use when storing an error (for use with context.Context).
```go
const ErrorCtxKey = caddy.CtxKey("handler_chain_error")
```

`errorEmptyHandler` 
an implicit suffix middleware that, if reached, sets the StatusCode to the
error stored in the `ErrorCtxKey`. This is to prevent situations where hte
Error chain does not actually handle the error(for instance, it matches only on some errors)

Convert the value from `request.context.value` to `HandlerError` and then use its `WriteHeader()`.
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

## `ResponseMatcher` (modules/caddyhttp/responsematchers.go)

`ResponseMatcher` is a type which can determine if an HTTP response matches some criteria
```go
type ResponseMatcher struct {
	StatusCode []int       `json:"status_code,omitempty"`
	Headers    http.Header `json:"headers,omitempty"`
}
```

__ResponseMatcher Fields__ : 
* `StatusCode` : 
If set, one of these status codes would be required.
A one-digit status can be used to represent all codes
in that class (e.g. 3 for all 3xx codes)

* `Headers` : 
If set, each headers specified must be one of the specified values, 
with the same logic used by the [request header matcher]

## `Route` (modules/caddyhttp/routes.go)

`Route` consists of a set of rules for matching HTTP requests,
a list of handler to execute, and optional flow control
parameters which customize the handling of HTTP requests
in highly flexible and performant manner.

```go
type Route struct {
	Group          string         `json:"group,omitempty"`
	MatcherSetsRaw RawMatcherSets `json:"matcher_sets,omitempty" caddy:"namespace=http.matchers"`
	HandlersRaw    []json.RawMessage
	Terminal       bool                `json:"terminal,omitempty"`
	MatcherSets    MatcherSets         `json:"-"`
	Handler        []MiddlewareHandler `json:"-"`
	middleware     []Middleware
}
```

__Route Fields__ : 

* `Group` :
Which is an optional name for a group to which this
route belongs. Grouping a route makes it mutually
exclusive with others in its group; if a route belongs
to a group, only the first matching route in that group 
will be executed.

* `MatcherSetsRaw` :
The matcher sets which will be used to qualify this 
route for a request (essentially the "if" statement of this route).
Each matcher set is OR'ed, but matchers within a set are AND'ed together.

---

__About `MatcherSetsRow`__ :

(modules/caddyhttp/routes.go)
```go
type RawMatcherSets []caddy.ModuleMap
```
Is a group of matcher sets in their raw, JSON from.

(caddy.go)
```go
type ModuleMap map[string]json.RawMessage
```
Is a map that can contain multiple modules, where the map key is the module's name. 
(The namespace is usually read from an associated field's struct tag.)
Because the module's name is given as the key in a module map,
the name does not have to be given in the `json.RawMessage`.

---

* `HandlerRaw` : 

The list of handlers for this route. Upon matching a request, they are chained
together in a middleware fashion: requests flow from the first handler to the last 
(top of the list to the bottom), 
with the possibility that any handler could stop
the chain and/or return an error. Response flow back through the chain 
(bottom of the list to the top) as they are written out to the client.

Not all handlers call the next handler in the chain. For example, the `reverse_proxy` handler always sends a request upstream or returns an error. Thus, configuring handlers after `reverse_proxy` int the same route is illogical, since they would never be executed. You will want to put handlers which originate the response at the very end of your route(s). 

Some handlers manipulate the response. Remember that requests flow down the list, and responses flow up the list.

For example, if you wanted to use both `templates` and `encode` handlers, you would need to put `templates` after `encode` in your route, because responses flow up. 
Thus, `templates` will be able to parse and execute the plain-text response as a template, and then return it up to the `encode` handler, which will then compress it into a binary format.

If `templates` came before `encode`, then `encode` would write a compressed, binary-encoded response to `templates` which would not be able to parse the response properly.

The correct order, then, is this:
```
[
	{"handler": "encode"},
	{"handler": "templates"},
	{"handler": "file_server"},
]
```

The request flow down (`encode` -> `templates` -> `file_server`).

1. First, `encode` will choose how to `encode` the response and warp the response.
2. Then, `templates` will warp the response with a buffer.
3. Finally, `file_server` will originate the content from a file.

The response flow up (`file_server` -> `templates` -> `encode`): 

1. First, `file_server` will write the file to the response.
2. That write will be buffered and then executed by `templates`.
3. Lastly, the write from `templates` will flow into `encode` which will compress the stream.


* `Terminal` :
If true, no more routes will be executed after this one.

* `MatcherSets` : 

```go
type MatcherSets []MatcherSet
```
```go
type MatcherSet  []RequestMatcher
```

`MatcherSet` is a set of matchers which
must all match in order for the request 
to be matched successfully.

__Router Functions__ :

* Empty() : check the route wether is empty.
* String() : Stringify the route.

---

## `RouteList` (modules/caddyhttp/routes.go)
A list of server routes that can create a middleware chain.
```go
type RouteList []Route
```

__RoutList Functions__ :

* Provision(ctx caddy.Context) error

This sets up both the matchers and handlers in the routes.

* ProvisionMatchers(ctx caddy.Context) error 

This sets up all the matchers by loading the
matcher modules. Only call this method directly if you need 
to set up matchers and handlers separately without having 
to provision a second time; otherwise use Provision instead.
