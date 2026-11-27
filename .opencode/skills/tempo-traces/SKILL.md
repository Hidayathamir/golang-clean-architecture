---
name: tempo-traces
description: Query distributed traces, spans, errors and latency from Tempo via the Grafana UI or the Tempo HTTP API. Use when auditing requests, debugging slow or failing traces, or inspecting trace data.
---

# Tempo Traces — Grafana + Tempo HTTP API

> **Enforcement**: ➖ **Not applicable**. Operational reference (query tooling), not a rule.

Traces are stored in **Tempo**, bundled inside the `grafana/otel-lgtm` container (no separate ClickHouse). There are two ways to query them:

- **Grafana UI**: http://localhost:3001 — login `admin` / `admin`. Explore → Tempo datasource. This is the interactive path (trace search, span waterfall, correlated logs/metrics).
- **Tempo HTTP API**: `http://localhost:3200` (query frontend) — scriptable, no login.
- **Required**: `docker-compose` must be running (`make docker-compose-up`).

## Tempo HTTP API

The API speaks **TraceQL** and the legacy search params. Useful endpoints:

```bash
# Search traces by service and time range (returns trace IDs)
curl -s "http://localhost:3200/api/search?q=%7B%20resource.service.name%3D%22webserver%22%20%7D&limit=20"

# Fetch one trace by ID (protobuf/JSON)
curl -s "http://localhost:3200/api/traces/<trace-id>"

# Search by tag (e.g. HTTP route)
curl -s "http://localhost:3200/api/search?tags=http.route%3D/api/v1/users"
```

## TraceQL essentials

Run TraceQL from Grafana Explore or pass it as `q=` to `/api/search`.

```traceql
-- Recent traces for a service
{ resource.service.name = "webserver" }

-- Spans with an error status
{ status = error }

-- Slow spans (> 1s)
{ duration > 1s }

-- Combine: slow spans in one service
{ resource.service.name = "workerconsumer" && duration > 500ms }

-- Spans belonging to a specific trace
{ trace:id = "<trace-id>" }
```

## Finding one request end to end

1. Make the request and read the `X-Trace-ID` response header (see the `tracing-logs` skill).
2. Open http://localhost:3001 → Explore → Tempo, or call `GET /api/traces/<trace-id>`.
3. The trace view shows the span waterfall; log lines sharing the same `trace_id` are one click away via the Loki datasource correlation.

## Migrating from the old ClickHouse queries

| Old (SigNoz ClickHouse) | New (Tempo) |
|---|---|
| `SELECT ... FROM signoz_traces.signoz_index_v2 WHERE serviceName = ...` | `{ resource.service.name = "..." }` |
| `WHERE traceID = '<id>'` | `{ trace:id = "<id>" }` or `GET /api/traces/<id>` |
| `WHERE statusCode = 2` | `{ status = error }` |
| `WHERE durationNano > 1000000000` | `{ duration > 1s }` |

For how spans are produced and correlated with logs, see the `otel-architecture` skill.