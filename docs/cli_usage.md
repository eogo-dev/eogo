# ZGO CLI Usage Guide

## Building the CLI

```bash
make build
# or
go build -o zgo cmd/zgo/main.go
```

This creates the `zgo` binary in the root directory.

## Available Commands

### General Commands

```bash
./zgo --help          # Show all commands
./zgo version         # Show version
./zgo env             # Display environment
./zgo serve           # Start HTTP server (same as `make server`)
```

### Database Migrations

```bash
# Run all pending migrations
./zgo migrate
./zgo db:migrate      # alias

# Rollback last migration batch
./zgo migrate:rollback
./zgo db:rollback     # alias

# Drop all tables and re-run migrations
./zgo migrate:fresh
./zgo db:fresh        # alias

# Show migration status
./zgo migrate:status
./zgo db:status       # alias
```

**Migration Status Output:**
```
Migration Status
──────────────────
Migration      Batch   Status   
───────        ─────   ──────   
create_users_table           1   Ran                    
create_teams_table           1   Ran                    
create_organizations_table   2   Ran
```

### Database Seeders

```bash
./zgo seed
./zgo db:seed
```

### Code Generation

#### Create Migration

```bash
./zgo make:migration create_posts_table
```

Creates: `database/migrations/YYYY_MM_DD_HHMMSS_create_posts_table.go`

Example generated file:
```go
package migrations

import (
    "github.com/go-gormigrate/gormigrate/v2"
    "gorm.io/gorm"
)

func init() {
    Migrations = append(Migrations, &gormigrate.Migration{
        ID: "create_posts_table",
        Migrate: func(tx *gorm.DB) error {
            // Add your migration here
            return nil
        },
        Rollback: func(tx *gorm.DB) error {
            // Add your rollback here
            return nil
        },
    })
}
```

#### Create Module

```bash
./zgo make:module Blog
```

Creates complete module structure:
```
internal/modules/blog/
├── model.go
├── dto.go
├── repository.go
├── service.go
├── handler.go
├── routes.go
└── provider.go
```

#### Create Individual Components

```bash
./zgo make:model Post
./zgo make:repository PostRepository
./zgo make:service PostService
./zgo make:handler PostHandler
./zgo make:seeder PostSeeder
```

### Routes

```bash
./zgo route:list
```

**Output:**
```
Registered Routes
───────────────────
Method    Path                  Handler
───────   ───────               ───────
GET       /                     Welcome
GET       /v1/users             user.List
GET       /v1/users/:id         user.Get
POST      /v1/login             user.Login
PUT       /v1/roles/:id         permission.UpdateRole
DELETE    /v1/roles/:id         permission.DeleteRole
```

## Migration Directory Structure

```
database/
├── migrations/
│   ├── 2025_06_18_000000_create_users_table.go
│   ├── 2025_06_18_000001_seed_default_users.go
│   ├── 2025_12_26_000000_create_roles_table.go
│   ├── 2025_12_26_000001_create_permissions_table.go
│   └── migrations.go              # Auto-registration
└── seeders/
    ├── default_users.go
    └── seeders.go
```

**Important:** The `cmd/migrate/` directory is empty and unused. Migrations are managed through `cmd/zgo/` CLI.

## Common Workflows

### 1. Create and Run Migration

```bash
# Create migration
./zgo make:migration create_posts_table

# Edit the file: database/migrations/YYYY_MM_DD_HHMMSS_create_posts_table.go

# Run migration
./zgo migrate

# Check status
./zgo migrate:status
```

### 2. Create Full Module

```bash
# Generate module
./zgo make:module Post

# Add migration
./zgo make:migration create_posts_table

# Update internal/modules/wire.go to include PostProviderSet
# Update routes/router.go to register routes

# Run migration
./zgo migrate

# Start server
make dev
```

### 3. Reset Database

```bash
# WARNING: This drops all tables!
./zgo migrate:fresh

# Optionally run seeders
./zgo seed
```

## Development Workflow

```bash
# 1. Build CLI
make build

# 2. Run migrations
./zgo migrate

# 3. Start dev server with hot reload
make dev

# 4. In another terminal, list routes
./zgo route:list

# 5. Create new module
./zgo make:module Product
```

## Notes

- **Binary location:** Root directory (`./zgo`)
- **Migration files:** `database/migrations/`
- **Auto-registration:** Migrations use `init()` to register
- **Naming:** Migration files are timestamped
- **`cmd/migrate/` is unused** - Use `./zgo migrate` instead

## See Also

- [Platform Guide](./platform_guide.md) - Development patterns
- [Dependency Injection](./dependency_injection.md) - Wire setup
- [Module README](../internal/modules/README.md) - Module architecture
