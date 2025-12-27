# EOGO Logging

Eogo logging system supporting multiple channels, level control, and file rotation.

## Quick Start

```go
import "github.com/eogo-dev/eogo/internal/infra/logger"

// Bootstrap (call once in main.go)
logger.Boot()
defer logger.Close()

// Basic usage
logger.Debug("debug message")
logger.Info("user logged in", map[string]any{"user_id": 123})
logger.Warning("disk space low")
logger.Error("failed to connect", map[string]any{"error": err.Error()})
```

## Log Levels

From lowest to highest severity:

| Level | Description | Production Output |
|-------|------|---------|
| `debug` | Detailed debug information | ❌ Silenced |
| `info` | Interesting events (Login, SQL) | ❌ Silenced |
| `notice` | Normal but significant events | ❌ Silenced |
| `warning` | Exceptional occurrences (Deprecated APIs) | ✅ Logged |
| `error` | Runtime errors | ✅ Logged |
| `critical` | Critical conditions (Component unavailable) | ✅ Logged |
| `alert` | Action must be taken immediately | ✅ Logged |
| `emergency` | System is unusable | ✅ Logged |

## Configuration

```env
# Log level
LOG_LEVEL=debug

# Log directory
LOG_PATH=storage/logs

# Filename pattern (supports date variables)
LOG_FILE={Y}-{m}-{d}.log

# Max size per file (MB)
LOG_MAX_SIZE=100

# Retention days
LOG_MAX_AGE=14

# JSON format (recommended for production)
LOG_JSON=false
```

## Channel Usage

```go
// Use pre-defined channels
logger.HTTP().Info("request received", map[string]any{
    "method": "POST",
    "path":   "/api/users",
})

logger.Database().Debug("query executed", map[string]any{
    "sql":      "SELECT * FROM users",
    "duration": "12ms",
})

logger.Auth().Warning("login failed", map[string]any{
    "email":  "[email]",
    "reason": "invalid password",
})

// Custom channel
logger.Channel("payment").Info("payment processed")
```

## Contextual Logging

```go
// Add global context to a logger instance
log := logger.WithContext(map[string]any{
    "request_id": "abc-123",
    "user_id":    456,
})

log.Info("processing order")  // Automatically includes request_id and user_id
log.Error("order failed")
```

## Usage in Gin Handlers

```go
func (h *Handler) Create(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Log using request context
    logger.Ctx(ctx).Info("creating user", map[string]any{
        "email": req.Email,
    })
    
    // Or use a specific channel
    logger.API().Info("user created", map[string]any{
        "user_id": user.ID,
    })
}
```

## Debug Mode Behavior

When `APP_DEBUG=true`:
- All log levels are output.
- Console output is colorized.
- Logs are simultaneously written to files.

When `APP_DEBUG=false` and `APP_ENV=production`:
- Only `warning` and higher are output.
- Console output is disabled.
- Logs are written to files in JSON format.

## Log File Location

```
storage/logs/
├── 2024-12-25.log      # Rotated by date
├── 2024-12-24.log
└── 2024-12-23.log
```

## AI & Debugging Friendly

Log formats are clean and easy for AI analysis:

```text
[2024-12-25 10:30:45] http.INFO: request received {"method":"POST","path":"/api/users","duration":"45ms"}
[2024-12-25 10:30:45] database.DEBUG: query executed {"sql":"INSERT INTO users...","duration":"12ms"}
[2024-12-25 10:30:46] http.ERROR: request failed {"error":"validation failed","code":400}
```

JSON Format (Production):
```json
{"time":"2024-12-25T10:30:45+08:00","level":"ERROR","channel":"http","message":"request failed","context":{"error":"validation failed"}}
```
