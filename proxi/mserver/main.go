package main

import (
	"github.com/caddyserver/caddy/v2"
)

type App struct{}

func (*App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "layer4",
		New: func() caddy.Module { return new(App) },
	}
}

func main() {}

// Interface guard
var _ caddy.App = (*App)(nil)
