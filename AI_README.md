# Golang Clean Architecture Template

## Purpose

Strict Clean Architecture implementation for Golang services. Enforces dependency inversion and layer isolation. Template only - discard sample implementation.

## Directory Intent

`cmd/`: Application entry points (web, worker, consumer, migrate)

`internal/`: Core business logic
- `entity/`: Domain models
- `usecase/`: Business logic
- `delivery/`: External adapters
- `infra/`: Concrete implementations. Database, cache, messaging, storage clients
- `config/`: Configuration
- `provider/`: External service clients initialization
- `dto/`: Data transfer objects between layers
- `converter/`: Entity-DTO transformations

`pkg/`: Shared utilities

`db/`: Database migrations

`test/`: Test suites

## Extension Rules

**New Feature**: Create usecase interface in `internal/usecase/feature/`. Implement in same package. Wire through dependency injection.

**New Entity**: Define in `internal/entity/`. Add repository interface. Implement in `internal/infra/repository/`.

**New Endpoint**: Create controller method in `internal/delivery/http/feature_controller.go`. Add route in `internal/delivery/http/route/`.

**New Consumer**: Create consumer in `internal/delivery/messaging/feature_consumer.go`. Add route in `internal/delivery/messaging/route/`.

**New External Service**: Add provider in `internal/provider/`. Reference in usecase through interface.

## Testing Strategy

Unit tests: Mock all dependencies. Test business logic isolation.

E2E tests: Test complete flows through delivery layer.

## Invariants

- Business logic stays in usecases
- All external calls through interfaces
- Configuration centralized
- Error handling standardized
- Context propagation mandatory
- Never edit `internal/mock/`, generate it with `make genereate`
- Never edit `api/`, generate it with `make swag`
