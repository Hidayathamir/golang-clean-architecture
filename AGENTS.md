# AGENTS.md

Always-on rules for this repo. Deep reference lives in skills — see **Skills** at the bottom.

## Code generation — run after interface changes

```bash
make generate         # Regenerate all mocks in internal/mock/ (deletes dir first, uses moq)
make swag             # Regenerate swagger docs in api/ (deletes dir first)
```

**Never edit `internal/mock/` or `api/` by hand.** Both are fully generated. Editing them is denied via `permissions` in `opencode.jsonc` (as are `logs/`).

**Enforced in CI** (`.github/workflows/rules.yml`) via `make check-codegen` — regenerates and fails if the result differs from what's committed.

## How to run the project

Configuration is `config.json` at the project root (Viper). Always run in this order:

1. **Start infrastructure** — `make docker-compose-up` (wait for it to finish). Postgres, Kafka, SigNoz, Redpanda Console.
2. **Run migrations** — `make migrate`.
3. **Run the app servers** (separate terminals, any order after migration):
   ```bash
   make run-webserver      # HTTP APIs       -> logs/webserver_log.jsonl
   make run-workerconsumer # Kafka consumers -> logs/workerconsumer_log.jsonl
   make run-workerproducer # Outbox producer -> logs/workerproducer_log.jsonl
   ```
   - `run-workerconsumer` handles async tasks (notifications, image uploads) from Kafka topics.
   - `run-workerproducer` polls the `outboxes` table and produces events to Kafka (at-least-once delivery with crash recovery via DB transactions).

**Always start services with the `make run-*` commands.** Never use `go run` or `nohup` directly — they bypass tuistory and cause stale-process conflicts. `run-*` commands are idempotent (rerunning kills the previous session).

## Database migrations

There is exactly **one** migration path:

```bash
make migrate   # cmd/migrate -> pkg/migrate -> sql-migrate (library), Up only
```

- Config comes from `config.json` (`database.*`) — the single source of truth for DB credentials.
- `sql-migrate` is used as a **Go library** (`pkg/migrate/migrate.go`), not as a CLI. There is no `dbconfig.yml` and no `new-migration`/`migrate-up`/`migrate-down`/`migrate-status` target.
- Migration files live in `db/migrations/`, named `<timestamp>_<name>.sql`, with `-- +migrate Up` / `-- +migrate Down` sections.
- `make migrate` only applies **Up** (all pending). It cannot roll back or show status; use `psql` against the `gorp_migrations` table for that.

## Architecture layers (top → bottom)

```
cmd/              Entry points: webserver, workerconsumer, consumer, migrate
internal/
  inbound/        HTTP controllers (package http), Kafka consumers (package messaging)
  dto/            Request/response DTOs between layers
  converter/      Entity ↔ DTO transforms
  usecase/*usecase/   Business logic (package userusecase, imageusecase, notifusecase)
  entity/         Domain models (GORM structs)
  outbound/       External implementations: repository, cache, messaging, storage, search
  provider/       Component initialization (DB, Redis, Kafka, S3, Elasticsearch, Fiber)
  dependency_injection/   Wire everything together
  config/         Viper-based config from config.json
pkg/              Shared utilities: logkit, errkit, telemetry, retrykit, validatorkit, caller, ctx, constant
```

## Decorator pattern — mandatory for all outbound interfaces

Every outbound interface (repository, cache, messaging, storage, search) **must** have a logger-middleware wrapper (`_mw_logger.go`) that:
- Wraps the real implementation in a `{Interface}MwLogger` struct with a `Next {Interface}` field
- Starts an OpenTelemetry span (`telemetry.Start`)
- Calls `retrykit.DBRetry` (for DB) or the real method
- Records errors to the span (`telemetry.RecordError`)
- Emits structured log via `logkit.LogMw(ctx, logrus.Fields{...}, err)`
- Satisfies the interface (`var _ Interface = &InterfaceMwLogger{}`)

Same pattern applies to usecases: `usecase → UserUsecaseMwLogger(UserUsecase)`.

Wire the decorator chain in `dependency_injection/` — real impl first, then logger wrapper:
```go
var userRepo repository.UserRepository
userRepo = repository.NewUserRepository(cfg)
userRepo = repository.NewUserRepositoryMwLogger(userRepo)  // decorates
```

**Enforced** via `make check-rules` (`scripts/checkrules`): every interface in `internal/outbound/` must have a sibling `<name>_mw_logger.go` declaring `<Name>MwLogger`. Runs in CI; run locally before pushing.

## File naming conventions

| Pattern | Example | Purpose |
|---------|---------|---------|
| `snake_case.go` | `user_controller.go` | All Go files |
| `*_mw_logger.go` | `user_repository_mw_logger.go` | Logger decorator |
| `*_test.go` | `user_usecase_test.go` | Tests (external test package: `package foo_test`) |
| `helper_test.go` | `helper_test.go` | Shared test fixtures per package |
| `mock/Mock{Interface}.go` | `mock/MockRepositoryUser.go` | Moq-generated mocks |

**Enforced** via `make check-rules` (`scripts/checkrules`): Go files under `internal/`, `pkg/`, `cmd/` must be snake_case (with optional `_test`), except generated mocks. Runs in CI.

## Error handling — always do both

```go
err = errkit.SetCode(err, http.StatusNotFound)            // set HTTP status
return errkit.AddFuncName(err, "package.(*Type).Method")  // enrich with call site
```

VSCode auto-runs `addfuncname` on `.go` file save. If you add new Go files outside VSCode, run `make add-func-name` manually.

## Comments — explain why, never what

- Comments are only for **why**: non-obvious intent, a tradeoff, a workaround, a subtle invariant.
- Never restate what the code does — names and structure already do that. A "what" comment stales as code changes.
- No obvious comments (`// increment i`), no section-header noise (`// === helpers ===`), no commented-out code.
- Godoc on exported types/funcs/methods **is** allowed — it's API surface, not narration.
- If body code needs a "why not obvious" note, keep it short and inline.

## Testing conventions

- **Test package**: `package {pkg}_test` (external, black-box)
- **DB mocking**: `go-sqlmock` + GORM with `PreferSimpleProtocol: true` (see `helper_test.go`)
- **Interface mocking**: Use moq-generated mocks from `internal/mock/` — set `.Func` fields on the mock struct to control behavior
- **Fixtures**: Place shared helpers in `helper_test.go` per package
- Unit tests: `make go-test` (tests `./internal/...`)
- E2E tests: `make go-e2e-test` (tests `./test/e2etest/...`)

## Route conventions (Fiber)

- Public routes: standard paths (`/api/users`)
- Authenticated routes: `_` prefix on the path segment (`/api/users/_current`, `/api/users/_follow`)
- Swagger annotations are on controller methods; Swagger served at `/swagger/*`

## Adding a new feature

1. Define entity in `internal/entity/`
2. Define DTOs in `internal/dto/`
3. Create interface + impl in `internal/usecase/<name>usecase/` — add `//go:generate moq` directive
4. Create outbound interface + impl in `internal/outbound/<layer>/` — add `//go:generate moq` directive
5. Add MwLogger decorators for both usecase and infra
6. Wire in `internal/dependency_injection/`
7. Create controller/consumer in `internal/inbound/<http|messaging>/`
8. Register routes in `internal/inbound/<http|messaging>/route/`
9. Run `make generate && make swag`

## Configuration

Single file: `config.json` at project root. Parsed via Viper in `internal/config/config.go`. Typed access through `config.Config` struct — always inject via constructor, never use a global.

## Observability & management tools

| Tool | URL | Description |
| :--- | :--- | :--- |
| **SigNoz** | http://localhost:3301 | Distributed traces, logs, application metrics. |
| **Redpanda Console** | http://localhost:8080 | Kafka topics, consumers, messages in real-time. |
| **Swagger** | http://localhost:3000/swagger | Interactive API documentation and testing. |
| **Kibana** | http://localhost:5601 | Explore/visualize data in Elasticsearch. |

## Skills — deep reference (loaded on demand)

- **`signoz-traces`** — ClickHouse MCP setup + trace/span/error SQL queries against SigNoz.
- **`otel-architecture`** — full OTel implementation: trace/log providers, decorator spans, HTTP + Kafka instrumentation, log correlation, OTLP export.
- **`kafka-error-handling`** — topic topology, retry/DLQ flow, error classification, consumer loop types, graceful shutdown.
- **`git-commit-convention`** — Conventional Commits 1.0.0: format, types, scopes, breaking changes, footers.