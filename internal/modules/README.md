# Modules

业务领域模块目录，遵循领域驱动设计 (DDD) 模式。

## 模块概览

| 模块 | 描述 | 类型 |
|------|------|------|
| `user` | 用户认证 (注册/登录/JWT) | 简单模块 |
| `permission` | RBAC (角色/权限) | 简单模块 |
| `llm` | LLM 供应商/模型/渠道管理 | 复合模块 |
| `finance` | 财务管理 (套餐/交易/充值) | 复合模块 |

## 标准模块结构 (8 文件)

```text
module_name/
├── model.go        # 数据库实体 (GORM)
├── dto.go          # DTO + Mapper 函数
├── repository.go   # 数据访问层 (接口+实现)
├── service.go      # 业务逻辑层 (接口+实现)
├── handler.go      # HTTP 处理器
├── routes.go       # 路由注册
├── provider.go     # Wire 依赖注入
└── service_test.go # 单元测试 (可选)
```

### 各文件职责

| 文件 | 职责 | 依赖 |
|------|------|------|
| `model.go` | 定义 `UserPO` 等数据库持久化对象 | GORM |
| `dto.go` | 请求/响应结构 + `toDomain()`/`toUserPO()` 转换 | domain |
| `repository.go` | 数据库 CRUD，返回 `domain.User` | domain, GORM |
| `service.go` | 业务逻辑，使用 `domain.User` | domain, repository |
| `handler.go` | HTTP 绑定，DTO ↔ Service 调用 | service, dto |
| `routes.go` | `Register(router)` 路由注册 | handler |
| `provider.go` | Wire `ProviderSet` 定义 | wire |

## Domain 层

`internal/domain/` 包含核心业务实体，**被所有模块共享**：

```go
// internal/domain/user.go
type User struct {
    ID        uint
    Username  string
    Email     string
    Password  string  // 内部使用，不暴露给 DTO
    // ...
}
```

### 数据流向

```
HTTP Request → [handler] → DTO
                   ↓
               [service] → domain.User (业务逻辑)
                   ↓
              [repository] → UserPO (数据库)
                   ↓
              [mapper] → domain.User ← 返回
```

## 复合模块结构

对于复杂领域，使用子模块：

```text
llm/
├── provider/       # 供应商子模块
│   ├── dto.go
│   ├── repository.go
│   ├── service.go
│   └── handler.go
├── channel/        # 渠道子模块
├── model.go        # 共享实体
└── router.go       # 统一路由注册
```

## 创建新模块

```bash
# 使用 CLI 生成
./eogo make:module Blog

# 生成后需要:
# 1. 在 routes/api.go 注册路由
# 2. 运行 cd internal/wiring && wire
```

## 命名规范

| 类型 | 命名模式 | 示例 |
|------|----------|------|
| 实体 (PO) | `{Entity}PO` | `UserPO` |
| 领域实体 | `domain.{Entity}` | `domain.User` |
| 请求 DTO | `{Action}{Entity}Request` | `CreateUserRequest` |
| 响应 DTO | `{Entity}Response` | `UserResponse` |
| 接口 | `{Entity}{Layer}` | `UserRepository`, `UserService` |

## 最佳实践

1. **DTO 包含 Mapper** - 转换函数放在 `dto.go`，不单独拆分
2. **接口定义在同文件** - Repository/Service 接口与实现在同一文件
3. **使用 Domain 层** - 业务逻辑使用 `domain.User`，不直接暴露 `UserPO`
4. **私有实现** - 实现类型首字母小写 (unexported)
5. **构造函数返回接口** - `NewService() Service`
