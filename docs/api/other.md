# Other API

> Generated: 2025-12-26 00:05:08

## Base URL

See [API Documentation](./api.md) for environment-specific base URLs.

## Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/health` | /health | 🔓 |
| `GET` | `/protected` | /protected | 🔓 |

---

## Details

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

