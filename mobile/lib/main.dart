import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:app_links/app_links.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';

void main() {
  runApp(
    const ProviderScope(
      child: LightShareApp(),
    ),
  );
}

class LightShareApp extends ConsumerStatefulWidget {
  const LightShareApp({super.key});

  @override
  ConsumerState<LightShareApp> createState() => _LightShareAppState();
}

class _LightShareAppState extends ConsumerState<LightShareApp> {
  late AppLinks _appLinks;
  StreamSubscription<Uri>? _linkSubscription;

  @override
  void initState() {
    super.initState();
    _initDeepLinks();
  }

  @override
  void dispose() {
    _linkSubscription?.cancel();
    super.dispose();
  }

  Future<void> _initDeepLinks() async {
    _appLinks = AppLinks();

    // Handle links when app is already running
    _linkSubscription = _appLinks.uriLinkStream.listen((uri) {
      _handleDeepLink(uri);
    }, onError: (err) {
      debugPrint('Deep link error: $err');
    });

    // Handle initial link if app was opened via deep link
    try {
      final initialLink = await _appLinks.getInitialLink();
      if (initialLink != null) {
        _handleDeepLink(initialLink);
      }
    } catch (e) {
      debugPrint('Failed to get initial link: $e');
    }
  }

  void _handleDeepLink(Uri uri) {
    debugPrint('Deep link received: $uri');

    // For custom URL schemes like lightshare://verify-email?token=xxx,
    // the "verify-email" part is the host, not the path
    final host = uri.host;
    final queryParams = uri.queryParameters;

    String? targetRoute;

    // Determine target route based on host
    if (host == 'verify-email') {
      final token = queryParams['token'];
      if (token != null) {
        targetRoute = '/auth/verify-email?token=$token';
      }
    } else if (host == 'magic-link') {
      final token = queryParams['token'];
      if (token != null) {
        targetRoute = '/auth/magic-link?token=$token';
      }
    }

    if (targetRoute != null) {
      debugPrint('Navigating to: $targetRoute');

      // Wait for the next frame to ensure router is ready
      WidgetsBinding.instance.addPostFrameCallback((_) {
        final router = ref.read(routerProvider);
        debugPrint('Executing navigation to: $targetRoute');
        router.go(targetRoute!); // Safe to use ! because we checked null above
      });
    } else {
      debugPrint('Unknown deep link host: $host');
    }
  }

  @override
  Widget build(BuildContext context) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: 'LightShare',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.darkTheme,
      themeMode: ThemeMode.dark,
      routerConfig: router,
    );
  }
}
