# EOGO 🚀
**Evolving Orchestration for Go**

Eogo is a modern, high-performance Go framework designed for an elite developer experience. It provides a "Vibe Enterprise" foundation for building robust, multi-tenant SaaS applications with zero friction.

![Eogo Banner](https://img.shields.io/badge/Vibe-Enterprise-blueviolet?style=for-the-badge)
![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=for-the-badge&logo=go)
![Architecture](https://img.shields.io/badge/Arch-Modular-success?style=for-the-badge)

---

## ✨ Features

- **Modular Architecture**: Isolated domain modules (`user`, `org`, `team`) for clean scaling.
- **Enterprise Core**: Pre-built multi-tenancy, RBAC, and API Key management.
- **Developer First CLI**: High-performance generator (`eogo make:module`).
- **Modern DI**: Type-safe dependency injection via Google Wire.
- **Testing Suite**: Comprehensive support for Unit, Integration, and Feature tests.

### Monitoring & Observability
- **[Monitor Dashboard](http://localhost:8025/monitor)**: Built-in health and stats monitoring.

---

## Documentation

- [Usage & Configuration](docs/usage_and_config.md)
- [Dependency Injection (Wire)](docs/dependency_injection.md) - Provider pattern similar to NestJS
- [Production Logging](docs/production_logging.md) - ClickHouse & Sentry integration
- [API Documentation](docs/api/)

---

## 🚀 Quick Start

### 1. Requirements
- Go 1.23+
- PostgreSQL (optional, defaults to SQLite for quick start)
- Redis (optional)

### 2. Installation
```bash
# Clone the repository
git clone https://github.com/eogo-dev/eogo.git
cd eogo

# Configure environment
cp .env.example .env

# Run migrations
make migrate
```

### 3. Run Development Server
```bash
make air
```
Visit: `http://localhost:8025`

---

## 📂 Project Structure

```text
├── cmd/eogo              # Framework CLI
├── internal/
│   ├── bootstrap/        # Lifecycle & Kernels
│   ├── modules/          # Business Domains (User, Org, Team)
│   └── platform/         # Framework Core (DB, Cache, Router)
├── routes/               # Global Route Definitions
└── tests/                # Integrated Test Platform
```

---

## 🏗️ Development SOP

- **Add Feature**: `./eogo make:module NewFeature`
- **Add Migration**: Add model to `internal/bootstrap/migrate.go`.
- **Run Tests**: `make test`

---

## 📜 License
MIT © 2025 Eogo Team
