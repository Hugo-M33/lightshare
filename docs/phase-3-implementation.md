# Phase 3 Implementation: Provider Connection

## Overview

Phase 3 implements the provider connection feature, allowing users to connect their LIFX smart lighting accounts to LightShare using personal access tokens. This phase establishes the foundation for controlling smart lights through the backend proxy.

## Implementation Status

✅ **Completed** - All Phase 3 components have been implemented

## Architecture

### Backend Components

#### 1. Database Schema

**Migration**: `backend/migrations/000003_create_accounts_table.up.sql`

```sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_account_id VARCHAR(255) NOT NULL,
    encrypted_token BYTEA NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(owner_user_id, provider, provider_account_id)
);
```

#### 2. Encryption Layer

**Location**: `backend/pkg/crypto/`

- **AES-256-GCM Encryption**: Provider tokens are encrypted using AES-256-GCM before storage
- **Local Dev KMS**: Uses `ENCRYPTION_KEY` environment variable (32-byte hex string)
- **Key Management**:
  - `LoadEncryptionKey()` - loads key from environment
  - `GenerateEncryptionKey()` - utility to generate new keys
  - `EncryptToken()` / `DecryptToken()` - encryption/decryption functions

**Security Features**:
- Random nonce for each encryption (prevents pattern analysis)
- Authentication tags (prevents tampering)
- Key validation (enforces 32-byte key length)

#### 3. Provider Abstraction

**Location**: `backend/pkg/providers/`

**Interface**:
```go
type Client interface {
    ValidateToken(token string) (*AccountInfo, error)
    GetAccountInfo(token string) (*AccountInfo, error)
}
```

**Implemented Providers**:
- ✅ LIFX (`backend/pkg/providers/lifx/client.go`)
- 🚧 Philips Hue (coming soon)

**Factory Pattern**:
```go
client, err := providers.NewClient(providers.ProviderLIFX)
```

#### 4. Business Logic

**Provider Service** (`backend/internal/services/provider.go`):
- `ConnectProvider()` - validates token → encrypts → stores account
- `ListAccounts()` - retrieves all accounts for a user
- `DisconnectAccount()` - deletes account (with ownership validation)

**Account Repository** (`backend/internal/repository/account.go`):
- CRUD operations for accounts
- Handles encrypted token storage/retrieval
- Enforces unique constraint on (user, provider, account)

#### 5. API Endpoints

**Provider Handler** (`backend/internal/handlers/provider.go`):

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/providers/connect` | ✅ | Connect a provider account |
| GET | `/api/v1/accounts` | ✅ | List connected accounts |
| DELETE | `/api/v1/accounts/:id` | ✅ | Disconnect an account |

**Request/Response Examples**:

```bash
# Connect provider
POST /api/v1/providers/connect
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "provider": "lifx",
  "token": "c87cec....."
}

# Response
{
  "id": "uuid",
  "provider": "lifx",
  "provider_account_id": "location-id",
  "metadata": {
    "lights_count": 5
  },
  "created_at": "2025-01-23T..."
}
```

### Mobile Components

#### 1. Models

**Location**: `mobile/lib/core/models/`

- `account.dart` - Account model with JSON serialization
- `provider.dart` - Provider enum, DTOs for API requests/responses

#### 2. Services

**Location**: `mobile/lib/core/services/provider_service.dart`

- API client for provider operations
- Error handling with DioException
- Maps HTTP responses to models

#### 3. State Management

**Location**: `mobile/lib/core/providers/accounts_provider.dart`

**Riverpod State**:
```dart
final accountsProvider = StateNotifierProvider<AccountsNotifier, AccountsState>
```

**State Operations**:
- `loadAccounts()` - fetches all accounts
- `connectProvider()` - connects new provider
- `disconnectAccount()` - removes account

#### 4. UI Screens

**Location**: `mobile/lib/features/providers/screens/`

1. **accounts_screen.dart** - Lists connected accounts with disconnect action
2. **provider_selection_screen.dart** - Choose provider (LIFX/Hue)
3. **token_entry_screen.dart** - Enter token with instructions

**Features**:
- Glass morphism UI matching app theme
- Inline help text with token generation instructions
- Error handling with snackbar notifications
- Loading states during API calls

## Setup Instructions

### Backend Setup

1. **Generate Encryption Key**:
```bash
openssl rand -hex 32
```

2. **Set Environment Variable**:
```bash
export ENCRYPTION_KEY="your-64-character-hex-string"
```

3. **Run Migrations**:
```bash
cd backend
golang-migrate -path migrations -database $DATABASE_URL up
```

4. **Start Server**:
```bash
go run cmd/server/main.go
```

### Mobile Setup

1. **Update Router** (if not already done):
Add routes in `mobile/lib/core/router/app_router.dart`:
```dart
GoRoute(
  path: '/accounts',
  builder: (context, state) => const AccountsScreen(),
),
GoRoute(
  path: '/providers/connect',
  builder: (context, state) => const ProviderSelectionScreen(),
),
GoRoute(
  path: '/providers/connect/token',
  builder: (context, state) {
    final provider = state.extra as Provider;
    return TokenEntryScreen(provider: provider);
  },
),
```

2. **Run App**:
```bash
cd mobile
flutter run
```

## Testing

### Backend Tests

**Encryption Tests** (`backend/pkg/crypto/crypto_test.go`):
- ✅ Encrypt/decrypt round trip
- ✅ Invalid key handling
- ✅ Encryption uniqueness (random nonce)
- ✅ Garbage data rejection

**Provider Service Tests** (`backend/internal/services/provider_test.go`):
- ✅ Invalid provider rejection
- ✅ Account listing
- ✅ Account disconnect
- ✅ Ownership validation

**Run Tests**:
```bash
cd backend
go test ./pkg/crypto/
go test ./internal/services/
```

## Usage Flow

### User Journey

1. **User navigates to "Connected Accounts"**
2. **Clicks "Connect Account"**
3. **Selects "LIFX"** from provider list
4. **Follows instructions** to get LIFX token from https://cloud.lifx.com/settings
5. **Pastes token** and clicks "Connect Account"
6. **Backend validates token** by calling LIFX API
7. **Token is encrypted** using AES-256-GCM
8. **Account is saved** to database
9. **Success message shown**, user redirected to accounts list

### LIFX Token Generation

Users need to:
1. Visit https://cloud.lifx.com/settings
2. Log in to LIFX account
3. Scroll to "Personal Access Tokens"
4. Click "Generate New Token"
5. Copy the generated token
6. Paste into LightShare app

## Security Considerations

### Token Encryption

- ✅ **AES-256-GCM** - Industry standard authenticated encryption
- ✅ **Random nonces** - Each encryption produces unique ciphertext
- ✅ **Authentication tags** - Prevents tampering
- ✅ **Key never exposed** - Stored only in backend environment

### API Security

- ✅ **JWT Authentication** - All provider endpoints require valid JWT
- ✅ **Ownership validation** - Users can only access their own accounts
- ✅ **Token validation** - Provider tokens validated before acceptance
- ✅ **No token exposure** - Encrypted tokens never sent to client

### Error Handling

- ✅ **Generic error messages** - Don't leak internal details
- ✅ **Input validation** - Tokens validated before processing
- ✅ **Rate limiting** - (inherited from Phase 2 middleware)

## Known Limitations

1. **LIFX Only**: Philips Hue support planned for future phase
2. **Token Method Only**: OAuth flow will be added in Phase 5
3. **No Token Refresh**: Users must manually update expired tokens
4. **Basic Validation**: Only verifies token works (list lights call)

## Future Enhancements (Phase 5 - OAuth)

- OAuth 2.0 flow for LIFX
- OAuth 2.0 flow for Philips Hue
- Automatic token refresh
- Token expiration tracking
- Multi-account support per provider

## Troubleshooting

### Backend Issues

**Error: "ENCRYPTION_KEY environment variable not set"**
- Generate key: `openssl rand -hex 32`
- Set in environment: `export ENCRYPTION_KEY="..."`

**Error: "invalid provider token"**
- Token may be expired
- Token may have insufficient permissions
- LIFX API may be down
- Check token at https://cloud.lifx.com/settings

### Mobile Issues

**"Unable to connect to server"**
- Check API_BASE_URL configuration
- Verify backend is running
- Check network connectivity

**"Account already connected"**
- Each LIFX account can only be connected once per user
- Disconnect existing account first if reconnecting

## File Structure

```
backend/
├── migrations/
│   ├── 000003_create_accounts_table.up.sql
│   └── 000003_create_accounts_table.down.sql
├── pkg/
│   ├── crypto/
│   │   ├── crypto.go (encryption functions)
│   │   ├── key.go (key management)
│   │   └── crypto_test.go
│   └── providers/
│       ├── provider.go (interface & factory)
│       └── lifx/
│           └── client.go
├── internal/
│   ├── models/
│   │   └── account.go
│   ├── repository/
│   │   └── account.go
│   ├── services/
│   │   ├── provider.go
│   │   └── provider_test.go
│   └── handlers/
│       └── provider.go

mobile/
├── lib/
│   ├── core/
│   │   ├── models/
│   │   │   ├── account.dart
│   │   │   └── provider.dart
│   │   ├── services/
│   │   │   └── provider_service.dart
│   │   └── providers/
│   │       └── accounts_provider.dart
│   └── features/
│       └── providers/
│           └── screens/
│               ├── accounts_screen.dart
│               ├── provider_selection_screen.dart
│               └── token_entry_screen.dart
```

## API Reference

See [api.md](./api.md) for detailed API documentation including:
- Request/response schemas
- Error codes
- Authentication requirements
- Rate limiting

## Next Steps

After Phase 3 completion, proceed to:
- **Phase 4**: Light Control - Actually controlling lights through backend proxy
- Implement device listing endpoints
- Implement light control endpoints (on/off, brightness, color)
- Add backend caching for device state
