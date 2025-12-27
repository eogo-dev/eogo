# Modules

This directory contains all business domain modules for the ZGI Console API.

## Module Overview

| Module | Description | Type |
|--------|-------------|------|
| `dashboard` | Admin dashboard statistics | Simple |
| `finance` | Financial management (plans, transactions, recharge) | Sub-modules |
| `llm` | LLM provider/model/channel management | Sub-modules |
| `permission` | RBAC (roles, permissions) | Simple |
| `tenant` | Tenant management | Simple |
| `user` | User authentication | Simple |

## Module Types

### Simple Module

For single-domain functionality with < 500 lines per file:

```text
module_name/
├── model.go        # Database entity
├── dto.go          # Request/Response DTOs
├── repository.go   # Data access (interface + impl)
├── service.go      # Business logic (interface + impl)
├── handler.go      # HTTP handlers
├── routes.go       # Route registration
└── module.go       # Module initialization
```

### Sub-module Structure

For complex domains with multiple independent sub-domains:

```text
module_name/
├── submodule1/     # Independent sub-domain
│   ├── dto.go
│   ├── repository.go
│   ├── service.go
│   ├── handler.go
│   └── service_test.go    # Private method tests
├── submodule2/
│   └── ...
├── model.go        # Shared models (all entities)
└── router.go       # Centralized route registration
```

## Key Files

### model.go (Parent Module Root)

All entity definitions shared across sub-modules:

```go
package llm

type Provider struct {
    ID          string `gorm:"primaryKey;size:36"`
    Name        string `gorm:"size:50;uniqueIndex"`
    DisplayName string `gorm:"size:100"`
    // ...
}

type Model struct {
    ID       string `gorm:"primaryKey;size:36"`
    Provider string `gorm:"size:50;index"`
    Name     string `gorm:"size:100"`
    // ...
}

type Channel struct {
    // ...
}
```

### router.go (Parent Module Root)

Centralized route registration for all sub-modules:

```go
package llm

func RegisterAdminRoutes(r *gin.RouterGroup, db *gorm.DB) {
    llm := r.Group("/llm")

    // Provider routes
    providerRepo := provider.NewRepository(db)
    providerSvc := provider.NewService(providerRepo)
    providerHandler := provider.NewHandler(providerSvc)
    
    providers := llm.Group("/providers")
    {
        providers.POST("", providerHandler.Create)
        providers.GET("", providerHandler.List)
        // ...
    }

    // Model routes
    // ... similar pattern
}
```

## Testing Guidelines

### In-Module Tests (Private Methods)

Place `*_test.go` files inside the module to test unexported functions:

```go
// internal/modules/llm/provider/service_test.go
package provider

func Test_detectModelType(t *testing.T) {
    // Test private helper function
    result := detectModelType("gpt-4")
    assert.Equal(t, "llm", result)
}
```

### External Tests (Public API)

Place in `tests/` directory:

- `tests/unit/[module]/` - Service tests with mocked repos
- `tests/e2e/` - Full integration tests

## Creating a New Module

### Simple Module

```bash
./eogo make:module Blog
```

### Sub-module Structure

1. Create parent directory: `internal/modules/mymodule/`
2. Create `model.go` with all entities
3. Create sub-directories for each sub-domain
4. Create `router.go` for centralized routing
5. Register in `routes/admin.go`:

```go
import "github.com/eogo-dev/eogo/internal/modules/mymodule"

// In RegisterAdminRoutes:
mymodule.RegisterAdminRoutes(admin, db)
```

## Current Module Structure

### LLM Module (`llm/`)

```text
llm/
├── provider/       # LLM provider management
├── llmmodel/       # LLM model management
├── channel/        # Channel management
├── credential/     # API credential management
├── sync/           # Sync from LMBase
├── redeem/         # Redeem code management
├── recyclebin/     # Soft-deleted items
├── statistics/     # Usage statistics
├── model.go        # Shared entities
└── router.go       # Centralized routes
```

### Finance Module (`finance/`)

```text
finance/
├── plan/           # Subscription plans
├── transaction/    # Transaction records
├── recharge/       # Recharge orders
└── router.go       # Centralized routes
```

## Best Practices

1. **Flat First** - Avoid unnecessary nesting
2. **One Package Per Directory** - All files share same package
3. **Interface + Impl Together** - No separate `_impl.go` files
4. **Centralized Routing** - One `router.go` per parent module
5. **Shared Models** - One `model.go` for all entities
6. **In-Module Tests** - Test private methods inside module
