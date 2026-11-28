# AGENTS.md

Always-on rules and workflow for this repo. Deep reference lives in skills — see **Skills** at the bottom.

---

# Rules (enforced by hooks)

These are checked automatically. Run `make setup` once after cloning to activate the hooks.

## Code generation — run after interface changes

```bash
make generate         # Regenerate all mocks in internal/mock/ (deletes dir first, uses moq)
make swag             # Regenerate swagger docs in api/ (deletes dir first)
```

**Never edit `internal/mock/` or `api/` by hand.** Both are fully generated. Editing them is denied via `permissions` in `opencode.jsonc` (as are `logs/`).

**Enforced** by a shared pre-commit hook (`.githooks/pre-commit` → `make check-rules`). Stale generated code is caught by the **pre-push** hook (`.githooks/pre-push` → `make check-codegen`): regenerate + assert no diff, so a push is blocked until the regenerated `internal/mock`/`api` is committed.

## Commit messages — Conventional Commits 1.0.0

Every commit message must follow Conventional Commits — `<type>[(scope)][!]: <description>`, using types `feat`, `fix`, `build`, `chore`, `ci`, `docs`, `style`, `refactor`, `perf`, `test`, `revert`. Merge, revert and fixup/squash commits are exempt.

**Enforced** by the shared commit-msg hook (`.githooks/commit-msg` → `make check-commit-msg`, `scripts/checkcommitmsg`). Also available manually: `make check-commit-msg MSG_FILE=<file>`. Full spec: the `git-commit-convention` skill.

## Decorator pattern — mandatory for all outbound interfaces

Every outbound interface (repository, cache, messaging, storage, search) **must** have a logger-middleware wrapper (`_mw_logger.go`) that wraps the real implementation in a `{Interface}MwLogger` struct with a `Next {Interface}` field and satisfies the interface (`var _ Interface = &InterfaceMwLogger{}`).

Same pattern applies to usecases: `usecase → UserUsecaseMwLogger(UserUsecase)`.

Wire the decorator chain in `dependency_injection/` — real impl first, then logger wrapper:
```go
var userRepo repository.UserRepository
userRepo = repository.NewUserRepository(cfg)
userRepo = repository.NewUserRepositoryMwLogger(userRepo)  // decorates
```

**Enforced** via `make check-rules` (`scripts/checkrules`): every interface in `internal/outbound/` must have a sibling `<name>_mw_logger.go` declaring `<Name>MwLogger`, and every interface in `internal/usecase/` must have a sibling `<name>_usecase_mw_logger.go` declaring `<Name>MwLogger`.

Every decorated method follows the same four steps — start span, call the real method (wrapped in retry), record the error on the span, emit the log (see `internal/outbound/repository/follow_repository_mw_logger.go:26`):

```go
func (r *FollowRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, follow *entity.Follow) error {
	ctx, span := telemetry.Start(ctx) // span named after the caller's func
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error { // no-value; DBRetryWithData for value-returning
		return r.Next.Create(ctx, db, follow)
	})
	telemetry.RecordError(span, err) // span.RecordError + status=ERROR
	fields := logrus.Fields{"follow": follow}
	logkit.LogMw(ctx, fields, err) // Debug on success, Error on failure
	return err
}
```

- The retry helper matches the outbound layer: `DBRetry`/`DBRetryWithData` (repository), `CacheRetry`/`CacheRetryWithData` (cache), `MessagingRetry`/`MessagingRetryWithData` (messaging), `SearchRetry`/`SearchRetryWithData` (search), `StorageRetry`/`StorageRetryWithData` (storage). Use the `*WithData` form for value-returning methods. Each layer has its own tuning in `pkg/retrykit`. Usecases call `Next` directly — no retry.
- **Enforced** by `make check-rules`: every error-returning method in an outbound `_mw_logger.go` must contain `telemetry.Start(ctx)`, `defer span.End()`, a `retrykit.` call, `telemetry.RecordError(span, err)`, and `logkit.LogMw(ctx, ...)`; usecase decorators are exempt from the retry marker.
- The telemetry mechanics (span naming, error recording, log correlation) live in the `otel-architecture` skill → "Decorator pattern".

## File naming conventions

Go files under `internal/`, `pkg/`, `cmd/` must be snake_case (with optional `_test`), except generated mocks. Decorators end in `_mw_logger.go`; tests end in `_test.go`; per-package fixtures are `helper_test.go`; moq mocks are `mock/Mock{Interface}.go`.

**Enforced** via `make check-rules` (`scripts/checkrules`). Full table: the `architecture` skill.

## Testing conventions

Test files under `internal/` must declare `package {dir}_test` (external, black-box).

**Enforced** via `make check-rules` (`scripts/checkrules`). DB mocking, fixtures, and how to run tests: the `architecture` skill.

## Migration file naming

Migration files under `db/migrations/` must be named `<14-digit-timestamp>_<snake_case>.sql` and carry both `-- +migrate Up` and `-- +migrate Down` sections, Up first.

**Enforced** via `make check-rules` (`scripts/checkrules`). How to run migrations: see **Database migrations** under Workflow.

## Repository table/column names — use the constants

Repository implementations (`internal/outbound/repository/`) must never hardcode table or column names as string literals. Use `pkg/constant/table` (e.g. `Table(table.User)`) and `pkg/constant/column` (e.g. `column.CreatedAt.Eq(t)`). Entities' `TableName()` return the same `table.*` constant.

**Enforced** via `make check-rules` (`scripts/checkrules`): any string literal equal to a known table/column name in a repository impl fails. Decorators (`_mw_logger.go`) and tests are exempt.

## Struct field names — mirror the outbound interface

A struct field whose type is an outbound interface (`cache.*`, `repository.*`, `messaging.*`, `search.*`, `storage.*`) must be named after that interface — `NotificationPublisher cache.NotificationPublisher`, not `Publisher cache.NotificationPublisher`. Controllers and consumers keep their single short `Usecase` field; that is intentional and out of scope.

**Enforced** via `make check-rules` (`scripts/checkrules`): a mismatched field name in `internal/outbound`, `internal/usecase`, `internal/inbound`, or `internal/dependency_injection` fails. Decorators (`_mw_logger.go`) and tests are exempt.

## Dependency injection containers — one layer per file

`internal/dependency_injection/` has one file per layer — `inbound_dependency_injection.go`, `usecase_dependency_injection.go`, `outbound_dependency_injection.go`, `middleware_dependency_injection.go`. Each declares plural container structs (`Controllers`, `Hubs`, `Consumers`, `Usecases`, `Caches`, `Messagings`, `Repositories`, `Searches`, `Storages`, `Middlewares`) that hold **only** that layer's type, each with a matching `Setup<Container>` func.

A container must not mix layers — e.g. no `Usecases` field in a `Consumers` struct, no raw client or repository in `Usecases`. Pass cross-layer dependencies as `SetupX` parameters instead.

**Enforced** via `make check-rules` (`scripts/checkrules`): a container file declaring a struct field from a different layer fails.

## Error wrapping — `errkit.AddFuncName`

A handled error must be wrapped with `errkit.AddFuncName(err, "<call site>")` before it is returned, in the impl files of every layer: `internal/outbound/*`, `internal/usecase`, `internal/inbound`.

Bare propagation of an `err` identifier (e.g. the outer `return err` after a `DB.Transaction` closure whose body already wrapped) is exempt. Decorators (`_mw_logger.go`) and tests are exempt — decorators record the error instead of wrapping it.

**Enforced** via `make check-rules` (`scripts/checkrules`): a return inside an `if err != nil` guard that returns an unwrapped error value fails.

## Error handling — tagged via `pkg/errkit`

Errors carry their HTTP status, a client-facing message, and the call site that produced them, all handled centrally by `pkg/errkit`. Non-retryable errors for Kafka are marked with `errkit.WrapNonRetryable` / read back with `errkit.IsNonRetryable` (see the `kafka-error-handling` skill → "Error Classification"). Details live in the `architecture` skill → "Error handling".

## Comments — explain why, never what

- Comments are only for **why**: non-obvious intent, a tradeoff, a workaround, a subtle invariant.
- Never restate what the code does — names and structure already do that. A "what" comment stales as code changes.
- No obvious comments (`// increment i`), no section-header noise (`// === helpers ===`), no commented-out code.
- Godoc on exported types/funcs/methods **is** allowed — it's API surface, not narration.
- If body code needs a "why not obvious" note, keep it short and inline.

---

# Workflow

## How to run the project

**After cloning, run `make setup` once** — it verifies all required tools (hard-fails if any are missing), downloads Go dependencies, activates the shared pre-commit hook, and prints the next steps. It is idempotent, so rerun it any time.

Configuration is `settings/config.json` (Viper). Always run in this order:

1. **Start infrastructure** — `make docker-compose-up` (wait for it to finish). Postgres, Kafka (+ Zookeeper), Redis, Elasticsearch, Kibana, MinIO (S3 on `:9100`, console on `:9001`), Redpanda Console, and Grafana LGTM (one container: OpenTelemetry Collector + Tempo + Loki + Prometheus + Grafana; OTLP on `:4317`/`:4318`, UI on `:3000`). The app creates the `image` bucket on startup, so no manual bucket setup is needed.
2. **Run migrations** — `make migrate`.
3. **Run the app servers** — one command starts all four in the background:
   ```bash
   make run-all            # starts webserver + 3 workers (idempotent, returns immediately)
   make stop-all           # stops everything
   ```
   All four are part of a working system:

   | Server | Purpose | Required? |
   | :--- | :--- | :--- |
   | `run-webserver` | Serves the HTTP API (`:3000`) | **Yes** |
   | `run-workerconsumer` | Handles async Kafka tasks (notifications, image uploads) | **Yes** |
   | `run-workerproducer` | Polls the `outboxes` table and produces events to Kafka (at-least-once with crash recovery) | **Yes** |
   | `run-workeridempotencycleaner` | Prunes old `message_idempotency` rows on a ticker | Production only |

   Without `run-workerproducer`, outbox events never reach Kafka; without `run-workerconsumer`, they never get processed. Each can also be started individually (`make run-webserver`, etc.) — they are idempotent (rerunning restarts the same named session).

   Inspect sessions with `tuistory sessions`, attach with `tuistory attach -s golang-clean-arch-webserver`, and read live logs with `tail -f logs/*.jsonl`.

**Always start services with the `make run-*` / `make run-all` commands.** Never use `go run` or `nohup` directly — they bypass tuistory and cause stale-process conflicts. The `run-*` and `stop-all` targets **require tuistory** — if it is missing they fail with a clear error (run `make check-tools`). One-shot targets (`make migrate`, `make go-test`, `make go-e2e-test`) call plain `go run`/`go test` and do **not** use tuistory.

## Tooling prerequisites

`make setup` runs `make check-tools` (which **fails** if anything is missing), downloads Go deps, then activates git hooks. Required tools:

- **go** — compiler/runtime.
- **docker** — infrastructure via `make docker-compose-up`.
- **swag** — `make swag` (swagger docs). `go install github.com/swaggo/swag/cmd/swag@latest`.
- **moq** — `make generate` (mocks). `go install github.com/matryer/moq@latest`.
- **golangci-lint** — `make format`. `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`.
- **npm / Node** — `make diagram` (mmdc) and installing tuistory.
- **tuistory** — required by `make run-*` for idempotent restarts. `npm install -g tuistory`.

## Database migrations

There is exactly **one** migration path:

```bash
make migrate   # cmd/migrate -> pkg/migrate -> sql-migrate (library), Up only
```

- Config comes from `settings/config.json` (`database.*`) — the single source of truth for DB credentials.
- `sql-migrate` is used as a **Go library** (`pkg/migrate/migrate.go`), not as a CLI. There is no `dbconfig.yml` and no `new-migration`/`migrate-up`/`migrate-down`/`migrate-status` target.
- Migration files live in `db/migrations/`, named `<timestamp>_<name>.sql`, with `-- +migrate Up` / `-- +migrate Down` sections.
- `make migrate` only applies **Up** (all pending). It cannot roll back or show status; use `psql` against the `gorp_migrations` table for that.

## Finding logs for one request

Each server writes its own `logs/<service>_log.jsonl`, and every line carries the request's `trace_id`. The trace id is returned to the caller as the **`X-Trace-ID`** response header.

```bash
# 1. Make a request and read the trace id from the response.
curl -i http://localhost:3000/api/users/_current
# -> X-Trace-ID: 6e49d1047e1d41e0a295caad7bbecaeb

# 2. Merge every matching line from all log files into one, time-sorted, tagged file.
make trace-logs TRACE_ID=6e49d1047e1d41e0a295caad7bbecaeb
# -> logs/trace_logs_search_result.jsonl   (override path with OUT=<file>)
```

Full detail — the tool, Grafana/Tempo correlation, and the async-work caveat: the `tracing-logs` skill.

---

# Reference

Pointers to the deep-dive skills. Load the skill when you need the detail.

## Observability & management tools

| Tool | URL | Description |
| :--- | :--- | :--- |
| **Grafana (LGTM)** | http://localhost:3001 | Distributed traces (Tempo), logs (Loki), metrics (Prometheus). Login `admin` / `admin`. |
| **Redpanda Console** | http://localhost:8080 | Kafka topics, consumers, messages in real-time. |
| **Swagger** | http://localhost:3000/swagger | Interactive API documentation and testing. |
| **Kibana** | http://localhost:5601 | Explore/visualize data in Elasticsearch. |
| **MinIO Console** | http://localhost:9001 | Browse/inspect the S3 (`image`) bucket. Login `test` / `testtesttest`. |

## Idempotency

Kafka consumers are **at-least-once**, so every message carries an **`x-idempotency-key` header** and nearly every handler in `route_messaging.go` is wrapped by a de-dup decorator (skipped duplicates commit the offset, no retry/DLQ; the wrapper is fail-open; run `make run-workeridempotencycleaner` in production).

**Enforced** via `make check-rules` (`scripts/checkrules`): every `ConsumeEventSingle`/`ConsumeEventBatch`/`ConsumeEventRetry` call in `route_messaging.go` must pass a handler wrapped by `messaging.IdempotencyHandlerSingle`/`Batch`. Inherently idempotent consumers (the Elasticsearch sync ones) are allowlisted in `checkConsumerIdempotency`.

Full detail — the table, wrapper functions, and behavior: the `idempotency` skill.

## Architecture, layout, and adding features

Layers (`cmd`/`internal`/`pkg`), route conventions, config, and the step-by-step for adding a feature: the `architecture` skill.

---

# Skills — deep reference (loaded on demand)

- **`architecture`** — layers, file layout, route/config conventions, and how to add a feature.
- **`tracing-logs`** — find all logs for one request by trace id (local files or Grafana/Tempo).
- **`tempo-traces`** — query traces/spans/errors from Tempo via the Grafana UI or the Tempo HTTP API (TraceQL).
- **`otel-architecture`** — full OTel implementation: trace/log providers, decorator spans, HTTP + Kafka instrumentation, log correlation, OTLP export.
- **`kafka-error-handling`** — topic topology, retry/DLQ flow, error classification, consumer loop types, graceful shutdown.
- **`idempotency`** — consumer de-duplication: `x-idempotency-key`, insert-if-absent usecase, inbound wrapper, `message_idempotency` table, cleaner binary.
- **`git-commit-convention`** — Conventional Commits 1.0.0: format, types, scopes, breaking changes, footers.