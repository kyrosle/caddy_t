package caddy

type Module interface {
	CaddyModule() ModuleInfo
}

type ModuleInfo struct {
	ID  ModuleID
	New func() Module
}

type ModuleID string