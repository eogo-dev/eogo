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

// ========== 方式1: Handler 直接查询 (最简单，AI推荐) ==========
func (h *Handler) List(c *gin.Context) {
    // 一行搞定！自动处理 page, per_page, path, query
    paginator, err := pagination.Auto[*domain.User](c, h.db.Model(&UserPO{}))
    if err != nil {
        response.HandleError(c, "Failed to list", err)
        return
    }
    response.Success(c, paginator)  // 自动输出 data + meta + links
}

// 带条件查询
paginator, err := pagination.AutoWithScope[*domain.User](c, h.db, func(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "active").Order("created_at DESC")
})

// ========== 方式2: 通过 Service 层 (DDD标准) ==========
// Service:
func (s *service) List(ctx context.Context, page, perPage int) (*pagination.Result[*domain.User], error) {
    users, total, err := s.repo.FindAll(ctx, page, perPage)
    if err != nil {
        return nil, err
    }
    return pagination.NewResult(users, total, page, perPage), nil
}

// Handler:
func (h *Handler) List(c *gin.Context) {
    req := pagination.FromContext(c)
    result, err := h.service.List(c.Request.Context(), req.GetPage(), req.GetPerPage())
    if err != nil {
        response.HandleError(c, "Failed to list", err)
        return
    }
    response.Success(c, result.ToPaginator(c))  // 自动设置 path 和 query
}
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

// List - paginated (最简方式)
func (h *Handler) List(c *gin.Context) {
    paginator, err := pagination.Auto[*domain.User](c, h.db.Model(&UserPO{}))
    if err != nil {
        response.HandleError(c, "Failed to list", err)
        return
    }
    response.Success(c, paginator)
}

// List - 通过 Service (DDD标准)
func (h *Handler) ListViaService(c *gin.Context) {
    req := pagination.FromContext(c)
    result, err := h.service.List(c.Request.Context(), req.GetPage(), req.GetPerPage())
    if err != nil {
        response.HandleError(c, "Failed to list", err)
        return
    }
    response.Success(c, result.ToPaginator(c))
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
