package protocol

type Capability struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Transport string   `json:"transport"`
	Status    string   `json:"status"`
	Notes     string   `json:"notes"`
	Engines   []string `json:"engines"`
}

func Catalog() []Capability {
	return []Capability{
		{"vless-reality", "VLESS + REALITY", "Xray transport", "planned-integration", "Primary MVP target; requires compatible engine and client validation.", []string{"Xray-core"}},
		{"hysteria2", "Hysteria 2", "QUIC/UDP", "planned-integration", "Requires UDP reachability and compatible engine.", []string{"Hysteria 2"}},
		{"shadowsocks", "Shadowsocks", "Engine-dependent", "planned-integration", "Select and pin a maintained implementation.", []string{"Xray-core", "Shadowsocks implementations"}},
		{"wireguard", "WireGuard", "L3 VPN / UDP", "planned-integration", "Needs interface, routing, key and firewall management.", []string{"wireguard-tools", "kernel WireGuard"}},
		{"vmess", "VMess", "Engine-dependent", "planned-integration", "Validate client interoperability.", []string{"Xray-core"}},
		{"trojan", "Trojan", "TLS/TCP", "planned-integration", "TLS and client interoperability required.", []string{"Xray-core"}},
		{"tuic", "TUIC", "QUIC/UDP", "planned-integration", "Requires compatible server and UDP reachability.", []string{"TUIC"}},
		{"socks5", "SOCKS5", "TCP/UDP", "planned-integration", "SOCKS5 alone does not encrypt traffic.", []string{"Xray-core", "compatible proxies"}},
	}
}
