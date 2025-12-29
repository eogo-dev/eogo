# Error Handling

> Unified error handling patterns for consistent API responses.

## Response Package

Use `pkg/response` for all API responses:

```go
import "github.com/zgiai/zgo/pkg/response"
```

## Success Responses

```go
// 200 OK with data
response.OK(c, user)

// 200 OK (alias)
response.Success(c, data)

// 201 Created
response.Created(c, newResource)

// 204 No Content
response.NoContent(c)
```

## Error Responses

```go
// 400 Bad Request
response.BadRequest(c, "Invalid input", err)

// 401 Unauthorized
response.Unauthorized(c, "Invalid credentials")

// 403 Forbidden
response.Forbidden(c, "Access denied")

// 404 Not Found
response.NotFound(c, "User not found", err)

// 409 Conflict
response.Conflict(c, "Email already exists", err)

// 422 Unprocessable Entity
response.UnprocessableEntity(c, "Validation failed", err)

// 500 Internal Server Error
response.InternalServerError(c, "Something went wrong", err)
```

## Response Format

All responses follow a consistent structure:

### Success

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "john"
  }
}
```

### Error

```json
{
  "code": 400,
  "message": "Invalid input",
  "error": "email format is invalid"
}
```

## Handler Pattern

```go
func (h *Handler) Create(c *gin.Context) {
    // 1. Bind and validate input
    var req CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request body", err)
        return
    }
    
    // 2. Call service
    result, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        // 3. Handle specific errors
        switch {
        case errors.Is(err, ErrDuplicateEmail):
            response.Conflict(c, "Email already exists", err)
        case errors.Is(err, ErrNotFound):
            response.NotFound(c, "Resource not found", err)
        default:
            response.InternalServerError(c, "Failed to create", err)
        }
        return
    }
    
    // 4. Return success
    response.Created(c, result)
}
```

## Custom Errors

Define domain-specific errors:

```go
// internal/modules/user/errors.go
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrDuplicateEmail   = errors.New("email already exists")
    ErrInvalidPassword  = errors.New("invalid password")
    ErrAccountLocked    = errors.New("account is locked")
)
```

## Error Wrapping

Wrap errors with context:

```go
func (s *service) GetUser(ctx context.Context, id uint) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return user, nil
}
```

## Validation Errors

Handle validation errors with details:

```go
func (h *Handler) Create(c *gin.Context) {
    var req CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Extract validation errors
        var ve validator.ValidationErrors
        if errors.As(err, &ve) {
            errors := make(map[string]string)
            for _, fe := range ve {
                errors[fe.Field()] = fe.Tag()
            }
            response.UnprocessableEntity(c, "Validation failed", errors)
            return
        }
        response.BadRequest(c, "Invalid JSON", err)
        return
    }
}
```

Response:

```json
{
  "code": 422,
  "message": "Validation failed",
  "error": {
    "Email": "required",
    "Password": "min"
  }
}
```

## Panic Recovery

The recovery middleware catches panics:

```go
// Registered in HTTP kernel
r.Use(gin.Recovery())
```

For custom panic handling:

```go
func CustomRecovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                // Log the panic
                log.Error().
                    Interface("panic", r).
                    Str("stack", string(debug.Stack())).
                    Msg("panic recovered")
                
                response.InternalServerError(c, "Internal server error", nil)
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

## Error Logging

Log errors with context:

```go
func (h *Handler) Create(c *gin.Context) {
    result, err := h.service.Create(ctx, &req)
    if err != nil {
        // Log with trace ID
        span := trace.SpanFromContext(c.Request.Context())
        log.Error().
            Err(err).
            Str("trace_id", span.SpanContext().TraceID().String()).
            Str("action", "create_user").
            Msg("failed to create user")
        
        response.InternalServerError(c, "Failed to create user", err)
        return
    }
}
```

## Error Middleware

Centralized error handling:

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        // Check for errors after handler execution
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            switch {
            case errors.Is(err, gorm.ErrRecordNotFound):
                response.NotFound(c, "Resource not found", err)
            case errors.Is(err, context.DeadlineExceeded):
                response.InternalServerError(c, "Request timeout", err)
            default:
                response.InternalServerError(c, "Internal error", err)
            }
        }
    }
}
```

## Best Practices

### DO ✅

- Use the response package consistently
- Define domain-specific errors
- Wrap errors with context
- Log errors with trace IDs
- Return appropriate HTTP status codes

### DON'T ❌

- Expose internal error details to clients
- Use generic error messages
- Ignore errors
- Log sensitive data in errors
- Return 500 for client errors

## Error Codes

Consider using error codes for client handling:

```go
const (
    ErrCodeValidation    = 1001
    ErrCodeUnauthorized  = 1002
    ErrCodeNotFound      = 1003
    ErrCodeDuplicate     = 1004
    ErrCodeInternal      = 5000
)

type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}
```

## Next Steps

- [API Design](./03-api-design.md)
- [Testing Guide](./01-testing.md)
