# Team API

> Generated: 2025-12-25 17:34:47

## Base URL

See [API Documentation](./api.md) for environment-specific base URLs.

## Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/v1/teams` | Create a new team | 🔒 |
| `GET` | `/v1/teams/:id` | Get team details | 🔒 |
| `PUT` | `/v1/teams/:id` | Update a team | 🔒 |
| `DELETE` | `/v1/teams/:id` | Delete a team | 🔒 |
| `GET` | `/v1/teams/:id/hierarchy` | Get team hierarchy | 🔒 |
| `GET` | `/v1/org-teams/:organization_id` | List teams by organization | 🔒 |

---

## Details

### POST `/v1/teams`

**Create a new team**

Creates a new team with the provided details. The name and organization_id are required fields.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.teams.store` |

#### Request Body

```json
{
  "description": "string",
  "display_name": "John Doe",
  "name": "John Doe",
  "organization_id": 1,
  "parent_team_id": null
}
```

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `name` | `string` | ✅ | Required, Min: 2, Max: 100 |
| `display_name` | `string` | ❌ | Max: 100 |
| `description` | `string` | ❌ | Max: 500 |
| `organization_id` | `uint` | ✅ | Required |
| `parent_team_id` | `*uint` | ❌ | - |

#### Response

```json
{
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "description": "string",
  "display_name": "John Doe",
  "id": 1,
  "name": "John Doe",
  "organization_id": 1,
  "parent_team_id": null,
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Example

```bash
curl -X POST 'http://localhost:6066/api/v1/v1/teams' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"description": "string","display_name": "John Doe","name": "John Doe","organization_id": 1,"parent_team_id": null}'
```

---

### GET `/v1/teams/:id`

**Get team details**

Retrieves detailed information about a specific team using its ID. Authentication is required.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.teams.show` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Response

```json
{
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "description": "string",
  "display_name": "John Doe",
  "id": 1,
  "name": "John Doe",
  "organization_id": 1,
  "parent_team_id": null,
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/teams/1' \
  -H 'Authorization: Bearer <token>'
```

---

### PUT `/v1/teams/:id`

**Update a team**

Updates one or more attributes of an existing team. Only specified fields will be modified.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.teams.update` |

#### Request Body

```json
{
  "description": "string",
  "display_name": "John Doe",
  "name": "John Doe",
  "parent_team_id": null,
  "status": null
}
```

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `name` | `string` | ❌ | Min: 2, Max: 100 |
| `display_name` | `string` | ❌ | Max: 100 |
| `description` | `string` | ❌ | Max: 500 |
| `parent_team_id` | `*uint` | ❌ | - |
| `status` | `*int` | ❌ | - |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Response

```json
{
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "description": "string",
  "display_name": "John Doe",
  "id": 1,
  "name": "John Doe",
  "organization_id": 1,
  "parent_team_id": null,
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Example

```bash
curl -X PUT 'http://localhost:6066/api/v1/v1/teams/1' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"description": "string","display_name": "John Doe","name": "John Doe","parent_team_id": null,"status": null}'
```

---

### DELETE `/v1/teams/:id`

**Delete a team**

Removes a team by its ID. Requires authentication and appropriate access rights.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.teams.destroy` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Response

```json
{
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "description": "string",
  "display_name": "John Doe",
  "id": 1,
  "name": "John Doe",
  "organization_id": 1,
  "parent_team_id": null,
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Example

```bash
curl -X DELETE 'http://localhost:6066/api/v1/v1/teams/1' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/teams/:id/hierarchy`

**Get team hierarchy**

Retrieves the hierarchical structure of a specific team by its ID. Requires authentication.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.teams.hierarchy` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Response

```json
{
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "description": "string",
  "display_name": "John Doe",
  "id": 1,
  "name": "John Doe",
  "organization_id": 1,
  "parent_team_id": null,
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/teams/1/hierarchy' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/org-teams/:organization_id`

**List teams by organization**

Returns a list of teams associated with the specified organization ID. Authentication is required.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.teams.by_org` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `organization_id` | `integer` | Resource identifier |

#### Response

```json
{
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "description": "string",
  "display_name": "John Doe",
  "id": 1,
  "name": "John Doe",
  "organization_id": 1,
  "parent_team_id": null,
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/org-teams/1' \
  -H 'Authorization: Bearer <token>'
```

---

