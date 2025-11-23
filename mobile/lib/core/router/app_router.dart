import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/auth_provider.dart';
import '../models/provider.dart' as models;
import '../../features/auth/screens/login_screen.dart';
import '../../features/auth/screens/signup_screen.dart';
import '../../features/auth/screens/email_verification_screen.dart';
import '../../features/auth/screens/magic_link_screen.dart';
import '../../features/home/screens/home_screen.dart';
import '../../features/providers/screens/accounts_screen.dart';
import '../../features/providers/screens/provider_selection_screen.dart';
import '../../features/providers/screens/token_entry_screen.dart';
import '../../features/devices/screens/devices_screen.dart';
import '../../features/devices/screens/device_detail_screen.dart';

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
      GoRoute(
        path: '/accounts',
        builder: (context, state) => const AccountsScreen(),
      ),
      GoRoute(
        path: '/providers/connect',
        builder: (context, state) => const ProviderSelectionScreen(),
      ),
      GoRoute(
        path: '/providers/connect/token',
        builder: (context, state) {
          final provider = state.extra as models.Provider;
          return TokenEntryScreen(provider: provider);
        },
      ),
      GoRoute(
        path: '/devices',
        builder: (context, state) => const DevicesScreen(),
      ),
      GoRoute(
        path: '/devices/:deviceId',
        builder: (context, state) {
          final deviceId = state.pathParameters['deviceId']!;
          return DeviceDetailScreen(deviceId: deviceId);
        },
      ),
    ],
  );
});
