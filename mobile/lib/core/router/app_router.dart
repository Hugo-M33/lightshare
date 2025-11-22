import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/auth_provider.dart';
import '../../features/auth/screens/login_screen.dart';
import '../../features/auth/screens/signup_screen.dart';
import '../../features/auth/screens/email_verification_screen.dart';
import '../../features/auth/screens/magic_link_screen.dart';
import '../../features/home/screens/home_screen.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authProvider);

  return GoRouter(
    initialLocation: '/',
    redirect: (context, state) {
      final isAuthenticated = authState.isAuthenticated;
      final isLoading = authState.isLoading;

      // Show loading during auth check
      if (isLoading) {
        return null;
      }

      final isAuthRoute = state.matchedLocation.startsWith('/auth');

      // Redirect to login if not authenticated and not on auth route
      if (!isAuthenticated && !isAuthRoute) {
        return '/auth/login';
      }

      // Redirect to home if authenticated and on auth route
      if (isAuthenticated && isAuthRoute) {
        return '/';
      }

      return null;
    },
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => const HomeScreen(),
      ),
      GoRoute(
        path: '/auth/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/auth/signup',
        builder: (context, state) => const SignupScreen(),
      ),
      GoRoute(
        path: '/auth/verify-email',
        builder: (context, state) {
          final token = state.uri.queryParameters['token'];
          return EmailVerificationScreen(token: token);
        },
      ),
      GoRoute(
        path: '/auth/magic-link',
        builder: (context, state) {
          final token = state.uri.queryParameters['token'];
          return MagicLinkScreen(token: token);
        },
      ),
    ],
  );
});
