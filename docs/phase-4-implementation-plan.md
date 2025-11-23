# Phase 4 Implementation Plan: Light Control

## Overview

Phase 4 implements the core functionality of LightShare - viewing and controlling smart lights through the app. This phase builds on Phase 3's provider connection system by adding device discovery and control capabilities.

**Status**: 🚧 Planning Phase

**Dependencies**: Phase 3 (Provider Connection) ✅ Complete

**Estimated Timeline**: This is a multi-step implementation with backend and mobile components that need to be developed in parallel.

---

## Goals

1. **Device Discovery**: List all lights from connected LIFX accounts
2. **Basic Controls**: Power on/off, brightness adjustment
3. **Advanced Controls**: Color selection, effects (LIFX-specific)
4. **Mobile UI**: Intuitive device list and control interface
5. **State Management**: Real-time state updates with optimistic UI
6. **Error Handling**: Graceful handling of offline devices and rate limits
7. **Caching**: Redis-based device state caching

---

## Architecture Overview

### Backend Flow
```
Mobile App → Backend API → Provider Client → LIFX/Hue API
               ↓
          Redis Cache
               ↓
          PostgreSQL (account tokens)
```

### Key Principles
- **Backend Proxy Pattern**: All provider API calls go through backend
- **Token Security**: Provider tokens never exposed to clients
- **Optimistic UI**: Mobile updates UI immediately, syncs with backend
- **Caching**: Device lists cached with short TTL (30-60 seconds)
- **Rate Limiting**: Per-account rate limiting to respect provider limits

---

## Implementation Tasks

### 1. Backend - Data Models & Types

#### 1.1 Device Model
**Location**: `backend/internal/models/device.go`

**Create**:
```go
// Device represents a smart light device
type Device struct {
    ID           string                 `json:"id"`           // Provider-specific device ID
    AccountID    string                 `json:"account_id"`   // Our account UUID
    Provider     string                 `json:"provider"`     // "lifx" or "hue"
    Label        string                 `json:"label"`        // User-friendly name
    Power        string                 `json:"power"`        // "on" or "off"
    Brightness   float64                `json:"brightness"`   // 0.0 - 1.0
    Color        *DeviceColor           `json:"color,omitempty"`
    Connected    bool                   `json:"connected"`
    Reachable    bool                   `json:"reachable"`
    Group        *DeviceGroup           `json:"group,omitempty"`
    Location     *DeviceLocation        `json:"location,omitempty"`
    Capabilities []string               `json:"capabilities"` // ["color", "temperature", "effects"]
    Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type DeviceColor struct {
    Hue        float64 `json:"hue"`        // 0-360
    Saturation float64 `json:"saturation"` // 0.0-1.0
    Kelvin     int     `json:"kelvin"`     // 1500-9000 (for white balance)
}

type DeviceGroup struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type DeviceLocation struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

#### 1.2 Control Action Types
**Location**: `backend/internal/models/action.go`

**Create**:
```go
// ActionRequest represents a control action request
type ActionRequest struct {
    Action     string                 `json:"action" validate:"required"` // "power", "brightness", "color", "effect"
    Parameters map[string]interface{} `json:"parameters" validate:"required"`
}

// Supported actions:
// - "power": {"state": "on"|"off"}
// - "brightness": {"level": 0.0-1.0, "duration": seconds}
// - "color": {"hue": 0-360, "saturation": 0-1, "brightness": 0-1, "duration": seconds}
// - "temperature": {"kelvin": 1500-9000, "duration": seconds}
// - "effect": {"name": "pulse"|"breathe", "color": {...}, "duration": seconds}
```

**Tasks**:
- [ ] Create `device.go` model with JSON tags
- [ ] Create `action.go` with action request/response types
- [ ] Add validation tags for all fields
- [ ] Add helper methods (e.g., `IsValidAction()`, `ValidateParameters()`)

---

### 2. Backend - Extended Provider Interface

#### 2.1 Update Provider Client Interface
**Location**: `backend/pkg/providers/provider.go`

**Modify Interface**:
```go
type Client interface {
    ValidateToken(token string) (*AccountInfo, error)
    GetAccountInfo(token string) (*AccountInfo, error)

    // New methods for Phase 4
    ListDevices(token string) ([]*Device, error)
    GetDevice(token, deviceID string) (*Device, error)
    SetPower(token, selector string, state bool, duration float64) error
    SetBrightness(token, selector string, level float64, duration float64) error
    SetColor(token, selector string, color *DeviceColor, duration float64) error
    SetColorTemperature(token, selector string, kelvin int, duration float64) error
}

// Selector format: "all", "id:d073d5", "group_id:xxx", "location_id:xxx"
```

**Tasks**:
- [ ] Extend `Client` interface with device and control methods
- [ ] Define `Device` struct in provider package
- [ ] Add selector pattern documentation
- [ ] Add error types for provider-specific errors

#### 2.2 LIFX Client - Device Listing
**Location**: `backend/pkg/providers/lifx/client.go`

**Implement**:
```go
// ListDevices returns all lights for the account
func (c *Client) ListDevices(token string) ([]*Device, error) {
    // GET https://api.lifx.com/v1/lights/all
    // Parse response and convert to Device structs
    // Handle pagination if needed
}

// GetDevice returns a specific light
func (c *Client) GetDevice(token, deviceID string) (*Device, error) {
    // GET https://api.lifx.com/v1/lights/id:{deviceID}
}
```

**LIFX API Response Example**:
```json
[
  {
    "id": "d073d5001234",
    "uuid": "8fa5f072-af97-44ed-ae54-e70fd7bd9d20",
    "label": "Living Room",
    "connected": true,
    "power": "on",
    "brightness": 0.5,
    "color": {
      "hue": 120.0,
      "saturation": 1.0,
      "kelvin": 3500
    },
    "group": {
      "id": "group123",
      "name": "Living Room"
    },
    "location": {
      "id": "location123",
      "name": "Home"
    },
    "product": {
      "name": "LIFX Color",
      "capabilities": {
        "has_color": true,
        "has_variable_color_temp": true
      }
    }
  }
]
```

**Tasks**:
- [ ] Implement `ListDevices()` method
- [ ] Implement `GetDevice()` method
- [ ] Map LIFX response to unified Device struct
- [ ] Extract capabilities from product info
- [ ] Handle connected/reachable status
- [ ] Add unit tests with mock HTTP responses

#### 2.3 LIFX Client - Control Actions
**Location**: `backend/pkg/providers/lifx/client.go`

**Implement**:
```go
// SetPower turns lights on or off
// PUT https://api.lifx.com/v1/lights/{selector}/state
func (c *Client) SetPower(token, selector string, state bool, duration float64) error

// SetBrightness adjusts brightness (0.0-1.0)
func (c *Client) SetBrightness(token, selector string, level float64, duration float64) error

// SetColor sets hue, saturation, brightness
func (c *Client) SetColor(token, selector string, color *DeviceColor, duration float64) error

// SetColorTemperature sets white balance (1500-9000K)
func (c *Client) SetColorTemperature(token, selector string, kelvin int, duration float64) error

// Pulse creates a pulse effect
func (c *Client) Pulse(token, selector string, color *DeviceColor, cycles int, period float64) error

// Breathe creates a breathing effect
func (c *Client) Breathe(token, selector string, color *DeviceColor, cycles int, period float64) error
```

**LIFX State API**:
```bash
PUT https://api.lifx.com/v1/lights/all/state
{
  "power": "on",
  "brightness": 0.5,
  "color": "hue:120 saturation:1.0",
  "duration": 1.0  // fade time in seconds
}
```

**Tasks**:
- [ ] Implement all control methods
- [ ] Add parameter validation (ranges, required fields)
- [ ] Handle LIFX-specific color format conversion
- [ ] Implement effects (Pulse, Breathe)
- [ ] Add rate limiting awareness
- [ ] Add comprehensive error handling
- [ ] Write unit tests for each control method

---

### 3. Backend - Device Service Layer

#### 3.1 Device Service
**Location**: `backend/internal/services/device.go`

**Create Service**:
```go
type DeviceService struct {
    accountRepo *repository.AccountRepository
    cache       *redis.Client
    cacheTTL    time.Duration
}

func NewDeviceService(accountRepo *repository.AccountRepository, cache *redis.Client) *DeviceService {
    return &DeviceService{
        accountRepo: accountRepo,
        cache:       cache,
        cacheTTL:    60 * time.Second, // 1 minute cache
    }
}

// ListDevices returns all devices for a user's accounts
func (s *DeviceService) ListDevices(ctx context.Context, userID string) ([]*models.Device, error) {
    // 1. Get all accounts for user
    // 2. Check cache for each account's devices
    // 3. If cache miss, fetch from provider and cache
    // 4. Merge devices from all accounts
    // 5. Return sorted by location/group
}

// GetDevice returns a specific device
func (s *DeviceService) GetDevice(ctx context.Context, userID, accountID, deviceID string) (*models.Device, error) {
    // 1. Verify user owns account
    // 2. Check cache
    // 3. If miss, fetch from provider
    // 4. Return device
}

// ExecuteAction performs a control action on device(s)
func (s *DeviceService) ExecuteAction(ctx context.Context, userID, accountID string, selector string, action *models.ActionRequest) error {
    // 1. Verify user owns account (or has grant access)
    // 2. Get decrypted token
    // 3. Get provider client
    // 4. Execute action via provider client
    // 5. Invalidate cache for affected devices
    // 6. Return result
}

// RefreshDevices forces a cache refresh for an account
func (s *DeviceService) RefreshDevices(ctx context.Context, userID, accountID string) ([]*models.Device, error) {
    // 1. Verify ownership
    // 2. Fetch from provider
    // 3. Update cache
    // 4. Return devices
}
```

**Caching Strategy**:
```go
// Cache key format: "devices:account:{account_id}"
// Cache TTL: 60 seconds (balance between freshness and API rate limits)

func (s *DeviceService) getCachedDevices(accountID string) ([]*models.Device, error) {
    key := fmt.Sprintf("devices:account:%s", accountID)
    data, err := s.cache.Get(context.Background(), key).Bytes()
    // ... unmarshal and return
}

func (s *DeviceService) setCachedDevices(accountID string, devices []*models.Device) error {
    key := fmt.Sprintf("devices:account:%s", accountID)
    data, _ := json.Marshal(devices)
    return s.cache.Set(context.Background(), key, data, s.cacheTTL).Err()
}

func (s *DeviceService) invalidateCache(accountID string) error {
    key := fmt.Sprintf("devices:account:%s", accountID)
    return s.cache.Del(context.Background(), key).Err()
}
```

**Tasks**:
- [ ] Create `DeviceService` struct with dependencies
- [ ] Implement `ListDevices()` with caching
- [ ] Implement `GetDevice()` with caching
- [ ] Implement `ExecuteAction()` with cache invalidation
- [ ] Implement `RefreshDevices()`
- [ ] Add access control checks (owner + granted users)
- [ ] Add rate limiting per account
- [ ] Write service tests with mocked repos and cache

---

### 4. Backend - API Endpoints

#### 4.1 Device Handler
**Location**: `backend/internal/handlers/device.go`

**Create Endpoints**:
```go
// GET /api/v1/devices
// List all devices across all user's connected accounts
func (h *DeviceHandler) ListDevices(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    devices, err := h.deviceService.ListDevices(c.Context(), userID)
    // Return devices grouped by account/location
}

// GET /api/v1/accounts/:accountId/devices
// List devices for a specific account
func (h *DeviceHandler) ListAccountDevices(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    accountID := c.Params("accountId")
    devices, err := h.deviceService.ListDevices(c.Context(), userID, accountID)
    // Return devices for specific account
}

// GET /api/v1/accounts/:accountId/devices/:deviceId
// Get a specific device
func (h *DeviceHandler) GetDevice(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    accountID := c.Params("accountId")
    deviceID := c.Params("deviceId")
    device, err := h.deviceService.GetDevice(c.Context(), userID, accountID, deviceID)
    // Return device details
}

// POST /api/v1/accounts/:accountId/devices/:selector/action
// Execute a control action on device(s)
// selector can be: "all", "id:abc123", "group_id:xxx", "location_id:xxx"
func (h *DeviceHandler) ExecuteAction(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    accountID := c.Params("accountId")
    selector := c.Params("selector")

    var action models.ActionRequest
    if err := c.BodyParser(&action); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
    }

    err := h.deviceService.ExecuteAction(c.Context(), userID, accountID, selector, &action)
    // Return success/failure
}

// POST /api/v1/accounts/:accountId/devices/refresh
// Force refresh device list (invalidate cache)
func (h *DeviceHandler) RefreshDevices(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    accountID := c.Params("accountId")
    devices, err := h.deviceService.RefreshDevices(c.Context(), userID, accountID)
    // Return refreshed devices
}
```

**Request/Response Examples**:

**List Devices**:
```bash
GET /api/v1/devices
Authorization: Bearer <jwt>

Response:
{
  "devices": [
    {
      "id": "d073d5001234",
      "account_id": "uuid",
      "provider": "lifx",
      "label": "Living Room",
      "power": "on",
      "brightness": 0.5,
      "color": {
        "hue": 120.0,
        "saturation": 1.0,
        "kelvin": 3500
      },
      "connected": true,
      "reachable": true,
      "group": {
        "id": "group123",
        "name": "Living Room"
      },
      "capabilities": ["color", "temperature", "effects"]
    }
  ]
}
```

**Execute Action**:
```bash
POST /api/v1/accounts/:accountId/devices/all/action
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "action": "power",
  "parameters": {
    "state": "on",
    "duration": 1.0
  }
}

Response:
{
  "success": true,
  "message": "Action executed successfully"
}
```

**Tasks**:
- [ ] Create `DeviceHandler` struct
- [ ] Implement all endpoint handlers
- [ ] Add request validation middleware
- [ ] Add rate limiting middleware (per-account)
- [ ] Add comprehensive error handling
- [ ] Return proper HTTP status codes
- [ ] Write handler tests

#### 4.2 Router Registration
**Location**: `backend/cmd/server/main.go`

**Add Routes**:
```go
// Device routes (protected)
deviceHandler := handlers.NewDeviceHandler(deviceService)
api.Get("/devices", authMiddleware, deviceHandler.ListDevices)
api.Get("/accounts/:accountId/devices", authMiddleware, deviceHandler.ListAccountDevices)
api.Get("/accounts/:accountId/devices/:deviceId", authMiddleware, deviceHandler.GetDevice)
api.Post("/accounts/:accountId/devices/:selector/action", authMiddleware, deviceHandler.ExecuteAction)
api.Post("/accounts/:accountId/devices/refresh", authMiddleware, deviceHandler.RefreshDevices)
```

**Tasks**:
- [ ] Register device routes
- [ ] Apply auth middleware
- [ ] Apply rate limiting middleware
- [ ] Update API documentation

---

### 5. Mobile - Data Models

#### 5.1 Device Model
**Location**: `mobile/lib/core/models/device.dart`

**Create**:
```dart
import 'package:freezed_annotation/freezed_annotation.dart';

part 'device.freezed.dart';
part 'device.g.dart';

@freezed
class Device with _$Device {
  const factory Device({
    required String id,
    required String accountId,
    required String provider,
    required String label,
    required String power,
    required double brightness,
    DeviceColor? color,
    required bool connected,
    required bool reachable,
    DeviceGroup? group,
    DeviceLocation? location,
    required List<String> capabilities,
    Map<String, dynamic>? metadata,
  }) = _Device;

  factory Device.fromJson(Map<String, dynamic> json) => _$DeviceFromJson(json);
}

@freezed
class DeviceColor with _$DeviceColor {
  const factory DeviceColor({
    required double hue,
    required double saturation,
    required int kelvin,
  }) = _DeviceColor;

  factory DeviceColor.fromJson(Map<String, dynamic> json) => _$DeviceColorFromJson(json);
}

@freezed
class DeviceGroup with _$DeviceGroup {
  const factory DeviceGroup({
    required String id,
    required String name,
  }) = _DeviceGroup;

  factory DeviceGroup.fromJson(Map<String, dynamic> json) => _$DeviceGroupFromJson(json);
}

@freezed
class DeviceLocation with _$DeviceLocation {
  const factory DeviceLocation({
    required String id,
    required String name,
  }) = _DeviceLocation;

  factory DeviceLocation.fromJson(Map<String, dynamic> json) => _$DeviceLocationFromJson(json);
}
```

**Tasks**:
- [ ] Create device model classes with Freezed
- [ ] Add JSON serialization
- [ ] Add helper methods (e.g., `hasCapability()`, `isOn()`)
- [ ] Generate code with `flutter pub run build_runner build`

---

### 6. Mobile - Device Service

#### 6.1 Device API Service
**Location**: `mobile/lib/core/services/device_service.dart`

**Create**:
```dart
import 'package:dio/dio.dart';
import '../models/device.dart';

class DeviceService {
  final Dio _dio;

  DeviceService(this._dio);

  // List all devices
  Future<List<Device>> listDevices() async {
    final response = await _dio.get('/devices');
    final data = response.data as Map<String, dynamic>;
    final devices = data['devices'] as List;
    return devices.map((d) => Device.fromJson(d)).toList();
  }

  // List devices for specific account
  Future<List<Device>> listAccountDevices(String accountId) async {
    final response = await _dio.get('/accounts/$accountId/devices');
    final data = response.data as Map<String, dynamic>;
    final devices = data['devices'] as List;
    return devices.map((d) => Device.fromJson(d)).toList();
  }

  // Get single device
  Future<Device> getDevice(String accountId, String deviceId) async {
    final response = await _dio.get('/accounts/$accountId/devices/$deviceId');
    return Device.fromJson(response.data);
  }

  // Execute action
  Future<void> executeAction(
    String accountId,
    String selector,
    String action,
    Map<String, dynamic> parameters,
  ) async {
    await _dio.post(
      '/accounts/$accountId/devices/$selector/action',
      data: {
        'action': action,
        'parameters': parameters,
      },
    );
  }

  // Convenience methods
  Future<void> setPower(String accountId, String deviceId, bool state) async {
    await executeAction(accountId, 'id:$deviceId', 'power', {
      'state': state ? 'on' : 'off',
      'duration': 0.5,
    });
  }

  Future<void> setBrightness(String accountId, String deviceId, double level) async {
    await executeAction(accountId, 'id:$deviceId', 'brightness', {
      'level': level,
      'duration': 0.5,
    });
  }

  Future<void> setColor(String accountId, String deviceId, double hue, double saturation) async {
    await executeAction(accountId, 'id:$deviceId', 'color', {
      'hue': hue,
      'saturation': saturation,
      'brightness': 1.0,
      'duration': 0.5,
    });
  }

  // Refresh devices
  Future<List<Device>> refreshDevices(String accountId) async {
    final response = await _dio.post('/accounts/$accountId/devices/refresh');
    final data = response.data as Map<String, dynamic>;
    final devices = data['devices'] as List;
    return devices.map((d) => Device.fromJson(d)).toList();
  }
}
```

**Tasks**:
- [ ] Create `DeviceService` class
- [ ] Implement all API methods
- [ ] Add error handling with DioException
- [ ] Add request/response logging
- [ ] Add timeout configuration
- [ ] Write unit tests with mocked Dio

---

### 7. Mobile - State Management

#### 7.1 Devices Provider
**Location**: `mobile/lib/core/providers/devices_provider.dart`

**Create**:
```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/device.dart';
import '../services/device_service.dart';

// Devices state
@freezed
class DevicesState with _$DevicesState {
  const factory DevicesState({
    @Default([]) List<Device> devices,
    @Default(false) bool isLoading,
    @Default(false) bool isRefreshing,
    String? error,
    DateTime? lastUpdated,
  }) = _DevicesState;
}

// Devices notifier
class DevicesNotifier extends StateNotifier<DevicesState> {
  final DeviceService _deviceService;

  DevicesNotifier(this._deviceService) : super(const DevicesState());

  // Load all devices
  Future<void> loadDevices() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final devices = await _deviceService.listDevices();
      state = state.copyWith(
        devices: devices,
        isLoading: false,
        lastUpdated: DateTime.now(),
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  // Refresh devices (pull-to-refresh)
  Future<void> refreshDevices() async {
    state = state.copyWith(isRefreshing: true, error: null);

    try {
      // Get unique account IDs from current devices
      final accountIds = state.devices
          .map((d) => d.accountId)
          .toSet()
          .toList();

      // Refresh each account
      final List<Device> allDevices = [];
      for (final accountId in accountIds) {
        final devices = await _deviceService.refreshDevices(accountId);
        allDevices.addAll(devices);
      }

      state = state.copyWith(
        devices: allDevices,
        isRefreshing: false,
        lastUpdated: DateTime.now(),
      );
    } catch (e) {
      state = state.copyWith(
        isRefreshing: false,
        error: e.toString(),
      );
    }
  }

  // Optimistic power toggle
  Future<void> togglePower(String accountId, String deviceId) async {
    // Find device
    final deviceIndex = state.devices.indexWhere((d) => d.id == deviceId);
    if (deviceIndex == -1) return;

    final device = state.devices[deviceIndex];
    final newPowerState = device.power == 'on' ? 'off' : 'on';

    // Optimistic update
    final updatedDevices = List<Device>.from(state.devices);
    updatedDevices[deviceIndex] = device.copyWith(power: newPowerState);
    state = state.copyWith(devices: updatedDevices);

    // Execute action
    try {
      await _deviceService.setPower(accountId, deviceId, newPowerState == 'on');
      // Success - state already updated
    } catch (e) {
      // Revert on error
      updatedDevices[deviceIndex] = device;
      state = state.copyWith(devices: updatedDevices, error: e.toString());
    }
  }

  // Optimistic brightness change
  Future<void> setBrightness(String accountId, String deviceId, double level) async {
    final deviceIndex = state.devices.indexWhere((d) => d.id == deviceId);
    if (deviceIndex == -1) return;

    final device = state.devices[deviceIndex];

    // Optimistic update
    final updatedDevices = List<Device>.from(state.devices);
    updatedDevices[deviceIndex] = device.copyWith(brightness: level);
    state = state.copyWith(devices: updatedDevices);

    // Execute action
    try {
      await _deviceService.setBrightness(accountId, deviceId, level);
    } catch (e) {
      // Revert on error
      updatedDevices[deviceIndex] = device;
      state = state.copyWith(devices: updatedDevices, error: e.toString());
    }
  }

  // Set color
  Future<void> setColor(String accountId, String deviceId, double hue, double saturation) async {
    final deviceIndex = state.devices.indexWhere((d) => d.id == deviceId);
    if (deviceIndex == -1) return;

    final device = state.devices[deviceIndex];

    // Optimistic update
    final updatedDevices = List<Device>.from(state.devices);
    final newColor = DeviceColor(hue: hue, saturation: saturation, kelvin: device.color?.kelvin ?? 3500);
    updatedDevices[deviceIndex] = device.copyWith(color: newColor);
    state = state.copyWith(devices: updatedDevices);

    // Execute action
    try {
      await _deviceService.setColor(accountId, deviceId, hue, saturation);
    } catch (e) {
      // Revert on error
      updatedDevices[deviceIndex] = device;
      state = state.copyWith(devices: updatedDevices, error: e.toString());
    }
  }
}

// Provider
final devicesProvider = StateNotifierProvider<DevicesNotifier, DevicesState>((ref) {
  final deviceService = ref.watch(deviceServiceProvider);
  return DevicesNotifier(deviceService);
});

// Helper providers
final deviceServiceProvider = Provider<DeviceService>((ref) {
  final dio = ref.watch(dioProvider);
  return DeviceService(dio);
});

// Filtered devices by location
final devicesByLocationProvider = Provider<Map<String, List<Device>>>((ref) {
  final devices = ref.watch(devicesProvider).devices;
  final Map<String, List<Device>> grouped = {};

  for (final device in devices) {
    final locationName = device.location?.name ?? 'Ungrouped';
    grouped.putIfAbsent(locationName, () => []).add(device);
  }

  return grouped;
});
```

**Tasks**:
- [ ] Create `DevicesState` with Freezed
- [ ] Create `DevicesNotifier` class
- [ ] Implement load/refresh methods
- [ ] Implement optimistic UI updates for controls
- [ ] Add error handling and revert logic
- [ ] Create helper providers (filtered, grouped)
- [ ] Write provider tests

---

### 8. Mobile - UI Screens

#### 8.1 Devices List Screen
**Location**: `mobile/lib/features/devices/screens/devices_screen.dart`

**Features**:
- List all devices grouped by location
- Show power state, brightness, color indicator
- Quick power toggle
- Pull-to-refresh
- Navigate to device detail
- Empty state when no devices
- Loading state
- Error handling

**UI Structure**:
```dart
Scaffold
├─ AppBar ("My Lights")
├─ RefreshIndicator
│  └─ ListView/CustomScrollView
│     ├─ LocationHeader ("Living Room")
│     ├─ DeviceCard
│     │  ├─ Device icon (color indicator)
│     │  ├─ Label
│     │  ├─ Status (on/off, brightness %)
│     │  └─ Power toggle switch
│     ├─ DeviceCard
│     └─ ...
└─ FloatingActionButton (Refresh)
```

**Tasks**:
- [ ] Create `DevicesScreen` widget
- [ ] Create `DeviceCard` widget
- [ ] Implement pull-to-refresh
- [ ] Add loading shimmer effect
- [ ] Add empty state UI
- [ ] Add error display with retry
- [ ] Add navigation to device detail
- [ ] Match glassmorphism design theme

#### 8.2 Device Detail Screen
**Location**: `mobile/lib/features/devices/screens/device_detail_screen.dart`

**Features**:
- Device info (name, type, connection status)
- Large power toggle
- Brightness slider
- Color picker (if capable)
- Temperature slider (if capable)
- Effects buttons (LIFX only)
- Last updated timestamp

**UI Structure**:
```dart
Scaffold
├─ AppBar (Device label)
├─ SingleChildScrollView
│  ├─ DeviceHeader (icon, status, location)
│  ├─ PowerCard (large toggle)
│  ├─ BrightnessCard (slider 0-100%)
│  ├─ ColorCard (color wheel picker)
│  ├─ TemperatureCard (slider 1500-9000K)
│  └─ EffectsCard (Pulse, Breathe buttons)
└─ ...
```

**Tasks**:
- [ ] Create `DeviceDetailScreen` widget
- [ ] Create control card widgets (Power, Brightness, Color, etc.)
- [ ] Implement color picker widget
- [ ] Implement temperature slider
- [ ] Add effects buttons (LIFX)
- [ ] Add capability-based UI rendering
- [ ] Add haptic feedback on interactions
- [ ] Match glassmorphism design theme

#### 8.3 Device Control Widgets
**Location**: `mobile/lib/features/devices/widgets/`

**Create**:
- `power_toggle.dart` - Large power switch widget
- `brightness_slider.dart` - Brightness control with percentage
- `color_picker.dart` - Hue/saturation color wheel
- `temperature_slider.dart` - Color temperature slider
- `effect_button.dart` - Effect action button (Pulse, Breathe, etc.)
- `device_status_indicator.dart` - Connection/reachability indicator

**Tasks**:
- [ ] Create all control widgets
- [ ] Add smooth animations
- [ ] Add haptic feedback
- [ ] Add debouncing for sliders (don't spam API)
- [ ] Add visual feedback for pending actions
- [ ] Style with glassmorphism theme

---

### 9. Testing

#### 9.1 Backend Tests

**Provider Client Tests** (`backend/pkg/providers/lifx/client_test.go`):
- [ ] Test `ListDevices()` with mock HTTP responses
- [ ] Test `GetDevice()` with valid/invalid device ID
- [ ] Test all control methods (SetPower, SetBrightness, etc.)
- [ ] Test error handling (network errors, 401, 429 rate limit)
- [ ] Test color conversion and formatting

**Device Service Tests** (`backend/internal/services/device_test.go`):
- [ ] Test `ListDevices()` with caching
- [ ] Test cache invalidation after actions
- [ ] Test access control (owner vs non-owner)
- [ ] Test error propagation from provider
- [ ] Test concurrent requests handling

**Handler Tests** (`backend/internal/handlers/device_test.go`):
- [ ] Test all endpoints with valid requests
- [ ] Test authentication requirements
- [ ] Test parameter validation
- [ ] Test error responses (400, 401, 403, 404, 500)
- [ ] Test rate limiting

**Integration Tests**:
- [ ] End-to-end flow: list devices → execute action → verify state
- [ ] Test with real LIFX test account (optional)

#### 9.2 Mobile Tests

**Model Tests** (`mobile/test/models/device_test.dart`):
- [ ] Test JSON serialization/deserialization
- [ ] Test helper methods

**Service Tests** (`mobile/test/services/device_service_test.dart`):
- [ ] Test all API methods with mocked Dio
- [ ] Test error handling
- [ ] Test request formatting

**Provider Tests** (`mobile/test/providers/devices_provider_test.dart`):
- [ ] Test load/refresh logic
- [ ] Test optimistic updates
- [ ] Test revert on error
- [ ] Test filtered providers

**Widget Tests** (`mobile/test/features/devices/screens/`):
- [ ] Test DevicesScreen renders correctly
- [ ] Test DeviceDetailScreen renders correctly
- [ ] Test control widgets work
- [ ] Test pull-to-refresh
- [ ] Test navigation

**Integration Tests** (`mobile/integration_test/`):
- [ ] E2E: Login → Navigate to devices → Toggle power
- [ ] E2E: Adjust brightness → Verify UI update

---

### 10. Documentation

#### 10.1 API Documentation
**Location**: `docs/api.md`

**Update**:
- [ ] Document all device endpoints
- [ ] Add request/response examples
- [ ] Document selector pattern
- [ ] Document action types and parameters
- [ ] Add error codes

#### 10.2 Implementation Guide
**Location**: `docs/phase-4-implementation.md`

**Create**:
- [ ] Architecture overview
- [ ] Setup instructions
- [ ] Usage guide
- [ ] Troubleshooting section
- [ ] Known limitations

---

## Success Criteria

Phase 4 is complete when:

- ✅ Backend can list devices from LIFX accounts
- ✅ Backend can execute control actions (power, brightness, color)
- ✅ Mobile app displays device list grouped by location
- ✅ Mobile app can toggle power with optimistic UI
- ✅ Mobile app can adjust brightness with slider
- ✅ Mobile app can change color (if device supports it)
- ✅ Device state is cached with reasonable TTL
- ✅ Pull-to-refresh works and invalidates cache
- ✅ Offline/unreachable devices are indicated
- ✅ All tests passing (backend + mobile)
- ✅ Documentation updated

---

## Known Limitations & Future Work

### Phase 4 Limitations
1. **LIFX Only**: Hue support deferred to future phase
2. **Basic Controls**: Advanced features (scenes, schedules) in later phases
3. **No Real-time Updates**: Manual refresh required (WebSocket in future)
4. **Simple Caching**: No conflict resolution for multi-user scenarios
5. **Rate Limiting**: Basic implementation, may need refinement

### Future Enhancements
- **Phase 5**: Add OAuth flow (better than token entry)
- **Phase 6**: Sharing system (invite users to control your lights)
- **Phase 7**: Advanced features (scenes, schedules, automation)
- **Phase 8**: Philips Hue support
- **Phase 9**: Real-time device state updates (WebSocket/SSE)
- **Phase 10**: Widgets, Apple Watch, voice control

---

## Implementation Order

**Recommended sequence**:

1. **Backend Foundation** (Days 1-2)
   - Models (Device, Action)
   - Extended LIFX client (ListDevices, control methods)
   - Provider interface updates

2. **Backend Service & API** (Days 3-4)
   - Device service with caching
   - Handler endpoints
   - Route registration
   - Basic testing

3. **Mobile Models & Services** (Day 5)
   - Device models with Freezed
   - Device service
   - State management setup

4. **Mobile UI - List** (Days 6-7)
   - Devices list screen
   - Device cards
   - Pull-to-refresh
   - Navigation

5. **Mobile UI - Controls** (Days 8-9)
   - Device detail screen
   - Control widgets (power, brightness, color)
   - Optimistic updates

6. **Testing & Polish** (Days 10-11)
   - Write/complete all tests
   - Fix bugs
   - UI polish
   - Documentation

7. **Integration Testing** (Day 12)
   - E2E tests
   - Manual QA
   - Performance testing
   - Bug fixes

---

## Notes

- **Start with happy path**: Get basic list + power toggle working first
- **Test incrementally**: Don't wait until end to test
- **Use mock data**: Create mock LIFX responses for faster development
- **Commit frequently**: Logical commits for each component
- **Document as you go**: Update docs when adding features

---

## Questions to Resolve

- [ ] Caching TTL: 60 seconds reasonable? Or make it configurable?
- [ ] Rate limiting: What limits per account? (LIFX: 120 req/min)
- [ ] Selector pattern: Support group/location selectors or just device IDs?
- [ ] Effects: Include in Phase 4 or defer to Phase 7?
- [ ] Concurrent control: Last-write-wins or queue actions?
- [ ] Device icons: Use generic icons or provider-specific?

---

## References

- [LIFX API Documentation](https://api.developer.lifx.com/)
- [LIFX State API](https://api.developer.lifx.com/docs/set-state)
- [LIFX Effects API](https://api.developer.lifx.com/docs/breathe-effect)
- [Phase 3 Implementation Docs](./phase-3-implementation.md)
- [Roadmap](./roadmap.md)
