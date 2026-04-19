# Flowforge API

Go HTTP API (Fiber, GORM, PostgreSQL, Clerk). Local development and production run via Docker Compose.

[![CI](https://github.com/clmntgbr/flowforge-api/actions/workflows/test.yml/badge.svg)](https://github.com/clmntgbr/flowforge-api/actions/workflows/test.yml)
[![Codecov](https://codecov.io/gh/clmntgbr/flowforge-api/graph/badge.svg)](https://app.codecov.io/gh/clmntgbr/flowforge-api)
[![Handler tests](https://img.shields.io/endpoint?url=https%3A%2F%2Fraw.githubusercontent.com%2Fclmntgbr%2Fflowforge-api%2Fmain%2Fbadges%2Fhandler-tests.json)](https://github.com/clmntgbr/flowforge-api/tree/main/tests/handler_test)

## Features

- REST API on **Fiber v3**
- **Clerk** JWT validation (JWKS) and user sync
- **Clerk webhooks** with Svix signature verification
- **GORM** + PostgreSQL
- Health endpoints (liveness, readiness, startup)
- Docker Compose for dev and prod

## Tech stack

| Layer        | Choice                          |
| ------------ | ------------------------------- |
| Language     | Go 1.25 (`module forgeflow-api`) |
| HTTP         | Fiber v3                        |
| Data         | GORM, PostgreSQL                |
| Auth         | Clerk (`clerk-sdk-go`)          |
| Webhooks     | Svix verification               |

## Repository layout

| Path | Role |
| ---- | ---- |
| `server.go` | Entrypoint, routes, middleware stack |
| `config/` | Env loading, DB connection, migrations |
| `deps/` | Dependency injection |
| `middleware/` | Auth, webhook guards |
| `handler/` | HTTP handlers |
| `service/` | Business logic |
| `repository/` | Data access |
| `domain/` | Models |
| `dto/` | Request/response and event payloads |
| `ctxutil/` | Fiber context helpers |
| `rules/` | Validation rules |
| `validator/` | Shared validators |
| `errors/` | API errors |
| `tests/handler_test/` | Handler tests (external test package) |

## API

### Health (no auth)

| Method | Path |
| ------ | ---- |
| GET | `/livez` |
| GET | `/readyz` |
| GET | `/startupz` |

### Authenticated (`/api/*`)

All routes below require `Authorization: Bearer <Clerk token>`.

**Users**

| Method | Path |
| ------ | ---- |
| GET | `/api/users/me` |

**Organizations**

| Method | Path |
| ------ | ---- |
| GET | `/api/organizations` |
| GET | `/api/organizations/:id` |
| POST | `/api/organizations` |
| PUT | `/api/organizations/:id` |
| PUT | `/api/organizations/:id/activate` |

**Endpoints**

| Method | Path |
| ------ | ---- |
| GET | `/api/endpoints` |
| GET | `/api/endpoints/:id` |
| POST | `/api/endpoints` |
| PUT | `/api/endpoints/:id` |

**Workflows**

| Method | Path |
| ------ | ---- |
| GET | `/api/workflows` |
| GET | `/api/workflows/:id` |
| POST | `/api/workflows` |
| PUT | `/api/workflows/:id` |
| PUT | `/api/workflows/:id/steps` |

**Connexions**

| Method | Path |
| ------ | ---- |
| POST | `/api/connexions` |
| DELETE | `/api/connexions/:id` |

**Steps**

| Method | Path |
| ------ | ---- |
| GET | `/api/steps/:id` |
| PUT | `/api/steps/:id` |

### Webhooks (signed, not Bearer JWT)

| Method | Path | Notes |
| ------ | ---- | ----- |
| POST | `/webhook/clerk` | Svix headers; events `user.created`, `user.updated`, `user.deleted` |

## Clerk behavior

1. **API** — Middleware validates the Bearer token against `${CLERK_FRONTEND_API}/.well-known/jwks.json`. Unknown Clerk users can be created locally; banned users are rejected.
2. **Webhooks** — `svix-id`, `svix-timestamp`, `svix-signature` are verified with `CLERK_WEBHOOK_SECRET` before updating the local `users` table.

## Environment

Copy the example env and adjust values:

```bash
cp .env.dist .env
```

| Variable | Purpose |
| -------- | ------- |
| `PORT` | API port inside the container (often `3000`) |
| `GO_ENV` | `development` or `production` |
| `DATABASE_URL` | PostgreSQL DSN for GORM |
| `POSTGRES_*` | DB service settings (Compose) |
| `CLERK_SECRET_KEY` | Clerk backend secret |
| `CLERK_FRONTEND_API` | Issuer / JWKS base URL |
| `CLERK_WEBHOOK_SECRET` | Webhook signing secret |
| `NGROK_AUTHTOKEN` | Only if you use the ngrok service from `compose.yaml` |

## Local development

**Prerequisites:** Docker, Docker Compose, Make.

```bash
make dev
```

Default mapping is host **4000** → container **3000** (see your Compose file).

```bash
make dev-logs      # follow logs
make dev-down      # stop stack
make dev-restart
make dev-rebuild
make shell         # shell in API container
make test          # go test inside API container
make lint          # golangci-lint in API container
```

Production-style Compose:

```bash
make prod          # foreground
make prod-d        # detached
make prod-logs
make prod-down
make prod-restart
make prod-rebuild
make shell-prod
```

Other targets: `make build-dev`, `make build-prod`, `make clean`, `make clean-all`.

### Compose files

This repository ships **`compose.yaml`** (dev) and **`compose.prod.yaml`** (prod). The Makefile invokes `docker-compose` without an explicit `-f` flag; Compose v2 resolves `compose.yaml` in the project directory. If your CLI looks for another filename, either rename the files or add `-f compose.yaml` / `-f compose.prod.yaml` to the Makefile targets.

## Tests and CI

- **Local (host):**  
  `go test ./tests/handler_test -v`  
  Coverage for the `handler` package from external tests:  
  `go test ./tests/handler_test -coverprofile=coverage.txt -coverpkg=./handler`

- **CI:** [`.github/workflows/test.yml`](.github/workflows/test.yml) runs the same tests, uploads **`coverage.txt`** to [Codecov](https://codecov.io), publishes the job summary, stores the coverage artifact, and refreshes [`badges/handler-tests.json`](badges/handler-tests.json) on pushes to `main` or `master`.

- **Codecov:** add a GitHub Actions secret **`CODECOV_TOKEN`** (from your Codecov project settings).

- **Badge “Handler tests”:** Shields reads the JSON on the default branch. If that branch is not `main`, update the badge URL in this file. Private repos may block Shields unless the JSON stays publicly readable.

## Data model

`users` is auto-migrated at startup (among others, depending on `AutoMigrate`):

- `id` (UUID, PK)  
- `clerk_id` (unique)  
- `first_name`, `last_name`  
- `banned`  
- `created_at`, `updated_at`

## Troubleshooting

| Symptom | What to check |
| ------- | --------------- |
| `401` on API | `CLERK_FRONTEND_API`, token issuer, JWKS reachability |
| `401` on webhook | `CLERK_WEBHOOK_SECRET`, Svix headers unchanged |
| DB errors | `DATABASE_URL`, `POSTGRES_*`, DB container health |
