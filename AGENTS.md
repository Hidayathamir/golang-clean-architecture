# Repository Guidelines

## Project Structure & Module Organization
The entrypoints under `cmd/` drive the deployables: `web` serves HTTP traffic, `worker` handles asynchronous jobs, and `migrate` manages database schema tasks. Domain logic follows a clean architecture split within `internal/`, where `delivery` exposes adapters, `usecase` coordinates workflows, `entity` defines models, `repository` persists data, and `gateway` integrates upstream services. Shared helpers reside in `pkg/`, database migrations live in `db/migrations`, and end-to-end flows are captured in `integrationtest/`.

### Folder Reference
- `cmd/`: Go entrypoints for runnable binaries.
  - `cmd/web`: HTTP server bootstrap, middlewares, and route wiring.
  - `cmd/worker`: Background job consumers plus scheduler wiring.
  - `cmd/migrate`: Database migration runner using files from `db/migrations`.
- `internal/`: Core business layers split by Clean Architecture roles.
  - `internal/delivery`: HTTP/worker handlers, request validation, presenter glue.
  - `internal/usecase`: Application services orchestrating repositories and gateways.
  - `internal/entity` & `internal/model`: Domain structs (entity) and transport-specific projections (model).
  - `internal/repository`: Interfaces and implementations over persistence engines.
  - `internal/gateway`: Outbound integrations such as third-party APIs or queues.
  - `internal/config`: Runtime configuration loaders/parsers.
  - `internal/mock`: Test doubles for usecases, repositories, and gateways.
- `pkg/`: Reusable utilities (logging, error helpers, middleware primitives) that stay framework-agnostic.
- `db/migrations`: SQL migration files consumed by `cmd/migrate`; structure mirrors goose/flyway style timestamps.
- `api/`: Generated Swagger specs and client/server stubs—never hand-edit.
- `integrationtest/`: Black-box flows that spin up real adapters via docker-compose.
- `signoz/`: Observability sandbox configs (docker-compose fragments, dashboards).
- Project-wide metadata lives at the repo root: `Makefile` (automation), `README.md`, `README_my_note.md`, and `run_app.md` (process docs), plus `.env` samples under `config/` when present.

## Build, Test, and Development Commands
Run `make run` to start `cmd/web/main.go` and `make run-worker` for background processing. `make go-test` executes unit tests across `internal/...`, while `make go-integration-test` exercises scenarios in `integrationtest/...`. Use `make run-clean` for a full lint, generate, Swagger, and boot cycle, `make docker-compose` to reset local dependencies, and `make check-tools` to confirm required CLIs before onboarding.

## Coding Style & Naming Conventions
All Go code must stay gofmt-clean with default tab indentation; run `make format` (golangci-lint) before submitting. Keep package names lowercase (`repository`, `usecase`), files in `snake_case.go`, and exported identifiers following Go’s PascalCase with preserved initialisms (e.g., `HTTPServer`). Swagger artifacts under `api/` are generated and should not be edited manually.

## Testing Guidelines
Prefer table-driven unit tests adjacent to the code under test with `_test.go` suffix and names like `TestUserRepository_FindByID`. Cover new use cases and update fixtures as schemas evolve. Execute both `make go-test` and `make go-integration-test` before opening a PR to validate unit and integration coverage.

## Commit & Pull Request Guidelines
Commits follow a short, imperative style (`caller info`, `update note`) under 72 characters. Group related changes per commit, and reference issues or tickets in the body when relevant. Pull requests should summarize architectural impact, list verification steps (tests run, migrations applied), and attach Swagger diffs or log snippets when behavior shifts. Include screenshots for any HTTP interface updates to aid review.

## Environment & Configuration Tips
Run `make docker-compose` after checkout or schema changes to refresh local services. Keep environment variables in sync with `.env` samples and avoid committing secrets. When introducing new CLIs or generators, document installation steps in `README.md` and add checks to `make check-tools` where possible.
