# ZGO Platform Guide

*Comprehensive guide for building applications with the ZGO framework*

## Table of Contents

- [Quick Start](#quick-start)
- [Core Packages](#core-packages)
  - [Pagination](#pagination)
  - [Response Helpers](#response-helpers)
  - [Resource Transformers](#resource-transformers)
  - [Validation](#validation)
  - [Error Handling](#error-handling)
- [Module Development](#module-development)
  - [Module Structure](#module-structure)
  - [Creating a Module](#creating-a-module)
  - [Dependency Injection](#dependency-injection)
- [Best Practices](#best-practices)

---

## Quick Start

### Basic CRUD Handler Example

```go
package blog

import (
    "github.com/zgiai/zgo/pkg/pagination"
    "github.com/zgiai/zgo/pkg/response"
    "github.com/zgiai/zgo/pkg/resource"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type Handler struct {
    db *gorm.DB
}

// List posts with pagination
func (h *Handler) List(c *gin.Context) {
    // One-liner pagination
    posts, paging, err := pagination.PaginateFromContext[Post](c, h.db)
    if err != nil {
        response.InternalServerError(c, "Failed to fetch posts", err)
        return
    }

    // Transform and respond
    resources := resource.NewPaginatedCollection(
        transformPosts(posts), 
        paging.Page, 
        paging.PageSize, 
        paging.Total,
    )
    c.JSON(200, resources.ToResponse())
}

// Create a post
func (h *Handler) Create(c *gin.Context) {
    var req CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request", err)
        return
    }

    post := &Post{Title: req.Title, Content: req.Content}
    if err := h.db.Create(post).Error; err != nil {
        response.InternalServerError(c, "Failed to create post", err)
        return
    }

    response.Created(c, transformPost(post))
}
```

---

## Core Packages

### Pagination

**Package:** `github.com/zgiai/zgo/pkg/pagination`

#### Quick Reference

```go
// Extract pagination from context (recommended)
req := pagination.FromContext(c)

// One-liner paginate (recommended)
items, paging, err := pagination.PaginateFromContext[User](c, db)

// Manual pagination
items, paging, err := pagination.Paginate[User](db, req)

// Build result manually
result := pagination.BuildResult(total, page, pageSize)
```

#### API Reference

##### `FromContext(c *gin.Context) *Request`

Extracts pagination parameters from query string.

**Query Parameters:**
- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 10, max: 100)
- `keyword` - Search keyword (optional)

**Example:**
```go
// GET /api/users?page=2&page_size=20&keyword=john
req := pagination.FromContext(c)
// req.Page = 2, req.PageSize = 20, req.Keyword = "john"
```

##### `PaginateFromContext[T](c *gin.Context, db *gorm.DB) ([]T, *Result, error)`

**Recommended** - One-liner pagination with automatic query extraction.

**Example:**
```go
func (h *Handler) ListUsers(c *gin.Context) {
    users, paging, err := pagination.PaginateFromContext[User](c, h.db)
    if err != nil {
        response.InternalServerError(c, "Failed to fetch users", err)
        return
    }
    
    c.JSON(200, gin.H{"data": users, "paging": paging})
}
```

##### `Paginate[T](db *gorm.DB, req *Request) ([]T, *Result, error)`

Generic pagination with manual request.

**Example:**
```go
req := &pagination.Request{Page: 1, PageSize: 20}
users, paging, err := pagination.Paginate[User](db.Where("status = ?", "active"), req)
```

#### Response Format

```json
{
  "total": 156,
  "page": 2,
  "page_size": 20,
  "last_page": 8,
  "from": 21,
  "to": 40
}
```

---

### Response Helpers

**Package:** `github.com/zgiai/zgo/pkg/response`

#### Standard Responses

```go
// Success (200 OK)
response.OK(c, user)
response.Success(c, data)

// Created (201)
response.Created(c, newUser)

// No Content (204)
response.NoContent(c)

// Bad Request (400)
response.BadRequest(c, "Invalid input", err)

// Unauthorized (401)
response.Unauthorized(c, "Token expired", err)

// Forbidden (403)
response.Forbidden(c, "Access denied", err)

// Not Found (404)
response.NotFound(c, "User not found", err)

// Internal Server Error (500)
response.InternalServerError(c, "Database error", err)
```

#### Response Structure

All responses follow this format:

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

Error responses:

```json
{
  "code": 400,
  "message": "Invalid input"
}
```

---

### Resource Transformers

**Package:** `github.com/zgiai/zgo/pkg/resource`

Resource transformers convert database models to API responses, hiding sensitive fields and formatting output.

#### Basic Resource

```go
// Define a resource
type UserResource struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // Password is NOT included (hidden)
}

func (r UserResource) ToMap() map[string]interface{} {
    return map[string]interface{}{
        "id":    r.ID,
        "name":  r.Name,
        "email": r.Email,
    }
}

// Transform model to resource
func transformUser(user *User) UserResource {
    return UserResource{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }
}

// Usage in handler
func (h *Handler) GetUser(c *gin.Context) {
    user, _ := h.service.GetUser(id)
    resource.Respond(c, 200, transformUser(user))
}
```

#### Collection

```go
// Multiple items
func (h *Handler) ListUsers(c *gin.Context) {
    users, _ := h.service.ListUsers()
    
    // Transform all users
    resources := make([]UserResource, len(users))
    for i, user := range users {
        resources[i] = transformUser(&user)
    }
    
    // Create collection
    collection := resource.NewCollection(resources)
    resource.RespondCollection(c, 200, collection)
}

// Response:
// {
//   "data": [
//     {"id": 1, "name": "John", "email": "john@example.com"},
//     {"id": 2, "name": "Jane", "email": "jane@example.com"}
//   ]
// }
```

#### Paginated Collection

```go
func (h *Handler) ListUsers(c *gin.Context) {
    users, paging, _ := pagination.PaginateFromContext[User](c, h.db)
    
    // Transform
    resources := make([]UserResource, len(users))
    for i, user := range users {
        resources[i] = transformUser(&user)
    }
    
    // Paginated collection
    collection := resource.NewPaginatedCollection(
        resources, 
        paging.Page, 
        paging.PageSize, 
        paging.Total,
    )
    
    c.JSON(200, collection.ToResponse())
}

// Response:
// {
//   "data": [...],
//   "meta": {
//     "current_page": 1,
//     "per_page": 10,
//     "total": 156,
//     "total_pages": 16
//   }
// }
```

#### Advanced: Meta and Links

```go
collection := resource.NewCollection(resources)
    .WithMeta(map[string]interface{}{
        "generated_at": time.Now(),
        "version": "v1",
    })
    .WithLinks(map[string]string{
        "self": "/api/users",
        "docs": "/docs/users",
    })

resource.RespondCollection(c, 200, collection)
```

---

### Validation

**Package:** `github.com/zgiai/zgo/pkg/validation`

Use struct tags and Gin's binding for validation:

```go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Age      int    `json:"age" binding:"gte=18,lte=120"`
}

func (h *Handler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Validation failed", err)
        return
    }
    // ...
}
```

**Common Validators:**
- `required` - Field is required
- `email` - Valid email format
- `min=n`, `max=n` - String length or number range
- `gte=n`, `lte=n` - Greater/less than or equal
- `oneof=a b c` - Value must be one of
- `url` - Valid URL

---

### Error Handling

**Package:** `github.com/zgiai/zgo/pkg/errors`

```go
import "github.com/zgiai/zgo/pkg/errors"

// Predefined errors
err := errors.ErrUserNotFound
err := errors.ErrInvalidCredentials
err := errors.ErrUnauthorized

// Custom errors
err := errors.New("CUSTOM_ERROR", "Something went wrong", 400)

// In handler
if err != nil {
    errors.Abort(c, err)
    return
}
```

---

## Module Development

### Module Structure

#### Simple Module

For single-domain functionality:

```
blog/
├── model.go        # Database entity
├── dto.go          # Request/Response DTOs
├── resource.go     # Resource transformers
├── repository.go   # Data access layer
├── service.go      # Business logic
├── handler.go      # HTTP handlers
├── routes.go       # Route registration
└── provider.go     # Dependency injection
```

#### Sub-module Structure

For complex domains with multiple sub-domains:

```
ecommerce/
├── product/
│   ├── dto.go
│   ├── repository.go
│   ├── service.go
│   └── handler.go
├── order/
│   ├── dto.go
│   ├── repository.go
│   ├── service.go
│   └── handler.go
├── model.go        # Shared models (Product, Order, etc.)
├── routes.go       # Centralized routing
└── provider.go     # Dependency injection
```

### Creating a Module

#### Using CLI (Recommended)

```bash
./zgo make:module Blog
```

This generates:
- `internal/modules/blog/` directory
- All necessary files (model, dto, repository, service, handler, routes)
- Wire provider registration
- Route registration stub

#### Manual Creation

1. **Create directory:** `internal/modules/mymodule/`

2. **Define model** (`model.go`):
```go
package mymodule

type Post struct {
    ID        uint      `gorm:"primaryKey"`
    Title     string    `gorm:"size:200;not null"`
    Content   string    `gorm:"type:text"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

3. **Create DTOs** (`dto.go`):
```go
type CreatePostRequest struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content" binding:"required"`
}

type PostResponse struct {
    ID      uint   `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}
```

4. **Implement Repository** (`repository.go`):
```go
type Repository interface {
    Create(ctx context.Context, post *Post) error
    FindByID(ctx context.Context, id uint) (*Post, error)
    List(ctx context.Context) ([]Post, error)
}

type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, post *Post) error {
    return r.db.WithContext(ctx).Create(post).Error
}
```

5. **Implement Service** (`service.go`):
```go
type Service interface {
    CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error)
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error) {
    post := &Post{
        Title:   req.Title,
        Content: req.Content,
    }
    return post, s.repo.Create(ctx, post)
}
```

6. **Implement Handler** (`handler.go`):
```go
type Handler struct {
    service Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
    var req CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request", err)
        return
    }

    post, err := h.service.CreatePost(c.Request.Context(), &req)
    if err != nil {
        response.InternalServerError(c, "Failed to create post", err)
        return
    }

    response.Created(c, post)
}
```

7. **Register Routes** (`routes.go`):
```go
package mymodule

import (
    "github.com/zgiai/zgo/internal/infra/router"
)

func Register(r *router.Router) {
    // Get dependencies from DI container
    handler := GetHandler() // Wire-generated

    posts := r.Group("/posts")
    {
        posts.POST("", handler.Create).Name("posts.create")
        posts.GET("", handler.List).Name("posts.list")
        posts.GET("/:id", handler.Get).Name("posts.get")
        posts.PUT("/:id", handler.Update).Name("posts.update")
        posts.DELETE("/:id", handler.Delete).Name("posts.delete")
    }
}
```

### Dependency Injection

Create `provider.go` for Wire:

```go
//go:build wireinject
// +build wireinject

package mymodule

import (
    "github.com/google/wire"
    "gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
    NewRepository,
    NewService,
    NewHandler,
)

func GetHandler(db *gorm.DB) *Handler {
    wire.Build(ProviderSet)
    return nil
}
```

Update `internal/modules/wire.go`:

```go
import mymodule "github.com/zgiai/zgo/internal/modules/mymodule"

var AllProviders = wire.NewSet(
    user.ProviderSet,
    permission.ProviderSet,
    mymodule.ProviderSet,  // Add this
)
```

---

## Best Practices

### 1. Use Syntax Sugar

```go
// ✅ Good - Use one-liner pagination
users, paging, _ := pagination.PaginateFromContext[User](c, db)

// ❌ Avoid - Manual parameter extraction
page, _ := strconv.Atoi(c.Query("page"))
pageSize, _ := strconv.Atoi(c.Query("page_size"))
```

### 2. Hide Sensitive Data

```go
// ✅ Good - Use resource transformers
type UserResource struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // Password NOT included
}

// ❌ Avoid - Returning raw models
c.JSON(200, user) // Exposes password hash!
```

### 3. Consistent Error Responses

```go
// ✅ Good
if err != nil {
    response.InternalServerError(c, "Failed to create user", err)
    return
}

// ❌ Avoid
if err != nil {
    c.JSON(500, gin.H{"error": err.Error()})
    return
}
```

### 4. Validate All Input

```go
// ✅ Good
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required,min=2"`
    Email string `json:"email" binding:"required,email"`
}

if err := c.ShouldBindJSON(&req); err != nil {
    response.BadRequest(c, "Validation failed", err)
    return
}
```

### 5. Use Dependency Injection

```go
// ✅ Good - Constructor injection
type Handler struct {
    service Service
    logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *Handler {
    return &Handler{service: service, logger: logger}
}

// ❌ Avoid - Global variables
var globalDB *gorm.DB
```

### 6. Keep Layers Separate

```
Handler  → Service  → Repository → Database
  ↓          ↓           ↓
HTTP      Business    Data Access
Layer      Logic       Layer
```

- **Handler**: HTTP request/response only
- **Service**: Business logic, orchestration
- **Repository**: Database queries only

### 7. Write Tests

```go
// Unit test
func TestService_CreateUser(t *testing.T) {
    mockRepo := &MockRepository{}
    service := NewService(mockRepo)
    
    user, err := service.CreateUser(ctx, &req)
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### 8. Use Meaningful Names

```go
// ✅ Good
func (h *Handler) ListActivePosts(c *gin.Context)
func (s *Service) CreateUserWithRole(ctx, req)

// ❌ Avoid
func (h *Handler) Get(c *gin.Context)
func (s *Service) Create(ctx, data)
```

---

## Next Steps

- Review [modules/README.md](../internal/modules/README.md) for module architecture
- Check [dependency_injection.md](./dependency_injection.md) for Wire setup
- Read [logging.md](./logging.md) for logging best practices
- Explore [API documentation](/swagger/index.html) for endpoint reference
