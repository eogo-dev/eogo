# API Documentation

> Generated: 2025-12-26 00:05:08

## Base URLs

| Environment | URL |
|-------------|-----|
| 🏠 Local | `http://localhost:8025/api/v1` |

## Authentication

Protected endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

## Overview

Total endpoints: **2**

## Table of Contents

- [Other](#other) (2 endpoints)

---

## Other

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/health` | /health | 🔓 |
| `GET` | `/protected` | /protected | 🔓 |

### GET `/health`

**/health**

| Property | Value |
|----------|-------|
| Auth | 🔓 Not required |

#### Example

```bash
curl -X GET 'http://localhost:8025/api/v1/health'
```

---

### GET `/protected`

**/protected**

| Property | Value |
|----------|-------|
| Auth | 🔓 Not required |

#### Example

```bash
curl -X GET 'http://localhost:8025/api/v1/protected'
```

---

