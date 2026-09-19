# AGENTS.md

## Commands

- Dev: `make dev` (or `go run ./cmd/api`).
- Database: `make postgres-up|postgres-down`.
- Migrations: `make migrate-up|down|status|version|reset`. **When a new migration is needed, always run `make migration name=<descriptive_name>` from the repository root. This Make target is the only permitted way to create a migration file; never create, timestamp, rename, copy, or invoke goose directly to create one. After the target generates the file, edit the generated SQL as needed. If the command fails, fix the command or environment and do not bypass it.**
- Checks: `gofmt -l .`, `go vet ./...`, `go build ./...`; frontend commands run from `web/` with `bun run dev|build|typecheck`.
- No tests or CI are configured.

## Layout

- `cmd/api` — server entrypoint; `cmd/migrate` — goose CLI.
- `internal/module/<domain>/` — one package per domain.
- `internal/httpx` — JSON envelope and middleware; `internal/router` — module composition.
- `internal/database` — `sqlx.DB` backed by the pgx driver; `internal/postgis` — PostGIS types and codecs.
- `migrations/` — goose SQL; `web/src` — React app.
- `web/` is the frontend directory; keep all frontend code there only.

## Conventions

- Keep code boring, explicit, conventional, and easy to read. Prefer 20 obvious lines over 5 clever lines.
- Prefer small, clear duplication over an abstraction that adds indirection or debt. Do not add speculative abstractions, optimizations, or extreme edge-case handling; handle edge cases only when the requirements demand them.
- Each `internal/module/<domain>/` is a self-contained package organized as `domain.go`, `repository.go`, `service.go`, `handler.go`, `route.go`, and `module.go`; group related resources with section banners.
- Keep layers strict: `module.go` wires repository → service → handler → route; repositories own SQL, services own validation and business rules, handlers own HTTP/input validation/responses, and routes own mux registration. Modules never import one another; compose them in `internal/router` with `http.StripPrefix`.
- Use Go 1.22 route patterns such as `GET /{$}` and method names such as `GetMany`, `GetManyByCategory`, and `GetCategories`.
- Return every response through `httpx.RespondJSON` or `httpx.RespondError`. Use `invalid_<field>` for bad input and `internal_error` for unexpected errors.
- Use plain model structs with `db` tags but no JSON tags, pointer fields for nullable columns, and string IDs.
- Repositories use `*sqlx.DB`, raw SQL with uppercase table names, context-aware `GetContext`/`SelectContext`/`ExecContext`, `BindNamed` for named writes, and an early `ctx.Err()` check.
- Run `gofmt`. Group imports as standard library, external packages, then `bearuang-go/...`; document exported types; retain existing section banners, wrapped long calls, and interface assertions.
- Migration files must be generated with `make migration name=<descriptive_name>`; do not manually create migration files. Generated migrations must contain goose `Up` and `Down` sections.
- Read configuration from environment variables with in-code defaults; do not load dotenv files in Go.
