---
name: seeding
description: The database seeding subsystem — the `make seed` target, cmd/seed -> pkg/seed, the db/seeds/ SQL files, how it differs from migrations, and the fresh-environment gotchas for MinIO and Kafka. Use when working on seed data, dev fixtures, resetting a local database, or when `make seed` fails right after a clean docker-compose reset.
---

# Seeding (dev data)

> **Enforcement**: ➖ **Reference only.** Nothing here is checked by `make check-rules`; seeds are a developer convenience, not a versioned system.

Seeding is **separate from migrations**. Migrations own the schema and are versioned; seeds own data and are meant to be re-runnable against a reset database. Keeping them apart stops dev fixtures from leaking into production migrations.

## Pieces

| Piece | Location |
|-------|----------|
| Entry point | `cmd/seed/main.go` |
| Logic | `pkg/seed/seed.go` — reads `cfg.GetDatabaseSeeds()`, runs each `*.sql` in filename order via `db.Exec` |
| Seed SQL | `db/seeds/` — one file per concern, e.g. `001_seed_users.sql` |
| Config | `database.seeds` in `settings/config.json` (default `db/seeds`) |
| Config key/getters | `internal/config/configkey.go` (`DatabaseSeeds`), `internal/config/config.go` (`GetDatabaseSeeds`) |
| Make target | `make seed` (`go run cmd/seed/main.go`) |

## Behavior

- `pkg/seed` executes only the `-- +migrate Up` section of each file; the `-- +migrate Down` section is ignored (kept for symmetry with migration files, but never run).
- Files are applied in **lexical filename order** — prefix with zero-padded numbers (`001_`, `002_`) to control ordering.
- **No tracking table.** Seeds re-run every time; each statement must be idempotent (`on conflict do nothing`) or the DB must be reset first.
- **Not built on `sql-migrate`** on purpose: `FileMigrationSource` shares the `gorp_migrations` table with real migrations and fails with `unknown migration in database` when it sees migration filenames it did not apply. `pkg/seed` uses a plain read-and-exec loop instead.

## Standard local workflow

```bash
make docker-compose-up-clean   # reset infra (down -v, then up)
make migrate                   # schema (34 migrations)
make seed                      # dev data
make run-all                   # start the 4 servers
```

## Fresh-environment gotchas

After `docker compose down -v`, two things are wiped and can make `make go-e2e-test` fail even though the seed worked:

- **MinIO bucket `image` missing.** `ensureBucketExists` (`internal/provider/aws_provider.go`) creates it on boot with 5 retries, then **panics** — but if the webserver started before MinIO was healthy it can fail silently. Symptom: presigned PUT returns `404`. Fix — create it and restart the webserver:
  ```bash
  AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=testtesttest AWS_DEFAULT_REGION=us-east-1 \
    aws --endpoint-url http://localhost:9100 s3 mb s3://image
  ```
- **Kafka `UNKNOWN_TOPIC_ID`.** Topics are auto-created on first produce, but a producer started before `down -v` caches stale metadata. Symptom: outbox rows stay `pending`, e2e times out on "outbox to drain". Fix — restart the servers so metadata is refreshed:
  ```bash
  make stop-all && make run-all
  ```

Order that avoids both: infra up → wait for Postgres/MinIO/Kafka healthy → `make migrate` → `make seed` → `make run-all`.