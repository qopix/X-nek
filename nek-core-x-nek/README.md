# NEK Core + X-NEK

Go management-plane MVP scaffold: REST API, dashboard, node/user/subscription records, protocol capability catalog, basic agent health endpoint, and CLI placeholder.

**This is not yet a working VPN server.** Protocols are listed as `planned-integration`; subscription URLs return metadata placeholders, not importable client configurations. Actual transport adapters, expiry enforcement, real usage metrics, Telegram bot, and payment integration remain roadmap items.

## Run locally

Requires Go 1.22+.

```sh
go test ./...
go run ./cmd/nek-core
```

Open http://127.0.0.1:8080. The first launch prints a generated admin token. Copy it into the panel. Data is stored at `data/nek.json` (override with `NEK_DATA`).

Environment:
- `NEK_ADDR` (default `127.0.0.1:8080`)
- `NEK_DATA` (default `data/nek.json`)
- `NEK_ADMIN_TOKEN` (optional; generated if unset)
- `NEK_PUBLIC_URL` (default `http://127.0.0.1:8080`)

Authenticated API endpoints use `Authorization: Bearer <token>`:
- `GET /api/health` (public health only)
- `GET /api/overview`, `GET /api/protocols`
- `GET|POST /api/nodes`, `GET|POST /api/users`
- `POST /api/users/{id}/disable`, `POST /api/users/{id}/enable`
- `GET|POST /api/subscriptions`
- `GET /sub/{token}` (placeholder metadata only)

## Before internet deployment

- Keep listening on localhost behind a TLS reverse proxy until production auth is added.
- Use a strong secret, firewall the admin API, and do not expose the scaffold agent publicly.
- Implement mutual TLS or signed commands for agents, least-privilege execution, secret redaction, rate limits, backups, and audit logs.
- Use SQLite/PostgreSQL before concurrent production use; the JSON store is only for the early MVP.
- Validate all protocol engines and clients before marking any protocol operational.
