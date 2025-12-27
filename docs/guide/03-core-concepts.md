# 核心概念

> 理解 Eogo 的架构理念和设计模式

## 设计哲学

Eogo 的设计遵循以下原则：

1. **简单性优先** - Go 的简洁哲学，避免过度设计
2. **关注点分离** - 每层只做一件事
3. **依赖倒置** - 高层不依赖低层，都依赖抽象
4. **可测试性** - 接口驱动，易于 Mock

## 架构模式

### 分层架构 (Layered Architecture)

```
┌─────────────────────────────────────────┐
│           Presentation Layer            │  ← HTTP Handlers
├─────────────────────────────────────────┤
│           Application Layer             │  ← Services (Use Cases)
├─────────────────────────────────────────┤
│             Domain Layer                │  ← Entities, Value Objects
├─────────────────────────────────────────┤
│          Infrastructure Layer           │  ← Database, Cache, External
└─────────────────────────────────────────┘
```

**各层职责：**

| 层 | 职责 | Eogo 对应 |
|---|------|----------|
| Presentation | 处理 HTTP 请求/响应 | `modules/*/handler.go` |
| Application | 编排业务流程 | `modules/*/service.go` |
| Domain | 核心业务逻辑 | `internal/domain/` |
| Infrastructure | 技术实现 | `internal/infra/` |

### 领域驱动设计 (DDD)

Eogo 采用 DDD 战术模式：

```
Domain Layer
├── Entity (实体)           - 有唯一标识的对象
├── Value Object (值对象)   - 无标识，按值比较
├── Aggregate Root (聚合根) - 事务边界
├── Domain Service (领域服务) - 跨实体逻辑
├── Domain Event (领域事件)  - 解耦通知
└── Repository (仓储接口)    - 数据访问抽象
```

### 依赖注入 (Dependency Injection)

使用 Google Wire 实现编译时依赖注入：

```go
// 定义 Provider
var ProviderSet = wire.NewSet(
    NewRepository,
    wire.Bind(new(Repository), new(*UserRepositoryImpl)),
    NewService,
    NewHandler,
)

// Wire 自动生成注入代码
func InitApp(db *gorm.DB) (*App, error) {
    wire.Build(
        user.ProviderSet,
        permission.ProviderSet,
        wire.Struct(new(App), "*"),
    )
    return nil, nil
}
```

**优势：**
- 编译时检查依赖
- 零运行时开销
- 显式依赖关系

## 设计模式

### 1. Repository 模式

**问题**：业务逻辑不应该知道数据如何存储

**解决**：定义接口，隐藏实现细节

```go
// domain/user.go - 接口定义
type UserRepository interface {
    FindByID(ctx context.Context, id uint) (*User, error)
    Create(ctx context.Context, user *User) error
}

// modules/user/repository.go - 实现
type UserRepositoryImpl struct {
    db *gorm.DB
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.User, error) {
    var model User
    if err := r.db.First(&model, id).Error; err != nil {
        return nil, err
    }
    return model.ToDomain(), nil
}
```

### 2. 值对象模式

**问题**：原始类型无法表达业务约束

**解决**：创建自验证的值对象

```go
// 不好 - 到处验证
func Register(email string) {
    if !isValidEmail(email) { ... }
}

// 好 - 值对象自验证
type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    if !emailRegex.MatchString(email) {
        return Email{}, errors.New("invalid email")
    }
    return Email{value: email}, nil
}

func Register(email Email) {
    // email 一定是有效的
}
```

### 3. 领域事件模式

**问题**：注册后要发邮件、记日志、通知其他系统

**解决**：发布事件，订阅者各自处理

```go
// 发布事件
domain.Dispatch(ctx, domain.NewUserRegisteredEvent(user.ID, user.Username, user.Email))

// 订阅事件
domain.Subscribe("user.registered", func(ctx context.Context, event domain.Event) error {
    e := event.(domain.UserRegisteredEvent)
    return emailService.SendWelcome(e.Email)
})
```

### 4. 聚合根模式

**问题**：多个相关对象的修改需要保持一致性

**解决**：通过聚合根统一管理

```go
type UserAggregate struct {
    AggregateRoot
    User *User
}

func (a *UserAggregate) ChangePassword(newHash string) {
    a.User.Password = newHash
    a.AddEvent(NewUserPasswordChangedEvent(a.User.ID))
}

// 使用
aggregate := NewUserAggregate(user)
aggregate.ChangePassword(hashedPassword)
repo.Update(ctx, aggregate.User)
aggregate.DispatchEventsAsync(ctx)  // 保存成功后发布事件
```

### 5. 规约模式 (Specification)

**问题**：复杂查询条件难以复用

**解决**：将条件封装为可组合的规约

```go
// 定义规约
type ActiveUserSpec struct{}
func (s *ActiveUserSpec) IsSatisfiedBy(user *User) bool {
    return user.IsActive()
}

type EmailDomainSpec struct {
    domain string
}
func (s *EmailDomainSpec) IsSatisfiedBy(user *User) bool {
    return strings.HasSuffix(user.Email, "@"+s.domain)
}

// 组合规约
spec := And(&ActiveUserSpec{}, NewEmailDomainSpec("company.com"))
if spec.IsSatisfiedBy(user) {
    // 活跃的公司邮箱用户
}
```

### 6. 熔断器模式

**问题**：下游服务故障导致级联失败

**解决**：快速失败，保护系统

```go
breaker := breaker.New(breaker.Config{
    Name:      "payment-service",
    Threshold: 5,           // 5 次失败后熔断
    Timeout:   10 * time.Second,  // 10 秒后尝试恢复
})

err := breaker.Do(func() error {
    return paymentService.Charge(amount)
})
if errors.Is(err, breaker.ErrServiceUnavailable) {
    // 服务熔断中，使用降级方案
}
```

### 7. Singleflight 模式

**问题**：缓存失效时大量请求同时查询数据库

**解决**：相同请求只执行一次

```go
sf := singleflight.New()

// 100 个并发请求，只有 1 个会真正查询数据库
result, err := sf.Do("user:123", func() (any, error) {
    return db.GetUser(123)
})
```

## 核心流程

### 请求处理流程

```
HTTP Request
    ↓
[Middleware] → Logger, Recovery, CORS, Auth, Tracing, Metrics
    ↓
[Handler] → 解析请求，调用 Service
    ↓
[Service] → 编排业务逻辑，调用 Repository
    ↓
[Repository] → 数据访问
    ↓
[Database]
    ↓
Response ← [Handler] ← [Service] ← [Repository]
```

### 依赖注入流程

```
1. 定义接口 (domain/)
    ↓
2. 实现接口 (modules/*/repository.go)
    ↓
3. 注册 Provider (modules/*/provider.go)
    ↓
4. Wire 生成代码 (wiring/wire_gen.go)
    ↓
5. 启动时自动注入 (bootstrap/http.go)
```

## 最佳实践

### DO ✅

- 在 Domain 层定义接口
- 使用值对象封装业务规则
- 通过事件解耦模块
- 编写单元测试 Mock 接口
- 使用 Context 传递请求上下文

### DON'T ❌

- 在 Domain 层引入框架依赖
- 在 Handler 中写业务逻辑
- 直接在 Service 中操作数据库
- 忽略错误处理
- 使用全局变量

## 下一步

- [分层架构详解](./architecture/01-layered-architecture.md)
- [领域驱动设计](./architecture/02-domain-driven-design.md)
- [依赖注入指南](./architecture/03-dependency-injection.md)
