# Testing Guide

> Unit, integration, and feature testing strategies for ZGO.

## Test Structure

```
tests/
├── unit/           # Mock-based unit tests
├── integration/    # Database integration tests
└── feature/        # End-to-end API tests
```

## Unit Testing

Test business logic in isolation using mocks.

### Service Test

```go
// internal/modules/user/service_test.go
func TestUserService_Register(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    // Create mocks
    mockRepo := NewMockRepository(ctrl)
    mockJWT := jwt.NewMockService(ctrl)
    
    // Create service with mocks
    svc := NewService(mockRepo, mockJWT)
    
    // Setup expectations
    mockRepo.EXPECT().
        FindByEmail(gomock.Any(), "test@example.com").
        Return(nil, gorm.ErrRecordNotFound)
    
    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        DoAndReturn(func(ctx context.Context, user *User) error {
            user.ID = 1
            return nil
        })
    
    // Execute
    req := &RegisterRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    user, err := svc.Register(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "testuser", user.Username)
    assert.Equal(t, uint(1), user.ID)
}
```

### Generate Mocks

```bash
# Install mockgen
go install github.com/golang/mock/mockgen@latest

# Generate mock
mockgen -source=internal/modules/user/repository.go \
        -destination=internal/modules/user/mock_repository.go \
        -package=user
```

### Table-Driven Tests

```go
func TestUserService_Validate(t *testing.T) {
    tests := []struct {
        name    string
        input   *RegisterRequest
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid input",
            input: &RegisterRequest{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "password123",
            },
            wantErr: false,
        },
        {
            name: "empty username",
            input: &RegisterRequest{
                Username: "",
                Email:    "test@example.com",
                Password: "password123",
            },
            wantErr: true,
            errMsg:  "username is required",
        },
        {
            name: "invalid email",
            input: &RegisterRequest{
                Username: "testuser",
                Email:    "invalid",
                Password: "password123",
            },
            wantErr: true,
            errMsg:  "invalid email format",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Integration Testing

Test with real database (SQLite in-memory).

### Setup

```go
// tests/integration/setup.go
func SetupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    // Run migrations
    err = db.AutoMigrate(&user.User{}, &permission.Role{})
    require.NoError(t, err)
    
    return db
}
```

### Repository Test

```go
func TestUserRepository_Create(t *testing.T) {
    db := SetupTestDB(t)
    repo := user.NewRepository(db)
    
    u := &user.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashed",
    }
    
    err := repo.Create(context.Background(), u)
    
    assert.NoError(t, err)
    assert.NotZero(t, u.ID)
    assert.NotZero(t, u.CreatedAt)
}

func TestUserRepository_FindByEmail(t *testing.T) {
    db := SetupTestDB(t)
    repo := user.NewRepository(db)
    
    // Seed data
    db.Create(&user.User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "hashed",
    })
    
    // Test
    found, err := repo.FindByEmail(context.Background(), "test@example.com")
    
    assert.NoError(t, err)
    assert.Equal(t, "testuser", found.Username)
}
```

## Feature Testing

End-to-end API tests with real HTTP requests.

### Test Framework

```go
// internal/infra/testing/testing.go
type TestCase struct {
    t      *testing.T
    engine *gin.Engine
    req    *http.Request
    resp   *httptest.ResponseRecorder
}

func NewTestCase(t *testing.T, engine *gin.Engine) *TestCase {
    return &TestCase{t: t, engine: engine}
}

func (tc *TestCase) Post(path string) *TestCase {
    tc.req = httptest.NewRequest("POST", path, nil)
    return tc
}

func (tc *TestCase) WithJSON(data any) *TestCase {
    body, _ := json.Marshal(data)
    tc.req = httptest.NewRequest(tc.req.Method, tc.req.URL.Path, bytes.NewReader(body))
    tc.req.Header.Set("Content-Type", "application/json")
    return tc
}

func (tc *TestCase) Call() *TestCase {
    tc.resp = httptest.NewRecorder()
    tc.engine.ServeHTTP(tc.resp, tc.req)
    return tc
}

func (tc *TestCase) AssertOk() *TestCase {
    assert.Equal(tc.t, http.StatusOK, tc.resp.Code)
    return tc
}

func (tc *TestCase) AssertJSONPath(path, expected string) *TestCase {
    // Extract nested JSON path like "data.user.name"
    var result map[string]any
    json.Unmarshal(tc.resp.Body.Bytes(), &result)
    
    value := extractPath(result, path)
    assert.Equal(tc.t, expected, value)
    return tc
}
```

### Feature Test Setup

```go
// tests/feature/setup.go
func SetupApp() *gin.Engine {
    cfg := &config.Config{}
    cfg.Database.Driver = "sqlite"
    cfg.Database.Memory = true
    cfg.JWT.Secret = "test-secret"
    
    db, _ := database.NewDB(cfg)
    bootstrap.RunMigrations(db)
    
    // Manual DI for tests
    jwtService := jwt.NewService(cfg)
    userRepo := user.NewRepository(db)
    userService := user.NewService(userRepo, jwtService)
    userHandler := user.NewHandler(userService)
    
    gin.SetMode(gin.TestMode)
    r := gin.New()
    routes.Setup(r, &app.Handlers{User: userHandler})
    
    return r
}

func NewTestCase(t *testing.T) *testing.TestCase {
    return testing.NewTestCase(t, SetupApp())
}
```

### API Tests

```go
// tests/feature/auth_test.go
func TestUserRegistration(t *testing.T) {
    tc := NewTestCase(t)
    
    tc.Post("/v1/register").
        WithJSON(map[string]any{
            "username": "testuser",
            "email":    "test@example.com",
            "password": "password123",
        }).
        Call().
        AssertOk().
        AssertJSONPath("data.username", "testuser").
        AssertJSONPath("data.email", "test@example.com")
}

func TestUserLogin(t *testing.T) {
    tc := NewTestCase(t)
    
    // Register first
    tc.Post("/v1/register").
        WithJSON(map[string]any{
            "username": "loginuser",
            "email":    "login@example.com",
            "password": "password123",
        }).
        Call().
        AssertOk()
    
    // Login
    tc.Post("/v1/login").
        WithJSON(map[string]any{
            "username": "login@example.com",
            "password": "password123",
        }).
        Call().
        AssertOk().
        AssertJSONStructure([]string{"data.access_token", "data.user"})
}

func TestProtectedRoute(t *testing.T) {
    tc := NewTestCase(t)
    
    // Get token
    token := tc.LoginAs("testuser")
    
    // Access protected route
    tc.Get("/v1/profile").
        WithHeader("Authorization", "Bearer "+token).
        Call().
        AssertOk()
}
```

## Running Tests

```bash
# All tests
make test

# Unit tests only
go test ./internal/modules/...

# Integration tests
go test ./tests/integration/...

# Feature tests
go test ./tests/feature/...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

## Test Helpers

### Assert Response Structure

```go
func (tc *TestCase) AssertJSONStructure(keys []string) *TestCase {
    var result map[string]any
    json.Unmarshal(tc.resp.Body.Bytes(), &result)
    
    for _, key := range keys {
        value := extractPath(result, key)
        assert.NotNil(tc.t, value, "missing key: %s", key)
    }
    return tc
}
```

### Login Helper

```go
func (tc *TestCase) LoginAs(username string) string {
    email := username + "@example.com"
    
    // Register
    tc.Post("/v1/register").
        WithJSON(map[string]any{
            "username": username,
            "email":    email,
            "password": "password123",
        }).
        Call()
    
    // Login
    tc.Post("/v1/login").
        WithJSON(map[string]any{
            "username": email,
            "password": "password123",
        }).
        Call()
    
    var resp map[string]any
    json.Unmarshal(tc.resp.Body.Bytes(), &resp)
    
    return resp["data"].(map[string]any)["access_token"].(string)
}
```

## Best Practices

### DO ✅

- Test behavior, not implementation
- Use table-driven tests
- Isolate unit tests with mocks
- Use in-memory DB for integration tests
- Test error cases

### DON'T ❌

- Test private methods directly
- Share state between tests
- Use real external services
- Skip error path testing
- Write flaky tests

## Next Steps

- [Error Handling](./02-error-handling.md)
- [API Design](./03-api-design.md)
