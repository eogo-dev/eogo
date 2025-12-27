# Eogo CLI Usage Guide

## Building the CLI

```bash
make build
# or
go build -o eogo cmd/eogo/main.go
```

This creates the `eogo` binary in the root directory.

## Available Commands

### General Commands

```bash
./eogo --help          # Show all commands
./eogo version         # Show version
./eogo env             # Display environment
./eogo serve           # Start HTTP server (same as `make server`)
```

### Database Migrations

```bash
# Run all pending migrations
./eogo migrate
./eogo db:migrate      # alias

# Rollback last migration batch
./eogo migrate:rollback
./eogo db:rollback     # alias

# Drop all tables and re-run migrations
./eogo migrate:fresh
./eogo db:fresh        # alias

# Show migration status
./eogo migrate:status
./eogo db:status       # alias
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
./eogo seed
./eogo db:seed
```

### Code Generation

#### Create Migration

```bash
./eogo make:migration create_posts_table
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
./eogo make:module Blog
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
./eogo make:model Post
./eogo make:repository PostRepository
./eogo make:service PostService
./eogo make:handler PostHandler
./eogo make:seeder PostSeeder
```

### Routes

```bash
./eogo route:list
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

**Important:** The `cmd/migrate/` directory is empty and unused. Migrations are managed through `cmd/eogo/` CLI.

## Common Workflows

### 1. Create and Run Migration

```bash
# Create migration
./eogo make:migration create_posts_table

# Edit the file: database/migrations/YYYY_MM_DD_HHMMSS_create_posts_table.go

# Run migration
./eogo migrate

# Check status
./eogo migrate:status
```

### 2. Create Full Module

```bash
# Generate module
./eogo make:module Post

# Add migration
./eogo make:migration create_posts_table

# Update internal/modules/wire.go to include PostProviderSet
# Update routes/router.go to register routes

# Run migration
./eogo migrate

# Start server
make dev
```

### 3. Reset Database

```bash
# WARNING: This drops all tables!
./eogo migrate:fresh

# Optionally run seeders
./eogo seed
```

## Development Workflow

```bash
# 1. Build CLI
make build

# 2. Run migrations
./eogo migrate

# 3. Start dev server with hot reload
make dev

# 4. In another terminal, list routes
./eogo route:list

# 5. Create new module
./eogo make:module Product
```

## Notes

- **Binary location:** Root directory (`./eogo`)
- **Migration files:** `database/migrations/`
- **Auto-registration:** Migrations use `init()` to register
- **Naming:** Migration files are timestamped
- **`cmd/migrate/` is unused** - Use `./eogo migrate` instead

## See Also

- [Platform Guide](./platform_guide.md) - Development patterns
- [Dependency Injection](./dependency_injection.md) - Wire setup
- [Module README](../internal/modules/README.md) - Module architecture
