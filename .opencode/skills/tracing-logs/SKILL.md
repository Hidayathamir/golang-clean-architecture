---
name: tracing-logs
description: Find every log line for one request by its trace id, from the per-service JSONL log files or Grafana/Tempo. Use when debugging a request end-to-end, correlating logs across servers, or when asked how to trace a request through the system.
---

# Tracing logs by trace id

> **Enforcement**: ➖ **Not applicable**. Operational reference (debugging workflow), not a rule.

Every log line the app writes carries the request's `trace_id`, and the same id is returned to the caller as the **`X-Trace-ID`** response header. That id is the handle for pulling one request's logs out of every server at once.

## The fast path — local JSONL files

Each server writes its own log file through `scripts/dev/runsvc.sh`:

| Server | Log file |
| :--- | :--- |
| webserver | `logs/webserver_log.jsonl` |
| workerconsumer | `logs/workerconsumer_log.jsonl` |
| workerproducer | `logs/workerproducer_log.jsonl` |
| workeridempotencycleaner | `logs/workeridempotencycleaner_log.jsonl` |

Workflow:

```bash
# 1. Make a request and read the trace id from the response header.
curl -i http://localhost:3000/api/users/_current
# -> X-Trace-ID: 6e49d1047e1d41e0a295caad7bbecaeb

# 2. Merge every matching line from all log files into one, time-sorted, tagged file.
make trace-logs TRACE_ID=6e49d1047e1d41e0a295caad7bbecaeb
# -> logs/trace_logs_search_result.jsonl   (override path with OUT=<file>)
```

`make trace-logs` runs `scripts/tracelogs`, which:

- scans every `logs/*_log.jsonl`;
- keeps lines whose `trace_id` matches (invalid JSON lines are skipped);
- adds a **`service`** field derived from the file name (`webserver`, `workerconsumer`, ...);
- sorts all lines by their `time` field;
- writes one JSONL file (`logs/trace_logs_search_result.jsonl` by default).

Exit codes: `2` when the trace id is missing, `1` when no line matches.

Why a tool and not `grep`: a single request spans several processes (webserver → outbox → workerproducer → Kafka → workerconsumer), each writing its own file. `grep '<id>' logs/*.jsonl` works, but does not merge, tag, or sort; the tool does. Prefer it whenever a request crosses services.

## The full path — Grafana / Tempo

When the local files are gone, rotated, or you want spans and logs in one view, open Grafana at http://localhost:3001 (login `admin` / `admin`) and search Tempo by the same `trace_id`, or query the Tempo HTTP API. See the `tempo-traces` skill for the TraceQL equivalents and log correlation via Loki.

## A caveat on async work

The producer and consumer **extract the original trace context** from the outbox row (see the `otel-architecture` skill), so a message triggered by a request keeps that request's trace id in its Kafka headers. But if you are looking at a *worker's own* poll cycle (workerproducer scanning pending outboxes), those lines carry a **new** trace id from the poll, not the original request's. Search by the id from the HTTP response to follow the request; search the worker's own id to inspect the poll.