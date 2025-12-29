# Production Logging Guide

ZGO supports high-volume log storage via ClickHouse and real-time error tracking via Sentry.

## Sentry Integration

Sentry captures exceptions and critical errors automatically.

### 1. Configuration
Add your Sentry DSN to `.env`:
```env
SENTRY_DSN=https://your-dsn@sentry.io/project
```

### 2. Usage
Any log with level `Error`, `Critical`, `Alert`, or `Emergency` will be automatically sent to Sentry.
```go
logger.Error("Database connection failed", map[string]any{"error": err.Error()})
```

---

## ClickHouse Integration

ClickHouse is used for high-performance, asynchronous log storage.

### 1. Configuration
Enable the ClickHouse handler in `.env`:
```env
LOG_CH_ENABLED=true
LOG_CH_LEVEL=info
LOG_CH_BATCH_SIZE=500
LOG_CH_INTERVAL_MS=1000
```

### 2. Strategy
ZGO uses an **Asynchronous Batching Strategy**:
- Logs are buffered in memory.
- Flushed to ClickHouse when buffer size reaches `LOG_CH_BATCH_SIZE` or after `LOG_CH_INTERVAL_MS`.

### 3. ClickHouse Schema (Recommended)
```sql
CREATE TABLE app_logs (
    time DateTime,
    level String,
    channel String,
    message String,
    request_id String,
    context String
) ENGINE = MergeTree()
ORDER BY (time, level)
```

## Recommended Production Setup

For massive scale, we recommend keeping log writing local to files and using **Vector** to ship logs to ClickHouse. However, for specialized needs or simpler deployments, the built-in `ClickHouseHandler` provides a direct, high-performance path.
