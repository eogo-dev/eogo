# Eogo Framework Guide

> EOGO (Evolving Orchestration for Go) - A modern Go web framework designed for high-velocity development.

## 📚 Documentation

### Getting Started

- [Quick Start](./01-getting-started.md) - Get started in 5 minutes
- [Project Structure](./02-project-structure.md) - Directory structure explained
- [Core Concepts](./03-core-concepts.md) - Architecture and design patterns

### Architecture

- [Layered Architecture](./architecture/01-layered-architecture.md) - Overall architecture design
- [Domain-Driven Design](./architecture/02-domain-driven-design.md) - DDD practices
- [Dependency Injection](./architecture/03-dependency-injection.md) - Wire usage guide

### Core Modules

- [Domain Layer](./modules/01-domain-layer.md) - Domain layer details
- [Infrastructure Layer](./modules/02-infrastructure-layer.md) - Infrastructure layer
- [Business Modules](./modules/03-business-modules.md) - Business module development

### Infrastructure

- [Database](./infrastructure/01-database.md) - GORM and database operations

### Production

- [Observability](./production/01-observability.md) - Tracing, Metrics, Logging
- [Resilience](./production/02-resilience.md) - Circuit breaker, rate limiting, retry
- [Health Checks](./production/03-health-checks.md) - K8s readiness probes

### Best Practices

- [Testing Guide](./best-practices/01-testing.md) - Unit, integration, feature tests
- [Error Handling](./best-practices/02-error-handling.md) - Unified error handling

---

## 🎯 Features

| Feature | Description |
| ------- | ----------- |
| Layered Architecture | Domain → Application → Infrastructure, clear separation |
| Compile-time DI | Google Wire, zero runtime overhead |
| DDD Support | Entities, Value Objects, Aggregates, Domain Events |
| Production Ready | Circuit breaker, rate limiting, tracing, metrics, health checks |
| Code Generation | CLI tools for automatic module generation |
| Test Friendly | Interface-driven, easy to mock |

## 🚀 Quick Start

```bash
# Clone project
git clone https://github.com/eogo-dev/eogo.git
cd eogo

# Install dependencies
go mod download

# Build CLI
make build

# Create module
./eogo make:module Blog

# Run migrations
./eogo migrate

# Start server
make air
```

## 📖 Recommended Reading Order

1. **Beginners**: Quick Start → Project Structure → Core Concepts
2. **Deep Dive**: Layered Architecture → Domain-Driven Design → Dependency Injection
3. **Development**: Business Modules → Database → Testing Guide
4. **Production**: Observability → Resilience → Health Checks

---

*Documentation is continuously updated. Contributions welcome!*
