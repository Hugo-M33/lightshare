# API Reference

Base URL: `https://api.lightshare.app/v1`

All endpoints require authentication unless noted otherwise. Include the access token in the Authorization header:

```
Authorization: Bearer <access_token>
```

## Authentication

### POST /auth/signup

Create a new user account.

**Request:**
```json
{
    "email": "user@example.com",
    "password": "securepassword123"
}
```

**Response:** `201 Created`
```json
{
    "user": {
        "id": "uuid",
        "email": "user@example.com",
        "created_at": "2024-01-15T10:30:00Z"
    },
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 3600
}
```

### POST /auth/login

Authenticate and receive tokens.

**Request:**
```json
{
    "email": "user@example.com",
    "password": "securepassword123"
}
```

**Response:** `200 OK`
```json
{
    "user": {
        "id": "uuid",
        "email": "user@example.com"
    },
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 3600
}
```

### POST /auth/refresh

Refresh an expired access token.

**Request:**
```json
{
    "refresh_token": "eyJ..."
}
```

**Response:** `200 OK`
```json
{
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 3600
}
```

### POST /auth/logout

Revoke refresh token.

**Request:**
```json
{
    "refresh_token": "eyJ..."
}
```

**Response:** `204 No Content`

---

## User Profile

### GET /me

Get current user profile and entitlements.

**Response:** `200 OK`
```json
{
    "id": "uuid",
    "email": "user@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "entitlements": {
        "tier": "pro",
        "max_shares": 10,
        "ads_enabled": false,
        "expires_at": "2025-01-15T10:30:00Z"
    }
}
```

### GET /me/entitlements

Get just entitlements (lightweight call for checking features).

**Response:** `200 OK`
```json
{
    "tier": "pro",
    "max_shares": 10,
    "ads_enabled": false,
    "expires_at": "2025-01-15T10:30:00Z"
}
```

---

## Provider Connection

### POST /providers/connect

Initiate provider connection (OAuth or token).

**Request (OAuth):**
```json
{
    "provider": "lifx",
    "method": "oauth"
}
```

**Response:** `200 OK`
```json
{
    "auth_url": "https://cloud.lifx.com/oauth/authorize?...",
    "state": "random-state-string"
}
```

**Request (Personal Token):**
```json
{
    "provider": "lifx",
    "method": "token",
    "token": "c1a2b3..."
}
```

**Response:** `201 Created`
```json
{
    "account": {
        "id": "uuid",
        "provider": "lifx",
        "provider_account_id": "user@lifx",
        "created_at": "2024-01-15T10:30:00Z"
    }
}
```

### GET /providers/oauth/callback

OAuth callback endpoint (called by provider).

**Query Parameters:**
- `code` - Authorization code
- `state` - State for CSRF protection

**Response:** Redirects to app deep link with success/error

---

## Accounts

### GET /accounts

List all connected accounts (owned and shared with user).

**Response:** `200 OK`
```json
{
    "accounts": [
        {
            "id": "uuid",
            "provider": "lifx",
            "provider_account_id": "user@lifx",
            "is_owner": true,
            "role": "owner",
            "device_count": 5,
            "share_count": 2,
            "created_at": "2024-01-15T10:30:00Z"
        },
        {
            "id": "uuid",
            "provider": "hue",
            "provider_account_id": "other@hue",
            "is_owner": false,
            "role": "controller",
            "device_count": 3,
            "shared_by": "friend@example.com",
            "created_at": "2024-01-10T08:00:00Z"
        }
    ]
}
```

### GET /accounts/:id

Get account details with devices.

**Response:** `200 OK`
```json
{
    "id": "uuid",
    "provider": "lifx",
    "provider_account_id": "user@lifx",
    "is_owner": true,
    "devices": [
        {
            "id": "d073d5xxxxxx",
            "name": "Living Room",
            "type": "bulb",
            "power": "on",
            "brightness": 0.8,
            "color": {
                "hue": 240,
                "saturation": 0.5,
                "kelvin": 3500
            },
            "online": true
        }
    ]
}
```

### DELETE /accounts/:id

Disconnect a provider account (owner only).

**Response:** `204 No Content`

### POST /accounts/:id/action

Perform an action on devices.

**Request:**
```json
{
    "selector": "all",
    "action": "set_state",
    "params": {
        "power": "on",
        "brightness": 0.5
    }
}
```

**Selector options:**
- `all` - All devices in account
- `id:d073d5xxxxxx` - Specific device
- `group:Living Room` - Device group
- `label:Desk Lamp` - Device by label

**Action options:**
- `set_state` - Set power, brightness, color
- `toggle` - Toggle power
- `breathe` - Breathe effect
- `pulse` - Pulse effect

**Response:** `200 OK`
```json
{
    "results": [
        {
            "id": "d073d5xxxxxx",
            "status": "ok",
            "label": "Living Room"
        }
    ]
}
```

---

## Sharing

### GET /accounts/:id/shares

List users with access to an account (owner only).

**Response:** `200 OK`
```json
{
    "shares": [
        {
            "id": "uuid",
            "user": {
                "id": "uuid",
                "email": "friend@example.com"
            },
            "role": "controller",
            "created_at": "2024-01-15T10:30:00Z"
        }
    ],
    "invitations": [
        {
            "id": "uuid",
            "email": "pending@example.com",
            "status": "pending",
            "expires_at": "2024-01-22T10:30:00Z"
        }
    ],
    "limits": {
        "current": 2,
        "max": 2
    }
}
```

### POST /accounts/:id/invite

Create an invitation to share access.

**Request:**
```json
{
    "email": "friend@example.com",
    "role": "controller"
}
```

**Response:** `201 Created`
```json
{
    "invitation": {
        "id": "uuid",
        "email": "friend@example.com",
        "expires_at": "2024-01-22T10:30:00Z"
    }
}
```

**Error (limit reached):** `403 Forbidden`
```json
{
    "error": "share_limit_reached",
    "message": "Upgrade to Pro to share with more users",
    "current": 2,
    "max": 2
}
```

### DELETE /accounts/:id/shares/:share_id

Remove a user's access (owner only).

**Response:** `204 No Content`

### DELETE /accounts/:id/invitations/:invitation_id

Cancel a pending invitation (owner only).

**Response:** `204 No Content`

---

## Invitations

### GET /invitations/pending

List pending invitations for current user.

**Response:** `200 OK`
```json
{
    "invitations": [
        {
            "id": "uuid",
            "account": {
                "provider": "lifx",
                "device_count": 5
            },
            "invited_by": "friend@example.com",
            "expires_at": "2024-01-22T10:30:00Z"
        }
    ]
}
```

### POST /invitations/:token/accept

Accept an invitation.

**Response:** `200 OK`
```json
{
    "account": {
        "id": "uuid",
        "provider": "lifx",
        "role": "controller"
    }
}
```

### POST /invitations/:token/decline

Decline an invitation.

**Response:** `204 No Content`

---

## Billing

### GET /billing/products

Get available subscription products.

**Response:** `200 OK`
```json
{
    "products": [
        {
            "id": "pro_monthly",
            "name": "LightShare Pro",
            "description": "Unlimited sharing, no ads",
            "prices": {
                "apple": "com.lightshare.pro.monthly",
                "google": "pro_monthly",
                "stripe": "price_xxxxx"
            },
            "features": [
                "Share with up to 10 users",
                "No advertisements",
                "Priority support"
            ]
        }
    ]
}
```

### POST /billing/validate

Validate a purchase receipt.

**Request (Apple):**
```json
{
    "platform": "apple",
    "receipt_data": "base64-encoded-receipt",
    "product_id": "com.lightshare.pro.monthly"
}
```

**Request (Google):**
```json
{
    "platform": "google",
    "purchase_token": "token-from-google",
    "product_id": "pro_monthly"
}
```

**Response:** `200 OK`
```json
{
    "valid": true,
    "subscription": {
        "product_id": "pro_monthly",
        "status": "active",
        "expires_at": "2025-01-15T10:30:00Z"
    },
    "entitlements": {
        "tier": "pro",
        "max_shares": 10,
        "ads_enabled": false
    }
}
```

### POST /billing/restore

Restore purchases (re-validate existing subscriptions).

**Request:**
```json
{
    "platform": "apple",
    "receipt_data": "base64-encoded-receipt"
}
```

**Response:** `200 OK`
```json
{
    "restored": true,
    "entitlements": {
        "tier": "pro",
        "max_shares": 10,
        "ads_enabled": false
    }
}
```

---

## Error Responses

All errors follow this format:

```json
{
    "error": "error_code",
    "message": "Human readable message",
    "details": {}
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `invalid_request` | 400 | Malformed request body |
| `unauthorized` | 401 | Missing or invalid token |
| `forbidden` | 403 | Insufficient permissions |
| `not_found` | 404 | Resource not found |
| `conflict` | 409 | Resource already exists |
| `rate_limited` | 429 | Too many requests |
| `provider_error` | 502 | Error from LIFX/Hue API |
| `internal_error` | 500 | Server error |

### Specific Error Codes

| Code | Description |
|------|-------------|
| `share_limit_reached` | User has reached sharing limit for their tier |
| `invitation_expired` | Invitation token has expired |
| `invitation_invalid` | Invitation token not found or already used |
| `provider_auth_failed` | Provider token validation failed |
| `receipt_invalid` | IAP receipt validation failed |

---

## Rate Limits

| Endpoint | Limit |
|----------|-------|
| `/auth/*` | 10 requests/minute |
| `/accounts/:id/action` | 60 requests/minute |
| All other endpoints | 120 requests/minute |

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1705312200
```

---

## Webhooks (Future)

Planned webhook support for:
- Subscription status changes
- Device state changes
- Share events
