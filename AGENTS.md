# Repository Guidelines

## Project Structure & Module Organization
The `cmd/` folder hosts entrypoints: `web` for the HTTP server, `worker` for async jobs, and `migrate` for schema tasks. Core business logic follows clean architecture within `internal/`: `delivery` adapters, `usecase` orchestrators, `entity` definitions, `repository` persistence, and `gateway` integrations. Shared helpers sit in `pkg/`, database artifacts in `db/migrations`, and end-to-end scenarios under `integrationtest/`.

## Build, Test, and Development Commands
Use `make run` to start the HTTP server via `cmd/web/main.go`; `make run-worker` runs background workers. `make go-test` executes unit tests in `internal/...`, while `make go-integration-test` targets `integrationtest/...`. Run `make run-clean` for a full cycle of linting, code generation, Swagger refresh, and app startup. `make docker-compose` resets and brings up local dependencies, and `make check-tools` verifies required CLIs.

## Coding Style & Naming Conventions
Code must stay `gofmt`-clean; rely on the default Go tab indentation. Run `make format` (golangci-lint) before submitting to enforce lint rules such as `errcheck` and `gocyclo`. Keep package names lowercase (e.g., `repository`, `usecase`), Go files in `snake_case.go`, and exported identifiers following Go’s PascalCase/Initialism guidance (e.g., `HTTPServer`). Swagger definitions are generated to `api/`, so avoid manual edits there.

## Testing Guidelines
Prefer table-driven unit tests next to the code under test with `_test.go` suffix and descriptive function names like `TestUserRepository_FindByID`. Integration tests should live in `integrationtest/` and call real adapters via Dockerized services. Maintain coverage for new use cases and update fixtures when schemas change; run both `make go-test` and `make go-integration-test` before opening a PR.

## Commit & Pull Request Guidelines
Follow the existing short, imperative commit style (`caller info`, `update note`) and keep messages under 72 characters. Group related changes per commit and reference issues in the body when applicable. Pull requests should summarize architecture impact, list verification steps (tests run, migrations applied), and include Swagger diffs or log snippets when behavior changes. Link tracking tickets and attach screenshots for HTTP interface updates.
