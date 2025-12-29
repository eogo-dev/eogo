# ZGO 🚀

**Enterprise Orchestration in Go**

The Orchestrable Go Framework for the Intelligent Era.

[![Website](https://img.shields.io/badge/Website-zgo.dev-blue?style=for-the-badge)](https://zgo.dev)

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
zgo/
├── cmd/
│   ├── zgo/              # CLI tool
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
git clone https://github.com/zgiai/zgo.git && cd zgo
cp .env.example .env

# Install dependencies
go mod download

# Install zgo CLI globally (recommended)
make install
# Now you can use 'zgo' command anywhere!

# Or just build locally
make build
./zgo serve

# Start development server with hot-reload
make air
```

Visit: `http://localhost:8025`

### Windows Users 🪟

Windows users can use the provided PowerShell or batch scripts:

```powershell
# PowerShell (Recommended)
.\make.ps1 setup
.\make.ps1 install
.\make.ps1 dev

# Or Command Prompt
make.bat setup
make.bat install
make.bat dev
```

See [Windows Development Guide](docs/WINDOWS.md) for detailed setup instructions.

### Global Installation

After `make install` (or `.\make.ps1 install` on Windows), use zgo from anywhere:

```bash
zgo version               # Show version
zgo serve                 # Start server
zgo make:module Blog      # Generate new module
zgo db:migrate            # Run migrations
zgo db:migrate --env=prod # Production migrations
```

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
- [**Windows Development Guide**](docs/WINDOWS.md) 🪟
- [Module Development](internal/modules/README.md)
- [Dependency Injection (Wire)](docs/dependency_injection.md)
- [AI Collaboration Guide](AGENTS.md)
- [API Documentation](docs/api/)

---

## 📜 License
MIT © 2025 ZGO Team
