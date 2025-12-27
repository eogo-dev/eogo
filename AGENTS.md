# AGENTS.md

AI 编码助手工作指南。

## 项目概述

EOGO 是一个现代化 Go 框架，采用领域驱动设计 (DDD) + 分层架构。

## 目录结构

```text
eogo/
├── cmd/
│   ├── eogo/              # CLI 工具
│   └── server/            # HTTP 服务入口
├── internal/
│   ├── bootstrap/         # 应用启动
│   ├── domain/            # 领域实体 (核心业务)
│   ├── modules/           # 业务模块
│   │   └── user/          # 示例: 8 个文件
│   │       ├── model.go       # 数据库实体 (UserPO)
│   │       ├── dto.go         # DTO + Mapper 函数
│   │       ├── repository.go  # 数据访问层
│   │       ├── service.go     # 业务逻辑层
│   │       ├── handler.go     # HTTP 处理器
│   │       ├── routes.go      # 路由注册
│   │       ├── provider.go    # Wire DI
│   │       └── service_test.go
│   ├── infra/             # 基础设施 (33+ 组件)
│   └── wiring/            # Wire 依赖注入
├── pkg/                   # 公共库
├── routes/                # 全局路由
└── tests/                 # 测试
```

## 常用命令

```bash
make build         # 构建 CLI
make test          # 运行测试
make lint          # 代码检查
make wire          # 生成 DI
make air           # 热重载开发
```

## 模块结构 (8 文件标准)

| 文件 | 职责 |
|------|------|
| `model.go` | 数据库实体 `UserPO` (GORM) |
| `dto.go` | 请求/响应 DTO + `toDomain()`/`toUserPO()` 转换 |
| `repository.go` | 数据访问，返回 `domain.User` |
| `service.go` | 业务逻辑，使用 `domain.User` |
| `handler.go` | HTTP 处理器 |
| `routes.go` | 路由注册 |
| `provider.go` | Wire ProviderSet |

## Domain 层

`internal/domain/` 包含核心业务实体：

```go
// internal/domain/user.go
type User struct {
    ID        uint
    Username  string
    Email     string
    Password  string
}
```

**数据流**: `Handler(DTO) → Service(domain.User) → Repository(UserPO)`

## 统一响应

```go
import "github.com/eogo-dev/eogo/pkg/response"

response.Success(c, data)
response.BadRequest(c, "error", err)
response.NotFound(c, "not found", err)
```

## 分页

```go
import "github.com/eogo-dev/eogo/pkg/pagination"

paginator, _ := pagination.PaginateFromContext[User](c, db)
c.JSON(200, paginator)
```

## Wire 依赖注入

```go
// internal/modules/user/provider.go
var ProviderSet = wire.NewSet(
    NewRepository,
    wire.Bind(new(Repository), new(*repository)),
    NewService,
    wire.Bind(new(Service), new(*service)),
    NewHandler,
)
```

运行 `cd internal/wiring && wire` 生成代码。

## 创建新模块

```bash
./eogo make:module Blog

# 然后:
# 1. 在 routes/api.go 注册路由
# 2. 运行 wire
```

## 开发规范

1. **DTO 包含 Mapper** - 转换函数放在 `dto.go`
2. **使用 Domain 层** - 业务逻辑使用 `domain.User`
3. **私有实现** - 结构体首字母小写
4. **构造函数返回接口** - `NewService() Service`
5. **snake_case JSON** - `json:"user_id"`
6. **英文代码注释** - 代码和注释用英文

## 测试

```bash
# 单元测试
go test ./internal/modules/user/...

# 集成测试
go test ./tests/integration/...

# 所有测试
make test
```
