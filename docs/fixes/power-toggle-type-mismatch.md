# Fix: Power Toggle Type Mismatch

## Issue
When attempting to toggle a light's power state (on/off), the API returned a 400 error:
```
"error":"missing or invalid 'state' parameter (must be string)","status":400
```

## Root Cause
Type mismatch between the mobile app and backend API:

- **Backend Expected**: `state` parameter as a string (`"on"` or `"off"`)
  - See `backend/internal/models/action.go:60-69`
- **Mobile Sent**: `state` parameter as a boolean (`true` or `false`)
  - See `mobile/lib/core/models/action_request.dart:35`

## Solution
Modified the `ActionRequest.power()` factory method to convert the boolean state to a string:

**Before:**
```dart
parameters: {
  'state': state,  // boolean
  'duration': duration,
}
```

**After:**
```dart
parameters: {
  'state': state ? 'on' : 'off',  // string
  'duration': duration,
}
```

## Additional Fix
Also standardized the `temperature` action to send `kelvin` as a double for consistency with other numeric parameters:

**Before:**
```dart
parameters: {
  'kelvin': kelvin,  // int
  'duration': duration,
}
```

**After:**
```dart
parameters: {
  'kelvin': kelvin.toDouble(),  // double
  'duration': duration,
}
```

## Files Modified
- `mobile/lib/core/models/action_request.dart`

## Testing
- ✅ Flutter analyzer: No issues
- ✅ All Flutter tests pass (8/8)

## Impact
- Power toggle functionality now works correctly
- Type consistency improved across action requests
