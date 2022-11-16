package caddyhttp

import "net/http"

type RequestMatcher interface {
	Match(*http.Request) bool
}

type Handler interface {
	ServeHttp(http.ResponseWriter, *http.Request) error
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (f HandlerFunc) ServeHttp(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type Middleware func(Handler) Handler

type MiddlewareHandler interface {
	ServeHttp(http.ResponseWriter, *http.Request, Handler) error
}

var emptyHandler Handler = HandlerFunc(func(http.ResponseWriter, *http.Request) error { return nil })

var errorEmptyHandler Handler = HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
	httpError := r.Context().Value(ErrorCtxKey)
	if handlerError, ok := httpError.(HandlerError); ok {
		w.WriteHeader(handlerError.StatusCode)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return nil
})

type ResponseHandler struct {
	Match      *ResponseMatcher `json:"match,omitempty"`
	StatusCode WeakString       `json:"statusCode,omitempty"`
}

type WeakString string
