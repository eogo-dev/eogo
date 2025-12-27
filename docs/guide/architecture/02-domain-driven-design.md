# 领域驱动设计 (DDD)

> Eogo 的 Domain 层采用 DDD 战术模式，将业务逻辑集中在领域层。

## 概述

DDD (Domain-Driven Design) 是一种软件设计方法，核心思想是将业务逻辑放在领域层，而不是散落在各处。

```
internal/domain/
├── user.go           # 实体 (Entity)
├── permission.go     # 实体 + 仓储接口
├── value_objects.go  # 值对象 (Value Object)
├── events.go         # 领域事件 (Domain Event)
├── services.go       # 领域服务 (Domain Service)
├── aggregate.go      # 聚合根 (Aggregate Root)
└── errors.go         # 领域错误
```

## 实体 (Entity)

**定义**：有唯一标识的对象，即使属性相同，ID 不同就是不同的实体。

```go
// internal/domain/user.go
type User struct {
    ID        uint      // 唯一标识 - 实体的核心特征
    Username  string
    Email     string
    Password  string
    Nickname  string
    Status    int
    CreatedAt time.Time
    UpdatedAt time.Time
}

// 业务方法
func (u *User) IsActive() bool {
    return u.Status == 1
}
```

**特点**：
- 有唯一标识 (ID)
- 可变的（属性可以改变）
- 包含业务方法
- 无框架依赖（无 GORM 标签、无 JSON 标签）

**对比 GORM Model**：

```go
// ❌ GORM Model - 有框架依赖
type User struct {
    gorm.Model
    Username string `json:"username" gorm:"unique"`
}

// ✅ Domain Entity - 纯净
type User struct {
    ID       uint
    Username string
}
```

**为什么分离？**
- Domain Entity 可以在任何地方使用
- 不受 ORM 框架升级影响
- 便于单元测试

## 值对象 (Value Object)

**定义**：无标识的对象，按值比较，不可变。

```go
// internal/domain/value_objects.go

// Email 值对象
type Email struct {
    value string  // 私有字段，不可变
}

func NewEmail(email string) (Email, error) {
    email = strings.TrimSpace(strings.ToLower(email))
    if email == "" {
        return Email{}, errors.New("email cannot be empty")
    }
    if !emailRegex.MatchString(email) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: email}, nil
}

func (e Email) String() string {
    return e.value
}

func (e Email) Equals(other Email) bool {
    return e.value == other.value
}

func (e Email) Domain() string {
    parts := strings.Split(e.value, "@")
    return parts[1]
}
```

**特点**：
- 无标识
- 不可变（创建后不能修改）
- 自我验证（创建时验证）
- 按值比较

**Eogo 中的值对象**：

| 值对象 | 用途 | 验证规则 |
|--------|------|----------|
| `Email` | 邮箱地址 | 格式验证 |
| `Username` | 用户名 | 3-32字符，字母数字 |
| `Password` | 密码哈希 | 封装哈希值 |
| `UserStatus` | 用户状态 | 枚举值 |
| `Money` | 金额 | 货币一致性 |

**使用示例**：

```go
// 创建时验证
email, err := domain.NewEmail("invalid")  // 返回错误
email, err := domain.NewEmail("user@example.com")  // 成功

// 业务方法
if email.Domain() == "company.com" {
    // 公司邮箱
}

// 比较
email1.Equals(email2)  // 按值比较
```

## 领域事件 (Domain Event)

**定义**：表示领域中发生的重要事情，用于解耦模块。

```go
// internal/domain/events.go

// 事件接口
type Event interface {
    EventName() string
    OccurredAt() time.Time
}

// 用户注册事件
type UserRegisteredEvent struct {
    BaseEvent
    UserID   uint
    Username string
    Email    string
}

func (e UserRegisteredEvent) EventName() string {
    return "user.registered"
}

// 创建事件
func NewUserRegisteredEvent(userID uint, username, email string) UserRegisteredEvent {
    return UserRegisteredEvent{
        BaseEvent: NewBaseEvent(),
        UserID:    userID,
        Username:  username,
        Email:     email,
    }
}
```

**事件调度器**：

```go
// 订阅事件
domain.Subscribe("user.registered", func(ctx context.Context, event domain.Event) error {
    e := event.(domain.UserRegisteredEvent)
    
    // 发送欢迎邮件
    return emailService.SendWelcome(e.Email)
})

domain.Subscribe("user.registered", func(ctx context.Context, event domain.Event) error {
    e := event.(domain.UserRegisteredEvent)
    
    // 记录审计日志
    return auditService.Log("user_registered", e.UserID)
})

// 发布事件
domain.Dispatch(ctx, domain.NewUserRegisteredEvent(user.ID, user.Username, user.Email))

// 异步发布
domain.DispatchAsync(ctx, event)
```

**Eogo 中的领域事件**：

| 事件 | 触发时机 | 典型处理 |
|------|----------|----------|
| `UserRegisteredEvent` | 用户注册成功 | 发送欢迎邮件 |
| `UserLoggedInEvent` | 用户登录成功 | 记录登录日志 |
| `UserPasswordChangedEvent` | 密码修改成功 | 发送安全通知 |
| `UserDeletedEvent` | 用户删除 | 清理关联数据 |
| `RoleAssignedEvent` | 角色分配 | 权限缓存更新 |

## 领域服务 (Domain Service)

**定义**：跨实体的业务逻辑，不属于任何单一实体。

```go
// internal/domain/services.go

// 认证服务
type AuthenticationService struct {
    userRepo UserRepository    // 依赖接口
    hasher   PasswordHasher    // 依赖接口
    tokens   TokenGenerator    // 依赖接口
}

func NewAuthenticationService(
    userRepo UserRepository,
    hasher PasswordHasher,
    tokens TokenGenerator,
) *AuthenticationService {
    return &AuthenticationService{
        userRepo: userRepo,
        hasher:   hasher,
        tokens:   tokens,
    }
}

func (s *AuthenticationService) Authenticate(ctx context.Context, identifier, password string) (*AuthResult, error) {
    // 1. 查找用户
    user, err := s.userRepo.FindByUsername(ctx, identifier)
    if err != nil {
        user, err = s.userRepo.FindByEmail(ctx, identifier)
        if err != nil {
            return nil, ErrInvalidCredentials
        }
    }
    
    // 2. 检查状态
    if !user.IsActive() {
        return nil, ErrAccountDisabled
    }
    
    // 3. 验证密码
    if !s.hasher.Verify(password, user.Password) {
        return nil, ErrInvalidCredentials
    }
    
    // 4. 生成令牌
    token, err := s.tokens.Generate(user.ID, user.Username)
    if err != nil {
        return nil, err
    }
    
    return &AuthResult{User: user, AccessToken: token}, nil
}
```

**Eogo 中的领域服务**：

| 服务 | 职责 |
|------|------|
| `AuthenticationService` | 用户认证（登录） |
| `RegistrationService` | 用户注册 |
| `PasswordService` | 密码修改/重置 |
| `AuthorizationService` | 权限检查 |

**为什么需要领域服务？**

```go
// ❌ 不好 - 认证逻辑放在 User 实体
func (u *User) Authenticate(password string, hasher PasswordHasher) bool {
    return hasher.Verify(password, u.Password)
}
// 问题：User 需要知道 PasswordHasher

// ✅ 好 - 认证逻辑放在领域服务
func (s *AuthenticationService) Authenticate(ctx context.Context, identifier, password string) (*AuthResult, error) {
    // 服务协调多个对象
}
```

## 聚合根 (Aggregate Root)

**定义**：一组相关对象的根，保证事务一致性。

```go
// internal/domain/aggregate.go

// 聚合根基类
type AggregateRoot struct {
    events []Event
    mu     sync.Mutex
}

func (a *AggregateRoot) AddEvent(event Event) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.events = append(a.events, event)
}

func (a *AggregateRoot) DispatchEvents(ctx context.Context) error {
    events := a.GetEvents()
    for _, event := range events {
        if err := Dispatch(ctx, event); err != nil {
            return err
        }
    }
    a.ClearEvents()
    return nil
}

// 用户聚合根
type UserAggregate struct {
    AggregateRoot
    User *User
}

func NewUserAggregate(user *User) *UserAggregate {
    return &UserAggregate{User: user}
}

func (a *UserAggregate) ChangePassword(newHashedPassword string) {
    a.User.Password = newHashedPassword
    a.AddEvent(NewUserPasswordChangedEvent(a.User.ID))
}

func (a *UserAggregate) Disable() error {
    if a.User.Status == int(UserStatusDisabled) {
        return ErrAccountDisabled
    }
    a.User.Status = int(UserStatusDisabled)
    return nil
}
```

**使用示例**：

```go
func (s *UserService) ChangePassword(ctx context.Context, userID uint, oldPwd, newPwd string) error {
    user, err := s.repo.FindByID(ctx, userID)
    if err != nil {
        return err
    }
    
    // 创建聚合
    aggregate := domain.NewUserAggregate(user)
    
    // 通过聚合修改
    aggregate.ChangePassword(hashedNewPassword)
    
    // 保存
    if err := s.repo.Update(ctx, aggregate.User); err != nil {
        return err
    }
    
    // 保存成功后发布事件
    aggregate.DispatchEventsAsync(ctx)
    return nil
}
```

## 仓储接口 (Repository)

**定义**：数据访问的抽象，定义在 Domain 层，实现在 Infrastructure 层。

```go
// internal/domain/user.go - 接口定义
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uint) error
    FindByID(ctx context.Context, id uint) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByUsername(ctx context.Context, username string) (*User, error)
    FindAll(ctx context.Context, page, pageSize int) ([]*User, int64, error)
}

// internal/modules/user/repository.go - 实现
type UserRepositoryImpl struct {
    db *gorm.DB
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    var model User
    if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
        return nil, err
    }
    return model.ToDomain(), nil
}
```

## 规约模式 (Specification)

**定义**：将业务规则封装为可组合的对象。

```go
// internal/domain/services.go

type Specification[T any] interface {
    IsSatisfiedBy(item T) bool
}

// 活跃用户规约
type ActiveUserSpec struct{}

func (s *ActiveUserSpec) IsSatisfiedBy(user *User) bool {
    return user.IsActive()
}

// 邮箱域名规约
type EmailDomainSpec struct {
    domain string
}

func (s *EmailDomainSpec) IsSatisfiedBy(user *User) bool {
    email, _ := NewEmail(user.Email)
    return email.Domain() == s.domain
}

// 组合规约
activeCompanyUsers := And(&ActiveUserSpec{}, NewEmailDomainSpec("company.com"))

for _, user := range users {
    if activeCompanyUsers.IsSatisfiedBy(user) {
        // 活跃的公司用户
    }
}
```

## 领域错误

```go
// internal/domain/errors.go
var (
    ErrUserNotFound       = errors.New("user not found")
    ErrEmailAlreadyExists = errors.New("email already registered")
    ErrInvalidCredentials = errors.New("invalid username or password")
    ErrAccountDisabled    = errors.New("account is disabled")
    ErrPermissionDenied   = errors.New("permission denied")
)
```

## 最佳实践

### DO ✅

```go
// 使用值对象验证
email, err := domain.NewEmail(req.Email)
if err != nil {
    return err
}

// 通过聚合根修改
aggregate := domain.NewUserAggregate(user)
aggregate.ChangePassword(newHash)

// 发布领域事件
domain.DispatchAsync(ctx, domain.NewUserRegisteredEvent(...))
```

### DON'T ❌

```go
// 在 Domain 层引入框架
type User struct {
    gorm.Model  // ❌ 不要这样
}

// 直接修改实体
user.Password = newHash  // ❌ 应该通过聚合根

// 在 Handler 中写业务逻辑
func (h *Handler) Register(c *gin.Context) {
    if existingUser != nil {  // ❌ 业务逻辑应该在 Service
        return
    }
}
```

## 下一步

- [依赖注入](./03-dependency-injection.md) - Wire 使用指南
- [业务模块开发](../modules/03-business-modules.md) - 实战开发
