# AGENTS.md

> **Startup instructions moved to [AGENTS_STARTUP.md](AGENTS_STARTUP.md).**

## Code generation — run after interface changes

```bash
make generate         # Regenerate all mocks in internal/mock/ (deletes dir first, uses moq)
make swag             # Regenerate swagger docs in api/ (deletes dir first)
```

**Never edit `internal/mock/` or `api/` by hand.** Both are fully generated. `.aiexclude` confirms this.

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

## SigNoz tracing — ClickHouse MCP

The opencode agent queries traces/spans directly from **ClickHouse** via an MCP server (no SigNoz API key needed). Configured in `~/.config/opencode/opencode.jsonc`.

- **ClickHouse**: `localhost:8123`, user `default`, no password
- **MCP server**: `@arvoretech/clickhouse-mcp` (read-only, runs via `npx`)
- **Tools**: `read_query`, `list_tables`, `describe_table`, `list_databases`
- **Required**: `docker-compose` must be running (`make docker-compose-up`)

### SigNoz ClickHouse tables for traces

| Table | Contents |
|---|---|
| `signoz_traces.signoz_index_v2` | Main traces/spans index (serviceName, name, durationNano, statusCode, tags, httpMethod, httpRoute, etc.) |
| `signoz_traces.signoz_spans` | Span details (traceID + model JSON blob) |
| `signoz_traces.signoz_error_index_v2` | Error traces |

### Common trace queries

```sql
-- Recent traces for a service
SELECT traceID, spanID, serviceName, name, durationNano,
       timestamp, statusCode, httpMethod, httpRoute
FROM signoz_traces.signoz_index_v2
WHERE serviceName = '<service>'
  AND timestamp > now() - INTERVAL 1 HOUR
ORDER BY timestamp DESC
LIMIT 20

-- Get full trace by ID
SELECT * FROM signoz_traces.signoz_index_v2
WHERE traceID = '<trace-id>'
ORDER BY timestamp

-- Error traces last hour
SELECT traceID, spanID, serviceName, name, durationNano,
       timestamp, statusMessage, exceptionType, exceptionMessage
FROM signoz_traces.signoz_index_v2
WHERE statusCode = 2
  AND timestamp > now() - INTERVAL 1 HOUR
ORDER BY timestamp DESC

-- Slow traces (p99+ duration)
SELECT traceID, spanID, serviceName, name, durationNano,
       timestamp, httpMethod, httpRoute
FROM signoz_traces.signoz_index_v2
WHERE durationNano > 1000000000  -- > 1s
  AND timestamp > now() - INTERVAL 1 HOUR
ORDER BY durationNano DESC
LIMIT 20
```

## Configuration

Single file: `config.json` at project root. Parsed via Viper in `internal/config/config.go`. Typed access through `config.Config` struct — always inject via constructor, never use a global.
