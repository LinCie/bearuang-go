# AGENTS.md

## Commands

- Dev: `make dev` (or `go run ./cmd/api`).
- Database: `make postgres-up|postgres-down`.
- Migrations: `make migrate-up|down|status|version|reset`; create one with `make migration name=create_users`.
- Checks: `gofmt -l .`, `go vet ./...`, `go build ./...`; frontend commands run from `web/` with `bun run dev|lint|build`.
- No tests or CI are configured.

## Layout

- `cmd/api` — server entrypoint; `cmd/migrate` — goose CLI.
- `internal/module/<domain>/` — one package per domain.
- `internal/httpx` — JSON envelope and middleware; `internal/router` — module composition.
- `internal/database` — `sqlx.DB` backed by the pgx driver; `internal/postgis` — PostGIS types and codecs.
- `migrations/` — goose SQL; `web/src` — React app.

## Conventions

- Keep code boring, explicit, conventional, and easy to read. Prefer 20 obvious lines over 5 clever lines.
- Prefer small, clear duplication over an abstraction that adds indirection or debt. Do not add speculative abstractions, optimizations, or extreme edge-case handling; handle edge cases only when the requirements demand them.
- Modules use `domain.go`, `repository.go`, `service.go`, `handler.go`, `route.go`, and `module.go`. Wire repository → service → handler; handlers do HTTP, repositories do SQL, services delegate.
- Modules must not import each other. Compose them in `internal/router` with `http.StripPrefix`.
- Use Go 1.22 route patterns such as `GET /{$}` and method names such as `GetMany`, `GetManyByCategory`, and `GetCategories`.
- Return every response through `httpx.RespondJSON` or `httpx.RespondError`. Use `invalid_<field>` for bad input and `internal_error` for unexpected errors.
- Use plain model structs with `db` tags but no JSON tags, pointer fields for nullable columns, and string IDs.
- Repositories use `*sqlx.DB`, raw SQL with uppercase table names, context-aware `GetContext`/`SelectContext`/`ExecContext`, `BindNamed` for named writes, and an early `ctx.Err()` check.
- Run `gofmt`. Group imports as standard library, external packages, then `bearuang-go/...`; document exported types; retain existing section banners, wrapped long calls, and interface assertions.
- Name migrations `<timestamp>_<name>.sql` and include goose `Up` and `Down` sections.
- Read configuration from environment variables with in-code defaults; do not load dotenv files in Go.
