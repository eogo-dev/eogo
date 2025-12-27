# EOGO Business Development Guide

## Module Structure

```text
internal/modules/post/
├── model.go       # Database models
├── dto.go         # Request/Response DTOs
├── repository.go  # Data access layer
├── service.go     # Business logic layer
├── handler.go     # HTTP handler layer
└── routes.go      # Route definitions
```

## Development Steps

### 1. Model (Database Structure)

Define your entities in `model.go` using GORM tags.

```go
package post

type Post struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `gorm:"size:200;not null" json:"title"`
    // ...
}
```

### 2. DTO (Data Transfer Objects)

Define request and response structures in `dto.go`.

```go
package post

type CreateRequest struct {
    Title   string `json:"title" binding:"required"`
}
```

### 3. Repository (Data Access)

Implement data persistence logic in `repository.go`.

### 4. Service (Business Logic)

Implement core business rules and domain logic in `service.go`.

### 5. Handler (HTTP Layer)

Handle request binding and response formatting in `handler.go`.

### 6. Routes (Self-Registration)

Define routes and wire dependencies in `routes.go`.

```go
package post

func Register(r *router.Router) {
    db := database.GetDB()
    repo := NewRepository(db)
    // ... wiring
    r.POST("/posts", handler.Create)
}
```

### 7. Global Registration

Add your module to `routes/api.go`:

```go
func Register(r *router.Router) {
    // ...
    post.Register(r)
}
```

---

## Layer Responsibilities

| Layer | Responsibility | Dependencies |
|---|------|------|
| Handler | Request validation, HTTP responses | Service |
| Service | Business logic, permissions, transactions | Repository |
| Repository | Database operations | gorm.DB |

## Coding Rules

1. **Handlers** must not perform database operations directly.
2. **Services** must not handle HTTP-specific logic (e.g., Gin context).
3. **Repositories** focus on data access, not business logic.
4. Every method must accept `ctx context.Context` as its first argument.
5. Depend on interfaces rather than concrete implementations for Services and Repositories.
