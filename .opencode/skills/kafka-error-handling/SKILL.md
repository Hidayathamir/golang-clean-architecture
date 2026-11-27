---
name: kafka-error-handling
description: The Kafka topic topology, retry/DLQ flow, error classification, consumer loop types, and graceful shutdown. Use when working on Kafka consumers, retry topics, DLQ handling, or consumer error paths.
---

# Kafka Consumer Error Handling

## Topic Topology

Every primary topic has a corresponding **retry** and **DLQ** topic, mapped in `pkg/constant/topic/topic.go`.

| Primary | Retry | DLQ |
|---------|-------|-----|
| `image.uploaded` | `image.uploaded.retry` | `image.uploaded.dlq` |
| `image.liked` | `image.liked.retry` | `image.liked.dlq` |
| `image.commented` | `image.commented.retry` | `image.commented.dlq` |
| `user.followed` | `user.followed.retry` | `user.followed.dlq` |
| `notif` | `notif.retry` | `notif.dlq` |

A `x-retry-count` header tracks retry attempts (stored on the record, incremented each hop). Max retries configured via `kafka.consumer.max_retries`.

## Consumer Loop Types

Three loop patterns in `internal/inbound/messaging/consumer.go`:

| Function | Poll Mode | Record Count | Error Routing |
|----------|-----------|-------------|---------------|
| `ConsumeEventSingle` | Poll 1 record | 1 | DLQ if `IsNonRetryable`, else retry |
| `ConsumeEventBatch` | Poll all buffered | N | All → retry (with count=1) |
| `ConsumeEventRetry` | Poll 1 record | 1 | DLQ if exhausted or non-retryable, else retry |

## Error Classification — `pkg/errkit/nonretryable.go`

Consumer handlers wrap **JSON unmarshal failures** with `WrapNonRetryable` — these are structural errors that will never succeed on retry. Business logic errors (from usecases) are returned as-is, so they are retryable.

| Error Source | Wrapped? | Classification |
|-------------|----------|---------------|
| `json.Unmarshal` failure | `WrapNonRetryable` | Non-retryable → DLQ |
| Usecase returns error | No wrapper | Retryable → retry topic |

## Flow Diagrams

### Single Consumer (`ConsumeEventSingle`)

```mermaid
flowchart TD
    A[Poll 1 record from primary topic] --> B{Poll error?}
    B -->|Yes| C[Log error, continue poll]
    B -->|No| D[Get record]
    D --> E[Call handler]
    E --> F{Handler returned error?}
    F -->|No| G[Commit offset]
    F -->|Yes| H{IsNonRetryable?}
    H -->|Yes| I[Produce record → DLQ topic]
    H -->|No| J[Increment x-retry-count<br>Produce record → Retry topic]
    I --> G
    J --> G
    G --> K{Context cancelled?}
    K -->|No| A
    K -->|Yes| L[Stop consumer, close client]
```

### Batch Consumer (`ConsumeEventBatch`)

```mermaid
flowchart TD
    A[Poll all buffered records] --> B{Poll error?}
    B -->|Yes| C[Log error, continue poll]
    B -->|No| D[Get records]
    D --> E{Records > 0?}
    E -->|No| H
    E -->|Yes| F[Call batch handler]
    F --> G{Handler returned error?}
    G -->|No| H[Commit offset]
    G -->|Yes| I[ALL records → retry topic<br>with x-retry-count=1]
    I --> H
    H --> J{Context cancelled?}
    J -->|No| A
    J -->|Yes| K[Stop consumer, close client]
```

> **Note**: Batch consumers currently send *all* records in the batch to retry even if only one failed. A TODO in the code acknowledges this should be per-record granularity.

### Retry Consumer (`ConsumeEventRetry`)

```mermaid
flowchart TD
    A[Poll 1 record from retry topic] --> B{Poll error?}
    B -->|Yes| C[Log error, continue poll]
    B -->|No| D[Get record,<br>parse x-retry-count]
    D --> E[Call handler]
    E --> F{Handler returned error?}
    F -->|No| G[Commit offset]
    F -->|Yes| H{IsNonRetryable?}
    H -->|Yes| I[Produce record → DLQ topic]
    H -->|No| J{retryCount >= maxRetries?}
    J -->|Yes| I
    J -->|No| K[Increment x-retry-count<br>Produce record → Retry topic]
    I --> G
    K --> G
    G --> L{Context cancelled?}
    L -->|No| A
    L -->|Yes| M[Stop consumer, close client]
```

### Full End-to-End Error Flow

```mermaid
flowchart LR
    subgraph Primary["Primary Topic"]
        P[(image.uploaded)]
    end

    subgraph Consumer["Primary Consumer (ConsumeEventSingle)"]
        C[Poll & Handle]
        C -->|error| D{Classify}
        D -->|NonRetryable| D1[→ DLQ]
        D -->|Retryable| D2[→ Retry]
    end

    subgraph RetryTopic["Retry Topic"]
        R[(image.uploaded.retry)]
    end

    subgraph RetryConsumer["Retry Consumer (ConsumeEventRetry)"]
        RC[Poll & Handle]
        RC -->|error| RD{Classify}
        RD -->|NonRetryable| RD1[→ DLQ]
        RD -->|count >= max| RD1
        RD -->|count < max| RD2[→ Retry again]
        RD2 --> R
    end

    subgraph DLQTopic["DLQ Topic"]
        DQ[(image.uploaded.dlq)]
    end

    P --> Consumer
    D2 --> R
    D1 --> DQ
    R --> RetryConsumer
    RD1 --> DQ
```

## Graceful Shutdown

```mermaid
sequenceDiagram
    participant Main as main()
    participant Ctx as context
    participant Wg as sync.WaitGroup
    participant Cli as consumer loop

    Main->>Ctx: context.WithCancel()
    Main->>Wg: wg.Go(consumer loop)
    Cli->>Cli: Poll & process records
    Note over Main: SIGINT/SIGTERM received
    Main->>Ctx: cancel()
    Ctx->>Cli: ctx.Done()
    Cli->>Cli: Break poll loop
    Cli->>Cli: client.Close()
    Cli->>Wg: goroutine exits
    Main->>Wg: wg.Wait() — all consumers done
    Main->>Main: exit
```

Key points:
- `cmd/workerconsumer/main.go` sets up a shared `context.WithCancel` and `sync.WaitGroup`.
- All consumer goroutines are started via `wg.Go` (syntactic sugar for `go func()` + `wg.Add`/`wg.Done`).
- Each consumer loop checks `ctx.Err()` after every poll-commit cycle.
- On OS signal, `cancel()` triggers `ctx.Done()` which causes all loops to exit, then `wg.Wait()` blocks until all consumers finish their current processing.