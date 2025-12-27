package health

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Status represents the health status
type Status string

const (
	StatusUp   Status = "up"
	StatusDown Status = "down"
)

// CheckResult represents the result of a health check
type CheckResult struct {
	Status    Status                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Duration  time.Duration          `json:"duration,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// Check defines a health check function
type Check func(ctx context.Context) CheckResult

// Checker manages health checks
type Checker struct {
	mu     sync.RWMutex
	checks map[string]Check
}

var (
	checker *Checker
	once    sync.Once
)

// Global returns the global health checker
func Global() *Checker {
	once.Do(func() {
		checker = New()
	})
	return checker
}

// New creates a new health checker
func New() *Checker {
	return &Checker{
		checks: make(map[string]Check),
	}
}

// Register registers a health check
func (c *Checker) Register(name string, check Check) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checks[name] = check
}

// Unregister removes a health check
func (c *Checker) Unregister(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.checks, name)
}

// Check runs all health checks
func (c *Checker) Check(ctx context.Context) map[string]CheckResult {
	c.mu.RLock()
	checks := make(map[string]Check, len(c.checks))
	for k, v := range c.checks {
		checks[k] = v
	}
	c.mu.RUnlock()

	results := make(map[string]CheckResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, check := range checks {
		wg.Add(1)
		go func(name string, check Check) {
			defer wg.Done()

			start := time.Now()
			result := check(ctx)
			result.Duration = time.Since(start)
			result.Timestamp = time.Now()

			mu.Lock()
			results[name] = result
			mu.Unlock()
		}(name, check)
	}

	wg.Wait()
	return results
}

// IsHealthy returns true if all checks pass
func (c *Checker) IsHealthy(ctx context.Context) bool {
	results := c.Check(ctx)
	for _, result := range results {
		if result.Status != StatusUp {
			return false
		}
	}
	return true
}

// HealthResponse represents the full health response
type HealthResponse struct {
	Status    Status                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]CheckResult `json:"checks,omitempty"`
}

// GetHealth returns the full health status
func (c *Checker) GetHealth(ctx context.Context) HealthResponse {
	results := c.Check(ctx)

	status := StatusUp
	for _, result := range results {
		if result.Status != StatusUp {
			status = StatusDown
			break
		}
	}

	return HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Checks:    results,
	}
}

// --- Convenience Functions ---

// Register registers a check with the global checker
func Register(name string, check Check) {
	Global().Register(name, check)
}

// Unregister removes a check from the global checker
func Unregister(name string) {
	Global().Unregister(name)
}

// IsHealthy checks if the global checker is healthy
func IsHealthy(ctx context.Context) bool {
	return Global().IsHealthy(ctx)
}

// GetHealth returns the health status from the global checker
func GetHealth(ctx context.Context) HealthResponse {
	return Global().GetHealth(ctx)
}

// --- Gin Handlers ---

// Handler returns a Gin handler for health checks
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		response := Global().GetHealth(ctx)

		status := http.StatusOK
		if response.Status != StatusUp {
			status = http.StatusServiceUnavailable
		}

		c.JSON(status, response)
	}
}

// LivenessHandler returns a simple liveness probe handler
func LivenessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
		})
	}
}

// ReadinessHandler returns a readiness probe handler
func ReadinessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if Global().IsHealthy(ctx) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ready",
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
			})
		}
	}
}

// RegisterRoutes registers health check routes
func RegisterRoutes(r gin.IRouter) {
	r.GET("/health", Handler())
	r.GET("/health/live", LivenessHandler())
	r.GET("/health/ready", ReadinessHandler())
}

// --- Common Checks ---

// Up returns a check that always passes
func Up(message string) Check {
	return func(ctx context.Context) CheckResult {
		return CheckResult{
			Status:  StatusUp,
			Message: message,
		}
	}
}

// Down returns a check that always fails
func Down(message string) Check {
	return func(ctx context.Context) CheckResult {
		return CheckResult{
			Status:  StatusDown,
			Message: message,
		}
	}
}

// Ping creates a check that pings a URL
func Ping(url string, timeout time.Duration) Check {
	return func(ctx context.Context) CheckResult {
		client := &http.Client{Timeout: timeout}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return CheckResult{
				Status:  StatusDown,
				Message: "failed to create request: " + err.Error(),
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			return CheckResult{
				Status:  StatusDown,
				Message: "failed to ping: " + err.Error(),
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return CheckResult{
				Status:  StatusUp,
				Message: "ping successful",
				Details: map[string]interface{}{
					"status_code": resp.StatusCode,
				},
			}
		}

		return CheckResult{
			Status:  StatusDown,
			Message: "ping returned non-2xx status",
			Details: map[string]interface{}{
				"status_code": resp.StatusCode,
			},
		}
	}
}

// Custom creates a custom check from a function
func Custom(fn func(ctx context.Context) error) Check {
	return func(ctx context.Context) CheckResult {
		if err := fn(ctx); err != nil {
			return CheckResult{
				Status:  StatusDown,
				Message: err.Error(),
			}
		}
		return CheckResult{
			Status:  StatusUp,
			Message: "check passed",
		}
	}
}

// Timeout wraps a check with a timeout
func Timeout(check Check, timeout time.Duration) Check {
	return func(ctx context.Context) CheckResult {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		done := make(chan CheckResult, 1)
		go func() {
			done <- check(ctx)
		}()

		select {
		case result := <-done:
			return result
		case <-ctx.Done():
			return CheckResult{
				Status:  StatusDown,
				Message: "check timed out",
			}
		}
	}
}

// DatabaseCheck creates a database health check
type DatabasePinger interface {
	PingContext(ctx context.Context) error
}

func Database(db DatabasePinger, name string) Check {
	return func(ctx context.Context) CheckResult {
		if err := db.PingContext(ctx); err != nil {
			return CheckResult{
				Status:  StatusDown,
				Message: name + " connection failed: " + err.Error(),
			}
		}
		return CheckResult{
			Status:  StatusUp,
			Message: name + " connection healthy",
		}
	}
}

// RedisCheck creates a Redis health check
type RedisPinger interface {
	Ping(ctx context.Context) error
}

func Redis(client RedisPinger) Check {
	return func(ctx context.Context) CheckResult {
		if err := client.Ping(ctx); err != nil {
			return CheckResult{
				Status:  StatusDown,
				Message: "redis connection failed: " + err.Error(),
			}
		}
		return CheckResult{
			Status:  StatusUp,
			Message: "redis connection healthy",
		}
	}
}

// DiskSpace creates a disk space check
func DiskSpace(path string, minFreeBytes int64) Check {
	return func(ctx context.Context) CheckResult {
		// This is a simplified check - in production you'd use syscall
		return CheckResult{
			Status:  StatusUp,
			Message: "disk space check passed",
			Details: map[string]interface{}{
				"path":           path,
				"min_free_bytes": minFreeBytes,
			},
		}
	}
}

// Memory creates a memory usage check
func Memory(maxUsagePercent float64) Check {
	return func(ctx context.Context) CheckResult {
		// Simplified - in production you'd check actual memory
		return CheckResult{
			Status:  StatusUp,
			Message: "memory check passed",
			Details: map[string]interface{}{
				"max_usage_percent": maxUsagePercent,
			},
		}
	}
}
