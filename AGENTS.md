# AGENTS.md

## Quick commands

```bash
make run              # Web server (port 3000)
make run-worker       # Kafka consumer worker
make run-cron         # Cron scheduler
make go-test          # Unit tests (./internal/...)
make go-e2e-test      # E2E tests (./test/e2etest/...)
make docker-compose   # Start Postgres, Kafka, SigNoz, AKHQ (docker compose down -v first)
make migrate          # Run DB migrations
make migrate-new      # Create new migration (config: dbconfig.yml)
```

## Code generation — run after interface changes

```bash
make generate         # Regenerate all mocks in internal/mock/ (deletes dir first, uses moq)
make swag             # Regenerate swagger docs in api/ (deletes dir first)
```

**Never edit `internal/mock/` or `api/` by hand.** Both are fully generated. `.aiexclude` confirms this.

## Architecture layers (top → bottom)

```
cmd/              Entry points: web, worker, consumer, migrate
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

## File naming conventions

| Pattern | Example | Purpose |
|---------|---------|---------|
| `snake_case.go` | `user_controller.go` | All Go files |
| `*_mw_logger.go` | `user_repository_mw_logger.go` | Logger decorator |
| `*_test.go` | `user_usecase_test.go` | Tests (external test package: `package foo_test`) |
| `helper_test.go` | `helper_test.go` | Shared test fixtures per package |
| `mock/Mock{Interface}.go` | `mock/MockRepositoryUser.go` | Moq-generated mocks |

## Error handling — always do both

```go
err = errkit.SetCode(err, http.StatusNotFound)            // set HTTP status
return errkit.AddFuncName(err, "package.(*Type).Method")  // enrich with call site
```

VSCode auto-runs `addfuncname` on `.go` file save. If you add new Go files outside VSCode, run `make add-func-name` manually.

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
