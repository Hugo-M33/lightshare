# Phase 2 Implementation - Authentication System

## Overview
Phase 2 has been successfully implemented with full authentication support including:
- Email/password authentication
- Email verification
- Magic link login (passwordless)
- JWT token management with refresh tokens
- Secure token storage

## Backend Implementation ✅

### Database
- ✅ Updated `users` table with email verification fields
- ✅ Created `refresh_tokens` table for secure token management
- ✅ Proper indexes for performance

### Core Services
- ✅ **Database Layer**: PostgreSQL with sqlx, connection pooling
- ✅ **Redis Client**: Session and cache management
- ✅ **JWT Service**: Token generation, validation, refresh
- ✅ **Crypto Service**: Bcrypt password hashing, SHA-256 token hashing
- ✅ **Email Service**: SMTP-based with HTML templates for verification and magic links

### Repositories
- ✅ **User Repository**: Full CRUD operations with email verification
- ✅ **Refresh Token Repository**: Token lifecycle management

### Auth Service
- ✅ Signup with email verification
- ✅ Login with password
- ✅ Email verification flow
- ✅ Magic link generation and verification
- ✅ Token refresh mechanism
- ✅ Logout (single and all devices)

### API Endpoints
All mounted under `/api/v1/auth`:
- `POST /signup` - Create new account
- `POST /login` - Authenticate with email/password
- `POST /verify-email` - Verify email with token
- `POST /magic-link` - Request magic link
- `POST /magic-link/verify` - Login with magic link
- `POST /refresh` - Refresh access token
- `POST /logout` - Logout from current device
- `GET /me` - Get current user (protected)
- `POST /logout-all` - Logout from all devices (protected)

### Middleware
- ✅ **Auth Middleware**: JWT validation, automatic token refresh
- ✅ **Role-based access control**: Support for different user roles

## Mobile Implementation ✅

### Models
- ✅ `User` model with JSON serialization
- ✅ `AuthResponse` and `SignupResponse` DTOs

### Services
- ✅ **ApiClient**: Dio-based HTTP client with:
  - Automatic token injection
  - Auto-refresh on 401
  - Error handling
- ✅ **AuthService**: Complete auth operations

### State Management (Riverpod)
- ✅ **App Providers**: Dependency injection setup
- ✅ **Auth Provider**: Comprehensive auth state management
  - Loading states
  - Error handling
  - User persistence

### Navigation
- ✅ **GoRouter setup** with:
  - Auth guards
  - Automatic redirects
  - Deep linking support

##TODO: Remaining UI Screens

The following screens need to be implemented (templates provided below):

### 1. Login Screen
Location: `mobile/lib/features/auth/screens/login_screen.dart`
- Email/password fields
- "Forgot password?" link → Magic link
- "Create account" link → Signup
- Loading state handling
- Error display

### 2. Signup Screen
Location: `mobile/lib/features/auth/screens/signup_screen.dart`
- Email/password fields
- Password strength indicator
- Terms acceptance
- Email verification message after signup
- Link to login

### 3. Email Verification Screen
Location: `mobile/lib/features/auth/screens/email_verification_screen.dart`
- Auto-verify on deep link
- Success/error messages
- Redirect to login after verification

### 4. Magic Link Screen
Location: `mobile/lib/features/auth/screens/magic_link_screen.dart`
- Email input
- "Check your email" confirmation
- Auto-login on deep link

### 5. Home Screen (Updated)
Location: `mobile/lib/features/home/screens/home_screen.dart`
- Display user email
- Logout button
- Placeholder for future features

### 6. Update main.dart
Location: `mobile/lib/main.dart`
- Wire up GoRouter
- Add app initialization

## Configuration

### Backend Environment Variables
```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/lightshare?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_EXPIRATION=1h
JWT_REFRESH_EXPIRATION=720h  # 30 days

# Email (for development, use Mailhog or similar)
SMTP_HOST=localhost
SMTP_PORT=1025
EMAIL_FROM=noreply@lightshare.com
EMAIL_FROM_NAME=LightShare
APP_BASE_URL=http://localhost:8080

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

### Mobile Configuration
Update `mobile/lib/core/providers/app_providers.dart`:
```dart
final apiBaseUrlProvider = Provider<String>((ref) {
  return 'http://10.0.2.2:8080'; // Android emulator
  // return 'http://localhost:8080'; // iOS simulator
  // return 'https://api.lightshare.com'; // Production
});
```

## Testing

### Backend
```bash
cd backend
go test ./...
```

### Mobile
```bash
cd mobile
flutter test
flutter run
```

## Deep Linking Setup

### Android (`android/app/src/main/AndroidManifest.xml`)
```xml
<intent-filter android:autoVerify="true">
    <action android:name="android.intent.action.VIEW" />
    <category android:name="android.intent.category.DEFAULT" />
    <category android:name="android.intent.category.BROWSABLE" />
    <data
        android:scheme="https"
        android:host="app.lightshare.com" />
    <data android:scheme="lightshare" />
</intent-filter>
```

### iOS (`ios/Runner/Info.plist`)
```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleTypeRole</key>
        <string>Editor</string>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>lightshare</string>
        </array>
    </dict>
</array>
```

## Security Features

1. **Password Security**
   - Bcrypt hashing with cost 12
   - Minimum 8 characters required
   - Stored hashes never exposed

2. **Token Security**
   - Short-lived access tokens (1 hour)
   - Long-lived refresh tokens (30 days)
   - SHA-256 hashed refresh tokens in DB
   - Automatic rotation on refresh

3. **Email Verification**
   - Cryptographically secure tokens
   - 24-hour expiration
   - One-time use

4. **Magic Links**
   - 15-minute expiration
   - One-time use
   - Cleared after successful login

## Next Steps (Phase 3)

After completing the UI screens, proceed to Phase 3:
1. Provider connection flow (LIFX, Hue)
2. Token encryption with KMS
3. OAuth2 flows
4. Provider API clients

## Notes

- The backend is fully implemented and ready for testing
- Run migrations before starting: `docker-compose up -d && cd backend && golang-migrate -path migrations -database $DATABASE_URL up`
- For development email testing, use [Mailhog](https://github.com/mailhog/MailHog) or similar SMTP server
- Mobile UI screens follow Material 3 design principles
