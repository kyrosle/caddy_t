package caddyhttp

import caddy "github.com/kyrosle/caddy_t/v2"

const ErrorCtxKey = caddy.CtxKey("handler_chain_error")

type HandlerError struct {
	Err        error
	StatusCode int
	ID         string
	Trace      string
}
