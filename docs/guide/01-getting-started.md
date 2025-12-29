# 快速开始

> 5 分钟上手 ZGO 框架

## 环境要求

- Go 1.21+
- MySQL 8.0+ 或 PostgreSQL 14+（可选，支持 SQLite）
- Redis 6+（可选）

## 安装

```bash
# 克隆项目
git clone https://github.com/zgiai/zgo.git
cd zgo

# 安装依赖
go mod download

# 构建 CLI 工具
make build
```

## 配置

复制环境配置文件：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```env
# 应用配置
APP_NAME=zgo
APP_ENV=development
APP_DEBUG=true

# 服务器配置
SERVER_PORT=8080
SERVER_MODE=debug

# 数据库配置
DB_DRIVER=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=zgo
DB_USERNAME=root
DB_PASSWORD=

# JWT 配置
JWT_SECRET=your-secret-key
JWT_EXPIRE=24h
```

## 运行

```bash
# 运行数据库迁移
./zgo migrate

# 启动开发服务器（热重载）
make air

# 或直接运行
go run cmd/server/main.go
```

访问 http://localhost:8080/health 验证服务是否正常。

## 创建第一个模块

使用 CLI 创建一个 Blog 模块：

```bash
./zgo make:module Blog
```

这会生成以下文件：

```
internal/modules/blog/
├── model.go       # 数据模型
├── dto.go         # 请求/响应 DTO
├── repository.go  # 数据访问层
├── service.go     # 业务逻辑层
├── handler.go     # HTTP 处理器
├── routes.go      # 路由注册
└── provider.go    # 依赖注入配置
```

## 注册路由

编辑 `routes/api.go`，添加模块路由：

```go
import "github.com/zgiai/zgo/internal/modules/blog"

func Setup(r *gin.Engine, handlers *app.Handlers) {
    // ... 其他路由
    
    // 注册 Blog 模块路由
    blog.RegisterRoutes(r, handlers.Blog)
}
```

## 重新生成 Wire

```bash
cd internal/wiring && wire
```

## 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 注册用户
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123"}'

# 登录
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"password123"}'
```

## 下一步

- [项目结构](./02-project-structure.md) - 了解目录结构
- [核心概念](./03-core-concepts.md) - 理解架构设计
- [业务模块开发](./modules/03-business-modules.md) - 开发自己的模块
