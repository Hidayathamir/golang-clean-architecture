# AGENTS_STARTUP.md

How to start this project from scratch. See [AGENTS.md](AGENTS.md) for architecture, codegen, and conventions.

## Quick commands

```bash
make docker-compose-up # Start all infra (Postgres, Kafka, ES, Redis, SigNoz...)
make migrate           # Run DB migrations
make run-webserver     # Web server (port 3000)
make run-workerconsumer        # Kafka consumer worker
make run-workerproducer # Outbox producer (polls outbox table, sends to Kafka)
make go-test           # Unit tests (./internal/...)
make go-e2e-test       # E2E tests (./test/e2etest/...)
make new-migration     # Create new migration (config: dbconfig.yml)
```

## Full startup flow (correct order)

Always run in this sequence:

1. **Start infra** — `make docker-compose-up` (wait for command to finish)
2. **Run migrations** — `make migrate`
3. **Start services** — `make run-webserver`, `make run-workerconsumer`, `make run-workerproducer` (any order after migration)

All `run-*` commands are **idempotent** — rerunning kills the previous session automatically. This works via `tuistory` (named background sessions). If `tuistory` is not installed, commands fall back to foreground `go run`.

## Startup script

To bring everything up from scratch in one shot:
```bash
make docker-compose-up && make migrate && make run-webserver & make run-workerconsumer & make run-workerproducer
```
