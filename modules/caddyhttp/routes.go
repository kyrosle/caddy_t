package caddyhttp

import (
	"encoding/json"
	"fmt"

	caddy "github.com/kyrosle/caddy_t/v2"
)

type RawMatcherSets []caddy.ModuleMap
type MatcherSets []MatcherSet
type MatcherSet []RequestMatcher

type Route struct {
	Group          string         `json:"group,omitempty"`
	MatcherSetsRaw RawMatcherSets `json:"matcher_sets,omitempty" caddy:"namespace=http.matchers"`
	HandlersRaw    []json.RawMessage
	Terminal       bool                `json:"terminal,omitempty"`
	MatcherSets    MatcherSets         `json:"-"`
	Handler        []MiddlewareHandler `json:"-"`
	middleware     []Middleware
}

func (r Route) Empty() bool {
	return len(r.MatcherSetsRaw) == 0 &&
		len(r.MatcherSets) == 0 &&
		len(r.HandlersRaw) == 0 &&
		len(r.Handler) == 0 &&
		!r.Terminal &&
		r.Group == ""
}

func (r Route) String() string {
	handlersRaw := "["
	for _, hr := range r.HandlersRaw {
		handlersRaw += " " + string(hr)
	}
	handlersRaw += "]"

	return fmt.Sprintf(`{Group:"%s" MatcherSetsRaw:%s HandlersRaw:%s Terminal:%t}`,
		r.Group, r.MatcherSetsRaw, handlersRaw, r.Terminal)
}

type RouterList []Route

func (routes RouterList) Provision(ctx caddy.Context) error {
	err := routes.ProvisionMatchers(ctx)
	if err != nil {
		return err
	}
	return routes.ProvisionHandlers(ctx, nil)
}

func (routes RouterList) ProvisionMatchers(ctx caddy.Context) error {
	for i := range routes {
		matchersIface, err := ctx.LoadMoudle(&routes[i], "MatcherSetsRaw")
		if err != nil {
			return fmt.Errorf("route %d: loading matcher moudles: %v", i, err)
		}
		err = routes[i].MatcherSets.FromInterface(matchersIface)
		if err != nil {
			return fmt.Errorf("route %d: %v", i, err)
		}
	}
}
