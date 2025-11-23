import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/device.dart';
import '../models/action_request.dart';
import 'app_providers.dart';

// Devices state class
class DevicesState {
  final List<Device> devices;
  final bool isLoading;
  final String? error;
  final DateTime? lastUpdated;

  const DevicesState({
    this.devices = const [],
    this.isLoading = false,
    this.error,
    this.lastUpdated,
  });

  DevicesState copyWith({
    List<Device>? devices,
    bool? isLoading,
    String? error,
    DateTime? lastUpdated,
  }) {
    return DevicesState(
      devices: devices ?? this.devices,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      lastUpdated: lastUpdated ?? this.lastUpdated,
    );
  }
}

// Devices state notifier with debouncing
class DevicesNotifier extends StateNotifier<DevicesState> {
  final Ref _ref;
  final Map<String, Timer> _debounceTimers = {};

  DevicesNotifier(this._ref) : super(const DevicesState());

  @override
  void dispose() {
    // Cancel all pending timers
    for (final timer in _debounceTimers.values) {
      timer.cancel();
    }
    _debounceTimers.clear();
    super.dispose();
  }

  /// Load all devices across all accounts
  Future<void> loadDevices() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final deviceService = _ref.read(deviceServiceProvider);
      final devices = await deviceService.listAllDevices();

      state = DevicesState(
        devices: devices,
        isLoading: false,
        lastUpdated: DateTime.now(),
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  /// Load devices for a specific account
  Future<void> loadAccountDevices(String accountId) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final deviceService = _ref.read(deviceServiceProvider);
      final devices = await deviceService.listAccountDevices(accountId);

      state = DevicesState(
        devices: devices,
        isLoading: false,
        lastUpdated: DateTime.now(),
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  /// Refresh devices from provider (bypasses cache)
  Future<void> refreshDevices(String accountId) async {
    try {
      final deviceService = _ref.read(deviceServiceProvider);
      final devices = await deviceService.refreshDevices(accountId);

      // Update only the devices for this account
      final updatedDevices = [
        ...state.devices.where((d) => d.accountId != accountId),
        ...devices,
      ];

      state = DevicesState(
        devices: updatedDevices,
        isLoading: false,
        lastUpdated: DateTime.now(),
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  /// Update device locally (optimistic UI)
  void _updateDeviceLocally(String deviceId, Device Function(Device) updater) {
    final updatedDevices = state.devices.map((device) {
      if (device.id == deviceId) {
        return updater(device);
      }
      return device;
    }).toList();

    state = state.copyWith(devices: updatedDevices);
  }

  /// Execute action with debouncing
  ///
  /// Immediately updates UI optimistically, then debounces API call by 500ms
  Future<void> _executeActionDebounced({
    required String accountId,
    required String selector,
    required ActionRequest action,
    required String deviceId,
    required Device Function(Device) optimisticUpdate,
  }) async {
    // Immediate optimistic UI update
    _updateDeviceLocally(deviceId, optimisticUpdate);

    // Cancel existing timer for this device
    final timerKey = '$deviceId-${action.action}';
    _debounceTimers[timerKey]?.cancel();

    // Create new debounced timer
    _debounceTimers[timerKey] = Timer(
      const Duration(milliseconds: 500),
      () async {
        try {
          final deviceService = _ref.read(deviceServiceProvider);
          await deviceService.executeAction(accountId, selector, action);

          // Optionally refresh device state from server after action
          // This ensures UI is in sync with actual device state
          try {
            final updatedDevice =
                await deviceService.getDevice(accountId, deviceId);
            _updateDeviceLocally(deviceId, (_) => updatedDevice);
          } catch (_) {
            // Silently fail refresh - optimistic update is still shown
          }
        } catch (e) {
          // Revert optimistic update on error
          await loadDevices();
          state = state.copyWith(error: e.toString());
        } finally {
          _debounceTimers.remove(timerKey);
        }
      },
    );
  }

  /// Set power state with optimistic UI
  Future<void> setPower(
    String accountId,
    String deviceId, {
    required bool state,
    double duration = 0.0,
  }) async {
    final selector = 'id:$deviceId';
    final action = ActionRequest.power(state: state, duration: duration);

    await _executeActionDebounced(
      accountId: accountId,
      selector: selector,
      action: action,
      deviceId: deviceId,
      optimisticUpdate: (device) => device.copyWith(
        power: state ? 'on' : 'off',
      ),
    );
  }

  /// Set brightness with debouncing and optimistic UI
  Future<void> setBrightness(
    String accountId,
    String deviceId, {
    required double level,
    double duration = 0.0,
  }) async {
    final selector = 'id:$deviceId';
    final action = ActionRequest.brightness(level: level, duration: duration);

    await _executeActionDebounced(
      accountId: accountId,
      selector: selector,
      action: action,
      deviceId: deviceId,
      optimisticUpdate: (device) => device.copyWith(
        brightness: level,
      ),
    );
  }

  /// Set color with debouncing and optimistic UI
  Future<void> setColor(
    String accountId,
    String deviceId, {
    required double hue,
    required double saturation,
    int? kelvin,
    double duration = 0.0,
  }) async {
    final selector = 'id:$deviceId';
    final action = ActionRequest.color(
      hue: hue,
      saturation: saturation,
      kelvin: kelvin,
      duration: duration,
    );

    await _executeActionDebounced(
      accountId: accountId,
      selector: selector,
      action: action,
      deviceId: deviceId,
      optimisticUpdate: (device) => device.copyWith(
        color: DeviceColor(
          hue: hue,
          saturation: saturation,
          kelvin: kelvin ?? device.color?.kelvin ?? 3500,
        ),
      ),
    );
  }

  /// Set temperature with debouncing and optimistic UI
  Future<void> setTemperature(
    String accountId,
    String deviceId, {
    required int kelvin,
    double duration = 0.0,
  }) async {
    final selector = 'id:$deviceId';
    final action = ActionRequest.temperature(kelvin: kelvin, duration: duration);

    await _executeActionDebounced(
      accountId: accountId,
      selector: selector,
      action: action,
      deviceId: deviceId,
      optimisticUpdate: (device) => device.copyWith(
        color: device.color?.copyWith(kelvin: kelvin) ??
            DeviceColor(hue: 0, saturation: 0, kelvin: kelvin),
      ),
    );
  }

  /// Trigger pulse effect (no debouncing - instant action)
  Future<void> pulseEffect(
    String accountId,
    String deviceId, {
    int cycles = 3,
    double period = 1.0,
    Map<String, dynamic>? color,
  }) async {
    try {
      final deviceService = _ref.read(deviceServiceProvider);
      final selector = 'id:$deviceId';

      await deviceService.pulseEffect(
        accountId,
        selector,
        cycles: cycles,
        period: period,
        color: color,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString());
      rethrow;
    }
  }

  /// Trigger breathe effect (no debouncing - instant action)
  Future<void> breatheEffect(
    String accountId,
    String deviceId, {
    int cycles = 3,
    double period = 1.0,
    Map<String, dynamic>? color,
  }) async {
    try {
      final deviceService = _ref.read(deviceServiceProvider);
      final selector = 'id:$deviceId';

      await deviceService.breatheEffect(
        accountId,
        selector,
        cycles: cycles,
        period: period,
        color: color,
      );
    } catch (e) {
      state = state.copyWith(error: e.toString());
      rethrow;
    }
  }

  /// Execute action on multiple devices (group or location)
  Future<void> executeGroupAction(
    String accountId,
    String selector,
    ActionRequest action,
  ) async {
    try {
      final deviceService = _ref.read(deviceServiceProvider);
      await deviceService.executeAction(accountId, selector, action);

      // Refresh devices after group action
      await loadDevices();
    } catch (e) {
      state = state.copyWith(error: e.toString());
      rethrow;
    }
  }

  void clearError() {
    state = state.copyWith(error: null);
  }
}

// Devices state provider
final devicesProvider =
    StateNotifierProvider<DevicesNotifier, DevicesState>((ref) {
  return DevicesNotifier(ref);
});

// Helper provider to get a specific device by ID
final deviceByIdProvider =
    Provider.family<Device?, String>((ref, deviceId) {
  final devicesState = ref.watch(devicesProvider);
  try {
    return devicesState.devices.firstWhere((device) => device.id == deviceId);
  } catch (_) {
    return null;
  }
});
