# Phase 3 Implementation Fixes

## Issue 1: Build Errors - Missing Theme Colors and Widget Parameters

### Problem
- `AppTheme.neonBlue` color was referenced but not defined in the theme
- `GradientButton` widget was called with an `icon` parameter that doesn't exist

### Solution
- Replaced all `AppTheme.neonBlue` references with `AppTheme.primaryPurple`
- Removed invalid `icon` parameter from `GradientButton` usage

### Files Modified
- `mobile/lib/features/providers/screens/accounts_screen.dart`
- `mobile/lib/features/providers/screens/provider_selection_screen.dart`
- `mobile/lib/features/providers/screens/token_entry_screen.dart`

---

## Issue 2: Infinite API Request Loop (Critical)

### Problem
When navigating to the accounts page, the app would spam the backend with repeated refresh token requests, creating an infinite loop:

```
POST /api/v1/auth/refresh -> 401
POST /api/v1/auth/refresh -> 401
POST /api/v1/auth/refresh -> 401
... (repeats indefinitely)
```

### Root Cause
The `ApiClient` interceptor had a critical flaw:
1. When any request received a 401 response, the interceptor would call `_refreshToken()`
2. The `_refreshToken()` method made a POST request to `/api/v1/auth/refresh` using `_dio.post()`
3. This refresh request would **also trigger the same interceptor**
4. If the refresh token was invalid/expired, it would get a 401
5. The 401 would trigger another refresh attempt
6. Loop continues infinitely

### Solution
Implemented three-layer protection against infinite loops:

1. **Skip interceptor for refresh endpoint**: Modified the `onRequest` handler to skip adding auth headers for the refresh endpoint
2. **Prevent recursive refresh attempts**: Added `_isRefreshing` flag to prevent concurrent refresh operations
3. **Clear tokens on refresh failure**: When refresh fails, clear all tokens to force re-login

### Code Changes

**File: `mobile/lib/core/services/api_client.dart`**

```dart
class ApiClient {
  bool _isRefreshing = false; // Added flag

  ApiClient(...) {
    _dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          // Skip auth header for refresh endpoint to prevent infinite loop
          if (options.path == '/api/v1/auth/refresh') {
            return handler.next(options);
          }
          // ... rest of onRequest
        },
        onError: (error, handler) async {
          // Don't try to refresh if:
          // 1. Already refreshing
          // 2. The failing request is the refresh endpoint itself
          if (error.response?.statusCode == 401 &&
              error.requestOptions.path != '/api/v1/auth/refresh' &&
              !_isRefreshing) {
            // ... refresh logic
          }
        },
      ),
    );
  }

  Future<bool> _refreshToken() async {
    // Prevent concurrent refresh attempts
    if (_isRefreshing) return false;
    _isRefreshing = true;

    try {
      // ... refresh logic
    } catch (e) {
      // Clear tokens on failure to force re-login
      await _secureStorage.delete(key: 'access_token');
      await _secureStorage.delete(key: 'refresh_token');
      return false;
    } finally {
      _isRefreshing = false;
    }
  }
}
```

### Additional Changes

**File: `mobile/lib/features/providers/screens/accounts_screen.dart`**

Added error handling in `initState` to gracefully handle authentication failures:

```dart
Future.microtask(() async {
  try {
    await ref.read(accountsProvider.notifier).loadAccounts();
  } catch (e) {
    // Handle error silently - will be shown in UI via error state
  }
});
```

### Testing
- User can navigate to accounts page without causing infinite API requests
- Invalid/expired tokens are properly cleared and force re-login
- Valid token refresh works correctly
- No more 401 spam in backend logs

---

---

## Issue 3: Navigation Stack Error - "There is nothing to pop"

### Problem
After successfully connecting a provider, navigating back from the accounts screen would crash with:
```
GoError: There is nothing to pop
```

### Root Cause
The token entry screen used `context.go('/accounts')` after successful connection, which **replaces the entire navigation stack** instead of popping back to the previous screen.

**Navigation flow:**
1. Home (/) → Push accounts (/accounts) - Stack: `[/, /accounts]`
2. Accounts → Push provider selection - Stack: `[/, /accounts, /providers/connect]`
3. Provider selection → Push token entry - Stack: `[/, /accounts, /providers/connect, /providers/connect/token]`
4. Success → **`context.go('/accounts')`** - Stack: `[/accounts]` ⚠️ Stack replaced!
5. Back button tries to pop → **CRASH** (nothing to pop)

### Solution

**File: `mobile/lib/features/providers/screens/token_entry_screen.dart`**

Changed from using `context.go()` to properly popping back:
```dart
// Before:
context.go('/accounts');  // Replaces entire stack

// After:
context.pop();  // Pop token entry screen
context.pop();  // Pop provider selection screen
// Now back at accounts screen with proper stack
```

**File: `mobile/lib/features/providers/screens/accounts_screen.dart`**

Made back button defensive to handle edge cases:
```dart
// Before:
onPressed: () => context.pop(),

// After:
onPressed: () {
  if (context.canPop()) {
    context.pop();
  } else {
    context.go('/');  // Fallback to home
  }
},
```

### Testing
- ✅ Navigate from Home → Accounts → Connect Provider → Enter Token → Success
- ✅ After success, back button on accounts screen works correctly
- ✅ Navigation stack is preserved properly
- ✅ No crashes on back button press

---

## Issue 4: Backend Middleware/Handler Key Mismatch

### Problem
Provider endpoints were returning 401 "unauthorized" even with valid authentication because handlers couldn't find the user ID.

### Root Cause
Auth middleware and handlers used **different context keys**:
- **Middleware** (auth.go:48): `c.Locals("user_id", claims.UserID)` (underscore)
- **Handlers** (provider.go): `c.Locals("userID")` (camelCase)

### Solution

**File: `backend/internal/handlers/provider.go`**

Fixed all three handlers to use consistent key:
```go
// Before:
userID, ok := c.Locals("userID").(uuid.UUID)

// After:
userID, ok := c.Locals("user_id").(uuid.UUID)
```

### Testing
- ✅ POST `/api/v1/providers/connect` works with valid auth
- ✅ GET `/api/v1/accounts` returns user's accounts
- ✅ DELETE `/api/v1/accounts/:id` successfully disconnects accounts

---

## Prevention Guidelines

### For Future API Client Development

1. **Never trigger interceptors recursively**: Authentication/refresh requests should bypass the same interceptor
2. **Use flags for stateful operations**: Prevent concurrent refresh/retry operations
3. **Fail gracefully**: Clear invalid state and force re-authentication rather than retry indefinitely
4. **Test token expiration**: Always test with expired/invalid tokens to catch infinite loops
5. **Monitor logs**: Backend request logs should never show repeated identical requests

### Testing Checklist for Token Refresh

- [ ] Valid token refresh succeeds
- [ ] Expired access token with valid refresh token refreshes successfully
- [ ] Expired refresh token clears tokens and redirects to login
- [ ] Invalid tokens don't cause infinite loops
- [ ] Concurrent requests don't trigger multiple refresh attempts
- [ ] Backend logs show clean request patterns
