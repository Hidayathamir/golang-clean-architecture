---
name: idempotency
description: The consumer-side idempotency (de-duplication) subsystem — the x-idempotency-key header, the insert-if-absent usecase, the inbound wrapper, the message_idempotency table, and the cleaner binary. Use when working on Kafka consumer de-duplication, duplicate message handling, the idempotency table, or the workeridempotencycleaner.
---

# Idempotency (consumer de-duplication)

> **Enforcement**: ✅ **Partially enforced.** `make check-rules` verifies every `ConsumeEvent*` call in `route_messaging.go` passes a handler wrapped by `IdempotencyHandlerSingle`/`Batch` (allowlisted: the Elasticsearch sync consumers). It does **not** verify that producers set the `x-idempotency-key` header — that stays a manual requirement.

Kafka consumers are **at-least-once**, so every message carries an **`x-idempotency-key` header** (a UUID) and nearly every handler in `internal/inbound/messaging/route/route_messaging.go` is wrapped by a de-dup decorator.

## Pieces

| Piece | Location |
|-------|----------|
| Entity | `internal/entity/message_idempotency.go` (table `message_idempotency`, PK `idempotency_key`) |
| Migration | `db/migrations/20260724000000_create_table_message_idempotency.sql` |
| Repository | `internal/outbound/repository/idempotency_repository.go` (+ `_mw_logger.go`) |
| Usecase | `internal/usecase/idempotencyusecase/` — `InsertIfNotExists` (insert-if-absent) and `DeleteOlderThan` |
| Inbound wrapper | `internal/inbound/messaging/idempotency.go` — `IdempotencyHandlerSingle` / `IdempotencyHandlerBatch` |
| Cleaner binary | `cmd/workeridempotencycleaner/` — deletes rows older than 7 days on a ticker |
| Config | `idempotency.cleanup_interval_seconds` (default `3600`) — the ticker interval only |

## Behavior

- A handler that sees an **already-inserted key** is skipped and returns `nil` — offset committed, no retry, no DLQ (de-dup is not error routing).
- The wrapper is **fail-open**: a missing header or a failed check logs and proceeds to the handler anyway.
- `IdempotencyHandlerBatch` filters per record; if all are duplicates the batch handler is never called.
- The cleaner runs via `make run-workeridempotencycleaner`. **Run it in production** — without it the table grows unbounded.
- Retention is **hardcoded to 7 days** in `cmd/workeridempotencycleaner/main.go` (`7*24*time.Hour`). Only the ticker interval is configurable (`idempotency.cleanup_interval_seconds`).

Deep reference for the consumer-side retry/DLQ flow around this wrapper: the `kafka-error-handling` skill.