# Project Usage & Configuration Guide

This document provides a step-by-step guide for configuring and using the Eogo project.

## 1. Overview

Eogo is a modern Go scaffold for AI-powered development, built with the "Vibe Enterprise" philosophy. It supports JWT authentication, multi-tenancy (Organizations/Teams), API Keys, and modular domain architecture out of the box.

## 2. Configuration Philosophy

- **Environment variables** are the primary way to configure sensitive and environment-specific values.
- **YAML config files** are used for non-sensitive, structured configuration.
- **Never commit real secrets or credentials** to the repository. Use placeholders in `.env.example`.

## 3. Configuration Workflow

### Step 1: Clone the Repository

```bash
git clone https://github.com/zgiai/eogo.git
cd eogo
```

### Step 2: Install Dependencies

```bash
go mod download
```

### Step 3: Prepare Environment Variables

- Copy the example file:
  ```bash
  cp .env.example .env
  ```
- Fill in your actual values for all required variables in `.env` (database, JWT secret, etc.).

### Step 4: Database Migration

```bash
make migrate
```

### Step 5: Start the Application

```bash
# Development with hot reload (recommended)
make air
```

The server will start on port **8025** by default.

## 4. Best Practices

- Always keep `.env.example` up to date with all required variables.
- Never commit `.env` with secrets to version control.
- Use strong, unique secrets for JWT and database.

## 5. Troubleshooting

- If the application fails to start, check for missing or incorrect environment variables.
- Ensure your database and Redis services are running and accessible.
- Review logs for detailed error messages.

---

For more details, see the main [README.md](../README.md) or open an issue on GitHub.
