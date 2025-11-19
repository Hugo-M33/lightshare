# Development Roadmap

This roadmap outlines the recommended order of implementation for LightShare. Each phase builds on the previous one, with clear dependencies noted.

---

## Phase 1: Project Foundation

**Goal:** Establish project structure, CI/CD, and basic infrastructure.

### 1.1 Repository Setup
- [ ] Initialize Flutter project structure (`mobile/`)
- [ ] Initialize Go module structure (`backend/`)
- [ ] Set up Docker and docker-compose for local dev
- [ ] Configure linting (golangci-lint, flutter analyze)
- [ ] Set up pre-commit hooks

### 1.2 CI/CD Pipeline
- [ ] GitHub Actions workflow for backend tests
- [ ] GitHub Actions workflow for Flutter tests
- [ ] Docker image build and push
- [ ] Code coverage reporting

### 1.3 Database Setup
- [ ] PostgreSQL container configuration
- [ ] Redis container configuration
- [ ] golang-migrate setup
- [ ] Initial migration: users table

### 1.4 Backend Skeleton
- [ ] Fiber app structure
- [ ] Configuration management (env vars)
- [ ] Structured logging setup
- [ ] Health check endpoint (`/health`)
- [ ] Basic middleware (CORS, request ID, logging)

**Deliverable:** Running backend with health endpoint, CI passing, Docker working.

---

## Phase 2: User Authentication

**Goal:** Complete user auth system with secure token handling.

**Depends on:** Phase 1

### 2.1 User Model & Storage
- [ ] Migration: complete users table
- [ ] User repository (create, find by email, find by ID)
- [ ] Password hashing with bcrypt

### 2.2 Auth Endpoints
- [ ] `POST /auth/signup` - user registration
- [ ] `POST /auth/login` - authentication
- [ ] `POST /auth/refresh` - token refresh
- [ ] `POST /auth/logout` - token revocation

### 2.3 JWT Implementation
- [ ] Access token generation (1 hour expiry)
- [ ] Refresh token generation (30 days expiry)
- [ ] Refresh token storage in database
- [ ] Auth middleware for protected routes
- [ ] Token revocation on logout

### 2.4 Mobile Auth UI
- [ ] Login screen
- [ ] Signup screen
- [ ] Secure token storage (flutter_secure_storage)
- [ ] Dio interceptor for auth headers
- [ ] Token refresh logic
- [ ] Auth state management (Riverpod)

**Deliverable:** Users can sign up, log in, and maintain sessions. Mobile app persists auth state.

---

## Phase 3: Provider Connection (Token Entry)

**Goal:** Allow users to connect LIFX/Hue accounts via personal access tokens.

**Depends on:** Phase 2

### 3.1 Token Encryption Setup
- [ ] KMS integration (AWS KMS or local dev alternative)
- [ ] DEK generation and storage
- [ ] AES-256-GCM encryption/decryption functions
- [ ] Secure key loading at boot

### 3.2 Account Model
- [ ] Migration: accounts table
- [ ] Account repository (create, find by user, find by ID)
- [ ] Encrypted token storage

### 3.3 Provider Clients
- [ ] LIFX API client
  - [ ] Token validation (list lights)
  - [ ] Get account info
- [ ] Hue API client (Remote API)
  - [ ] Token validation
  - [ ] Get account info

### 3.4 Connection Endpoints
- [ ] `POST /providers/connect` (token method)
  - [ ] Accept token from client
  - [ ] Validate with provider API
  - [ ] Encrypt and store token
  - [ ] Create account record
- [ ] `GET /accounts` - list connected accounts
- [ ] `DELETE /accounts/:id` - disconnect account

### 3.5 Mobile Connection UI
- [ ] Provider selection screen
- [ ] Token entry screen with instructions
- [ ] Account list screen
- [ ] Connection success/error handling

**Deliverable:** Users can connect LIFX/Hue accounts using personal tokens. Tokens stored encrypted.

---

## Phase 4: Light Control

**Goal:** Users can view and control their lights through the app.

**Depends on:** Phase 3

### 4.1 Device Fetching
- [ ] LIFX: list all lights endpoint
- [ ] Hue: list all lights endpoint
- [ ] Normalize device data across providers
- [ ] Cache device lists (Redis, short TTL)

### 4.2 Control Actions
- [ ] LIFX actions: power, brightness, color, effects
- [ ] Hue actions: power, brightness, color
- [ ] `POST /accounts/:id/action` endpoint
- [ ] Action validation and rate limiting

### 4.3 Mobile Control UI
- [ ] Device list screen (grouped by account/room)
- [ ] Device detail/control screen
- [ ] Power toggle
- [ ] Brightness slider
- [ ] Color picker (if supported)
- [ ] Optimistic UI updates
- [ ] Pull-to-refresh for device state

### 4.4 Error Handling
- [ ] Provider API error handling
- [ ] Offline device handling
- [ ] Rate limit feedback to user

**Deliverable:** Users can view all their lights and control power/brightness/color.

---

## Phase 5: OAuth Provider Connection

**Goal:** Add OAuth flow as alternative to token entry (better UX).

**Depends on:** Phase 4

### 5.1 OAuth Configuration
- [ ] Register app with LIFX developer portal
- [ ] Register app with Hue developer portal
- [ ] Store client IDs/secrets securely

### 5.2 OAuth Flow Backend
- [ ] `POST /providers/connect` (oauth method) - generate auth URL
- [ ] `GET /providers/oauth/callback` - handle callback
- [ ] Exchange code for tokens (with PKCE)
- [ ] Store access + refresh tokens

### 5.3 Token Refresh
- [ ] Automatic refresh on 401 from provider
- [ ] Refresh token rotation
- [ ] Handle refresh failure (re-auth required)

### 5.4 Mobile OAuth Flow
- [ ] Open auth URL in browser/webview
- [ ] Deep link handling for callback
- [ ] PKCE implementation

**Deliverable:** Users can connect via OAuth for seamless authentication.

---

## Phase 6: Sharing System

**Goal:** Account owners can share light access with other users.

**Depends on:** Phase 4

### 6.1 Sharing Models
- [ ] Migration: access_grants table
- [ ] Migration: invitations table
- [ ] Grant repository
- [ ] Invitation repository

### 6.2 Invitation System
- [ ] `POST /accounts/:id/invite` - create invitation
- [ ] Generate secure invitation tokens
- [ ] Store in Redis with TTL (7 days)
- [ ] Email sending (SendGrid/SES/etc.)
- [ ] `GET /invitations/pending` - list user's invitations
- [ ] `POST /invitations/:token/accept`
- [ ] `POST /invitations/:token/decline`

### 6.3 Access Control
- [ ] Modify action endpoint to check grants
- [ ] Owner vs grantee permissions
- [ ] `GET /accounts/:id/shares` - list shares
- [ ] `DELETE /accounts/:id/shares/:id` - revoke access

### 6.4 Share Limits (Free Tier)
- [ ] Default limit: 2 shares per account
- [ ] Check limit on invite creation
- [ ] Return appropriate error when exceeded

### 6.5 Mobile Sharing UI
- [ ] Share management screen (owner view)
- [ ] Invite user form
- [ ] Pending invitations list
- [ ] Accept/decline invitation screen
- [ ] Shared accounts indicator
- [ ] Revoke access confirmation

### 6.6 Notifications
- [ ] Push notification setup (Firebase)
- [ ] Notify on invitation received
- [ ] Notify on invitation accepted

**Deliverable:** Users can invite others to control their lights with free tier limit enforced.

---

## Phase 7: Subscriptions & Payments

**Goal:** Implement Pro tier with IAP and entitlement system.

**Depends on:** Phase 6

### 7.1 Subscription Model
- [ ] Migration: subscriptions table
- [ ] Subscription repository
- [ ] Entitlement calculation logic

### 7.2 Apple IAP
- [ ] Configure products in App Store Connect
- [ ] Server-side receipt validation
- [ ] `POST /billing/validate` (Apple)
- [ ] Handle subscription renewals
- [ ] Handle cancellations/expirations

### 7.3 Google Play Billing
- [ ] Configure products in Play Console
- [ ] Service account setup
- [ ] Server-side purchase validation
- [ ] `POST /billing/validate` (Google)
- [ ] Real-time developer notifications (RTDN)

### 7.4 Entitlements System
- [ ] `GET /me/entitlements` endpoint
- [ ] Update share limits based on subscription
- [ ] Track subscription status changes

### 7.5 Mobile IAP Integration
- [ ] in_app_purchase plugin setup
- [ ] Subscription screen with product info
- [ ] Purchase flow
- [ ] Restore purchases
- [ ] Subscription status display

### 7.6 Stripe (Web - Optional)
- [ ] Stripe product/price setup
- [ ] Checkout session creation
- [ ] Webhook handling
- [ ] Customer portal for management

**Deliverable:** Users can subscribe to Pro, unlocking increased share limits.

---

## Phase 8: Advertisements

**Goal:** Display ads for free tier users.

**Depends on:** Phase 7 (needs entitlements to know when to hide ads)

### 8.1 Ad Integration
- [ ] AdMob account setup
- [ ] Configure ad units (banner, possibly interstitial)
- [ ] google_mobile_ads plugin setup

### 8.2 Ad Placement
- [ ] Banner ad on home/device list screen
- [ ] Respect entitlements (hide for Pro)
- [ ] Handle ad loading failures gracefully

### 8.3 Privacy Compliance
- [ ] ATT prompt for iOS
- [ ] GDPR consent for EU users
- [ ] Update privacy policy

**Deliverable:** Free users see ads, Pro users don't.

---

## Phase 9: Polish & Production Prep

**Goal:** Prepare for app store submission.

**Depends on:** Phase 8

### 9.1 UI/UX Polish
- [ ] Consistent design system
- [ ] Loading states
- [ ] Empty states
- [ ] Error states
- [ ] Animations and transitions
- [ ] Accessibility review

### 9.2 Testing
- [ ] Backend unit tests (>80% coverage)
- [ ] Flutter unit tests
- [ ] Integration tests
- [ ] E2E tests for critical flows
- [ ] Manual QA for IAP flows
- [ ] Device testing matrix

### 9.3 Security Audit
- [ ] Dependency vulnerability scan
- [ ] Static analysis (SAST)
- [ ] Penetration testing
- [ ] Security checklist review

### 9.4 Documentation
- [ ] Privacy policy
- [ ] Terms of service
- [ ] App store descriptions
- [ ] Screenshots and previews
- [ ] Support documentation

### 9.5 Monitoring Setup
- [ ] Sentry for error tracking
- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] Alerting rules

### 9.6 Production Infrastructure
- [ ] Production database (managed PostgreSQL)
- [ ] Production Redis
- [ ] KMS key setup
- [ ] Domain and SSL certificates
- [ ] CDN for any static assets

**Deliverable:** App ready for submission with monitoring in place.

---

## Phase 10: Launch

**Goal:** Submit to app stores and launch.

**Depends on:** Phase 9

### 10.1 App Store Submission
- [ ] iOS app archive and upload
- [ ] App Store Connect metadata
- [ ] App review preparation
- [ ] Address review feedback

### 10.2 Play Store Submission
- [ ] Android app bundle
- [ ] Play Console listing
- [ ] Content rating questionnaire
- [ ] Review and launch

### 10.3 Launch Tasks
- [ ] Monitoring dashboards live
- [ ] On-call rotation setup
- [ ] Rollback procedures documented
- [ ] Customer support ready

**Deliverable:** App live on both stores.

---

## Post-Launch Phases

### Phase 11: OAuth Improvements
- [ ] Additional providers (if applicable)
- [ ] Improved token refresh handling
- [ ] Better error recovery

### Phase 12: Advanced Features
- [ ] Scenes/routines
- [ ] Schedules
- [ ] Widgets (iOS/Android)
- [ ] Apple Watch / Wear OS

### Phase 13: Analytics & Optimization
- [ ] Usage analytics
- [ ] A/B testing
- [ ] Performance optimization
- [ ] Conversion optimization

---

## Dependency Graph

```
Phase 1 (Foundation)
    │
    ▼
Phase 2 (Auth)
    │
    ▼
Phase 3 (Token Connection)
    │
    ▼
Phase 4 (Light Control)
    │
    ├────────────────┐
    ▼                ▼
Phase 5 (OAuth)   Phase 6 (Sharing)
                     │
                     ▼
                 Phase 7 (Payments)
                     │
                     ▼
                 Phase 8 (Ads)
                     │
                     ▼
                 Phase 9 (Polish)
                     │
                     ▼
                 Phase 10 (Launch)
```

---

## Quick Reference: Critical Path to MVP

For fastest path to a working MVP:

1. **Phase 1** - Foundation (1-2 weeks)
2. **Phase 2** - Auth (1 week)
3. **Phase 3** - Token Connection (1 week)
4. **Phase 4** - Light Control (1-2 weeks)
5. **Phase 6** - Sharing (1-2 weeks)
6. **Phase 7** - Payments (1-2 weeks)
7. **Phase 8** - Ads (few days)
8. **Phase 9** - Polish (1-2 weeks)
9. **Phase 10** - Launch (1 week + review time)

OAuth (Phase 5) can be done in parallel or deferred post-launch as enhancement.

---

## Notes

- **Start simple:** Token entry before OAuth, basic controls before effects
- **Test early:** Set up test accounts with LIFX/Hue from Phase 3
- **Security first:** Never skip encryption, always validate server-side
- **Platform rules:** IAP is mandatory for iOS/Android - no workarounds
