package caddy

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caddyserver/certmagic"
)

type Duration time.Duration

type CtxKey string

type ModuleMap map[string]json.RawMessage

type Config struct {
	Admin   *AdminConfig `json:"admin,omitempty"`
	Logging *Logging     `json:"logging,omitempty"`

	StorageRaw json.RawMessage `json:"storage,omitempty" caddy:"namespace=caddy.storage inline_key=module"`
	AppsRaw    ModuleMap       `json:"apps,omitempty" caddy:"namespace="`

	apps       map[string]App
	storage    certmagic.Storage
	cancelFunc context.CancelFunc
}


type App interface {
	Start() error
	Stop() error
}
