# CLAUDE.md - LightShare Project Context

## Project Overview

LightShare is a mobile application that enables users to connect and control smart lighting systems (LIFX, Philips Hue) and share access with other users. The app uses a freemium model with subscription-based premium features.

## Tech Stack

### Mobile (Flutter)
- **State Management**: Riverpod
- **HTTP Client**: Dio
- **Secure Storage**: flutter_secure_storage (for session tokens only)
- **OAuth**: oauth2_client or openid_client (prefer backend-handled OAuth)
- **In-App Purchases**: in_app_purchase plugin
- **Ads**: google_mobile_ads
- **Push Notifications**: firebase_messaging
- **Testing**: flutter_test, integration_test

### Backend (Go)
- **Framework**: Fiber
- **Database**: PostgreSQL with sqlx (prefer hand-written queries)
- **Migrations**: golang-migrate
- **Cache/Sessions**: Redis
- **HTTP Client**: resty or standard http.Client
- **Auth**: jwt-go with DB-stored refresh tokens
- **Encryption**: KMS + AES-GCM for token encryption

## Key Concepts

### Authentication Flow
- Backend handles OAuth flows with providers (LIFX/Hue)
- Client receives only session tokens (JWT access + refresh)
- Provider tokens are NEVER sent to clients - backend proxies all API calls

### Token Security
- Provider tokens encrypted at rest using DEK (Data Encryption Key)
- DEK encrypted by KMS master key (AWS KMS/GCP KMS/Vault)
- Tokens validated on acceptance via provider API call
- Only backend service can access decrypted tokens

### Sharing Model
- Owners can invite users via email (one-time invite tokens)
- Free tier: share with up to 2 users
- Pro tier: increased sharing limit (10+ or unlimited)
- Roles: viewer/controller (start with controller only)

### Subscription Model
- Free: 2 user sharing, ads shown, basic features
- Pro: increased sharing, no ads
- iOS: Apple IAP (required by Apple)
- Android: Google Play Billing (required by Google)
- Web: Stripe payments
- Server validates all receipts and maintains canonical entitlement state

## Database Schema (Core Tables)

```sql
users (id, email, password_hash, stripe_customer_id, role, created_at)
accounts (id, owner_user_id, provider, provider_account_id, encrypted_token, metadata)
access_grants (id, account_id, grantee_user_id, role, created_by, created_at)
invitations (id, account_id, invitee_email, invite_token, expires_at, status)
```

## Development Guidelines

### Security First
- Never log secrets or tokens
- Always use TLS/HTTPS
- Validate all receipts server-side
- Rate-limit API calls (both user-facing and to providers)
- Run SAST and dependency scanning in CI

### API Design
- All provider interactions go through backend proxy
- Client calls `/accounts/:id/action` - backend calls LIFX/Hue
- Entitlements checked via `/me/entitlements` endpoint

### Testing Requirements
- Unit tests for all backend logic
- Integration tests with mocked providers
- E2E tests for Flutter app
- Manual QA required for IAP flows

## Common Tasks

### Adding a New Provider
1. Implement OAuth flow handler in backend
2. Add provider-specific API client
3. Create token validation logic
4. Add provider option in mobile UI
5. Update docs/architecture.md

### Implementing New Features
1. Check if feature requires entitlement (pro vs free)
2. Add server-side entitlement check
3. Update mobile UI with entitlement-aware rendering
4. Add appropriate tests

### Working with Subscriptions
1. Mobile handles purchase flow via platform SDK
2. Send receipt to backend `/billing/validate`
3. Backend validates with Apple/Google servers
4. Backend updates user entitlements in DB
5. Client refreshes entitlements via `/me/entitlements`

## File Structure (Planned)

```
lightshare/
├── mobile/                 # Flutter app
│   ├── lib/
│   │   ├── features/       # Feature modules
│   │   ├── core/           # Shared utilities
│   │   └── main.dart
│   └── test/
├── backend/                # Go backend
│   ├── cmd/
│   │   └── server/
│   ├── internal/
│   │   ├── handlers/
│   │   ├── services/
│   │   ├── models/
│   │   └── middleware/
│   ├── migrations/
│   └── pkg/
├── docs/                   # Documentation
└── docker-compose.yml
```

## Environment Variables (Backend)

```
DATABASE_URL=postgres://...
REDIS_URL=redis://...
JWT_SECRET=...
KMS_KEY_ID=...              # For token encryption
LIFX_CLIENT_ID=...
LIFX_CLIENT_SECRET=...
HUE_CLIENT_ID=...
HUE_CLIENT_SECRET=...
APPLE_SHARED_SECRET=...     # For receipt validation
GOOGLE_SERVICE_ACCOUNT=...  # For Play Developer API
```

## Important Constraints

1. **Apple/Google IAP Required**: In-app digital features MUST use platform IAP - cannot use external payments for unlocking features in mobile apps

2. **Provider Rate Limits**: LIFX and Hue have different rate limits - implement per-provider throttling

3. **Hue Bridge**: Some Hue operations require local bridge - verify cloud API availability for each feature

4. **Token Storage**: ONLY store backend session tokens on device - provider tokens stay server-side

5. **Concurrent Access**: Multiple users controlling same device = last-write-wins with optimistic state updates

## MVP Milestones

1. Project skeleton (repo, CI, Docker, basic auth)
2. Provider connection flow (token entry + validation + secure storage)
3. Light controls UI (list devices, on/off via backend proxy)
4. Sharing invites + ACL + free tier limit enforcement
5. IAP integration (receipt validation, pro features)
6. Ads integration + UI polish
7. QA + privacy policy + store submission

## Useful Commands

```bash
# Backend
cd backend && go run cmd/server/main.go
go test ./...
golang-migrate -path migrations -database $DATABASE_URL up

# Mobile
cd mobile && flutter run
flutter test
flutter build apk --release
flutter build ios --release

# Docker
docker-compose up -d
docker-compose logs -f backend
```

## Resources

- [LIFX API Docs](https://api.developer.lifx.com/)
- [Philips Hue API](https://developers.meethue.com/)
- [Apple Receipt Validation](https://developer.apple.com/documentation/storekit/in-app_purchase/validating_receipts_with_the_app_store)
- [Google Play Billing](https://developer.android.com/google/play/billing)
- Always run gofmt, goimports, golint-cli, alignment after modifying the go project before considering the task finished. Some tools are in /Users/hmartin/go/bin/
- Always run flutter analyze and test before considering the task finished after modifying the flutter project