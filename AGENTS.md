# AGENTS.md

Always-on rules for this repo. Deep reference lives in skills — see **Skills** at the bottom.

## Code generation — run after interface changes

```bash
make generate         # Regenerate all mocks in internal/mock/ (deletes dir first, uses moq)
make swag             # Regenerate swagger docs in api/ (deletes dir first)
```

**Never edit `internal/mock/` or `api/` by hand.** Both are fully generated. Editing them is denied via `permissions` in `opencode.jsonc` (as are `logs/`).

**Enforced** by a shared pre-commit hook (`.githooks/pre-commit` → `make check-rules`). Run `make setup` once after cloning to activate it. `make check-codegen` (regenerate + assert no diff) is run manually when generated code changes.

## How to run the project

**After cloning, run `make setup` once** — activates the shared pre-commit hook that enforces project rules.

Configuration is `settings/config.json` (Viper). Always run in this order:

1. **Start infrastructure** — `make docker-compose-up` (wait for it to finish). Postgres, Kafka (+ Zookeeper), Redis, Elasticsearch, Kibana, Localstack, Redpanda Console, and SigNoz (with its own ClickHouse + Zookeeper).
2. **Run migrations** — `make migrate`.
3. **Run the app servers** (separate terminals, any order after migration):
   ```bash
   make run-webserver      # HTTP APIs       -> logs/webserver_log.jsonl
   make run-workerconsumer # Kafka consumers -> logs/workerconsumer_log.jsonl
   make run-workerproducer # Outbox producer -> logs/workerproducer_log.jsonl
   make run-workeridempotencycleaner # Idempotency cleanup -> logs/workeridempotencycleaner_log.jsonl
   ```
   - `run-workerconsumer` handles async tasks (notifications, image uploads) from Kafka topics.
   - `run-workerproducer` polls the `outboxes` table and produces events to Kafka (at-least-once delivery with crash recovery via DB transactions).
   - `run-workeridempotencycleaner` prunes old `message_idempotency` rows on a ticker (see **Idempotency** below).

**Always start services with the `make run-*` commands.** Never use `go run` or `nohup` directly — they bypass tuistory and cause stale-process conflicts. `run-*` commands are idempotent (rerunning kills the previous session).

## Database migrations

There is exactly **one** migration path:

```bash
make migrate   # cmd/migrate -> pkg/migrate -> sql-migrate (library), Up only
```

- Config comes from `settings/config.json` (`database.*`) — the single source of truth for DB credentials.
- `sql-migrate` is used as a **Go library** (`pkg/migrate/migrate.go`), not as a CLI. There is no `dbconfig.yml` and no `new-migration`/`migrate-up`/`migrate-down`/`migrate-status` target.
- Migration files live in `db/migrations/`, named `<timestamp>_<name>.sql`, with `-- +migrate Up` / `-- +migrate Down` sections.
- `make migrate` only applies **Up** (all pending). It cannot roll back or show status; use `psql` against the `gorp_migrations` table for that.

## Architecture layers (top → bottom)

```
cmd/              Entry points: webserver, workerconsumer, workerproducer, workeridempotencycleaner, migrate
internal/
  inbound/        HTTP controllers (package http), Kafka consumers (package messaging)
  dto/            Request/response DTOs between layers
  converter/      Entity ↔ DTO transforms
  usecase/*usecase/   Business logic (package userusecase, imageusecase, notifusecase, idempotencyusecase)
  entity/         Domain models (GORM structs, e.g. message_idempotency)
  outbound/       External implementations: repository, cache, messaging, storage, search
  provider/       Component initialization (DB, Redis, Kafka, S3, Elasticsearch, Fiber)
  dependency_injection/   Wire everything together
  config/         Viper-based config from settings/config.json
pkg/              Shared utilities: logkit, errkit, telemetry, retrykit, validatorkit, caller, ctx, constant
```

## Decorator pattern — mandatory for all outbound interfaces

Every outbound interface (repository, cache, messaging, storage, search) **must** have a logger-middleware wrapper (`_mw_logger.go`) that wraps the real implementation in a `{Interface}MwLogger` struct with a `Next {Interface}` field and satisfies the interface (`var _ Interface = &InterfaceMwLogger{}`).

Same pattern applies to usecases: `usecase → UserUsecaseMwLogger(UserUsecase)`.

Wire the decorator chain in `dependency_injection/` — real impl first, then logger wrapper:
```go
var userRepo repository.UserRepository
userRepo = repository.NewUserRepository(cfg)
userRepo = repository.NewUserRepositoryMwLogger(userRepo)  // decorates
```

**Enforced** via `make check-rules` (`scripts/checkrules`): every interface in `internal/outbound/` must have a sibling `<name>_mw_logger.go` declaring `<Name>MwLogger`. Runs in the pre-commit hook (activate with `make setup`).

The telemetry mechanics inside each generated method (span, `retrykit`, error recording, `logkit.LogMw`) live in the `otel-architecture` skill → "Decorator pattern".

## Idempotency

Kafka consumers are **at-least-once**, so every message carries an **`x-idempotency-key` header** and nearly every handler in `route_messaging.go` is wrapped by a de-dup decorator (skipped duplicates commit the offset, no retry/DLQ; the wrapper is fail-open; run `make run-workeridempotencycleaner` in production).

Full detail — the table, wrapper functions, and behavior: the `idempotency` skill.

## File naming conventions

| Pattern | Example | Purpose |
|---------|---------|---------|
| `snake_case.go` | `user_controller.go` | All Go files |
| `*_mw_logger.go` | `user_repository_mw_logger.go` | Logger decorator |
| `*_test.go` | `user_usecase_test.go` | Tests (external test package: `package foo_test`) |
| `helper_test.go` | `helper_test.go` | Shared test fixtures per package |
| `mock/Mock{Interface}.go` | `mock/MockRepositoryUser.go` | Moq-generated mocks |

**Enforced** via `make check-rules` (`scripts/checkrules`): Go files under `internal/`, `pkg/`, `cmd/` must be snake_case (with optional `_test`), except generated mocks. Runs in the pre-commit hook (activate with `make setup`).

## Error handling — always do both

```go
err = errkit.SetCode(err, http.StatusNotFound)            // set HTTP status
return errkit.AddFuncName(err, "package.(*Type).Method")  // enrich with call site
```

`make add-func-name` fills in the caller's function name for `errkit.AddFuncName`. The pre-commit hook runs it automatically (on any staged `.go` file) and re-stages the result, so you normally do nothing. Run it manually only if you commit with `--no-verify`. (For instant feedback on save, `.vscode/settings.json` runs the same command — optional, not required.)

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

Single file: `settings/config.json`. Parsed via Viper in `internal/config/config.go`. Typed access through `config.Config` struct — always inject via constructor, never use a global.

## Diagrams

`docs/architecture.mermaid` is the source; `docs/architecture.png` is generated (used by README).

- Regenerate: `make diagram` (fetches `mmdc` via `npx`, needs Node).
- The pre-commit hook auto-regenerates the PNG and stages it whenever a `docs/*.mermaid` change is committed.
- `docs/puppeteer-config.json` disables the Chromium sandbox (needed on Linux).

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
- **`idempotency`** — consumer de-duplication: `x-idempotency-key`, insert-if-absent usecase, inbound wrapper, `message_idempotency` table, cleaner binary.
- **`git-commit-convention`** — Conventional Commits 1.0.0: format, types, scopes, breaking changes, footers.