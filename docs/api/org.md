# Org API

> Generated: 2025-12-25 17:34:47

## Base URL

See [API Documentation](./api.md) for environment-specific base URLs.

## Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/v1/organizations` | Create organization | 🔒 |
| `GET` | `/v1/organizations` | List organizations | 🔒 |
| `GET` | `/v1/organizations/me` | Get my organizations | 🔒 |
| `GET` | `/v1/organizations/:id` | Get organization by ID | 🔒 |
| `PUT` | `/v1/organizations/:id` | Update organization | 🔒 |
| `DELETE` | `/v1/organizations/:id` | Delete an organization | 🔒 |

---

## Details

### POST `/v1/organizations`

**Create organization**

Creates a new organization with the provided details. Requires authentication and appropriate permissions.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.organizations.store` |

#### Example

```bash
curl -X POST 'http://localhost:6066/api/v1/v1/organizations' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/organizations`

**List organizations**

Retrieves a paginated list of all organizations. Requires authentication.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.organizations.index` |

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/organizations' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/organizations/me`

**Get my organizations**

Retrieves all organizations associated with the authenticated user. Requires authentication.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.organizations.me` |

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/organizations/me' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/organizations/:id`

**Get organization by ID**

Retrieves detailed information about a specific organization using its unique identifier. Requires authentication.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.organizations.show` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/organizations/1' \
  -H 'Authorization: Bearer <token>'
```

---

### PUT `/v1/organizations/:id`

**Update organization**

Updates the details of an existing organization identified by its ID. Requires authentication and ownership or admin rights.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.organizations.update` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Example

```bash
curl -X PUT 'http://localhost:6066/api/v1/v1/organizations/1' \
  -H 'Authorization: Bearer <token>'
```

---

### DELETE `/v1/organizations/:id`

**Delete an organization**

Removes an organization by its ID. Requires authentication and appropriate permissions.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.organizations.destroy` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Example

```bash
curl -X DELETE 'http://localhost:6066/api/v1/v1/organizations/1' \
  -H 'Authorization: Bearer <token>'
```

---

