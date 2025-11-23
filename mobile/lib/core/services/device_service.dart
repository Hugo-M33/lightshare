import 'package:dio/dio.dart';
import '../models/device.dart';
import '../models/action_request.dart';
import 'api_client.dart';

class DeviceService {
  final ApiClient _apiClient;

  DeviceService({
    required ApiClient apiClient,
  }) : _apiClient = apiClient;

  /// List all devices across all user accounts
  Future<List<Device>> listAllDevices() async {
    try {
      final response = await _apiClient.get('/api/v1/devices');

      final data = response.data as Map<String, dynamic>;
      final devices = (data['devices'] as List<dynamic>)
          .map((json) => Device.fromJson(json as Map<String, dynamic>))
          .toList();

      return devices;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  /// List devices for a specific account
  Future<List<Device>> listAccountDevices(String accountId) async {
    try {
      final response =
          await _apiClient.get('/api/v1/accounts/$accountId/devices');

      final data = response.data as Map<String, dynamic>;
      final devices = (data['devices'] as List<dynamic>)
          .map((json) => Device.fromJson(json as Map<String, dynamic>))
          .toList();

      return devices;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  /// Get a specific device by ID
  Future<Device> getDevice(String accountId, String deviceId) async {
    try {
      final response = await _apiClient
          .get('/api/v1/accounts/$accountId/devices/$deviceId');

      return Device.fromJson(response.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  /// Execute a control action on device(s)
  ///
  /// [accountId] - The account ID
  /// [selector] - Device selector (id:xxx, all, group_id:xxx, location_id:xxx)
  /// [action] - The action to execute
  Future<void> executeAction(
    String accountId,
    String selector,
    ActionRequest action,
  ) async {
    try {
      await _apiClient.post(
        '/api/v1/accounts/$accountId/devices/$selector/action',
        data: action.toJson(),
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  /// Force refresh devices from provider (bypasses cache)
  Future<List<Device>> refreshDevices(String accountId) async {
    try {
      final response =
          await _apiClient.post('/api/v1/accounts/$accountId/devices/refresh');

      final data = response.data as Map<String, dynamic>;
      final devices = (data['devices'] as List<dynamic>)
          .map((json) => Device.fromJson(json as Map<String, dynamic>))
          .toList();

      return devices;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  /// Set power state for device(s)
  Future<void> setPower(
    String accountId,
    String selector, {
    required bool state,
    double duration = 0.0,
  }) async {
    final action = ActionRequest.power(state: state, duration: duration);
    await executeAction(accountId, selector, action);
  }

  /// Set brightness level for device(s)
  Future<void> setBrightness(
    String accountId,
    String selector, {
    required double level,
    double duration = 0.0,
  }) async {
    final action = ActionRequest.brightness(level: level, duration: duration);
    await executeAction(accountId, selector, action);
  }

  /// Set color for device(s)
  Future<void> setColor(
    String accountId,
    String selector, {
    required double hue,
    required double saturation,
    int? kelvin,
    double duration = 0.0,
  }) async {
    final action = ActionRequest.color(
      hue: hue,
      saturation: saturation,
      kelvin: kelvin,
      duration: duration,
    );
    await executeAction(accountId, selector, action);
  }

  /// Set color temperature for device(s)
  Future<void> setTemperature(
    String accountId,
    String selector, {
    required int kelvin,
    double duration = 0.0,
  }) async {
    final action =
        ActionRequest.temperature(kelvin: kelvin, duration: duration);
    await executeAction(accountId, selector, action);
  }

  /// Trigger a pulse effect
  Future<void> pulseEffect(
    String accountId,
    String selector, {
    int cycles = 3,
    double period = 1.0,
    Map<String, dynamic>? color,
  }) async {
    final action = ActionRequest.effect(
      name: DeviceEffect.pulse,
      cycles: cycles,
      period: period,
      color: color,
    );
    await executeAction(accountId, selector, action);
  }

  /// Trigger a breathe effect
  Future<void> breatheEffect(
    String accountId,
    String selector, {
    int cycles = 3,
    double period = 1.0,
    Map<String, dynamic>? color,
  }) async {
    final action = ActionRequest.effect(
      name: DeviceEffect.breathe,
      cycles: cycles,
      period: period,
      color: color,
    );
    await executeAction(accountId, selector, action);
  }

  Exception _handleError(DioException error) {
    if (error.response?.data != null && error.response!.data is Map) {
      final errorMessage = error.response!.data['error'] as String? ??
          'An unexpected error occurred';
      return Exception(errorMessage);
    }

    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.sendTimeout:
        return Exception('Connection timeout. Please try again.');
      case DioExceptionType.connectionError:
        return Exception(
          'Unable to connect to server. Please check your internet connection.',
        );
      default:
        return Exception('An unexpected error occurred. Please try again.');
    }
  }
}
