import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../models/auth_response.dart';
import '../models/user.dart';
import 'api_client.dart';

class AuthService {
  final ApiClient _apiClient;
  final FlutterSecureStorage _secureStorage;

  AuthService({
    required ApiClient apiClient,
    required FlutterSecureStorage secureStorage,
  })  : _apiClient = apiClient,
        _secureStorage = secureStorage;

  Future<SignupResponse> signup({
    required String email,
    required String password,
  }) async {
    try {
      final response = await _apiClient.post(
        '/api/v1/auth/signup',
        data: {
          'email': email,
          'password': password,
        },
      );

      return SignupResponse.fromJson(response.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<AuthResponse> login({
    required String email,
    required String password,
  }) async {
    try {
      final response = await _apiClient.post(
        '/api/v1/auth/login',
        data: {
          'email': email,
          'password': password,
        },
      );

      final authResponse =
          AuthResponse.fromJson(response.data as Map<String, dynamic>);

      // Store tokens
      await _storeTokens(authResponse);

      return authResponse;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> verifyEmail(String token) async {
    try {
      await _apiClient.post(
        '/api/v1/auth/verify-email',
        data: {'token': token},
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> requestMagicLink(String email) async {
    try {
      await _apiClient.post(
        '/api/v1/auth/magic-link',
        data: {'email': email},
      );
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<AuthResponse> loginWithMagicLink(String token) async {
    try {
      final response = await _apiClient.post(
        '/api/v1/auth/magic-link/verify',
        data: {'token': token},
      );

      final authResponse =
          AuthResponse.fromJson(response.data as Map<String, dynamic>);

      // Store tokens
      await _storeTokens(authResponse);

      return authResponse;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> logout() async {
    try {
      final refreshToken = await _secureStorage.read(key: 'refresh_token');
      if (refreshToken != null) {
        await _apiClient.post(
          '/api/v1/auth/logout',
          data: {'refresh_token': refreshToken},
        );
      }
    } catch (e) {
      // Ignore logout errors
    } finally {
      await _clearTokens();
    }
  }

  Future<void> logoutAll() async {
    try {
      await _apiClient.post('/api/v1/auth/logout-all');
    } catch (e) {
      // Ignore logout errors
    } finally {
      await _clearTokens();
    }
  }

  Future<User?> getCurrentUser() async {
    try {
      final accessToken = await _secureStorage.read(key: 'access_token');
      if (accessToken == null) return null;

      final response = await _apiClient.get('/api/v1/auth/me');
      return User.fromJson(response.data as Map<String, dynamic>);
    } catch (e) {
      return null;
    }
  }

  Future<bool> isLoggedIn() async {
    final accessToken = await _secureStorage.read(key: 'access_token');
    return accessToken != null;
  }

  Future<void> _storeTokens(AuthResponse authResponse) async {
    await _secureStorage.write(
      key: 'access_token',
      value: authResponse.accessToken,
    );
    await _secureStorage.write(
      key: 'refresh_token',
      value: authResponse.refreshToken,
    );
  }

  Future<void> _clearTokens() async {
    await _secureStorage.delete(key: 'access_token');
    await _secureStorage.delete(key: 'refresh_token');
  }

  String _handleError(DioException error) {
    if (error.response != null) {
      final data = error.response!.data;
      if (data is Map<String, dynamic> && data.containsKey('error')) {
        return data['error'] as String;
      }
    }

    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return 'Connection timeout. Please check your internet connection.';
      case DioExceptionType.badResponse:
        return 'Server error. Please try again later.';
      case DioExceptionType.cancel:
        return 'Request cancelled.';
      default:
        return 'An unexpected error occurred.';
    }
  }
}
