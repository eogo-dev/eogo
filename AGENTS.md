# AGENTS.md

Instructions for AI coding agents working on the Eogo framework.

## Project Overview

EOGO is a modern Go framework using Domain-Driven Design (DDD) + layered architecture.

## Directory Structure

```text
eogo/
├── cmd/
│   ├── eogo/              # CLI tool
│   └── server/            # HTTP server entry
├── internal/
│   ├── bootstrap/         # Application startup
│   ├── domain/            # Domain entities (core business)
│   ├── modules/           # Business modules
│   │   └── user/          # Example: 8 files
│   │       ├── model.go       # Database entity (UserPO)
│   │       ├── dto.go         # DTO + Mapper functions
│   │       ├── repository.go  # Data access layer
│   │       ├── service.go     # Business logic layer
│   │       ├── handler.go     # HTTP handlers
│   │       ├── routes.go      # Route registration
│   │       ├── provider.go    # Wire DI
│   │       └── service_test.go
│   ├── infra/             # Infrastructure (33+ components)
│   └── wiring/            # Wire dependency injection
├── pkg/                   # Public libraries
├── routes/                # Global routes
└── tests/                 # Tests
```

## Common Commands

```bash
make build         # Build CLI
make test          # Run tests
make lint          # Code linting
make wire          # Generate DI
make air           # Hot-reload dev server
```

## Module Structure (8-file standard)

| File | Responsibility |
|------|----------------|
| `model.go` | Database entity `UserPO` (GORM) |
| `dto.go` | Request/Response DTO + `toDomain()`/`toUserPO()` mappers |
| `repository.go` | Data access, returns `domain.User` |
| `service.go` | Business logic, uses `domain.User` |
| `handler.go` | HTTP handlers |
| `routes.go` | Route registration |
| `provider.go` | Wire ProviderSet |

## Domain Layer

`internal/domain/` contains core business entities with JSON tags:

```go
// internal/domain/user.go
type User struct {
    ID        uint       `json:"id"`
    Username  string     `json:"username"`
    Email     string     `json:"email"`
    Password  string     `json:"-"`  // Always hidden!
    CreatedAt time.Time  `json:"created_at"`
}
```

**Data Flow**: `Handler(DTO) → Service(domain.User) → Repository(UserPO)`

## Handler Utilities

```go
import "github.com/eogo-dev/eogo/pkg/handler"

// Parse URL parameters (auto sends error response)
id, ok := handler.ParseID(c, "id")
if !ok {
    return  // 400 already sent
}

// Get authenticated user (auto sends 401)
userID, ok := handler.GetUserID(c)
if !ok {
    return  // 401 already sent
}

// Bind JSON request (auto sends error response)
var req CreateRequest
if !handler.BindJSON(c, &req) {
    return  // 400 already sent
}

// Query helpers
page := handler.QueryInt(c, "page", 1)
active := handler.QueryBool(c, "active", true)
```

## Unified Response

```go
import "github.com/eogo-dev/eogo/pkg/response"

// All responses use Success() - it auto-detects pagination!
response.Success(c, data)           // 200 with data
response.Success(c, paginator)      // 200 with data + meta + links (auto-detected)
response.Created(c, data)           // 201 created
response.NoContent(c)               // 204 no content

// Error responses
response.BadRequest(c, "message", err)
response.NotFound(c, "message", err)
response.Unauthorized(c)
response.Forbidden(c)
response.HandleError(c, "message", err)  // Auto-maps error to status code

// With inline transformation (optional, rare)
response.Transform(c, user, func(u *User) any {
    return map[string]any{"id": u.ID, "name": u.Username}
})
```

## Pagination

```go
import "github.com/eogo-dev/eogo/pkg/pagination"

// Simple pagination - just pass paginator to Success()
users, paginator, err := pagination.New[User](c, db.Model(&User{}))
response.Success(c, paginator)  // Auto-outputs data + meta + links

// With custom scope
users, paginator, err := pagination.NewWithScope[User](c, db, func(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "active").Order("created_at DESC")
})
response.Success(c, paginator)

// Manual pagination
req := pagination.FromContext(c)
users, total, _ := service.List(ctx, req.GetPage(), req.GetPerPage())
paginator := pagination.NewPaginator(users, total, req.GetPage(), req.GetPerPage())
response.Success(c, paginator)
```

## Complete Handler Example

```go
package user

import (
    "github.com/eogo-dev/eogo/pkg/handler"
    "github.com/eogo-dev/eogo/pkg/pagination"
    "github.com/eogo-dev/eogo/pkg/response"
    "github.com/gin-gonic/gin"
)

// Get - single resource
func (h *Handler) Get(c *gin.Context) {
    id, ok := handler.ParseID(c, "id")
    if !ok {
        return
    }

    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        response.HandleError(c, "User not found", err)
        return
    }

    response.Success(c, user)  // Domain直接输出
}

// List - paginated
func (h *Handler) List(c *gin.Context) {
    req := pagination.FromContext(c)
    users, total, err := h.service.List(c.Request.Context(), req.GetPage(), req.GetPerPage())
    if err != nil {
        response.HandleError(c, "Failed to list", err)
        return
    }

    paginator := pagination.NewPaginator(users, total, req.GetPage(), req.GetPerPage())
    paginator.SetPath(c.Request.URL.Path)
    response.Success(c, paginator)  // 自动输出 data + meta + links
}

// Create - with binding
func (h *Handler) Create(c *gin.Context) {
    var req CreateRequest
    if !handler.BindJSON(c, &req) {
        return
    }

    user, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        response.HandleError(c, "Create failed", err)
        return
    }

    response.Created(c, user)
}

// Update - authenticated user
func (h *Handler) UpdateProfile(c *gin.Context) {
    userID, ok := handler.GetUserID(c)
    if !ok {
        return
    }

    var req UpdateRequest
    if !handler.BindJSON(c, &req) {
        return
    }

    user, err := h.service.Update(c.Request.Context(), userID, &req)
    if err != nil {
        response.HandleError(c, "Update failed", err)
        return
    }

    response.Success(c, user)
}
```

## Wire Dependency Injection

```go
// internal/modules/user/provider.go
var ProviderSet = wire.NewSet(
    NewRepository,
    wire.Bind(new(Repository), new(*repository)),
    NewService,
    wire.Bind(new(Service), new(*service)),
    NewHandler,
)
```

Run `cd internal/wiring && wire` to generate code.

## Creating New Modules

```bash
./eogo make:module Blog

# Then:
# 1. Register routes in routes/api.go
# 2. Run wire
```

## Development Guidelines

1. **DTO includes Mapper** - Mapper functions go in `dto.go`
2. **Use Domain Layer** - Business logic uses `domain.User`
3. **Private implementations** - Struct names are unexported
4. **Constructors return interfaces** - `NewService() Service`
5. **snake_case JSON** - `json:"user_id"`
6. **English comments** - All code and comments in English
7. **Use handler package** - For ParseID, GetUserID, BindJSON
8. **Domain has JSON tags** - Sensitive fields use `json:"-"`

## Testing

```bash
# Unit tests
go test ./internal/modules/user/...

# Integration tests
go test ./tests/integration/...

# All tests
make test
```
