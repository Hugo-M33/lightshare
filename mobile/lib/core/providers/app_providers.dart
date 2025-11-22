import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../services/api_client.dart';
import '../services/auth_service.dart';

// Secure storage provider
final secureStorageProvider = Provider<FlutterSecureStorage>((ref) {
  return const FlutterSecureStorage(
    aOptions: AndroidOptions(
      encryptedSharedPreferences: true,
    ),
  );
});

// API base URL provider - can be overridden for different environments
final apiBaseUrlProvider = Provider<String>((ref) {
  // Default to localhost for development
  // TODO: Change this to your production API URL
  return 'http://localhost:8080';
});

// API client provider
final apiClientProvider = Provider<ApiClient>((ref) {
  final baseUrl = ref.watch(apiBaseUrlProvider);
  final secureStorage = ref.watch(secureStorageProvider);

  return ApiClient(
    baseUrl: baseUrl,
    secureStorage: secureStorage,
  );
});

// Auth service provider
final authServiceProvider = Provider<AuthService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  final secureStorage = ref.watch(secureStorageProvider);

  return AuthService(
    apiClient: apiClient,
    secureStorage: secureStorage,
  );
});
