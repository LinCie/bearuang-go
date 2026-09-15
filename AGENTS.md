# AGENTS.md

## Commands

- Dev server (air hot reload): `make dev`; direct: `go run ./cmd/api`
- Postgres 17 + PostGIS: `make postgres-up` / `make postgres-down`
- Migrations: `make migrate-up|down|status|version|reset`; new file: `make migration name=create_users`
- Backend checks: `gofmt -l .` (expect empty), `go vet ./...`, `go build ./...`
- Frontend (in `web/`, bun): `bun run dev|lint|build`
- No tests or CI configured.
- Note: Makefile `start`/`build` targets still reference the stale `./cmd/web` path; the entrypoint is `./cmd/api`.

## Layout

- `cmd/api` — server entry (currently a `main()` stub); `cmd/migrate` — goose CLI
- `internal/module/<domain>/` — one package per domain
- `internal/httpx` — JSON envelope + middleware `Chain`; `internal/router` — mounts modules
- `internal/database` — pgx pool; `internal/postgis` — `Point` alias, EWKB codec, pgx type registration
- `migrations/` — goose SQL; `web/src` — React app

## Backend conventions

- Module files: `model.go`, `repository.go` (`Repository` interface + `PostgresRepository`), `service.go`, `handler.go`, `route.go`, `module.go` (wires repo → service → handler and owns a sub-`ServeMux`).
- Layering: handler → service → repository. Handlers do HTTP only, SQL only in repositories; services delegate.
- Modules don't import each other; compose them in `internal/router` via `mux.Handle("/places/", http.StripPrefix("/places/", module.Mux))`.
- Routes use Go 1.22 patterns in `route.go`: `mux.HandleFunc("GET /{$}", h.GetMany)`. Method names: `GetMany`, `GetManyByCategory`, `GetCategories`.
- All responses go through `httpx.RespondJSON` / `httpx.RespondError`: `{"data": ...}` or `{"code", "message"}`. Bad input → 400 `invalid_<field>`; unexpected errors → 500 `internal_error` with `err.Error()`.
- Models are plain structs with no JSON tags (keys are Go field names), pointer fields for nullable columns, `int64` IDs.
- Queries: raw SQL through pgxpool, uppercase table names, `pgx.CollectRows` with a `scanX(ctx, row)` helper that checks `ctx.Err()` first.
- Style: gofmt; imports grouped std / external / `bearuang-go/...`; doc comments on exported types; `// ----` section banners; long calls wrapped one argument per line; interface assertions (`var _ Repository = (*PostgresRepository)(nil)`).
- Migrations: goose SQL `<timestamp>_<name>.sql` with `-- +goose Up` / `-- +goose Down`, created via `make migration`.
- Config comes from env (`DATABASE_URL`, `MIGRATIONS_DIR`, `POSTGRES_*`) with in-code defaults; no dotenv loading in Go.