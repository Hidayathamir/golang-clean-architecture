---
name: architecture
description: The layers, file layout, and conventions of this clean-architecture Go project — cmd/internal/pkg layout, inbound/usecase/outbound flow, route and config conventions, and the step-by-step for adding a feature. Use when navigating the codebase, deciding where code belongs, or adding a new feature or layer.
---

# Architecture and conventions

> **Enforcement**: ➖ **Mostly reference.** A few items here are enforced by `make check-rules`; those are called out inline. The enforced rules are summarized in AGENTS.md.

This is the deep reference for how the project is laid out. AGENTS.md carries the always-on rules; this skill explains the structure behind them.

## Layers (top → bottom)

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

The flow of a request is always **inbound → usecase → outbound**, with `converter` translating between `entity` and `dto` at the boundaries. Nothing skips a layer.

## Configuration

Single file: `settings/config.json`. Parsed via Viper in `internal/config/config.go`. Typed access through the `config.Config` struct — **always inject via constructor, never use a global**. The same file is the single source of truth for DB credentials, Kafka brokers, log level, and the OTLP endpoint.

## Route conventions (Fiber)

- Public routes: standard paths (`/api/users`).
- Authenticated routes: `_` prefix on the path segment (`/api/users/_current`, `/api/users/_follow`).
- Swagger annotations live on controller methods; Swagger is served at `/swagger/*`.
- Middleware chain is applied at group level: `OtelFiberMiddleware` then `TraceIDMiddleware` (sets the `X-Trace-ID` response header) then auth. See the `otel-architecture` skill for the middleware details.

## Adding a new feature

1. Define entity in `internal/entity/`.
2. Define DTOs in `internal/dto/`.
3. Create interface + impl in `internal/usecase/<name>usecase/` — add the `//go:generate moq` directive.
4. Create outbound interface + impl in `internal/outbound/<layer>/` — add the `//go:generate moq` directive.
5. Add MwLogger decorators for **both** usecase and infra (see the `otel-architecture` skill → "Decorator pattern"; the structure itself is an enforced rule in AGENTS.md).
6. Wire in `internal/dependency_injection/` (real impl first, then logger wrapper).
7. Create controller/consumer in `internal/inbound/<http|messaging>/`.
8. Register routes in `internal/inbound/<http|messaging>/route/`.
9. Run `make generate && make swag`.

## File naming (enforced)

| Pattern | Example | Purpose |
|---------|---------|---------|
| `snake_case.go` | `user_controller.go` | All Go files |
| `*_mw_logger.go` | `user_repository_mw_logger.go` | Logger decorator |
| `*_test.go` | `user_usecase_test.go` | Tests (external test package: `package foo_test`) |
| `helper_test.go` | `helper_test.go` | Shared test fixtures per package |
| `mock/Mock{Interface}.go` | `mock/MockRepositoryUser.go` | Moq-generated mocks |

Enforced by `make check-rules` (`scripts/checkrules`) in the pre-commit hook: Go files under `internal/`, `pkg/`, `cmd/` must be snake_case (with optional `_test`), except generated mocks.

## Testing conventions

- **Test package**: `package {pkg}_test` (external, black-box). Enforced by `make check-rules`.
- **DB mocking**: `go-sqlmock` + GORM with `PreferSimpleProtocol: true` (see `helper_test.go`).
- **Interface mocking**: moq-generated mocks from `internal/mock/` — set `.Func` fields on the mock struct to control behavior.
- **Fixtures**: shared helpers in `helper_test.go` per package.
- Unit tests: `make go-test` (tests `./internal/...`).
- E2E tests: `make go-e2e-test` (tests `./test/e2etest/...`).

## Error handling

Always do both:

```go
err = errkit.SetCode(err, http.StatusNotFound)            // set HTTP status
return errkit.AddFuncName(err, "package.(*Type).Method")  // enrich with call site
```

`make add-func-name` fills in the caller's function name for `errkit.AddFuncName`. The pre-commit hook runs it automatically on any staged `.go` file and re-stages the result.

## Diagrams

`docs/architecture.mermaid` is the source; `docs/architecture.png` is generated (used by README).

- Regenerate: `make diagram` (fetches `mmdc` via `npx`, needs Node).
- The pre-commit hook auto-regenerates the PNG and stages it whenever a `docs/*.mermaid` change is committed.
- `docs/puppeteer-config.json` disables the Chromium sandbox (needed on Linux).