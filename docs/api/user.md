# User API

> Generated: 2025-12-25 17:34:47

## Base URL

See [API Documentation](./api.md) for environment-specific base URLs.

## Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/v1/register` | Create user account | 🔓 |
| `POST` | `/v1/login` | Authenticate user | 🔓 |
| `POST` | `/v1/password/reset` | Request password reset | 🔓 |
| `GET` | `/v1/users/profile` | Get user profile | 🔒 |
| `PUT` | `/v1/users/profile` | Update user profile | 🔒 |
| `PUT` | `/v1/users/password` | Change user password | 🔒 |
| `DELETE` | `/v1/users/account` | Delete user account | 🔒 |
| `GET` | `/v1/users` | List all users | 🔒 |
| `GET` | `/v1/users/:id` | Get user by ID | 🔒 |
| `GET` | `/v1/users/:id/info` | Get user info | 🔒 |

---

## Details

### POST `/v1/register`

**Create user account**

Registers a new user with required fields: username, password, and email. Optional fields include nickname and phone. No authentication required.

| Property | Value |
|----------|-------|
| Auth | 🔓 Not required |
| Route Name | `v1.auth.register` |

#### Request Body

```json
{
  "email": "user@example.com",
  "nickname": "John Doe",
  "password": "********",
  "phone": "+1234567890",
  "username": "John Doe"
}
```

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `username` | `string` | ✅ | Required, Min: 3, Max: 50 |
| `password` | `string` | ✅ | Required, Min: 6, Max: 50 |
| `email` | `string` | ✅ | Required, Email format |
| `nickname` | `string` | ❌ | Max: 50 |
| `phone` | `string` | ❌ | Max: 20 |

#### Response

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "email": "user@example.com",
  "id": 1,
  "last_login": "2024-01-01T00:00:00Z",
  "nickname": "John Doe",
  "phone": "+1234567890",
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z",
  "username": "John Doe"
}
```

#### Example

```bash
curl -X POST 'http://localhost:6066/api/v1/v1/register' \
  -H 'Content-Type: application/json' \
  -d '{"email": "user@example.com","nickname": "John Doe","password": "********","phone": "+1234567890","username": "John Doe"}'
```

---

### POST `/v1/login`

**Authenticate user**

Logs in a user by validating username and password, returning an authentication token upon success. No authentication required.

| Property | Value |
|----------|-------|
| Auth | 🔓 Not required |
| Route Name | `v1.auth.login` |

#### Request Body

```json
{
  "password": "********",
  "username": "John Doe"
}
```

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `username` | `string` | ✅ | Required |
| `password` | `string` | ✅ | Required |

#### Response

```json
{
  "token": "string",
  "user": null
}
```

#### Example

```bash
curl -X POST 'http://localhost:6066/api/v1/v1/login' \
  -H 'Content-Type: application/json' \
  -d '{"password": "********","username": "John Doe"}'
```

---

### POST `/v1/password/reset`

**Request password reset**

Initiates a password reset flow by sending a reset link to the registered email address. No authentication required.

| Property | Value |
|----------|-------|
| Auth | 🔓 Not required |
| Route Name | `v1.auth.password.reset` |

#### Request Body

```json
{
  "email": "user@example.com"
}
```

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `email` | `string` | ✅ | Required, Email format |

#### Response

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "email": "user@example.com",
  "id": 1,
  "last_login": "2024-01-01T00:00:00Z",
  "nickname": "John Doe",
  "phone": "+1234567890",
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z",
  "username": "John Doe"
}
```

#### Example

```bash
curl -X POST 'http://localhost:6066/api/v1/v1/password/reset' \
  -H 'Content-Type: application/json' \
  -d '{"email": "user@example.com"}'
```

---

### GET `/v1/users/profile`

**Get user profile**

Retrieves the authenticated user's profile information. Requires valid authentication.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.users.profile` |

#### Response

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "email": "user@example.com",
  "id": 1,
  "last_login": "2024-01-01T00:00:00Z",
  "nickname": "John Doe",
  "phone": "+1234567890",
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z",
  "username": "John Doe"
}
```

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/users/profile' \
  -H 'Authorization: Bearer <token>'
```

---

### PUT `/v1/users/profile`

**Update user profile**

Updates the authenticated user's profile with provided fields such as nickname, avatar, phone, and bio. All fields are optional.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.users.profile.update` |

#### Request Body

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "string",
  "nickname": "John Doe",
  "phone": "+1234567890"
}
```

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `nickname` | `string` | ❌ | Max: 50 |
| `avatar` | `string` | ❌ | Max: 255 |
| `phone` | `string` | ❌ | Max: 20 |
| `bio` | `string` | ❌ | Max: 500 |

#### Response

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "email": "user@example.com",
  "id": 1,
  "last_login": "2024-01-01T00:00:00Z",
  "nickname": "John Doe",
  "phone": "+1234567890",
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z",
  "username": "John Doe"
}
```

#### Example

```bash
curl -X PUT 'http://localhost:6066/api/v1/v1/users/profile' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"avatar": "https://example.com/avatar.jpg","bio": "string","nickname": "John Doe","phone": "+1234567890"}'
```

---

### PUT `/v1/users/password`

**Change user password**

Changes the authenticated user's password using the provided old and new passwords. Both fields are required.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.users.password.update` |

#### Request Body

```json
{
  "new_password": "********",
  "old_password": "********"
}
```

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `old_password` | `string` | ✅ | Required |
| `new_password` | `string` | ✅ | Required, Min: 6, Max: 50 |

#### Response

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "email": "user@example.com",
  "id": 1,
  "last_login": "2024-01-01T00:00:00Z",
  "nickname": "John Doe",
  "phone": "+1234567890",
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z",
  "username": "John Doe"
}
```

#### Example

```bash
curl -X PUT 'http://localhost:6066/api/v1/v1/users/password' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"new_password": "********","old_password": "********"}'
```

---

### DELETE `/v1/users/account`

**Delete user account**

Permanently deletes the authenticated user's account. This action is irreversible and requires valid authentication.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.users.account.delete` |

#### Response

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "email": "user@example.com",
  "id": 1,
  "last_login": "2024-01-01T00:00:00Z",
  "nickname": "John Doe",
  "phone": "+1234567890",
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z",
  "username": "John Doe"
}
```

#### Example

```bash
curl -X DELETE 'http://localhost:6066/api/v1/v1/users/account' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/users`

**List all users**

Retrieves a list of all registered users. Accessible only to authenticated users with appropriate permissions.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.users.index` |

#### Response

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "email": "user@example.com",
  "id": 1,
  "last_login": "2024-01-01T00:00:00Z",
  "nickname": "John Doe",
  "phone": "+1234567890",
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z",
  "username": "John Doe"
}
```

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/users' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/users/:id`

**Get user by ID**

Retrieves the details of a specific user using their unique identifier. Authentication is required to access this endpoint.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.users.show` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Response

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "email": "user@example.com",
  "id": 1,
  "last_login": "2024-01-01T00:00:00Z",
  "nickname": "John Doe",
  "phone": "+1234567890",
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z",
  "username": "John Doe"
}
```

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/users/1' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/users/:id/info`

**Get user info**

Fetches extended information about a user identified by their user ID. This endpoint requires authentication and provides supplementary user data.

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `v1.users.info` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Response

```json
{
  "avatar": "https://example.com/avatar.jpg",
  "bio": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "deleted_at": "object",
  "email": "user@example.com",
  "id": 1,
  "last_login": "2024-01-01T00:00:00Z",
  "nickname": "John Doe",
  "phone": "+1234567890",
  "status": 1,
  "updated_at": "2024-01-01T00:00:00Z",
  "username": "John Doe"
}
```

#### Example

```bash
curl -X GET 'http://localhost:6066/api/v1/v1/users/1/info' \
  -H 'Authorization: Bearer <token>'
```

---

