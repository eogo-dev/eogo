# 项目结构

> Eogo 采用清晰的分层目录结构，遵循 Go 社区最佳实践。

## 目录总览

```
eogo/
├── cmd/                    # 应用入口
│   ├── eogo/              # CLI 工具
│   ├── server/            # HTTP 服务器
│   └── migrate/           # 数据库迁移工具
│
├── internal/              # 私有代码（核心）
│   ├── app/              # 应用容器
│   ├── bootstrap/        # 启动引导
│   ├── domain/           # 领域层（DDD）
│   ├── infra/            # 基础设施层
│   ├── modules/          # 业务模块
│   └── wiring/           # Wire 依赖注入
│
├── pkg/                   # 公共库（可被外部引用）
│   ├── errors/           # 错误处理
│   ├── logger/           # 日志
│   ├── pagination/       # 分页
│   ├── response/         # 响应格式
│   └── validation/       # 验证
│
├── routes/               # 路由定义
├── database/             # 数据库相关
│   ├── migrations/       # 迁移文件
│   └── seeders/          # 数据填充
│
├── tests/                # 测试文件
│   ├── unit/            # 单元测试
│   ├── integration/     # 集成测试
│   └── feature/         # 功能测试
│
├── storage/              # 存储目录
│   └── logs/            # 日志文件
│
└── docs/                 # 文档
```

## 核心目录详解

### `cmd/` - 应用入口

每个子目录是一个可执行程序：

```
cmd/
├── eogo/           # 框架 CLI 工具
│   └── main.go    # ./eogo make:module, ./eogo migrate
├── server/         # HTTP 服务器
│   └── main.go    # 主服务入口
└── migrate/        # 迁移工具
    └── main.go    # 独立迁移命令
```

### `internal/` - 私有代码

Go 的 `internal` 目录有特殊含义：只能被同一模块内的代码导入。

#### `internal/app/` - 应用容器

```go
// app.go - 应用依赖容器
type Application struct {
    Config       *config.Config
    DB           *gorm.DB
    JWTService   *jwt.Service
    EmailService *email.Service
    Handlers     *Handlers
}

type Handlers struct {
    User       *user.Handler
    Permission *permission.Handler
}
```

#### `internal/bootstrap/` - 启动引导

```
bootstrap/
├── app.go      # 应用初始化
└── http.go     # HTTP 服务器启动、中间件配置、优雅关闭
```

#### `internal/domain/` - 领域层 ⭐

**这是框架的核心，包含纯业务逻辑：**

```
domain/
├── user.go           # 用户实体 + Repository 接口
├── permission.go     # 权限实体 + Repository 接口
├── value_objects.go  # 值对象（Email, Username, Money）
├── events.go         # 领域事件 + EventDispatcher
├── services.go       # 领域服务（认证、注册、密码）
├── aggregate.go      # 聚合根
└── errors.go         # 领域错误
```

**特点：**
- 无框架依赖（无 GORM、无 Gin）
- 纯 Go 代码，可独立测试
- 定义接口，不实现

#### `internal/infra/` - 基础设施层

```
infra/
├── config/        # 配置加载
├── database/      # 数据库连接
├── cache/         # 缓存
├── redis/         # Redis 客户端
├── jwt/           # JWT 服务
├── email/         # 邮件服务
├── queue/         # 消息队列
├── middleware/    # HTTP 中间件
├── router/        # 路由封装
├── tracing/       # OpenTelemetry 追踪
├── metrics/       # Prometheus 指标
├── health/        # 健康检查
├── breaker/       # 熔断器
├── singleflight/  # 防缓存击穿
├── retry/         # 重试机制
├── ratelimit/     # 限流
└── testing/       # 测试工具
```

#### `internal/modules/` - 业务模块

每个模块是一个独立的业务领域：

```
modules/
├── user/              # 用户模块
│   ├── model.go      # GORM 模型
│   ├── dto.go        # 请求/响应 DTO
│   ├── repository.go # Repository 实现
│   ├── service.go    # 业务服务
│   ├── handler.go    # HTTP Handler
│   ├── routes.go     # 路由注册
│   └── provider.go   # Wire Provider
│
└── permission/        # 权限模块
    └── ...
```

#### `internal/wiring/` - 依赖注入

```
wiring/
├── wire.go       # Wire 配置（手写）
└── wire_gen.go   # Wire 生成（自动）
```

### `pkg/` - 公共库

可被外部项目导入的通用代码：

```
pkg/
├── errors/       # 结构化错误
├── logger/       # 日志封装
├── pagination/   # 分页工具
├── response/     # 统一响应格式
├── validation/   # 验证器
├── hash/         # 哈希工具
└── utils/        # 通用工具
```

### `routes/` - 路由定义

```
routes/
├── api.go        # API 路由注册
└── router.go     # 路由器封装
```

### `database/` - 数据库

```
database/
├── migrations/           # 迁移文件
│   ├── migrations.go    # 迁移注册
│   ├── user_migrations.go
│   └── permission_migrations.go
└── seeders/              # 数据填充
    └── seeder.go
```

### `tests/` - 测试

```
tests/
├── unit/          # 单元测试（Mock）
├── integration/   # 集成测试（真实 DB）
└── feature/       # 功能测试（HTTP + SQLite）
```

## 依赖方向

```
cmd → internal/bootstrap → internal/wiring → internal/modules
                                                    ↓
                                            internal/domain
                                                    ↑
                                            internal/infra
                                                    ↓
                                                  pkg
```

**规则：**
1. `domain` 不依赖任何层
2. `infra` 实现 `domain` 定义的接口
3. `modules` 组合 `domain` 和 `infra`
4. `pkg` 是独立的工具库

## 文件命名规范

| 文件 | 用途 |
|------|------|
| `model.go` | GORM 数据模型 |
| `dto.go` | 请求/响应数据传输对象 |
| `repository.go` | 数据访问层实现 |
| `service.go` | 业务逻辑层 |
| `handler.go` | HTTP 处理器 |
| `routes.go` | 路由注册 |
| `provider.go` | Wire 依赖提供者 |
| `*_test.go` | 测试文件 |

## 下一步

- [核心概念](./03-core-concepts.md) - 理解架构设计理念
- [分层架构](./architecture/01-layered-architecture.md) - 深入了解分层设计
