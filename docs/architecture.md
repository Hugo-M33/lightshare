# Architecture

## System Overview

LightShare follows a client-server architecture where the mobile app communicates with a backend API, which in turn proxies requests to smart lighting providers.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Mobile    │────▶│   Backend   │────▶│   LIFX/Hue  │
│   (Flutter) │◀────│   (Go)      │◀────│   APIs      │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                    ┌──────┴──────┐
                    │             │
              ┌─────▼─────┐ ┌─────▼─────┐
              │ PostgreSQL│ │   Redis   │
              └───────────┘ └───────────┘
```

## Components

### Mobile Application (Flutter)

**Responsibilities:**
- User interface for light control
- Account management and settings
- Subscription purchase flow (IAP)
- Local session token storage

**Key Packages:**
- `riverpod` - State management
- `dio` - HTTP client
- `flutter_secure_storage` - Secure token storage
- `in_app_purchase` - IAP handling
- `google_mobile_ads` - Advertisement display

### Backend API (Go/Fiber)

**Responsibilities:**
- User authentication and authorization
- OAuth flow handling with providers
- Secure token storage and encryption
- Proxy API calls to LIFX/Hue
- Subscription/receipt validation
- Sharing and invitation management

**Key Libraries:**
- `fiber` - Web framework
- `sqlx` - Database queries
- `jwt-go` - JWT handling
- `resty` - HTTP client for providers

### Database (PostgreSQL)

**Core Tables:**

```sql
-- Users table
users (
    id              UUID PRIMARY KEY,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    stripe_customer_id VARCHAR(255),
    role            VARCHAR(50) DEFAULT 'user',
    created_at      TIMESTAMP DEFAULT NOW()
)

-- Connected provider accounts
accounts (
    id                  UUID PRIMARY KEY,
    owner_user_id       UUID REFERENCES users(id),
    provider            VARCHAR(50) NOT NULL,  -- 'lifx', 'hue'
    provider_account_id VARCHAR(255) NOT NULL,
    encrypted_token     BYTEA NOT NULL,
    metadata            JSONB,
    created_at          TIMESTAMP DEFAULT NOW()
)

-- Access grants for sharing
access_grants (
    id              UUID PRIMARY KEY,
    account_id      UUID REFERENCES accounts(id),
    grantee_user_id UUID REFERENCES users(id),
    role            VARCHAR(50) DEFAULT 'controller',
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(account_id, grantee_user_id)
)

-- Pending invitations
invitations (
    id              UUID PRIMARY KEY,
    account_id      UUID REFERENCES accounts(id),
    invitee_email   VARCHAR(255) NOT NULL,
    invite_token    VARCHAR(255) UNIQUE NOT NULL,
    expires_at      TIMESTAMP NOT NULL,
    status          VARCHAR(50) DEFAULT 'pending',
    created_at      TIMESTAMP DEFAULT NOW()
)

-- Subscription entitlements
subscriptions (
    id              UUID PRIMARY KEY,
    user_id         UUID REFERENCES users(id),
    platform        VARCHAR(50) NOT NULL,  -- 'apple', 'google', 'stripe'
    product_id      VARCHAR(255) NOT NULL,
    status          VARCHAR(50) NOT NULL,
    expires_at      TIMESTAMP,
    receipt_data    TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
)
```

### Cache (Redis)

**Use Cases:**
- Session management
- Rate limiting (per-user and per-provider)
- Invitation token storage (with TTL)
- Temporary OAuth state storage

## Data Flows

### OAuth Provider Connection

```
1. User initiates connection in app
2. App calls POST /providers/connect
3. Backend generates OAuth URL with state
4. App opens URL in browser/webview
5. User authenticates with provider
6. Provider redirects to backend callback
7. Backend exchanges code for tokens
8. Backend validates token with test API call
9. Backend encrypts and stores token
10. Backend creates account record
11. App receives success response
```

### Personal Access Token Connection

```
1. User pastes token in app
2. App sends POST /providers/connect with token
3. Backend validates token with provider API
4. Backend encrypts and stores token
5. Backend creates account record
6. App receives success response
```

### Light Control Action

```
1. User taps control in app
2. App calls POST /accounts/:id/action
3. Backend verifies user has access (owner or grantee)
4. Backend decrypts provider token
5. Backend calls provider API
6. Backend returns result to app
7. App updates UI state
```

### Sharing Invitation Flow

```
1. Owner creates invite via POST /accounts/:id/invite
2. Backend checks share limit (free: 2, pro: 10+)
3. Backend creates invitation with token
4. Backend sends email to invitee
5. Invitee clicks link
   a. If registered: Accept screen in app
   b. If not registered: Sign up, then accept
6. Invitee accepts via POST /invitations/:token/accept
7. Backend creates access_grant record
8. Both users can now control lights
```

### Subscription Purchase Flow

```
1. User initiates purchase in app
2. App uses platform IAP SDK (Apple/Google)
3. Platform processes payment
4. App receives receipt
5. App sends receipt to POST /billing/validate
6. Backend validates with Apple/Google servers
7. Backend updates subscription record
8. Backend returns updated entitlements
9. App updates UI (remove ads, unlock features)
```

## Token Encryption

### Encryption Architecture

```
┌─────────────────────────────────┐
│          KMS Master Key         │
│    (AWS KMS / GCP KMS / Vault)  │
└────────────────┬────────────────┘
                 │ encrypts
                 ▼
┌─────────────────────────────────┐
│     Data Encryption Key (DEK)   │
│      (stored encrypted in DB)   │
└────────────────┬────────────────┘
                 │ encrypts
                 ▼
┌─────────────────────────────────┐
│       Provider Tokens           │
│    (AES-256-GCM encrypted)      │
└─────────────────────────────────┘
```

### Process

1. **At Boot**: Backend calls KMS to decrypt DEK, stores in memory
2. **On Store**: Token encrypted with DEK using AES-256-GCM
3. **On Read**: Ciphertext decrypted with in-memory DEK
4. **Rotation**: New DEK can be created, re-encrypt all tokens

## Provider Integration

### LIFX

- **Auth**: OAuth2 or Personal Access Token
- **API**: REST-based cloud API
- **Rate Limits**: 120 requests per minute
- **Features**: All operations available via cloud

### Philips Hue

- **Auth**: OAuth2 (Remote API) or Local API (bridge)
- **API**: REST-based
- **Rate Limits**: Varies by endpoint
- **Considerations**:
  - Remote API requires Hue developer account
  - Some features may require local bridge access
  - Bridge discovery needed for local control

## Deployment Architecture

### Development

```
Docker Compose:
- postgres:15
- redis:7
- backend (hot reload)
```

### Production

```
Cloud Provider (AWS/GCP/DO):
├── Load Balancer (TLS termination)
├── Backend Service (containerized, auto-scaling)
├── Managed PostgreSQL (with backups)
├── Managed Redis (ElastiCache/Memorystore)
└── KMS (for key management)
```

### CI/CD Pipeline

```
GitHub Actions:
1. Run tests (backend + mobile)
2. Build Docker image
3. Push to container registry
4. Deploy to staging
5. Run E2E tests
6. Deploy to production
7. Mobile: Fastlane for store releases
```

## Scaling Considerations

### Horizontal Scaling
- Backend is stateless - scale behind load balancer
- Use Redis for shared session state
- Database connection pooling

### Rate Limiting Strategy
- Per-user limits (prevent abuse)
- Per-provider limits (respect API quotas)
- Implement backoff and queuing for bulk operations

### Caching Strategy
- Cache device lists (invalidate on change)
- Cache user entitlements (short TTL)
- Don't cache light state (changes frequently)
