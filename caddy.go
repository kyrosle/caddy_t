package caddy

import (
	"encoding/json"
	"time"
)

type Duration time.Duration

type CtxKey string

type ModuleMap map[string]json.RawMessage

type Config struct {
	Admin   *AdminConfig `json:"admin,omitempty"`
	Logging *Logging     `json:"logging,omitempty"`

	StorageRaw json.RawMessage `json:"storage,omitempty" caddy:"namespace=caddy.storage inline_key=module"`
	AppsRaw ModuleMap `json:"admin,omitempty" caddy:"namespace="`
}
