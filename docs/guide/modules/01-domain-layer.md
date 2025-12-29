# Domain 层详解

> `internal/domain/` 是 ZGO 的核心，包含纯业务逻辑，无框架依赖。

## 目录结构

```
internal/domain/
├── user.go           # 用户实体 + Repository 接口
├── permission.go     # 权限实体 + Repository 接口
├── value_objects.go  # 值对象（Email, Username, Money 等）
├── events.go         # 领域事件 + EventDispatcher
├── services.go       # 领域服务（认证、注册、密码）
├── aggregate.go      # 聚合根
└── errors.go         # 领域错误定义
```

## user.go - 用户实体

```go
package domain

import (
    "context"
    "time"
)

// User 是纯领域实体
// 特点：无 GORM 标签、无 JSON 标签、无框架依赖
type User struct {
    ID        uint
    Username  string
    Email     string
    Password  string  // 哈希后的密码
    Nickname  string
    Avatar    string
    Phone     string
    Bio       string
    Status    int     // 1: active, 0: disabled
    LastLogin *time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}

// IsActive 是业务方法，属于实体本身
func (u *User) IsActive() bool {
    return u.Status == 1
}

// UserRepository 定义数据访问接口
// 接口定义在 Domain 层，实现在 Infrastructure 层
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uint) error
    FindByID(ctx context.Context, id uint) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByUsername(ctx context.Context, username string) (*User, error)
    FindAll(ctx context.Context, page, pageSize int) ([]*User, int64, error)
}
```

**设计要点**：
1. **纯净**：不依赖任何框架
2. **业务方法**：`IsActive()` 是业务逻辑
3. **接口定义**：Repository 接口定义在这里，遵循依赖倒置

## value_objects.go - 值对象

值对象是不可变的、按值比较的对象。

### Email 值对象

```go
type Email struct {
    value string  // 私有字段
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail 创建时验证
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

// String 返回字符串值
func (e Email) String() string {
    return e.value
}

// Equals 按值比较
func (e Email) Equals(other Email) bool {
    return e.value == other.value
}

// Domain 提取域名（业务方法）
func (e Email) Domain() string {
    parts := strings.Split(e.value, "@")
    if len(parts) != 2 {
        return ""
    }
    return parts[1]
}
```

### Username 值对象

```go
type Username struct {
    value string
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

func NewUsername(username string) (Username, error) {
    username = strings.TrimSpace(username)
    if username == "" {
        return Username{}, errors.New("username cannot be empty")
    }
    if !usernameRegex.MatchString(username) {
        return Username{}, errors.New("username must be 3-32 characters, alphanumeric with _ or -")
    }
    return Username{value: username}, nil
}
```

### UserStatus 枚举

```go
type UserStatus int

const (
    UserStatusDisabled UserStatus = 0
    UserStatusActive   UserStatus = 1
    UserStatusPending  UserStatus = 2
    UserStatusBanned   UserStatus = 3
)

func (s UserStatus) String() string {
    switch s {
    case UserStatusDisabled:
        return "disabled"
    case UserStatusActive:
        return "active"
    case UserStatusPending:
        return "pending"
    case UserStatusBanned:
        return "banned"
    default:
        return "unknown"
    }
}

func (s UserStatus) IsActive() bool {
    return s == UserStatusActive
}

func (s UserStatus) CanLogin() bool {
    return s == UserStatusActive
}
```

### Money 值对象

```go
type Money struct {
    amount   int64   // 最小单位（分）
    currency string  // ISO 4217 货币代码
}

func NewMoney(amount int64, currency string) (Money, error) {
    currency = strings.ToUpper(strings.TrimSpace(currency))
    if len(currency) != 3 {
        return Money{}, errors.New("currency must be a 3-letter ISO code")
    }
    return Money{amount: amount, currency: currency}, nil
}

// Add 相加（必须同币种）
func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, errors.New("cannot add different currencies")
    }
    return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}
```

## events.go - 领域事件

### 事件定义

```go
// Event 接口
type Event interface {
    EventName() string
    OccurredAt() time.Time
}

// BaseEvent 基类
type BaseEvent struct {
    occurredAt time.Time
}

func (e BaseEvent) OccurredAt() time.Time {
    return e.occurredAt
}

func NewBaseEvent() BaseEvent {
    return BaseEvent{occurredAt: time.Now()}
}

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
    BaseEvent
    UserID   uint
    Username string
    Email    string
}

func (e UserRegisteredEvent) EventName() string {
    return "user.registered"
}

func NewUserRegisteredEvent(userID uint, username, email string) UserRegisteredEvent {
    return UserRegisteredEvent{
        BaseEvent: NewBaseEvent(),
        UserID:    userID,
        Username:  username,
        Email:     email,
    }
}
```

### 事件调度器

```go
type EventHandler func(ctx context.Context, event Event) error

type EventDispatcher struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

func NewEventDispatcher() *EventDispatcher {
    return &EventDispatcher{
        handlers: make(map[string][]EventHandler),
    }
}

// Subscribe 订阅事件
func (d *EventDispatcher) Subscribe(eventName string, handler EventHandler) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.handlers[eventName] = append(d.handlers[eventName], handler)
}

// Dispatch 同步发布
func (d *EventDispatcher) Dispatch(ctx context.Context, event Event) error {
    d.mu.RLock()
    handlers := d.handlers[event.EventName()]
    d.mu.RUnlock()

    for _, handler := range handlers {
        if err := handler(ctx, event); err != nil {
            return err
        }
    }
    return nil
}

// DispatchAsync 异步发布
func (d *EventDispatcher) DispatchAsync(ctx context.Context, event Event) {
    d.mu.RLock()
    handlers := d.handlers[event.EventName()]
    d.mu.RUnlock()

    for _, handler := range handlers {
        go func(h EventHandler) {
            _ = h(ctx, event)
        }(handler)
    }
}

// 全局实例
var globalDispatcher = NewEventDispatcher()

func Subscribe(eventName string, handler EventHandler) {
    globalDispatcher.Subscribe(eventName, handler)
}

func Dispatch(ctx context.Context, event Event) error {
    return globalDispatcher.Dispatch(ctx, event)
}

func DispatchAsync(ctx context.Context, event Event) {
    globalDispatcher.DispatchAsync(ctx, event)
}
```

### 使用示例

```go
// 订阅事件（通常在 bootstrap 时）
domain.Subscribe("user.registered", func(ctx context.Context, event domain.Event) error {
    e := event.(domain.UserRegisteredEvent)
    return emailService.SendWelcome(e.Email, e.Username)
})

domain.Subscribe("user.registered", func(ctx context.Context, event domain.Event) error {
    e := event.(domain.UserRegisteredEvent)
    return auditService.Log("user_registered", e.UserID)
})

// 发布事件（在 Service 中）
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
    // ... 创建用户
    
    // 异步发布事件
    domain.DispatchAsync(ctx, domain.NewUserRegisteredEvent(user.ID, user.Username, user.Email))
    
    return user, nil
}
```

## services.go - 领域服务

### 认证服务

```go
// 依赖接口定义
type PasswordHasher interface {
    Hash(password string) (string, error)
    Verify(password, hash string) bool
}

type TokenGenerator interface {
    Generate(userID uint, username string) (string, error)
    Validate(token string) (userID uint, username string, err error)
}

// AuthenticationService 认证服务
type AuthenticationService struct {
    userRepo UserRepository
    hasher   PasswordHasher
    tokens   TokenGenerator
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

type AuthResult struct {
    User        *User
    AccessToken string
}

func (s *AuthenticationService) Authenticate(ctx context.Context, identifier, password string) (*AuthResult, error) {
    // 1. 查找用户（支持用户名或邮箱）
    user, err := s.userRepo.FindByUsername(ctx, identifier)
    if err != nil {
        user, err = s.userRepo.FindByEmail(ctx, identifier)
        if err != nil {
            return nil, ErrInvalidCredentials
        }
    }

    // 2. 检查账户状态
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

### 注册服务

```go
type RegistrationService struct {
    userRepo UserRepository
    hasher   PasswordHasher
}

type RegistrationRequest struct {
    Username string
    Email    string
    Password string
    Nickname string
}

func (s *RegistrationService) Register(ctx context.Context, req *RegistrationRequest) (*User, error) {
    // 1. 验证邮箱（使用值对象）
    email, err := NewEmail(req.Email)
    if err != nil {
        return nil, err
    }

    // 2. 验证用户名（使用值对象）
    username, err := NewUsername(req.Username)
    if err != nil {
        return nil, err
    }

    // 3. 检查邮箱是否已存在
    existing, _ := s.userRepo.FindByEmail(ctx, email.String())
    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }

    // 4. 哈希密码
    hashedPassword, err := s.hasher.Hash(req.Password)
    if err != nil {
        return nil, err
    }

    // 5. 创建用户
    user := &User{
        Username: username.String(),
        Email:    email.String(),
        Password: hashedPassword,
        Nickname: req.Nickname,
        Status:   int(UserStatusActive),
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}
```

## aggregate.go - 聚合根

```go
// AggregateRoot 基类
type AggregateRoot struct {
    events []Event
    mu     sync.Mutex
}

func (a *AggregateRoot) AddEvent(event Event) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.events = append(a.events, event)
}

func (a *AggregateRoot) GetEvents() []Event {
    a.mu.Lock()
    defer a.mu.Unlock()
    events := make([]Event, len(a.events))
    copy(events, a.events)
    return events
}

func (a *AggregateRoot) ClearEvents() {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.events = nil
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

// UserAggregate 用户聚合根
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

## errors.go - 领域错误

```go
package domain

import "errors"

var (
    // 用户错误
    ErrUserNotFound       = errors.New("user not found")
    ErrEmailAlreadyExists = errors.New("email already registered")
    ErrInvalidCredentials = errors.New("invalid username or password")
    ErrAccountDisabled    = errors.New("account is disabled")

    // 权限错误
    ErrPermissionDenied = errors.New("permission denied")
    ErrRoleNotFound     = errors.New("role not found")

    // 通用错误
    ErrNotFound     = errors.New("resource not found")
    ErrConflict     = errors.New("resource already exists")
    ErrInvalidInput = errors.New("invalid input")
)
```

## 最佳实践

### DO ✅

- 在 Domain 层定义所有业务规则
- 使用值对象封装验证逻辑
- 通过事件解耦模块
- 定义接口，不实现

### DON'T ❌

- 在 Domain 层引入框架依赖
- 在 Domain 层直接操作数据库
- 在 Domain 层处理 HTTP 请求

## 下一步

- [Infrastructure 层](./02-infrastructure-layer.md)
- [业务模块开发](./03-business-modules.md)
