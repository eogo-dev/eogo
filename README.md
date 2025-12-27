# EOGO 🚀
**Evolving Orchestration for Go**

现代化高性能 Go 框架，专为企业级 SaaS 应用设计。

![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go)
![Architecture](https://img.shields.io/badge/Arch-DDD-success?style=for-the-badge)

---

## ✨ 特性

- **领域驱动设计 (DDD)**: 清晰的领域层 + 模块化业务
- **企业级基础设施**: 熔断器、限流器、链路追踪、配置热更新
- **开发者优先**: CLI 代码生成、Wire 依赖注入、完善测试
- **生产就绪**: CI/CD、代码质量检查、OpenAPI 文档

---

## 📂 项目结构

```text
eogo/
├── cmd/
│   ├── eogo/              # CLI 工具
│   └── server/            # HTTP 服务入口
├── internal/
│   ├── bootstrap/         # 应用启动与生命周期
│   ├── domain/            # 核心领域实体 (DDD)
│   ├── modules/           # 业务模块 (user, permission, llm)
│   ├── infra/             # 基础设施 (33+ 组件)
│   │   ├── breaker/       # 熔断器
│   │   ├── ratelimit/     # 限流器 (内存/Redis)
│   │   ├── config/        # 配置管理 (热更新)
│   │   ├── tracing/       # OpenTelemetry
│   │   └── ...
│   └── wiring/            # Wire 依赖注入
├── pkg/                   # 可复用公共库
├── routes/                # 路由注册
├── tests/                 # 测试 (unit/integration/e2e)
├── docs/                  # 文档
└── .github/workflows/     # CI/CD
```

---

## 🚀 快速开始

```bash
# 克隆并配置
git clone https://github.com/eogo-dev/eogo.git && cd eogo
cp .env.example .env

# 安装依赖
go mod download

# 启动开发服务器
make air
```

访问: `http://localhost:8025`

---

## �️ 常用命令

```bash
make help          # 查看所有命令
make build         # 构建 CLI
make test          # 运行测试
make lint          # 代码检查
make cover         # 覆盖率报告
make wire          # 生成依赖注入
make docs          # 生成 API 文档
```

---

## 📖 文档

- [开发指南](docs/guide/README.md)
- [模块开发](internal/modules/README.md)
- [依赖注入 (Wire)](docs/dependency_injection.md)
- [AI 协作指南](AGENTS.md)
- [API 文档](docs/api/)

---

## 📜 License
MIT © 2025 Eogo Team
