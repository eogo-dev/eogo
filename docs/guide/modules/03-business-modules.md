# Business Module Development

> How to create and structure business modules in Eogo.

## Module Structure

Each module follows a consistent structure:

```
internal/modules/blog/
├── model.go       # Database entity
├── dto.go         # Request/Response DTOs
├── repository.go  # Data access layer
├── service.go     # Business logic
├── handler.go     # HTTP handlers
├── routes.go      # Route registration
└── provider.go    # Wire DI configuration
```

## Creating a Module

Use the CLI to generate a new module:

```bash
./eogo make:module Blog
```

This generates all files with proper structure and Wire configuration.

## Layer Responsibilities

### Model (model.go)

Database entity with GORM tags:

```go
type Blog struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
    
    Title     string `gorm:"size:255;not null" json:"title"`
    Content   string `gorm:"type:text" json:"content"`
    AuthorID  uint   `gorm:"index" json:"author_id"`
    Published bool   `gorm:"default:false" json:"published"`
}
```

### DTO (dto.go)

Request validation and response shaping:

```go
// CreateBlogRequest for creating a blog
type CreateBlogRequest struct {
    Title   string `json:"title" binding:"required,min=1,max=255"`
    Content string `json:"content" binding:"required"`
}

// BlogResponse for API responses
type BlogResponse struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Published bool      `json:"published"`
    CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts model to response
func (b *Blog) ToResponse() *BlogResponse {
    return &BlogResponse{
        ID:        b.ID,
        Title:     b.Title,
        Content:   b.Content,
        Published: b.Published,
        CreatedAt: b.CreatedAt,
    }
}
```

### Repository (repository.go)

Data access with interface abstraction:

```go
// Repository defines blog data operations
type Repository interface {
    Create(ctx context.Context, blog *Blog) error
    FindByID(ctx context.Context, id uint) (*Blog, error)
    FindAll(ctx context.Context) ([]*Blog, error)
    Update(ctx context.Context, blog *Blog) error
    Delete(ctx context.Context, id uint) error
}

type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, blog *Blog) error {
    return r.db.WithContext(ctx).Create(blog).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*Blog, error) {
    var blog Blog
    err := r.db.WithContext(ctx).First(&blog, id).Error
    if err != nil {
        return nil, err
    }
    return &blog, nil
}
```

### Service (service.go)

Business logic with interface abstraction:

```go
// Service defines blog business operations
type Service interface {
    Create(ctx context.Context, req *CreateBlogRequest, authorID uint) (*Blog, error)
    GetByID(ctx context.Context, id uint) (*Blog, error)
    List(ctx context.Context) ([]*Blog, error)
    Publish(ctx context.Context, id uint) error
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *CreateBlogRequest, authorID uint) (*Blog, error) {
    blog := &Blog{
        Title:    req.Title,
        Content:  req.Content,
        AuthorID: authorID,
    }
    
    if err := s.repo.Create(ctx, blog); err != nil {
        return nil, err
    }
    return blog, nil
}

func (s *service) Publish(ctx context.Context, id uint) error {
    blog, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    
    blog.Published = true
    return s.repo.Update(ctx, blog)
}
```

### Handler (handler.go)

HTTP handlers using the response package:

```go
type Handler struct {
    service Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
    var req CreateBlogRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request", err)
        return
    }
    
    userID := c.GetUint("user_id")
    blog, err := h.service.Create(c.Request.Context(), &req, userID)
    if err != nil {
        response.InternalServerError(c, "Failed to create blog", err)
        return
    }
    
    response.Created(c, blog.ToResponse())
}

func (h *Handler) Show(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        response.BadRequest(c, "Invalid ID", err)
        return
    }
    
    blog, err := h.service.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            response.NotFound(c, "Blog not found", err)
            return
        }
        response.InternalServerError(c, "Failed to get blog", err)
        return
    }
    
    response.OK(c, blog.ToResponse())
}
```

### Routes (routes.go)

Route registration:

```go
func Register(r *router.Router, h *Handler) {
    blogs := r.Group("/blogs")
    {
        blogs.GET("", h.Index)
        blogs.GET("/:id", h.Show)
        
        // Protected routes
        auth := blogs.Group("").Use(middleware.Auth())
        auth.POST("", h.Create)
        auth.PUT("/:id", h.Update)
        auth.DELETE("/:id", h.Delete)
        auth.POST("/:id/publish", h.Publish)
    }
}
```

### Provider (provider.go)

Wire dependency injection:

```go
var ProviderSet = wire.NewSet(
    NewRepository,
    wire.Bind(new(Repository), new(*repository)),
    NewService,
    wire.Bind(new(Service), new(*service)),
    NewHandler,
)
```

## Registering the Module

### 1. Add to Wire

```go
// internal/wiring/wire.go
func InitApp(db *gorm.DB) (*app.Application, error) {
    wire.Build(
        // ... existing providers
        blog.ProviderSet,  // Add module
        wire.Struct(new(app.Application), "*"),
    )
    return nil, nil
}
```

### 2. Add to App Handlers

```go
// internal/app/app.go
type Handlers struct {
    User       *user.Handler
    Permission *permission.Handler
    Blog       *blog.Handler  // Add handler
}
```

### 3. Register Routes

```go
// routes/api.go
func Setup(r *gin.Engine, h *app.Handlers) {
    api := r.Group("/v1")
    
    user.Register(api, h.User)
    permission.Register(api, h.Permission)
    blog.Register(api, h.Blog)  // Add routes
}
```

### 4. Add Migration

```go
// database/migrations/blog_migrations.go
func CreateBlogsTable() *gormigrate.Migration {
    return &gormigrate.Migration{
        ID: "202512270001_create_blogs",
        Migrate: func(db *gorm.DB) error {
            return db.AutoMigrate(&blog.Blog{})
        },
        Rollback: func(db *gorm.DB) error {
            return db.Migrator().DropTable("blogs")
        },
    }
}

// database/migrations/migrations.go
func All() []*gormigrate.Migration {
    return []*gormigrate.Migration{
        // ... existing
        CreateBlogsTable(),
    }
}
```

### 5. Regenerate Wire

```bash
cd internal/wiring && wire
```

## Pagination

Use the pagination package for list endpoints:

```go
func (h *Handler) Index(c *gin.Context) {
    paginator, err := pagination.PaginateFromContext[Blog](c, h.db)
    if err != nil {
        response.InternalServerError(c, "Failed to fetch blogs", err)
        return
    }
    
    c.JSON(http.StatusOK, paginator)
}
```

Response format:

```json
{
  "items": [...],
  "total": 100,
  "per_page": 15,
  "current_page": 1,
  "last_page": 7
}
```

## Testing

### Unit Test (Service)

```go
func TestBlogService_Create(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := NewMockRepository(ctrl)
    svc := NewService(mockRepo)
    
    req := &CreateBlogRequest{
        Title:   "Test Blog",
        Content: "Content",
    }
    
    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(nil)
    
    blog, err := svc.Create(context.Background(), req, 1)
    
    assert.NoError(t, err)
    assert.Equal(t, "Test Blog", blog.Title)
}
```

### Feature Test (API)

```go
func TestBlogCreate(t *testing.T) {
    tc := NewTestCase(t)
    
    // Login first
    token := tc.LoginAs("testuser")
    
    tc.Post("/v1/blogs").
        WithHeader("Authorization", "Bearer "+token).
        WithJSON(map[string]any{
            "title":   "My Blog",
            "content": "Hello World",
        }).
        Call().
        AssertCreated().
        AssertJSONPath("data.title", "My Blog")
}
```

## Best Practices

### DO ✅

- Use interfaces for Repository and Service
- Pass `context.Context` to all methods
- Use DTOs for request/response
- Validate input in handlers
- Use the response package for consistency

### DON'T ❌

- Put business logic in handlers
- Access database directly in handlers
- Return models directly (use DTOs)
- Ignore errors
- Skip context propagation

## Next Steps

- [Database Operations](../infrastructure/01-database.md)
- [Testing Guide](../best-practices/01-testing.md)
