# EOGO 🚀
**Evolving Orchestration for Go**

A modern, high-performance Go framework designed for enterprise-grade SaaS applications.

![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go)
![Architecture](https://img.shields.io/badge/Arch-DDD-success?style=for-the-badge)

---

## ✨ Features

- **Domain-Driven Design (DDD)**: Clean domain layer + modular business logic
- **Enterprise Infrastructure**: Circuit breaker, rate limiter, tracing, config hot-reload
- **Developer First**: CLI code generation, Wire DI, comprehensive testing
- **Production Ready**: CI/CD, code quality checks, OpenAPI documentation

---

## 📂 Project Structure

```text
eogo/
├── cmd/
│   ├── eogo/              # CLI tool
│   └── server/            # HTTP server entry
├── internal/
│   ├── bootstrap/         # Application lifecycle
│   ├── domain/            # Core domain entities (DDD)
│   ├── modules/           # Business modules (user, permission, llm)
│   ├── infra/             # Infrastructure (33+ components)
│   │   ├── breaker/       # Circuit breaker
│   │   ├── ratelimit/     # Rate limiter (memory/Redis)
│   │   ├── config/        # Config management (hot-reload)
│   │   ├── tracing/       # OpenTelemetry
│   │   └── ...
│   └── wiring/            # Wire dependency injection
├── pkg/                   # Reusable public libraries
├── routes/                # Route registration
├── tests/                 # Tests (unit/integration/e2e)
├── docs/                  # Documentation
└── .github/workflows/     # CI/CD
```

---

## 🚀 Quick Start

```bash
# Clone and configure
git clone https://github.com/eogo-dev/eogo.git && cd eogo
cp .env.example .env

# Install dependencies
go mod download

# Start development server
make air
```

Visit: `http://localhost:8025`

---

## 🛠️ Common Commands

```bash
make help          # Show all commands
make build         # Build CLI
make test          # Run tests
make lint          # Code linting
make cover         # Coverage report
make wire          # Generate DI code
make docs          # Generate API docs
```

---

## 📖 Documentation

- [Development Guide](docs/guide/README.md)
- [Module Development](internal/modules/README.md)
- [Dependency Injection (Wire)](docs/dependency_injection.md)
- [AI Collaboration Guide](AGENTS.md)
- [API Documentation](docs/api/)

---

## 📜 License
MIT © 2025 Eogo Team
