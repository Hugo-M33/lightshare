import 'package:dio/dio.dart';
import '../models/account.dart';
import '../models/provider.dart';
import 'api_client.dart';

class ProviderService {
  final ApiClient _apiClient;

  ProviderService({
    required ApiClient apiClient,
  }) : _apiClient = apiClient;

  Future<Account> connectProvider({
    required String provider,
    required String token,
  }) async {
    try {
      final response = await _apiClient.post(
        '/api/v1/providers/connect',
        data: {
          'provider': provider,
          'token': token,
        },
      );

      return Account.fromJson(response.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<List<Account>> listAccounts() async {
    try {
      final response = await _apiClient.get('/api/v1/accounts');

      final listResponse = ListAccountsResponse.fromJson(
        response.data as Map<String, dynamic>,
      );

      return listResponse.accounts;
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  Future<void> disconnectAccount(String accountId) async {
    try {
      await _apiClient.delete('/api/v1/accounts/$accountId');
    } on DioException catch (e) {
      throw _handleError(e);
    }
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
