# 分层架构

> ZGO 采用经典的分层架构，结合六边形架构的思想。

## 架构总览

```
┌─────────────────────────────────────────────────────────────────┐
│                        External World                            │
│                   (HTTP, gRPC, CLI, Message Queue)               │
└─────────────────────────────────────────────────────────────────┘
                                ↓ ↑
┌─────────────────────────────────────────────────────────────────┐
│                     Presentation Layer                           │
│                                                                  │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│   │   Handler   │  │  Middleware │  │   Router    │            │
│   └─────────────┘  └─────────────┘  └─────────────┘            │
│                                                                  │
│   职责：HTTP 请求解析、响应格式化、认证授权                        │
│   位置：internal/modules/*/handler.go, internal/infra/middleware │
└─────────────────────────────────────────────────────────────────┘
                                ↓ ↑
┌─────────────────────────────────────────────────────────────────┐
│                     Application Layer                            │
│                                                                  │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│   │   Service   │  │     DTO     │  │  Use Case   │            │
│   └─────────────┘  └─────────────┘  └─────────────┘            │
│                                                                  │
│   职责：编排业务流程、事务管理、调用领域服务                        │
│   位置：internal/modules/*/service.go                            │
└─────────────────────────────────────────────────────────────────┘
                                ↓ ↑
┌─────────────────────────────────────────────────────────────────┐
│                       Domain Layer                               │
│                                                                  │
│   ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ │
│   │ Entity  │ │  Value  │ │ Domain  │ │Aggregate│ │  Event  │ │
│   │         │ │ Object  │ │ Service │ │  Root   │ │         │ │
│   └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ │
│                                                                  │
│   职责：核心业务逻辑、业务规则、领域模型                           │
│   位置：internal/domain/                                         │
└─────────────────────────────────────────────────────────────────┘
                                ↓ ↑
┌─────────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                          │
│                                                                  │
│   ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ │
│   │Database │ │  Cache  │ │  Queue  │ │ Email   │ │ Tracing │ │
│   └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ │
│                                                                  │
│   职责：技术实现、外部服务集成、数据持久化                         │
│   位置：internal/infra/                                          │
└─────────────────────────────────────────────────────────────────┘
```

## 各层详解

### Presentation Layer (表现层)

**职责**：
- 接收 HTTP 请求
- 解析请求参数
- 调用 Application 层
- 格式化响应

**代码示例**：

```go
// internal/modules/user/handler.go
type Handler struct {
    service Service  // 依赖 Application 层
}

func (h *Handler) Register(c *gin.Context) {
    // 1. 解析请求
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid input", err)
        return
    }
    
    // 2. 调用 Service
    user, err := h.service.Register(c.Request.Context(), &req)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    // 3. 返回响应
    response.Created(c, user.ToResponse())
}
```

**规则**：
- ✅ 只处理 HTTP 相关逻辑
- ✅ 使用 DTO 传输数据
- ❌ 不包含业务逻辑
- ❌ 不直接访问数据库

### Application Layer (应用层)

**职责**：
- 编排业务流程
- 管理事务
- 调用领域服务
- 发布领域事件

**代码示例**：

```go
// internal/modules/user/service.go
type Service struct {
    repo       Repository      // 数据访问
    hasher     hash.Hasher     // 密码哈希
    jwt        *jwt.Service    // JWT 服务
    email      *email.Service  // 邮件服务
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
    // 1. 验证业务规则
    existing, _ := s.repo.FindByEmail(ctx, req.Email)
    if existing != nil {
        return nil, domain.ErrEmailAlreadyExists
    }
    
    // 2. 创建领域对象
    hashedPassword, _ := s.hasher.Hash(req.Password)
    user := &User{
        Username: req.Username,
        Email:    req.Email,
        Password: hashedPassword,
        Status:   1,
    }
    
    // 3. 持久化
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // 4. 发布事件
    domain.DispatchAsync(ctx, domain.NewUserRegisteredEvent(user.ID, user.Username, user.Email))
    
    // 5. 发送欢迎邮件
    go s.email.SendWelcome(user.Email, user.Username)
    
    return user, nil
}
```

**规则**：
- ✅ 编排多个领域服务
- ✅ 管理事务边界
- ✅ 发布领域事件
- ❌ 不包含核心业务规则
- ❌ 不直接操作数据库

### Domain Layer (领域层)

**职责**：
- 定义核心业务实体
- 封装业务规则
- 定义领域服务
- 定义仓储接口

**代码示例**：

```go
// internal/domain/user.go
type User struct {
    ID        uint
    Username  string
    Email     string
    Password  string
    Status    int
}

func (u *User) IsActive() bool {
    return u.Status == 1
}

// 仓储接口（不是实现）
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uint) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
}
```

```go
// internal/domain/services.go
type AuthenticationService struct {
    userRepo UserRepository
    hasher   PasswordHasher
    tokens   TokenGenerator
}

func (s *AuthenticationService) Authenticate(ctx context.Context, identifier, password string) (*AuthResult, error) {
    user, err := s.userRepo.FindByUsername(ctx, identifier)
    if err != nil {
        return nil, ErrInvalidCredentials
    }
    
    if !user.IsActive() {
        return nil, ErrAccountDisabled
    }
    
    if !s.hasher.Verify(password, user.Password) {
        return nil, ErrInvalidCredentials
    }
    
    token, _ := s.tokens.Generate(user.ID, user.Username)
    return &AuthResult{User: user, AccessToken: token}, nil
}
```

**规则**：
- ✅ 纯业务逻辑
- ✅ 无框架依赖
- ✅ 定义接口，不实现
- ❌ 不依赖任何外部层

### Infrastructure Layer (基础设施层)

**职责**：
- 实现仓储接口
- 数据库访问
- 缓存操作
- 外部服务集成

**代码示例**：

```go
// internal/modules/user/repository.go
type UserRepositoryImpl struct {
    db *gorm.DB
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
    model := &User{
        Username: user.Username,
        Email:    user.Email,
        Password: user.Password,
        Status:   user.Status,
    }
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    var model User
    if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
        return nil, err
    }
    return model.ToDomain(), nil
}
```

**规则**：
- ✅ 实现 Domain 层定义的接口
- ✅ 处理技术细节
- ✅ 可替换实现
- ❌ 不包含业务逻辑

## 依赖方向

```
Presentation → Application → Domain ← Infrastructure
                               ↑
                          (接口定义)
```

**关键点**：
1. 依赖方向向内（指向 Domain）
2. Domain 层不依赖任何层
3. Infrastructure 实现 Domain 定义的接口
4. 通过接口实现依赖倒置

## 数据流

### 请求流程

```
HTTP Request
    ↓
Handler.Register(c *gin.Context)
    ↓ RegisterRequest (DTO)
Service.Register(ctx, req)
    ↓ domain.User
Repository.Create(ctx, user)
    ↓ SQL
Database
```

### 响应流程

```
Database
    ↓ Row
Repository → domain.User
    ↓
Service → domain.User
    ↓
Handler → UserResponse (DTO)
    ↓
HTTP Response (JSON)
```

## 模型转换

```go
// DTO → Domain
func (req *RegisterRequest) ToDomain() *domain.User {
    return &domain.User{
        Username: req.Username,
        Email:    req.Email,
    }
}

// Domain → Response
func (u *domain.User) ToResponse() *UserResponse {
    return &UserResponse{
        ID:       u.ID,
        Username: u.Username,
        Email:    u.Email,
    }
}

// GORM Model → Domain
func (m *User) ToDomain() *domain.User {
    return &domain.User{
        ID:       m.ID,
        Username: m.Username,
        Email:    m.Email,
        Password: m.Password,
        Status:   m.Status,
    }
}
```

## 优势

| 优势 | 说明 |
|------|------|
| **可测试性** | 每层可独立测试，Mock 接口 |
| **可维护性** | 职责清晰，修改影响范围小 |
| **可扩展性** | 新增功能不影响现有代码 |
| **技术无关** | Domain 层不依赖框架 |
| **团队协作** | 不同层可并行开发 |

## 下一步

- [领域驱动设计](./02-domain-driven-design.md) - 深入 Domain 层
- [依赖注入](./03-dependency-injection.md) - Wire 使用指南
