package caddy

import (
	"context"
	"log"
)

type Context struct {
	context.Context
	moduleInstances map[string][]Module
	cfg             *Config
	cleanupFuncs    []func()
	ancestry        []Module
}

func NewContext(ctx Context) (Context, context.CancelFunc) {
	newCtx := Context{moduleInstances: make(map[string][]Module), cfg: ctx.cfg}
	c, cancel := context.WithCancel(ctx.Context)
	wrappedCancel := func() {
		cancel()
		for _, f := range ctx.cleanupFuncs {
			f()
		}

		for modName, modInstances := range newCtx.moduleInstances {
			for _, inst := range modInstances {
				if cu, ok := inst.(CleanerUpper); ok {
					err := cu.Cleanup()
					if err != nil {
						log.Printf("[ERROR] %s (%p): cleanup: %v", modName, inst, err)
					}
				}
			}
		}
	}
	newCtx.Context = c
	return newCtx, wrappedCancel
}
