import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../services/api_client.dart';
import '../services/auth_service.dart';
import '../services/provider_service.dart';
import '../services/device_service.dart';

// Secure storage provider
final secureStorageProvider = Provider<FlutterSecureStorage>((ref) {
  return const FlutterSecureStorage(
    aOptions: AndroidOptions(encryptedSharedPreferences: true),
  );
});

// API base URL provider - reads from environment variables
// Use --dart-define=API_BASE_URL=http://your-api-url to override
final apiBaseUrlProvider = Provider<String>((ref) {
  const baseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );
  return baseUrl;
});

// API client provider
final apiClientProvider = Provider<ApiClient>((ref) {
  final baseUrl = ref.watch(apiBaseUrlProvider);
  final secureStorage = ref.watch(secureStorageProvider);

  return ApiClient(baseUrl: baseUrl, secureStorage: secureStorage);
});

// Auth service provider
final authServiceProvider = Provider<AuthService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  final secureStorage = ref.watch(secureStorageProvider);

  return AuthService(apiClient: apiClient, secureStorage: secureStorage);
});

// Provider service provider
final providerServiceProvider = Provider<ProviderService>((ref) {
  final apiClient = ref.watch(apiClientProvider);

  return ProviderService(apiClient: apiClient);
});

// Device service provider
final deviceServiceProvider = Provider<DeviceService>((ref) {
  final apiClient = ref.watch(apiClientProvider);

  return DeviceService(apiClient: apiClient);
});
