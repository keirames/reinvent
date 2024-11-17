package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/mholt/caddy-l4/layer4"
	_ "github.com/mholt/caddy-l4/modules/l4clock"
	_ "github.com/mholt/caddy-l4/modules/l4dns"
	_ "github.com/mholt/caddy-l4/modules/l4echo"
	_ "github.com/mholt/caddy-l4/modules/l4http"
	_ "github.com/mholt/caddy-l4/modules/l4openvpn"
	_ "github.com/mholt/caddy-l4/modules/l4postgres"
	_ "github.com/mholt/caddy-l4/modules/l4proxy"
	_ "github.com/mholt/caddy-l4/modules/l4proxyprotocol"
	_ "github.com/mholt/caddy-l4/modules/l4quic"
	_ "github.com/mholt/caddy-l4/modules/l4rdp"
	_ "github.com/mholt/caddy-l4/modules/l4regexp"
	_ "github.com/mholt/caddy-l4/modules/l4socks"
	_ "github.com/mholt/caddy-l4/modules/l4ssh"
	_ "github.com/mholt/caddy-l4/modules/l4subroute"
	_ "github.com/mholt/caddy-l4/modules/l4tee"
	_ "github.com/mholt/caddy-l4/modules/l4throttle"
	_ "github.com/mholt/caddy-l4/modules/l4tls"
	_ "github.com/mholt/caddy-l4/modules/l4winbox"
	_ "github.com/mholt/caddy-l4/modules/l4wireguard"
	_ "github.com/mholt/caddy-l4/modules/l4xmpp"
)

func main() {
	caddycmd.Main()
}
