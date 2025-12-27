# Database Operations

> GORM database operations, migrations, and best practices.

## Configuration

```go
// internal/infra/config/config.go
type DatabaseConfig struct {
    Enabled      bool   `env:"DB_ENABLED" envDefault:"true"`
    Driver       string `env:"DB_DRIVER" envDefault:"mysql"`
    Host         string `env:"DB_HOST" envDefault:"localhost"`
    Port         int    `env:"DB_PORT" envDefault:"3306"`
    Database     string `env:"DB_DATABASE" envDefault:"eogo"`
    Username     string `env:"DB_USERNAME" envDefault:"root"`
    Password     string `env:"DB_PASSWORD" envDefault:""`
    MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
    MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"100"`
    Memory       bool   `env:"DB_MEMORY" envDefault:"false"` // SQLite in-memory
}
```

Environment variables:

```bash
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=eogo
DB_USERNAME=root
DB_PASSWORD=secret
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
```

## Supported Drivers

| Driver | DSN Format |
|--------|------------|
| MySQL | `user:pass@tcp(host:port)/db?charset=utf8mb4&parseTime=True` |
| PostgreSQL | `host=localhost user=user password=pass dbname=db port=5432` |
| SQLite | `file.db` or `:memory:` |

## Model Definition

```go
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`  // Auto-set
    UpdatedAt time.Time      `json:"updated_at"`  // Auto-updated
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
    
    Username string `gorm:"size:100;uniqueIndex;not null" json:"username"`
    Email    string `gorm:"size:255;uniqueIndex;not null" json:"email"`
    Password string `gorm:"size:255;not null" json:"-"`
    Status   int    `gorm:"default:1" json:"status"`
}

// TableName overrides the table name
func (User) TableName() string {
    return "users"
}
```

## Repository Pattern

```go
// Repository interface
type Repository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uint) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, page, perPage int) ([]*User, int64, error)
}

// Implementation
type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *repository) Update(ctx context.Context, user *User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&User{}, id).Error
}

func (r *repository) List(ctx context.Context, page, perPage int) ([]*User, int64, error) {
    var users []*User
    var total int64
    
    r.db.WithContext(ctx).Model(&User{}).Count(&total)
    
    offset := (page - 1) * perPage
    err := r.db.WithContext(ctx).
        Offset(offset).
        Limit(perPage).
        Find(&users).Error
    
    return users, total, err
}
```

## Query Builder

### Basic Queries

```go
// Find by ID
db.First(&user, 1)

// Find by condition
db.Where("email = ?", email).First(&user)

// Find all
db.Find(&users)

// Select specific columns
db.Select("id", "username").Find(&users)

// Order
db.Order("created_at DESC").Find(&users)

// Limit & Offset
db.Offset(10).Limit(20).Find(&users)
```

### Advanced Queries

```go
// Multiple conditions
db.Where("status = ? AND role = ?", 1, "admin").Find(&users)

// OR condition
db.Where("status = ?", 1).Or("role = ?", "admin").Find(&users)

// IN clause
db.Where("id IN ?", []uint{1, 2, 3}).Find(&users)

// LIKE
db.Where("username LIKE ?", "%john%").Find(&users)

// Raw SQL
db.Raw("SELECT * FROM users WHERE status = ?", 1).Scan(&users)
```

### Preloading Relations

```go
// Preload single relation
db.Preload("Profile").Find(&users)

// Preload multiple relations
db.Preload("Profile").Preload("Roles").Find(&users)

// Preload with conditions
db.Preload("Orders", "status = ?", "completed").Find(&users)

// Nested preload
db.Preload("Orders.Items").Find(&users)
```

## Transactions

```go
func (r *repository) Transfer(ctx context.Context, fromID, toID uint, amount int) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Deduct from sender
        if err := tx.Model(&Account{}).
            Where("id = ?", fromID).
            Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
            return err
        }
        
        // Add to receiver
        if err := tx.Model(&Account{}).
            Where("id = ?", toID).
            Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

### Manual Transaction

```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&profile).Error; err != nil {
    tx.Rollback()
    return err
}

return tx.Commit().Error
```

## Migrations

### Migration File

```go
// database/migrations/user_migrations.go
func CreateUsersTable() *gormigrate.Migration {
    return &gormigrate.Migration{
        ID: "202512270001_create_users",
        Migrate: func(db *gorm.DB) error {
            return db.AutoMigrate(&user.User{})
        },
        Rollback: func(db *gorm.DB) error {
            return db.Migrator().DropTable("users")
        },
    }
}

func AddAvatarToUsers() *gormigrate.Migration {
    return &gormigrate.Migration{
        ID: "202512270002_add_avatar_to_users",
        Migrate: func(db *gorm.DB) error {
            return db.Migrator().AddColumn(&user.User{}, "Avatar")
        },
        Rollback: func(db *gorm.DB) error {
            return db.Migrator().DropColumn(&user.User{}, "Avatar")
        },
    }
}
```

### Register Migrations

```go
// database/migrations/migrations.go
func All() []*gormigrate.Migration {
    return []*gormigrate.Migration{
        CreateUsersTable(),
        CreateRolesTable(),
        AddAvatarToUsers(),
    }
}
```

### Run Migrations

```bash
./eogo migrate
```

## Pagination

```go
import "github.com/eogo-dev/eogo/internal/infra/pagination"

// From Gin context (auto-extracts ?page=1&per_page=15)
paginator, err := pagination.PaginateFromContext[User](c, db)

// Manual pagination
paginator, err := pagination.Paginate[User](db, page, perPage)

// With query conditions
query := db.Where("status = ?", 1)
paginator, err := pagination.PaginateFromContext[User](c, query)
```

Response:

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

## Soft Delete

Models with `DeletedAt` field support soft delete:

```go
// Soft delete
db.Delete(&user, 1)

// Query excludes soft deleted by default
db.Find(&users) // Only non-deleted

// Include soft deleted
db.Unscoped().Find(&users)

// Permanently delete
db.Unscoped().Delete(&user, 1)

// Restore soft deleted
db.Model(&user).Update("deleted_at", nil)
```

## Hooks

```go
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Hash password before create
    hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashed)
    return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
    // Send welcome email
    return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
    // Validate before update
    return nil
}
```

## Connection Pool

```go
sqlDB, _ := db.DB()

// Set connection pool settings
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
sqlDB.SetConnMaxIdleTime(10 * time.Minute)
```

## Tracing

Enable database tracing:

```go
import "github.com/eogo-dev/eogo/internal/infra/tracing"

// Add tracing plugin
tracing.WithTracing(db, "eogo")
```

## Best Practices

### DO ✅

- Use context for all queries
- Use transactions for multi-step operations
- Use soft delete for important data
- Index frequently queried columns
- Use pagination for list endpoints

### DON'T ❌

- Use `db.Exec` for user input (SQL injection)
- Ignore errors
- Query without limits
- Use `SELECT *` when not needed
- Skip migrations

## Next Steps

- [Business Modules](../modules/03-business-modules.md)
- [Testing Guide](../best-practices/01-testing.md)
