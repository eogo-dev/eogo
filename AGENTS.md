# AGENTS.md

Instructions for AI coding agents working on the Eogo framework.

## Project Overview

EOGO (Evolving Orchestration for Go) is a modern Go framework designed for high-velocity development. It features a clean, modular architecture with automated tooling.

## Core Modules

The framework provides two essential modules as examples:
- **`user`**: Authentication (register, login, JWT)
- **`permission`**: RBAC (roles, permissions, user-role assignments)

**Note**: Business-specific modules should be created using `eogo make:module` and kept separate from the framework core.

## Setup & Commands

```bash
# Install dependencies
go mod download

# Build the framework CLI
make build

# Install globally (optional)
go install ./cmd/eogo

# Run tests (Unit, Integration, Feature)
make test

# Start development server (Air)
make air
```

## Directory Structure

```text
├── cmd/eogo              # Framework CLI
├── cmd/server            # Application entry point
├── internal/
│   ├── bootstrap/        # App/HTTP/Console Kernels & Migrations
│   ├── modules/          # Domain Modules (user, permission)
│   └── platform/         # Framework Core (database, router, pagination, response)
├── routes/
│   ├── api.go            # API registration
│   └── router.go         # HTTP kernel
└── tests/
    ├── unit/             # Mock-based unit tests
    ├── integration/      # Database integration tests
    └── feature/          # End-to-end API tests (In-memory DB)
```

## Framework Conventions

### 1. Auto Timestamps (GORM)
All models automatically get `CreatedAt` and `UpdatedAt` timestamps:

```go
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`  // Auto-set on create
    UpdatedAt time.Time      `json:"updated_at"`  // Auto-updated on save
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
    // ... other fields
}
```

**GORM handles this automatically** - no manual intervention needed.

### 2. Unified Response Structure
Use `internal/platform/response` for consistent API responses:

```go
import "github.com/eogo-dev/eogo/internal/platform/response"

// Success response
response.Success(c, data)
response.OK(c, data)
response.Created(c, data)
response.NoContent(c)

// Error responses
response.BadRequest(c, "Invalid input", err)
response.Unauthorized(c, "Not authenticated")
response.NotFound(c, "Resource not found", err)
response.InternalServerError(c, "Server error", err)
```

**Standard Response Format:**
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 3. Pagination
Use `internal/platform/pagination` for paginated results:

```go
import "github.com/eogo-dev/eogo/internal/platform/pagination"

// From Gin context (auto-extracts ?page=1&per_page=15)
paginator, err := pagination.PaginateFromContext[User](c, db)

// Manual pagination
paginator, err := pagination.Paginate[User](db, page, perPage)

// Return paginated response
c.JSON(http.StatusOK, paginator)
```

**Pagination Response:**
```json
{
  "items": [...],
  "total": 100,
  "per_page": 15,
  "current_page": 1,
  "last_page": 7,
  "from": 1,
  "to": 15,
  "next_page_url": "/api/v1/users?page=2&per_page=15",
  "prev_page_url": null
}
```

**Alternative Pagination Methods:**
- `SimplePaginate[T]()` - Faster, no total count
- `CursorPaginate[T]()` - For infinite scroll

## Dependency Injection (Wire)

Eogo uses **Google Wire** for compile-time dependency injection, similar to NestJS but type-safe with zero runtime overhead.

### Provider Pattern
Each module defines a `ProviderSet` in `provider.go`:

```go
// internal/modules/user/provider.go
var ProviderSet = wire.NewSet(
    NewRepository,                                    // Constructor
    wire.Bind(new(Repository), new(*UserRepositoryImpl)), // Interface → Implementation
    NewService,
    wire.Bind(new(Service), new(*UserServiceImpl)),
    NewHandler,
)
```

**What this means:**
- `NewRepository` creates `*UserRepositoryImpl`
- `wire.Bind` says: "when someone needs `Repository` interface, provide `*UserRepositoryImpl`"
- `NewService` depends on `Repository` interface (not concrete type)
- Wire generates the wiring code at compile-time

### Central Aggregation
```go
// internal/modules/wire.go
func InitApp(db *gorm.DB) (*App, error) {
    wire.Build(
        config.MustLoad,
        jwt.NewService,
        user.ProviderSet,       // Import module providers
        permission.ProviderSet,
        wire.Struct(new(App), "*"),
    )
    return nil, nil
}
```

### Comparison with NestJS
| NestJS | Eogo (Wire) |
|--------|-------------|
| `@Injectable()` | Constructor (`NewService`) |
| `@Module({ providers })` | `wire.NewSet(...)` |
| Runtime DI | Compile-time generation |

**Benefits:**
- ✅ Type-safe (compile errors for missing deps)
- ✅ Zero runtime cost (no reflection)
- ✅ Explicit dependencies
- ✅ Easy testing (mock interfaces)

**After creating a module**, run:
```bash
cd internal/modules && wire
```

See [docs/dependency_injection.md](file:///Users/stark/item/eogo/eogo/docs/dependency_injection.md) for details.

## Development SOP

### 1. Creating a Module
Always use the CLI to maintain consistency:
```bash
./eogo make:module Blog
```
This generates `internal/modules/blog/` with:
- `model.go`: Database entity with auto timestamps
- `dto.go`: Request/Response DTOs
- `repository.go`: Data access layer
- `service.go`: Business logic
- `handler.go`: HTTP handlers
- `routes.go`: Route registration
- `provider.go`: Dependency injection

### 2. Route Registration
Modules self-register their routes:
- **Module**: Define `Register(r *router.Router)` in `routes.go`
- **Global**: Add `blog.Register(r)` to `routes/api.go`

### 3. Database Migrations
Migrations are in `database/migrations/` (Laravel-style):
- Each module has its own migration file (e.g., `user_migrations.go`)
- All migrations are registered in `migrations.go`
- Run with `./eogo migrate`

**Adding a new migration:**
```go
// database/migrations/blog_migrations.go
func CreateBlogsTable() *gormigrate.Migration {
    return &gormigrate.Migration{
        ID: "202512270_create_blogs",
        Migrate: func(db *gorm.DB) error {
            return db.AutoMigrate(&blog.Blog{})
        },
        Rollback: func(db *gorm.DB) error {
            return db.Migrator().DropTable("blogs")
        },
    }
}
```

Then register in `migrations.go`:
```go
func All() []*gormigrate.Migration {
    return []*gormigrate.Migration{
        // ... existing
        CreateBlogsTable(), // Add here
    }
}
```

See [database/README.md](file:///Users/stark/item/eogo/eogo/database/README.md) for details.

## Testing SOP

### Unit Testing (Service Layer)
- Use mocks for repositories
- Place in `internal/modules/[module]/service_test.go`

### Feature Testing (API Layer)
Feature tests use real HTTP + in-memory SQLite:
```go
func TestUserRegistration(t *testing.T) {
    tc := NewTestCase(t)
    tc.Post("/v1/register").
       WithJSON(gin.H{"username": "test"}).
       Call().
       AssertOk()
}
```

## Guardrails
- **No Import Cycles**: Platform packages are pure infrastructure
- **Interfaces**: Services and Handlers depend on interfaces
- **Context**: Every Service/Repo method accepts `ctx context.Context`
- **Naming**: Use `snake_case` for JSON, `access_token` for JWT
- **Language**: All code and comments in English
- **Timestamps**: Use GORM's auto timestamps (`CreatedAt`, `UpdatedAt`)
- **Responses**: Use `response` package for consistency
- **Pagination**: Use `pagination` package for list endpoints
