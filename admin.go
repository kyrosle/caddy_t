package caddy

import (
	"crypto"
	"encoding/json"
	"net/http"

	"github.com/caddyserver/certmagic"
)

type AdminConfig struct {
	Disabled      bool            `json:"disabled,omitempty"`
	Listen        string          `json:"listen,omitempty"`
	EnforceOrigin bool            `json:"enforce_origin,omitempty"`
	Origins       []string        `json:"origins,omitempty"`
	Config        *ConfigSettings `json:"config,omitempty"`
	Identity      *IdentityConfig
	Remote        *RemoteAdmin
	routers       []AdminRouter
}

type ConfigSettings struct {
	Persist   *bool           `json:"persist,omitempty"`
	LoadRaw   json.RawMessage `json:"load,omitempty" caddy:"namespace=caddy.config_loaders inline_key=module"`
	LoadDelay Duration        `json:load_delay,omitempty"`
}

type IdentityConfig struct {
	Identifiers []string
	IssuersRaw  []json.RawMessage
	issuers     []certmagic.Issuer
}

type RemoteAdmin struct {
	Listen        string         `json:"listen,omitempty"`
	AccessControl []*AdminAccess `json:"access_control,omitempty"`
}

type AdminAccess struct {
	PublicKeys  []string           `json:public_keys,omitempty`
	Permissions []AdminPermissions `json:"permissions,omitempty"`
	publicKeys  []crypto.PublicKey
}

type AdminPermissions struct {
	Path    []string `json:"path,omitempty"`
	Methods []string `json:"methods,omitempty"`
}

type AdminRouter interface {
	Routes() []AdminRoute
}
type AdminRoute struct {
	Pattern string
	Handler AdminHandler
}

type AdminHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}
