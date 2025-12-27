# Apikey API

> Generated: 2025-12-25 17:34:47

## Base URL

See [API Documentation](./api.md) for environment-specific base URLs.

## Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/v1/api-keys` | Create API key | 🔒 |
| `GET` | `/v1/api-keys` | List API keys | 🔒 |
| `DELETE` | `/v1/api-keys/:id` | Delete API key | 🔒 |

---

## Details

### POST `/v1/api-keys`

**Create API key**

Generates a new API key with specified attributes including name, permissions, and expiration settings. Authentication is required to create a key.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.apikeys.store` |

#### Request Body

```json
{
  "expires_at": "2024-01-01T00:00:00Z",
  "name": "John Doe",
  "never_expire": true,
  "permissions": [
    "item1",
    "item2"
  ]
}
```

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `name` | `string` | ✅ | Required, Max: 100 |
| `permissions` | `[]string` | ❌ | Optional |
| `expires_at` | `time.Time` | ❌ | Optional |
| `never_expire` | `bool` | ❌ | Optional |

#### Response

```json
{
  "created_at": "2024-01-01T00:00:00Z",
  "expires_at": "2024-01-01T00:00:00Z",
  "id": 1,
  "key": "string",
  "last_used_at": "2024-01-01T00:00:00Z",
  "name": "John Doe",
  "permissions": [
    "item1",
    "item2"
  ],
  "prefix": "string",
  "user_id": 1
}
```

#### Example

```bash
curl -X POST 'http://localhost:6066/api/v1/v1/api-keys' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"expires_at": "2024-01-01T00:00:00Z","name": "John Doe","never_expire": true,"permissions": ["item1","item2"]}'
```

---

### GET `/v1/api-keys`

**List API keys**

Returns a list of all API keys associated with the authenticated user. This endpoint requires valid authentication for access.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.apikeys.index` |

#### Response

```json
{
  "data": [],
  "page": 1,
  "per_page": 1,
  "total": 1
}
```

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/api-keys' \
  -H 'Authorization: Bearer <token>'
```

---

### DELETE `/v1/api-keys/:id`

**Delete API key**

Removes an API key identified by its unique ID. Authentication is required, and the key must belong to the authenticated user.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.apikeys.destroy` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Response

```json
{
  "created_at": "2024-01-01T00:00:00Z",
  "expires_at": "2024-01-01T00:00:00Z",
  "id": 1,
  "key": "string",
  "last_used_at": "2024-01-01T00:00:00Z",
  "name": "John Doe",
  "permissions": [
    "item1",
    "item2"
  ],
  "prefix": "string",
  "user_id": 1
}
```

#### Example

```bash
curl -X DELETE 'http://localhost:6066/api/v1/v1/api-keys/1' \
  -H 'Authorization: Bearer <token>'
```

---

