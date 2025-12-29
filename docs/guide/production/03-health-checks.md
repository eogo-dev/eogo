# Health Checks

> Kubernetes-ready health probes for production deployments.

## Overview

ZGO provides health check infrastructure for:

- **Liveness Probe**: Is the application alive?
- **Readiness Probe**: Is the application ready to serve traffic?
- **Detailed Health**: Full system health status

## Endpoints

| Endpoint | Purpose | K8s Probe |
|----------|---------|-----------|
| `/health` | Full health status | - |
| `/health/live` | Liveness check | livenessProbe |
| `/health/ready` | Readiness check | readinessProbe |

## Quick Start

```go
import "github.com/zgiai/zgo/internal/infra/health"

// Register routes
health.RegisterRoutes(r)

// Or use instance
h := health.New()
h.Register("database", health.DatabaseChecker(db))
h.RegisterRoutes(r)
```

## Health Checkers

### Database

```go
health.Register("database", health.DatabaseChecker(db))
```

Response:
```json
{
  "status": "up",
  "details": {
    "open_connections": 10,
    "in_use": 2,
    "idle": 8,
    "max_open": 100
  }
}
```

### Redis

```go
health.Register("redis", health.RedisChecker(func(ctx context.Context) error {
    return redisClient.Ping(ctx).Err()
}))
```

### Disk Space

```go
// Alert if less than 1GB available
health.Register("disk", health.DiskSpace("/", 1<<30))
```

Response:
```json
{
  "status": "up",
  "details": {
    "available_bytes": 50000000000,
    "total_bytes": 100000000000,
    "used_percent": 50.0
  }
}
```

### Custom Checker

```go
health.Register("external-api", func(ctx context.Context) health.CheckResult {
    resp, err := http.Get("https://api.example.com/health")
    if err != nil {
        return health.CheckResult{
            Status:  health.StatusDown,
            Message: "API unreachable",
        }
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return health.CheckResult{
            Status:  health.StatusDegraded,
            Message: "API returned non-200",
        }
    }
    
    return health.CheckResult{Status: health.StatusUp}
})
```

### With Timeout

```go
// Timeout after 3 seconds
health.Register("slow-service", 
    health.Timeout(slowChecker, 3*time.Second))
```

## Health Status

Three possible states:

| Status | Meaning | HTTP Code |
|--------|---------|-----------|
| `up` | Healthy | 200 |
| `degraded` | Partially healthy | 200 |
| `down` | Unhealthy | 503 |

## Response Format

### Full Health (`/health`)

```json
{
  "status": "up",
  "timestamp": "2024-01-15T10:30:00Z",
  "checks": {
    "database": {
      "status": "up",
      "latency_ms": 5,
      "details": {
        "open_connections": 10
      }
    },
    "redis": {
      "status": "up",
      "latency_ms": 2
    },
    "disk": {
      "status": "up",
      "details": {
        "available_bytes": 50000000000
      }
    }
  }
}
```

### Liveness (`/health/live`)

```json
{
  "status": "up"
}
```

### Readiness (`/health/ready`)

```json
{
  "status": "up"
}
```

Or when unhealthy (HTTP 503):
```json
{
  "status": "down"
}
```

## Kubernetes Configuration

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: zgo
spec:
  template:
    spec:
      containers:
        - name: zgo
          image: zgo:latest
          ports:
            - containerPort: 8080
          livenessProbe:
            httpGet:
              path: /health/live
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /health/ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
            failureThreshold: 3
          startupProbe:
            httpGet:
              path: /health/live
              port: 8080
            initialDelaySeconds: 0
            periodSeconds: 5
            failureThreshold: 30
```

### Probe Timing

| Probe | Initial Delay | Period | Failure Threshold |
|-------|--------------|--------|-------------------|
| Startup | 0s | 5s | 30 (2.5 min max) |
| Liveness | 5s | 10s | 3 |
| Readiness | 5s | 5s | 3 |

## Programmatic Usage

### Check Health

```go
h := health.Global()

// Run all checks
results := h.Check(ctx)

// Check if healthy
if h.IsHealthy(ctx) {
    // All systems go
}

// Get overall status
status := h.OverallStatus(ctx)
switch status {
case health.StatusUp:
    // Healthy
case health.StatusDegraded:
    // Partially healthy
case health.StatusDown:
    // Unhealthy
}
```

### Dynamic Registration

```go
// Register at runtime
health.Register("new-service", checker)

// Unregister
health.Unregister("old-service")
```

## Integration with Graceful Shutdown

```go
// During shutdown, mark as not ready
func shutdown() {
    // Stop accepting new requests
    health.Register("shutdown", health.Down("shutting down"))
    
    // Wait for in-flight requests
    time.Sleep(5 * time.Second)
    
    // Shutdown server
    server.Shutdown(ctx)
}
```

## Best Practices

### DO ✅

- Check all critical dependencies
- Use timeouts for external checks
- Return detailed status for debugging
- Use readiness for dependency checks
- Use liveness for deadlock detection

### DON'T ❌

- Make liveness checks too complex
- Include non-critical services in readiness
- Set timeouts too short
- Ignore health check failures in CI/CD

## Monitoring Health

### Prometheus Metrics

```go
// Add health metrics
healthGauge := prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
        Name: "health_check_status",
        Help: "Health check status (1=up, 0=down)",
    },
    []string{"check"},
)

// Update periodically
go func() {
    for range time.Tick(30 * time.Second) {
        results := health.Check(ctx)
        for name, result := range results {
            value := 0.0
            if result.Status == health.StatusUp {
                value = 1.0
            }
            healthGauge.WithLabelValues(name).Set(value)
        }
    }
}()
```

### Alerting

```yaml
# Prometheus alert rule
groups:
  - name: health
    rules:
      - alert: ServiceUnhealthy
        expr: health_check_status == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Health check {{ $labels.check }} is failing"
```

## Next Steps

- [Observability](./01-observability.md)
- [Resilience Patterns](./02-resilience.md)
