import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/user.dart';
import 'app_providers.dart';

// Auth state class
class AuthState {
  final User? user;
  final bool isLoading;
  final bool isAuthenticated;
  final String? error;

  const AuthState({
    this.user,
    this.isLoading = false,
    this.isAuthenticated = false,
    this.error,
  });

  AuthState copyWith({
    User? user,
    bool? isLoading,
    bool? isAuthenticated,
    String? error,
  }) {
    return AuthState(
      user: user ?? this.user,
      isLoading: isLoading ?? this.isLoading,
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      error: error,
    );
  }
}

// Auth state notifier
class AuthNotifier extends StateNotifier<AuthState> {
  final Ref _ref;

  AuthNotifier(this._ref) : super(const AuthState()) {
    _checkAuthStatus();
  }

  Future<void> _checkAuthStatus() async {
    state = state.copyWith(isLoading: true);

    try {
      final authService = _ref.read(authServiceProvider);
      final user = await authService.getCurrentUser();

      if (user != null) {
        state = AuthState(
          user: user,
          isAuthenticated: true,
          isLoading: false,
        );
      } else {
        state = const AuthState(
          isAuthenticated: false,
          isLoading: false,
        );
      }
    } catch (e) {
      state = const AuthState(
        isAuthenticated: false,
        isLoading: false,
      );
    }
  }

  Future<void> signup({
    required String email,
    required String password,
  }) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final authService = _ref.read(authServiceProvider);
      final response = await authService.signup(
        email: email,
        password: password,
      );

      state = AuthState(
        user: response.user,
        isAuthenticated: false, // Email not verified yet
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  Future<void> login({
    required String email,
    required String password,
  }) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final authService = _ref.read(authServiceProvider);
      final response = await authService.login(
        email: email,
        password: password,
      );

      state = AuthState(
        user: response.user,
        isAuthenticated: true,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  Future<void> verifyEmail(String token) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final authService = _ref.read(authServiceProvider);
      await authService.verifyEmail(token);

      // Update user email verified status
      if (state.user != null) {
        state = state.copyWith(
          user: state.user!.copyWith(emailVerified: true),
          isLoading: false,
        );
      } else {
        state = state.copyWith(isLoading: false);
      }
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  Future<void> requestMagicLink(String email) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final authService = _ref.read(authServiceProvider);
      await authService.requestMagicLink(email);

      state = state.copyWith(isLoading: false);
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  Future<void> loginWithMagicLink(String token) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final authService = _ref.read(authServiceProvider);
      final response = await authService.loginWithMagicLink(token);

      state = AuthState(
        user: response.user,
        isAuthenticated: true,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
      rethrow;
    }
  }

  Future<void> logout() async {
    state = state.copyWith(isLoading: true);

    try {
      final authService = _ref.read(authServiceProvider);
      await authService.logout();

      state = const AuthState(
        isAuthenticated: false,
        isLoading: false,
      );
    } catch (e) {
      // Clear state even if logout fails
      state = const AuthState(
        isAuthenticated: false,
        isLoading: false,
      );
    }
  }

  Future<void> logoutAll() async {
    state = state.copyWith(isLoading: true);

    try {
      final authService = _ref.read(authServiceProvider);
      await authService.logoutAll();

      state = const AuthState(
        isAuthenticated: false,
        isLoading: false,
      );
    } catch (e) {
      // Clear state even if logout fails
      state = const AuthState(
        isAuthenticated: false,
        isLoading: false,
      );
    }
  }

  Future<void> refreshUser() async {
    try {
      final authService = _ref.read(authServiceProvider);
      final user = await authService.getCurrentUser();

      if (user != null) {
        state = state.copyWith(user: user);
      }
    } catch (e) {
      // Ignore errors when refreshing user
    }
  }
}

// Auth state provider
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(ref);
});
