# Mobile Deep Linking Guide

## Overview

LightShare implements custom URL scheme deep linking to provide a seamless email verification and magic link login experience. When users click email links, they're automatically redirected to the mobile app and logged in without manual intervention.

## How It Works

### Flow Diagram

```
User signs up
  ↓
Backend sends email with: lightshare://verify-email?token=abc123
  ↓
User clicks link in email
  ↓
OS opens LightShare app via deep link
  ↓
App receives deep link → routes to verification screen
  ↓
Screen auto-verifies email → receives JWT tokens
  ↓
User is logged in → redirected to home dashboard
```

## Configuration

### Backend Configuration

**Environment Variable**
```bash
# backend/.env
MOBILE_DEEP_LINK_SCHEME=lightshare
```

**Email Templates**
- Verification emails use: `lightshare://verify-email?token={token}`
- Magic link emails use: `lightshare://magic-link?token={token}`

**API Response**
The `POST /api/v1/auth/verify-email` endpoint now returns JWT tokens (access + refresh) upon successful verification, enabling immediate auto-login.

### Android Configuration

**File**: `mobile/android/app/src/main/AndroidManifest.xml`

```xml
<intent-filter>
    <action android:name="android.intent.action.VIEW"/>
    <category android:name="android.intent.category.DEFAULT"/>
    <category android:name="android.intent.category.BROWSABLE"/>
    <data android:scheme="lightshare"/>
</intent-filter>
```

This registers the app to handle URLs starting with `lightshare://`.

### iOS Configuration

**File**: `mobile/ios/Runner/Info.plist`

```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleTypeRole</key>
        <string>Editor</string>
        <key>CFBundleURLName</key>
        <string>com.lightshare.app</string>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>lightshare</string>
        </array>
    </dict>
</array>
```

This registers the app to handle the `lightshare://` URL scheme.

### Flutter Configuration

**Package**: `app_links: ^6.3.2` (in `pubspec.yaml`)

**Deep Link Handler** (`mobile/lib/main.dart`):
```dart
void _handleDeepLink(Uri uri) {
  // For custom URL schemes like lightshare://verify-email?token=xxx,
  // the "verify-email" part is the host, not the path
  final host = uri.host;
  final queryParams = uri.queryParameters;

  if (host == 'verify-email') {
    final token = queryParams['token'];
    if (token != null) {
      router.go('/auth/verify-email?token=$token');
    }
  } else if (host == 'magic-link') {
    final token = queryParams['token'];
    if (token != null) {
      router.go('/auth/magic-link?token=$token');
    }
  }
}
```

**Note**: In custom URL schemes, the format `scheme://host?query` is used, where:
- `lightshare` = scheme
- `verify-email` = host (not path)
- `token=xxx` = query parameters

## Authentication Flow

### Email Verification with Auto-Login

1. **User signs up** → Backend generates verification token
2. **Email sent** with link: `lightshare://verify-email?token=abc123`
3. **User clicks link** → OS opens app via deep link
4. **App receives URI** → Routes to `/auth/verify-email?token=abc123`
5. **EmailVerificationScreen**:
   - Calls `POST /api/v1/auth/verify-email` with token
   - Backend verifies email, marks user as verified
   - Backend generates JWT token pair
   - Returns: `{ user, access_token, refresh_token, expires_at, token_type }`
6. **AuthService**:
   - Stores tokens in secure storage
   - Updates auth state to authenticated
7. **Screen redirects** → Home dashboard (`/`)
8. **User is logged in** ✅

### Magic Link Login

Same flow as email verification, but uses:
- Endpoint: `POST /api/v1/auth/magic-link/verify`
- Deep link: `lightshare://magic-link?token=xyz789`

## Testing

### Testing on iOS Simulator

```bash
# Terminal command to simulate deep link
xcrun simctl openurl booted "lightshare://verify-email?token=test-token-123"
```

### Testing on Android Emulator

```bash
# Terminal command to simulate deep link
adb shell am start -W -a android.intent.action.VIEW \
  -d "lightshare://verify-email?token=test-token-123" \
  com.lightshare.app
```

### Testing on Physical Devices

#### Option 1: Send Test Email
1. Use a real email address during signup
2. Check inbox for verification email
3. Click the link on your phone
4. App should open automatically

#### Option 2: Create Test Link
1. Create a note/message with the deep link URL
2. Send to yourself (iMessage, WhatsApp, etc.)
3. Tap the link
4. App should open

#### Option 3: Browser Test
1. Open Safari/Chrome on the device
2. Type the deep link URL in the address bar
3. Press enter
4. iOS/Android will prompt to open the app

### Testing in Development

**Local Backend Setup**:
```bash
# Ensure backend is running
cd backend
go run cmd/server/main.go

# Backend should be accessible at:
# - iOS Simulator: http://localhost:8080
# - Android Emulator: http://10.0.2.2:8080
# - Physical Device: http://<your-local-ip>:8080
```

**Mobile App Setup**:
```bash
cd mobile

# For iOS simulator (uses localhost)
flutter run -d "iPhone 15 Pro"

# For Android emulator (needs special IP)
flutter run -d emulator-5554 --dart-define=API_BASE_URL=http://10.0.2.2:8080

# For physical device (use your machine's local IP)
flutter run -d <device-id> --dart-define=API_BASE_URL=http://192.168.1.XXX:8080
```

## Debugging

### Check Deep Link Registration

**iOS**:
```bash
# Check if app is registered for the scheme
xcrun simctl openurl booted "lightshare://test"
# Should open the app or show "No application knows how to open URL"
```

**Android**:
```bash
# Check intent filters
adb shell dumpsys package com.lightshare.app | grep -A 10 "scheme"
```

### Common Issues

#### Issue: Link Opens in Browser Instead of App

**Cause**: Deep link not properly registered

**Solution**:
1. Verify AndroidManifest.xml has the intent-filter
2. Verify Info.plist has CFBundleURLTypes
3. Rebuild and reinstall the app
4. Clear app data and restart

#### Issue: App Opens but Doesn't Navigate

**Cause**: Deep link handler not working

**Solution**:
1. Check console logs for "Deep link received: ..."
2. Verify the deep link format matches expected pattern
3. Check router configuration for the route paths

#### Issue: "Token expired" Error

**Cause**: Verification/magic link tokens expire after 24h/15min

**Solution**:
1. Request a new verification email
2. Use the new link immediately
3. Check backend logs for token validation errors

#### Issue: Auto-login Not Working

**Cause**: Backend not returning JWT tokens

**Solution**:
1. Check backend endpoint `/api/v1/auth/verify-email` response
2. Verify response includes: `access_token`, `refresh_token`
3. Check secure storage for saved tokens
4. Review AuthService._storeTokens() method

## Security Considerations

### Token Validation

- Verification tokens expire after 24 hours
- Magic link tokens expire after 15 minutes
- Tokens are single-use (cleared after verification)
- Backend validates token before issuing JWT

### URL Scheme Security

⚠️ **Custom URL schemes can be hijacked** by other apps on the device. For production, consider:

1. **Universal Links (iOS)** / **App Links (Android)**:
   - Requires verified domain (e.g., `https://app.lightshare.com`)
   - More secure than custom schemes
   - Fallback to web if app not installed

2. **Hybrid Approach**:
   - Development: `lightshare-dev://`
   - Production: `https://app.lightshare.com` with associated domains

### Token Storage

- JWT tokens stored in `flutter_secure_storage`
- Uses platform-native encryption (Keychain on iOS, EncryptedSharedPreferences on Android)
- Tokens never logged or exposed in URLs

## Future Enhancements

### Universal Links / App Links

For production, implement verified domain-based deep linking:

**iOS Associated Domains**:
```json
// .well-known/apple-app-site-association
{
  "applinks": {
    "apps": [],
    "details": [{
      "appID": "TEAM_ID.com.lightshare.app",
      "paths": ["/verify-email", "/magic-link"]
    }]
  }
}
```

**Android App Links**:
```json
// .well-known/assetlinks.json
[{
  "relation": ["delegate_permission/common.handle_all_urls"],
  "target": {
    "namespace": "android_app",
    "package_name": "com.lightshare.app",
    "sha256_cert_fingerprints": ["..."]
  }
}]
```

### Deferred Deep Linking

Handle scenarios where:
1. User clicks link but doesn't have app installed
2. User installs app
3. App opens and completes the verification

This requires additional attribution tracking.

## References

- [app_links package](https://pub.dev/packages/app_links)
- [iOS Universal Links](https://developer.apple.com/ios/universal-links/)
- [Android App Links](https://developer.android.com/training/app-links)
- [Flutter Deep Linking](https://docs.flutter.dev/ui/navigation/deep-linking)

## Related Documentation

- [Email Configuration](./backend-email-configuration.md)
- [Authentication Flow](./authentication-flow.md)
- [Mobile App Configuration](./mobile-configuration.md)
