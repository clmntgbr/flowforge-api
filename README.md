# Go API Template (Fiber + Clerk + PostgreSQL)

This organization is a Go API starter built with Fiber, GORM, PostgreSQL, and Clerk authentication.
It includes:

- JWT authentication with Clerk JWKS validation
- Clerk webhook verification (Svix signature)
- User synchronization between Clerk and local database
- Docker-based local development and production compose setup
- Makefile commands for daily workflows

## Tech Stack

- Go `1.25`
- Fiber `v3`
- GORM + PostgreSQL
- Clerk (`clerk-sdk-go`)
- Svix webhook signature verification
- Docker / Docker Compose

## Organization Structure

- `server.go`: application entry point and route registration
- `config/`: environment and database configuration
- `deps/`: dependency wiring
- `middleware/`: authentication and webhook protection
- `handler/`: HTTP handlers
- `service/`: business logic (auth, users, Clerk, webhooks)
- `repository/`: database access layer
- `domain/`: data models
- `dto/`: request/event DTOs
- `ctxutil/`: Fiber context helpers

## API Overview

### Health checks

- `GET /livez`
- `GET /readyz`
- `GET /startupz`

### Protected API

All `/api/*` routes require a `Bearer` token issued by Clerk.

- `GET /api/users/me`: returns the authenticated local user

### Webhooks

- `POST /webhook/clerk`: Clerk webhook endpoint

Supported Clerk events:

- `user.created`
- `user.updated`
- `user.deleted`

## Clerk Integration

This template uses Clerk in two different paths:

1. **API Authentication**
   - The middleware reads the `Authorization: Bearer <token>` header.
   - The token is validated against Clerk JWKS from:
     - `${CLERK_FRONTEND_API}/.well-known/jwks.json`
   - If the Clerk user is not yet in the local DB, it is created automatically.
   - If the local user is banned, access is denied.

2. **Webhook Synchronization**
   - `/webhook/clerk` is protected using Svix signature headers:
     - `svix-id`
     - `svix-timestamp`
     - `svix-signature`
   - Signature is verified with `CLERK_WEBHOOK_SECRET`.
   - User lifecycle events update the local `users` table.

## Environment Variables

Copy `.env.dist` to `.env` and fill values:

```bash
cp .env.dist .env
```

Main variables:

- `PORT`: API port inside the container (`3000` by default)
- `GO_ENV`: `development` or `production`
- `DATABASE_URL`: PostgreSQL DSN used by GORM
- `POSTGRES_*`: database service configuration
- `CLERK_SECRET_KEY`: Clerk backend secret key
- `CLERK_FRONTEND_API`: Clerk frontend API URL (issuer + JWKS base URL)
- `CLERK_WEBHOOK_SECRET`: signing secret for webhook verification
- `NGROK_AUTHTOKEN`: required if using ngrok service in `compose.yaml`

## Local Development

### Prerequisites

- Docker
- Docker Compose
- Make

### Start development stack

```bash
make dev
```

Then the API is exposed on `http://localhost:4000` (mapped to container `3000`).

Useful commands:

```bash
make dev-logs
make dev-down
make dev-restart
make dev-rebuild
make shell
make test
make lint
```

## Production Compose Commands

```bash
make prod
make prod-d
make prod-logs
make prod-down
make prod-restart
make prod-rebuild
make shell-prod
```

## Makefile Commands

### Development

- `make dev`: start development services in background
- `make dev-logs`: stream development logs
- `make dev-down`: stop development services
- `make dev-restart`: restart development services
- `make dev-rebuild`: rebuild and restart development services

### Production

- `make prod`: run production compose in foreground
- `make prod-d`: run production compose in background
- `make prod-logs`: stream production logs
- `make prod-down`: stop production services
- `make prod-restart`: restart production services
- `make prod-rebuild`: rebuild and restart production services

### Images

- `make build-dev`: build development image target
- `make build-prod`: build production image target

### Utilities

- `make shell`: shell inside development API container
- `make shell-prod`: shell inside production API container
- `make test`: run Go tests in API container
- `make lint`: run `golangci-lint --fix` in API container
- `make clean`: remove containers/volumes and prune temp/docker cache
- `make clean-all`: remove all related images, volumes, and prune more aggressively

## Docker Notes

- Development compose file in this repository is `compose.yaml`.
- Production compose file in this repository is `compose.prod.yaml`.
- The Makefile currently references `docker-compose.yml` and `docker-compose.prod.yml`.

If your environment does not resolve these names automatically, either:

- rename files to match the Makefile, or
- update Makefile paths to `compose.yaml` and `compose.prod.yaml`.

## Data Model

The `users` table is auto-migrated at startup with:

- `id` (UUID primary key)
- `clerk_id` (unique)
- `first_name`
- `last_name`
- `banned`
- `created_at`
- `updated_at`

## Troubleshooting

- **401 Invalid token**
  - Verify `CLERK_FRONTEND_API` is correct and reachable.
  - Verify the token issuer matches `CLERK_FRONTEND_API`.
- **401 Invalid signature on webhook**
  - Verify `CLERK_WEBHOOK_SECRET` matches your Clerk endpoint secret.
  - Ensure Clerk sends the Svix headers untouched.
- **DB connection failures**
  - Check `DATABASE_URL` and `POSTGRES_*` values.
  - Ensure database container is healthy.
