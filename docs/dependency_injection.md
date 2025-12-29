# Dependency Injection with Wire

## Overview

ZGO uses [Google Wire](https://github.com/google/wire) for compile-time dependency injection, similar to NestJS's DI system but **type-safe and zero runtime overhead**.

## Philosophy: Distributed Providers

Instead of a monolithic DI configuration, each module defines its own **ProviderSet** - a self-contained unit that declares:
- What it provides (constructors)
- What interfaces it implements (bindings)

This is conceptually similar to NestJS modules with `@Injectable()` decorators.

## Provider Pattern

### Module Structure
```
internal/modules/user/
├── model.go       # Domain entities
├── repository.go  # Data access (interface + impl)
├── service.go     # Business logic (interface + impl)
├── handler.go     # HTTP handlers
├── provider.go    # ⭐ DI configuration
└── routes.go      # Route registration
```

### Provider Example
```go
// internal/modules/user/provider.go
package user

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
    NewRepository,                                    // Constructor
    wire.Bind(new(Repository), new(*UserRepositoryImpl)), // Interface binding
    NewService,
    wire.Bind(new(Service), new(*UserServiceImpl)),
    NewHandler,
)
```

**What this does:**
1. `NewRepository` creates `*UserRepositoryImpl`
2. `wire.Bind` tells Wire: "when someone needs `Repository`, provide `*UserRepositoryImpl`"
3. Same pattern for Service
4. `NewHandler` depends on `Service` interface (not implementation)

### Central Aggregation
```go
// internal/modules/wire.go
package app

import (
    "github.com/zgiai/zgo/internal/modules/user"
    "github.com/zgiai/zgo/internal/modules/permission"
    "github.com/google/wire"
)

type App struct {
    User       *user.Handler
    Permission *permission.Handler
}

func InitApp(db *gorm.DB) (*App, error) {
    wire.Build(
        config.MustLoad,
        jwt.NewService,
        user.ProviderSet,       // ⭐ Import module providers
        permission.ProviderSet,
        wire.Struct(new(App), "*"),
    )
    return nil, nil
}
```

## Comparison with NestJS

| NestJS | ZGO (Wire) |
|--------|-------------|
| `@Injectable()` | Constructor function (`NewService`) |
| `@Module({ providers: [...] })` | `wire.NewSet(...)` |
| Runtime DI container | Compile-time code generation |
| `useClass`, `useValue` | `wire.Bind`, `wire.Value` |
| Circular dependency detection at runtime | Compile-time error |

## Creating a New Module

### 1. Generate Scaffold
```bash
./zgo make:module Blog
```

This auto-generates `provider.go`:
```go
package blog

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
    NewRepository,
    wire.Bind(new(Repository), new(*RepositoryImpl)),
    NewService,
    wire.Bind(new(Service), new(*ServiceImpl)),
    NewHandler,
)
```

### 2. Register in Central Wire
Add to `internal/modules/wire.go`:
```go
import "github.com/zgiai/zgo/internal/modules/blog"

type App struct {
    User       *user.Handler
    Permission *permission.Handler
    Blog       *blog.Handler  // ⭐ Add here
}

func InitApp(db *gorm.DB) (*App, error) {
    wire.Build(
        config.MustLoad,
        jwt.NewService,
        user.ProviderSet,
        permission.ProviderSet,
        blog.ProviderSet,  // ⭐ Add here
        wire.Struct(new(App), "*"),
    )
    return nil, nil
}
```

### 3. Generate DI Code
```bash
cd internal/modules && wire
```

This generates `wire_gen.go` with all the wiring code.

## Benefits

✅ **Type Safety**: Compile-time errors for missing dependencies  
✅ **Zero Runtime Cost**: No reflection, pure Go code  
✅ **Explicit Dependencies**: Clear constructor signatures  
✅ **Testability**: Easy to mock interfaces  
✅ **Scalability**: Each module is self-contained  

## Common Patterns

### Interface Binding
```go
wire.Bind(new(Service), new(*ServiceImpl))
```
"When someone needs `Service`, provide `*ServiceImpl`"

### Value Providers
```go
wire.Value(&Config{Port: 8080})
```

### Conditional Providers
```go
func ProvideLogger(cfg *Config) Logger {
    if cfg.Debug {
        return NewDebugLogger()
    }
    return NewProductionLogger()
}
```

## Troubleshooting

### "no provider found"
- Ensure constructor is in `ProviderSet`
- Check interface bindings are correct
- Verify central `wire.go` imports the module

### "unused provider"
- Remove from `ProviderSet` if not needed
- Or add to `App` struct if it should be exposed

### Circular dependencies
Wire will error at compile-time. Refactor to break the cycle.
