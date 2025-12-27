# Resilience Patterns

> Circuit breaker, rate limiting, retry, and singleflight for production stability.

## Overview

Eogo provides production-grade resilience patterns:

| Pattern | Package | Purpose |
|---------|---------|---------|
| Circuit Breaker | `internal/infra/breaker` | Prevent cascade failures |
| Singleflight | `internal/infra/singleflight` | Prevent cache stampede |
| Retry | `internal/infra/retry` | Handle transient failures |
| Rate Limiting | `internal/infra/ratelimit` | Protect resources |

## Circuit Breaker

Prevents cascade failures by stopping requests to failing services.

### States

```
┌─────────┐    failures >= threshold    ┌─────────┐
│ CLOSED  │ ─────────────────────────▶ │  OPEN   │
│(normal) │                            │(failing)│
└─────────┘                            └─────────┘
     ▲                                      │
     │                                      │ timeout
     │         ┌───────────┐                │
     └─────────│ HALF-OPEN │◀───────────────┘
   success     │ (testing) │
               └───────────┘
```

### Usage

```go
import "github.com/eogo-dev/eogo/internal/infra/breaker"

// Create breaker
cb := breaker.New(breaker.Config{
    Name:                "payment-service",
    Threshold:           5,              // Open after 5 failures
    Timeout:             10 * time.Second, // Stay open for 10s
    MaxHalfOpenRequests: 1,              // Allow 1 test request
})

// Execute with protection
err := cb.Do(func() error {
    return paymentService.Charge(amount)
})

if errors.Is(err, breaker.ErrServiceUnavailable) {
    // Circuit is open, use fallback
}
```

### With Fallback

```go
err := cb.DoWithFallback(
    func() error {
        return primaryService.Call()
    },
    func(err error) error {
        // Fallback when circuit is open
        return backupService.Call()
    },
)
```

### Breaker Group

Manage multiple breakers by name:

```go
group := breaker.NewGroup(breaker.Config{
    Threshold: 5,
    Timeout:   10 * time.Second,
})

// Auto-creates breaker for each service
err := group.Do("user-service", func() error {
    return userService.GetUser(id)
})

err = group.Do("order-service", func() error {
    return orderService.GetOrder(id)
})
```

### Check State

```go
state := cb.State()
switch state {
case breaker.StateClosed:
    // Normal operation
case breaker.StateOpen:
    // Failing, requests blocked
case breaker.StateHalfOpen:
    // Testing recovery
}
```

## Singleflight

Prevents cache stampede by deduplicating concurrent requests.

### Problem

Without singleflight:
```
100 concurrent requests for user:123
  → 100 database queries
  → Database overload
```

With singleflight:
```
100 concurrent requests for user:123
  → 1 database query
  → 99 requests wait and share result
```

### Usage

```go
import "github.com/eogo-dev/eogo/internal/infra/singleflight"

sf := singleflight.New()

// All concurrent calls with same key share one execution
result, err := sf.Do("user:123", func() (any, error) {
    return db.GetUser(123)
})

user := result.(*User)
```

### Type-Safe Version

```go
sf := singleflight.NewTyped[*User]()

user, err := sf.Do("user:123", func() (*User, error) {
    return db.GetUser(123)
})
// user is *User, no type assertion needed
```

### With Context

```go
// Respects context cancellation
user, err := sf.DoCtx(ctx, "user:123", func() (any, error) {
    return db.GetUser(123)
})
```

### Check if Fresh

```go
result, fresh, err := sf.DoEx("user:123", func() (any, error) {
    return db.GetUser(123)
})

if fresh {
    // This call executed the function
} else {
    // This call shared another's result
}
```

## Retry with Backoff

Handles transient failures with exponential backoff.

### Basic Usage

```go
import "github.com/eogo-dev/eogo/internal/infra/retry"

err := retry.Do(ctx, func(ctx context.Context) error {
    return externalService.Call()
},
    retry.WithMaxAttempts(3),
    retry.WithInitialDelay(100*time.Millisecond),
)
```

### Exponential Backoff

```go
// Delays: 100ms → 200ms → 400ms → 800ms (capped at MaxDelay)
err := retry.Do(ctx, fn,
    retry.WithMaxAttempts(5),
    retry.WithInitialDelay(100*time.Millisecond),
    retry.WithMaxDelay(10*time.Second),
    retry.WithMultiplier(2.0),
    retry.WithJitter(0.1),  // ±10% randomness
)
```

### Convenience Functions

```go
// Simple exponential backoff
retry.ExponentialBackoff(ctx, fn, 3)

// Fixed retries with no delay
retry.Times(ctx, 3, fn)

// Retry forever until success
retry.Forever(ctx, fn, time.Second)
```

### With Result

```go
user, err := retry.DoWithResult(ctx, func(ctx context.Context) (*User, error) {
    return userService.GetUser(id)
},
    retry.WithMaxAttempts(3),
)
```

### Conditional Retry

```go
err := retry.Do(ctx, fn,
    retry.WithShouldRetry(func(err error) bool {
        // Only retry on specific errors
        var netErr net.Error
        if errors.As(err, &netErr) && netErr.Temporary() {
            return true
        }
        return false
    }),
)
```

### Reusable Retrier

```go
retrier := retry.NewRetrier(
    retry.WithMaxAttempts(3),
    retry.WithInitialDelay(100*time.Millisecond),
)

// Use same config for multiple calls
err := retrier.Do(ctx, fn1)
err = retrier.Do(ctx, fn2)
```

## Rate Limiting

Protects resources from overload.

### Middleware

```go
import "github.com/eogo-dev/eogo/internal/infra/ratelimit"

// Global rate limit
r.Use(ratelimit.Middleware(ratelimit.Config{
    Max:      100,
    Duration: time.Minute,
}))

// Per-route rate limit
r.POST("/login", 
    ratelimit.PerMinute(10),  // 10 requests/minute
    handler.Login,
)
```

### By User

```go
r.Use(ratelimit.ByUser(100, time.Minute, func(c *gin.Context) string {
    return c.GetString("user_id")
}))
```

### By IP

```go
r.Use(ratelimit.ByIP(60, time.Minute))
```

## Combining Patterns

### Cache with Singleflight + Circuit Breaker

```go
type CachedUserService struct {
    cache   Cache
    db      *gorm.DB
    sf      singleflight.SingleFlight
    breaker breaker.Breaker
}

func (s *CachedUserService) GetUser(ctx context.Context, id uint) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    
    // Check cache first
    if user, ok := s.cache.Get(key); ok {
        return user.(*User), nil
    }
    
    // Singleflight prevents cache stampede
    result, err := s.sf.Do(key, func() (any, error) {
        // Circuit breaker protects database
        var user *User
        err := s.breaker.Do(func() error {
            return s.db.First(&user, id).Error
        })
        if err != nil {
            return nil, err
        }
        
        // Cache the result
        s.cache.Set(key, user, 5*time.Minute)
        return user, nil
    })
    
    if err != nil {
        return nil, err
    }
    return result.(*User), nil
}
```

### External API with Retry + Circuit Breaker

```go
func (s *PaymentService) Charge(ctx context.Context, amount int) error {
    return s.breaker.Do(func() error {
        return retry.Do(ctx, func(ctx context.Context) error {
            return s.gateway.Charge(amount)
        },
            retry.WithMaxAttempts(3),
            retry.WithInitialDelay(100*time.Millisecond),
        )
    })
}
```

## Best Practices

### DO ✅

- Use circuit breakers for external services
- Use singleflight for cache population
- Use retry for transient failures
- Set appropriate thresholds and timeouts
- Monitor breaker state changes

### DON'T ❌

- Retry non-idempotent operations blindly
- Set retry delays too short
- Ignore circuit breaker state in monitoring
- Use singleflight for write operations

## Next Steps

- [Health Checks](./03-health-checks.md)
- [Observability](./01-observability.md)
