# FreeFSM — Free Field Service Manager

Self-hosted, open-source field service management for FreeBSD and Linux.
Single static Go binary, PostgreSQL backend, zero npm/NPM dependencies.

## Features

- **Dashboard** — overview with quick access to all modules
- **Customers** — full CRUD with search, status filter, HTMX pagination
- **Jobs** — work orders with status workflow, scheduling, arrival windows
- **Schedule Calendar** — month view with clickable job cards
- **Estimates** — line items editor (Alpine.js), status workflow, tax calculation
- **Invoices** — line items editor, payment recording, status workflow
- **Items / Pricebook** — service and product catalog with SKU and pricing
- **Auth** — setup token for initial admin, bcrypt + HTTP-only session cookies

## Tech Stack

| Layer | Choice |
|-------|--------|
| Language | Go |
| Router | chi |
| Database | PostgreSQL (JSONB, TIMESTAMPTZ) |
| ORM | ent (type-safe codegen) |
| Templates | Templ (compile-time safety) |
| Interactivity | HTMX 2 + Alpine.js |
| CSS | Pico CSS |
| Deploy | Single binary, systemd + rc.d |

## Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL 15+
- `ent` CLI: `go install entgo.io/ent/cmd/ent@latest`
- `templ` CLI: `go install github.com/a-h/templ/cmd/templ@latest`
- Ensure `$HOME/go/bin` is in `$PATH`

### Database

```sql
CREATE USER freefsm WITH PASSWORD 'changeme';
CREATE DATABASE freefsm OWNER freefsm;
GRANT ALL PRIVILEGES ON DATABASE freefsm TO freefsm;
```

If using Fedora or any system with ident/peer auth, edit `pg_hba.conf` and change
`local` and `127.0.0.1` entries from `ident`/`peer` to `md5`, then restart
PostgreSQL.

### Build & Run

```bash
git clone https://github.com/MartialM1nd/freefsm.git
cd freefsm
cp .env.example .env
# Edit .env: set DB_PASSWORD, SESSION_SECRET, SETUP_TOKEN

make run
# → http://localhost:3000
```

### First-Time Setup

1. Visit `http://localhost:3000` — you'll be redirected to `/setup`
2. Enter your `FREEFSM_SETUP_TOKEN` value (from `.env`)
3. Create an admin account (name, email, password)
4. You're logged in and on the Dashboard

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `FREEFSM_DB_HOST` | `localhost` | PostgreSQL host |
| `FREEFSM_DB_PORT` | `5432` | PostgreSQL port |
| `FREEFSM_DB_NAME` | `freefsm` | Database name |
| `FREEFSM_DB_USER` | `freefsm` | Database user |
| `FREEFSM_DB_PASSWORD` | *(required)* | Database password |
| `FREEFSM_DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `FREEFSM_ADDR` | `:3000` | HTTP listen address |
| `FREEFSM_LOG_LEVEL` | `info` | `debug` / `info` / `warn` / `error` |
| `FREEFSM_SESSION_SECRET` | *(required)* | Cookie encryption key |
| `FREEFSM_SETUP_TOKEN` | *(required)* | Initial admin registration token |

## Project Structure

```
freefsm/
├── cmd/freefsm/           # Entry point + static file embed
│   └── static/            # CSS, JS (Pico, HTMX, Alpine)
├── internal/
│   ├── config/            # Env loading + DSN builder
│   ├── database/          # pgxpool connection + SQL migration runner
│   │   └── migrations/    # 7 SQL migration files
│   ├── ent/
│   │   └── schema/        # 11 ent schema definitions
│   ├── handlers/          # HTTP handlers (chi routes)
│   ├── middleware/         # Auth, Flash, user context
│   ├── services/          # Business logic (ent queries)
│   └── templates/         # 22 Templ files (pages + partials)
├── deploy/
│   ├── freebsd/           # rc.d service script
│   └── linux/             # systemd unit + config sample
├── Makefile               # build, install, fmt, lint, test
├── PLAN.md                # Full roadmap + architecture
└── go.mod
```

## Development

```bash
make build       # ent generate → templ generate → go build → dist/freefsm
make run         # build + run
make fmt         # go fmt ./...
make lint        # go vet ./...
make clean       # remove dist/
make install     # install to /usr/local/bin
```

### Adding a New Entity

1. Create a SQL migration in `internal/database/migrations/`
2. Define an ent schema in `internal/ent/schema/`
3. Run `ent generate ./internal/ent/schema`
4. Create a service in `internal/services/`
5. Create a handler in `internal/handlers/`
6. Create templates in `internal/templates/`
7. Register routes in `internal/handlers/router.go`

## Deployment

### Linux (systemd)

```bash
make install-linux   # installs binary + systemd unit
systemctl start freefsm
```

### FreeBSD (rc.d)

```bash
make install-freebsd # installs binary + rc.d script
service freefsm start
```

## License

AGPL-3.0
