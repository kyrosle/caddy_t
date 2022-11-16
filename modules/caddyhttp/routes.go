package caddyhttp

type Route struct {
	Group string `json:"group,omitempty"`
	MatcherSetsRaw
}
