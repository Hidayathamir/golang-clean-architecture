---
name: signoz-traces
description: Query distributed traces, spans, errors and latency from SigNoz via the ClickHouse MCP server. Use when auditing requests, debugging slow or failing traces, or inspecting trace data.
---

# SigNoz Traces — ClickHouse MCP

> **Enforcement**: ➖ **Not applicable**. Operational reference (query tooling), not a rule.

The opencode agent queries traces/spans directly from **ClickHouse** via an MCP server (no SigNoz API key needed). Configured in `~/.config/opencode/opencode.jsonc`.

- **ClickHouse**: `localhost:8123`, user `default`, no password
- **MCP server**: `@arvoretech/clickhouse-mcp` (read-only, runs via `npx`)
- **Tools**: `read_query`, `list_tables`, `describe_table`, `list_databases`
- **Required**: `docker-compose` must be running (`make docker-compose-up`)

## Tables for traces

| Table | Contents |
|---|---|
| `signoz_traces.signoz_index_v2` | Main traces/spans index (serviceName, name, durationNano, statusCode, tags, httpMethod, httpRoute, etc.) |
| `signoz_traces.signoz_spans` | Span details (traceID + model JSON blob) |
| `signoz_traces.signoz_error_index_v2` | Error traces |

## Common queries

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

For how spans are produced and correlated with logs, see the `otel-architecture` skill.