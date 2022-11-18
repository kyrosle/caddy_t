package caddy

import "context"

type Context struct {
	context.Context
	moduleInstances map[string][]Module
	cfg *Config
	cleanupFuncs []func()
	ancestry []Module
}