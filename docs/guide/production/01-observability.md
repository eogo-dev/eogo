# Observability

> Tracing, Metrics, and Logging for production monitoring.

## Overview

Eogo provides three pillars of observability:

| Pillar | Implementation | Purpose |
|--------|---------------|---------|
| Tracing | OpenTelemetry | Request flow tracking |
| Metrics | Prometheus | Performance monitoring |
| Logging | Structured logs | Debugging & auditing |

## Distributed Tracing

### Configuration

```go
// internal/infra/tracing/tracing.go
type Config struct {
    Enabled     bool
    ServiceName string
    Environment string
    Endpoint    string  // OTLP endpoint
    Insecure    bool
    SampleRate  float64
    Debug       bool    // Use stdout exporter
}
```

Environment variables:

```bash
TRACING_ENABLED=true
TRACING_ENDPOINT=localhost:4317
TRACING_SAMPLE_RATE=1.0
```

### HTTP Middleware

Automatically traces all HTTP requests:

```go
// Registered in HTTP kernel
r.Use(tracing.Middleware(cfg.App.Name))
r.Use(tracing.InjectTraceID())  // Adds X-Trace-ID header
```

Captured attributes:
- `http.method`, `http.route`, `http.status_code`
- `http.client_ip`, `http.user_agent`
- Request/response size
- Errors

### Database Tracing

GORM plugin for SQL tracing:

```go
// Enable database tracing
tracing.WithTracing(db, "eogo")
```

Captured attributes:
- `db.operation` (SELECT, INSERT, UPDATE, DELETE)
- `db.table`
- `db.rows_affected`
- Query duration

### Manual Spans

Create custom spans for business operations:

```go
func (s *service) ProcessOrder(ctx context.Context, orderID uint) error {
    ctx, span := s.tracer.Start(ctx, "ProcessOrder")
    defer span.End()
    
    // Add attributes
    span.SetAttributes(attribute.Int("order.id", int(orderID)))
    
    // Add events
    tracing.AddEvent(ctx, "validation_started")
    
    if err := s.validate(ctx, orderID); err != nil {
        tracing.RecordError(ctx, err)
        return err
    }
    
    tracing.AddEvent(ctx, "validation_completed")
    return nil
}
```

### Trace Context Propagation

Traces propagate across services via HTTP headers:

```go
// Automatic propagation via W3C Trace Context
// Headers: traceparent, tracestate
```

## Prometheus Metrics

### Built-in Metrics

HTTP metrics (auto-collected):

```
http_requests_total{method, path, status}
http_request_duration_seconds{method, path}
http_request_size_bytes{method, path}
http_response_size_bytes{method, path}
```

Database metrics:

```
db_queries_total{operation, table}
db_query_duration_seconds{operation, table}
db_connections_open
db_connections_in_use
```

Cache metrics:

```
cache_hits_total{cache}
cache_misses_total{cache}
```

Business metrics:

```
user_registrations_total
user_logins_total
active_users
```

### HTTP Middleware

```go
// Registered in HTTP kernel
r.Use(metrics.Middleware())

// Metrics endpoint
r.GET("/metrics", metrics.Handler())
```

### Recording Metrics

```go
// HTTP request
metrics.RecordHTTPRequest("POST", "/v1/orders", "201", duration)

// Database query
metrics.RecordDBQuery("SELECT", "users", duration)

// Cache
metrics.RecordCacheHit("user_cache")
metrics.RecordCacheMiss("user_cache")

// Business
metrics.RecordUserRegistration()
metrics.RecordUserLogin()
```

### Custom Metrics

```go
// Counter
orderCounter := metrics.Counter(
    "orders_total",
    "Total number of orders",
    "status", "type",
)
orderCounter.WithLabelValues("completed", "standard").Inc()

// Gauge
activeOrders := metrics.Gauge(
    "active_orders",
    "Number of active orders",
)
activeOrders.Set(42)

// Histogram
orderDuration := metrics.Histogram(
    "order_processing_seconds",
    "Order processing duration",
    prometheus.DefBuckets,
    "type",
)
orderDuration.WithLabelValues("standard").Observe(1.5)

// Timer helper
timer := metrics.NewTimer(orderDuration, "standard")
// ... do work
timer.ObserveDuration()
```

### Prometheus Scrape Config

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'eogo'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

## Structured Logging

### Log Format

JSON structured logs for easy parsing:

```json
{
  "level": "info",
  "time": "2024-01-15T10:30:00Z",
  "msg": "request completed",
  "trace_id": "abc123",
  "method": "POST",
  "path": "/v1/users",
  "status": 201,
  "duration_ms": 45
}
```

### Trace Correlation

Include trace ID in logs for correlation:

```go
func (h *Handler) Create(c *gin.Context) {
    span := trace.SpanFromContext(c.Request.Context())
    traceID := span.SpanContext().TraceID().String()
    
    log.Info().
        Str("trace_id", traceID).
        Str("action", "create_user").
        Msg("processing request")
}
```

## Grafana Dashboards

### HTTP Overview

```
# Request rate
rate(http_requests_total[5m])

# Error rate
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])

# P99 latency
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))
```

### Database Performance

```
# Query rate by operation
rate(db_queries_total[5m])

# Slow queries (>100ms)
histogram_quantile(0.95, rate(db_query_duration_seconds_bucket[5m]))

# Connection pool usage
db_connections_in_use / db_connections_open
```

## Jaeger Integration

### Docker Compose

```yaml
services:
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"  # UI
      - "4317:4317"    # OTLP gRPC
    environment:
      - COLLECTOR_OTLP_ENABLED=true
```

### Configuration

```bash
TRACING_ENABLED=true
TRACING_ENDPOINT=localhost:4317
TRACING_INSECURE=true
```

## Best Practices

### DO ✅

- Enable tracing in production
- Use meaningful span names
- Add business-relevant attributes
- Correlate logs with trace IDs
- Set appropriate sample rates

### DON'T ❌

- Log sensitive data
- Create too many custom metrics
- Ignore high cardinality labels
- Sample at 100% in production (unless needed)

## Next Steps

- [Resilience Patterns](./02-resilience.md)
- [Health Checks](./03-health-checks.md)
