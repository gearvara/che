package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/webteleport/caddy-gos"
	_ "github.com/webteleport/caddy-wasm"
	_ "github.com/webteleport/caddy-webteleport"
	_ "k0s.io/third_party/pkg/module/hub"
	_ "k0s.io/third_party/pkg/plugin/hello"
)

func main() {
	caddycmd.Main()
}
