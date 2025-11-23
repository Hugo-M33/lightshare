# Phase 4: Light Control - Implementation Tracker

**Status**: 🚧 In Progress
**Started**: 2025-11-23
**Branch**: `claude/phase-4-lights-control-015aGAVRbhPfvEMLYVczoAYe`

---

## Backend Implementation

### 1. Data Models ⏳ In Progress
- [ ] Create `backend/internal/models/device.go` with Device, DeviceColor, DeviceGroup, DeviceLocation
- [ ] Create `backend/internal/models/action.go` with ActionRequest and action types
- [ ] Add JSON tags and validation
- [ ] Add helper methods (IsValidAction, ValidateParameters)
- [ ] Write model unit tests

### 2. Provider Interface Extensions ⏳ In Progress
- [ ] Extend `backend/pkg/providers/provider.go` Client interface
- [ ] Add Device struct to provider package
- [ ] Add effect methods (Pulse, Breathe)
- [ ] Document selector pattern
- [ ] Add provider-specific error types

### 3. LIFX Client - Device Listing
- [ ] Implement `ListDevices()` in `backend/pkg/providers/lifx/client.go`
- [ ] Implement `GetDevice()` method
- [ ] Map LIFX response to unified Device struct
- [ ] Extract capabilities from product info
- [ ] Handle connected/reachable status
- [ ] Add unit tests with mock HTTP responses

### 4. LIFX Client - Control Actions
- [ ] Implement `SetPower()` method
- [ ] Implement `SetBrightness()` method
- [ ] Implement `SetColor()` method
- [ ] Implement `SetColorTemperature()` method
- [ ] Add parameter validation (ranges, required fields)
- [ ] Handle LIFX color format conversion

### 5. LIFX Client - Effects
- [ ] Implement `Pulse()` effect method
- [ ] Implement `Breathe()` effect method
- [ ] Add effect parameter validation
- [ ] Add unit tests for effects

### 6. Configuration Updates
- [ ] Add `DeviceCacheTTL` to `backend/internal/config/config.go`
- [ ] Add `RateLimitPerMin` to config
- [ ] Update config loading with defaults
- [ ] Document new env vars in README

### 7. Device Service Layer
- [ ] Create `backend/internal/services/device.go`
- [ ] Implement `ListDevices()` with caching
- [ ] Implement `GetDevice()` with caching
- [ ] Implement `ExecuteAction()` with cache invalidation
- [ ] Implement `RefreshDevices()`
- [ ] Add Redis caching helpers (get/set/invalidate)
- [ ] Add rate limiting with Redis
- [ ] Add access control checks (owner + granted users)
- [ ] Write service tests with mocked repos and cache

### 8. Device Handler & Routes
- [ ] Create `backend/internal/handlers/device.go`
- [ ] Implement `ListDevices` endpoint handler
- [ ] Implement `ListAccountDevices` endpoint handler
- [ ] Implement `GetDevice` endpoint handler
- [ ] Implement `ExecuteAction` endpoint handler
- [ ] Implement `RefreshDevices` endpoint handler
- [ ] Add request validation middleware
- [ ] Add rate limiting middleware
- [ ] Register routes in `backend/cmd/server/main.go`
- [ ] Write handler tests

### 9. Backend Testing
- [ ] Unit tests for LIFX client methods
- [ ] Unit tests for device service
- [ ] Unit tests for device handlers
- [ ] Integration test: list devices flow
- [ ] Integration test: execute action flow
- [ ] Test rate limiting behavior
- [ ] Test caching behavior

---

## Mobile Implementation

### 10. Data Models
- [ ] Create `mobile/lib/core/models/device.dart` with Freezed
- [ ] Create DeviceColor model
- [ ] Create DeviceGroup model
- [ ] Create DeviceLocation model
- [ ] Add JSON serialization
- [ ] Add helper methods (hasCapability, isOn)
- [ ] Run code generation: `flutter pub run build_runner build`
- [ ] Write model tests

### 11. Device Service
- [ ] Create `mobile/lib/core/services/device_service.dart`
- [ ] Implement `listDevices()` method
- [ ] Implement `listAccountDevices()` method
- [ ] Implement `getDevice()` method
- [ ] Implement `executeAction()` method
- [ ] Implement convenience methods (setPower, setBrightness, setColor)
- [ ] Implement `refreshDevices()` method
- [ ] Add error handling with DioException
- [ ] Write service tests with mocked Dio

### 12. State Management
- [ ] Create `mobile/lib/core/providers/devices_provider.dart`
- [ ] Create DevicesState with Freezed
- [ ] Create DevicesNotifier class
- [ ] Implement `loadDevices()` method
- [ ] Implement `refreshDevices()` method
- [ ] Implement `togglePower()` with optimistic update
- [ ] Implement `setBrightness()` (debounced)
- [ ] Implement `updateBrightnessLocally()` (instant UI)
- [ ] Implement `setColor()` (debounced)
- [ ] Implement `updateColorLocally()` (instant UI)
- [ ] Create helper providers (devicesByLocation, etc.)
- [ ] Run code generation
- [ ] Write provider tests

### 13. UI - Devices List Screen
- [ ] Create `mobile/lib/features/devices/screens/devices_screen.dart`
- [ ] Create `DeviceCard` widget
- [ ] Implement device list with grouping by location
- [ ] Add pull-to-refresh functionality
- [ ] Add loading shimmer effect
- [ ] Add empty state UI
- [ ] Add error display with retry
- [ ] Add navigation to device detail
- [ ] Style with glassmorphism theme
- [ ] Write widget tests

### 14. UI - Device Detail Screen
- [ ] Create `mobile/lib/features/devices/screens/device_detail_screen.dart`
- [ ] Create DeviceHeader widget
- [ ] Create PowerCard widget
- [ ] Create BrightnessCard widget
- [ ] Create ColorCard widget (if capable)
- [ ] Create TemperatureCard widget (if capable)
- [ ] Create EffectsCard widget (LIFX only)
- [ ] Add capability-based rendering
- [ ] Add haptic feedback
- [ ] Style with glassmorphism theme
- [ ] Write widget tests

### 15. UI - Control Widgets
- [ ] Create `mobile/lib/features/devices/widgets/power_toggle.dart`
- [ ] Create `brightness_slider.dart` with debouncing (500ms)
- [ ] Create `color_picker.dart` with debouncing (500ms)
- [ ] Create `temperature_slider.dart` with debouncing (500ms)
- [ ] Create `effect_button.dart` (Pulse, Breathe)
- [ ] Create `device_status_indicator.dart`
- [ ] Implement Timer-based debouncing
- [ ] Add haptic feedback
- [ ] Add visual feedback for pending API calls
- [ ] Test rapid slider movements only trigger one API call
- [ ] Write widget tests

### 16. Navigation & Routing
- [ ] Add `/devices` route to router
- [ ] Add `/devices/:id` route for device detail
- [ ] Update home screen with "My Lights" navigation
- [ ] Test navigation flow

### 17. Mobile Testing
- [ ] Unit tests for device models
- [ ] Unit tests for device service
- [ ] Unit tests for devices provider
- [ ] Widget tests for DevicesScreen
- [ ] Widget tests for DeviceDetailScreen
- [ ] Widget tests for control widgets
- [ ] Integration test: Login → View devices
- [ ] Integration test: Toggle power
- [ ] Integration test: Adjust brightness

---

## Documentation & Polish

### 18. API Documentation
- [ ] Update `docs/api.md` with device endpoints
- [ ] Add request/response examples
- [ ] Document selector pattern
- [ ] Document action types and parameters
- [ ] Add error codes

### 19. Implementation Guide
- [ ] Create `docs/phase-4-implementation.md`
- [ ] Document architecture
- [ ] Add setup instructions
- [ ] Add usage guide
- [ ] Add troubleshooting section
- [ ] Document known limitations

### 20. Final Polish
- [ ] Run backend linting: `golangci-lint run`
- [ ] Run mobile linting: `flutter analyze`
- [ ] Run all backend tests: `go test ./...`
- [ ] Run all mobile tests: `flutter test`
- [ ] Fix any failing tests
- [ ] Update README.md with Phase 4 info
- [ ] Manual QA with real LIFX account
- [ ] Performance testing (check API rate limits)

---

## Progress Summary

**Total Tasks**: ~110
**Completed**: 0
**In Progress**: 2
**Remaining**: 108

**Current Focus**: Backend models and provider interface extensions

---

## Notes

- Backend and mobile can be developed in parallel
- Test incrementally, don't wait until the end
- Commit frequently with logical groupings
- Update this file as tasks are completed

---

## Environment Setup Checklist

Backend:
- [ ] Set `DEVICE_CACHE_TTL_SECONDS=60` (optional, has default)
- [ ] Set `RATE_LIMIT_PER_MIN=30` (optional, has default)
- [ ] Ensure Redis is running for caching
- [ ] Run migrations (if any new ones)

Mobile:
- [ ] Update API base URL if needed
- [ ] Run `flutter pub get`
- [ ] Run `flutter pub run build_runner build --delete-conflicting-outputs`

Testing:
- [ ] Get LIFX personal access token for testing
- [ ] Connect at least one LIFX device to test account
