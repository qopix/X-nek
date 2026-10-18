# Roadmap

## Stage 0: foundation — scaffold included
- Go module, core/agent/CLI entry points
- API token auth and local JSON persistence
- Dashboard, protocol catalog

## Stage 1: real first node
- Pin Xray-core version
- Implement VLESS + REALITY adapter
- Generate actual client configs and validate/reload engine
- Test with Happ/Incy or another compatible client

## Stage 2: panel — partial scaffold included
- Node and user records; enable/disable user
- Role-based auth, sessions, CSRF protections as applicable
- Replace JSON storage with SQLite/PostgreSQL

## Stage 3: subscriptions/multi-node
- Enforce expiry in each engine
- Mutual authentication and TLS for agent/controller
- Real node health and traffic accounting

## Stage 4: transports
- Shadowsocks, Hysteria 2, WireGuard
- VMess, Trojan, TUIC, SOCKS5
- Integration tests for each client/protocol combination

## Stage 5: Telegram and billing
- Telegram bot with webhook secret
- Telegram Stars for digital goods sold inside Telegram
- External payment integrations where permitted; idempotent webhooks and audit log
