import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lightshare/main.dart';
import 'package:lightshare/features/auth/screens/login_screen.dart';
import 'package:lightshare/features/auth/screens/signup_screen.dart';
import 'package:lightshare/features/home/screens/home_screen.dart';

void main() {
  group('LightShareApp', () {
    testWidgets('App initializes with MaterialApp.router',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: LightShareApp(),
        ),
      );

      // The app should render without errors
      expect(find.byType(LightShareApp), findsOneWidget);
      expect(find.byType(MaterialApp), findsOneWidget);
    });

    testWidgets('App has correct title', (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: LightShareApp(),
        ),
      );
      await tester.pumpAndSettle();

      // Verify the app title is set correctly
      final MaterialApp app =
          tester.widget(find.byType(MaterialApp).first);
      expect(app.title, equals('LightShare'));
    });
  });

  group('LoginScreen', () {
    testWidgets('Login screen renders correctly',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: LoginScreen(),
          ),
        ),
      );

      // Check for main branding elements
      expect(find.text('LightShare'), findsOneWidget);
      expect(find.text('Control your lights, anywhere'), findsOneWidget);
      expect(find.byIcon(Icons.lightbulb_outline), findsOneWidget);

      // Check for form elements
      expect(find.text('Welcome Back'), findsOneWidget);
      expect(find.text('Sign in to continue'), findsOneWidget);
      expect(find.byType(TextFormField), findsNWidgets(2));

      // Check for buttons
      expect(find.text('Sign In'), findsOneWidget);
      expect(find.text('Login with Magic Link'), findsOneWidget);
      expect(find.text('Sign Up'), findsOneWidget);
    });

    testWidgets('Email and password fields are present',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: LoginScreen(),
          ),
        ),
      );

      // Find email field
      final emailField = find.ancestor(
        of: find.text('Email'),
        matching: find.byType(TextFormField),
      );
      expect(emailField, findsOneWidget);

      // Find password field
      final passwordField = find.ancestor(
        of: find.text('Password'),
        matching: find.byType(TextFormField),
      );
      expect(passwordField, findsOneWidget);

      // Verify email field has proper keyboard type
      final emailTextFormField =
          tester.widget<TextFormField>(emailField);
      expect(emailTextFormField.keyboardType,
          equals(TextInputType.emailAddress));

      // Verify password field is obscured
      final passwordTextFormField =
          tester.widget<TextFormField>(passwordField);
      expect(passwordTextFormField.obscureText, isTrue);
    });

    testWidgets('Password visibility toggle works',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: LoginScreen(),
          ),
        ),
      );

      // Find the password visibility toggle button
      final visibilityButton = find.byIcon(Icons.visibility_outlined);
      expect(visibilityButton, findsOneWidget);

      // Tap the visibility toggle
      await tester.tap(visibilityButton);
      await tester.pump();

      // Verify the icon changed
      expect(find.byIcon(Icons.visibility_off_outlined),
          findsOneWidget);
      expect(find.byIcon(Icons.visibility_outlined), findsNothing);
    });
  });

  group('SignupScreen', () {
    testWidgets('Signup screen renders correctly',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: SignupScreen(),
          ),
        ),
      );

      // Check for branding
      expect(find.text('LightShare'), findsOneWidget);
      expect(find.byIcon(Icons.lightbulb_outline), findsOneWidget);

      // Check for signup form elements
      expect(find.text('Create Account'), findsOneWidget);
      expect(find.byType(TextFormField), findsNWidgets(2));
      expect(find.text('Sign Up'), findsAtLeastNWidgets(1));
    });
  });

  group('HomeScreen', () {
    testWidgets('Home screen renders correctly for authenticated user',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: HomeScreen(),
          ),
        ),
      );

      // Check for main elements
      expect(find.text('LightShare'), findsOneWidget);
      expect(find.byIcon(Icons.lightbulb_outline), findsAtLeastNWidgets(1));
      expect(find.text('Welcome back! 👋'), findsOneWidget);

      // Check for stat cards
      expect(find.text('Devices'), findsOneWidget);
      expect(find.text('Shared'), findsOneWidget);

      // Check for coming soon section
      expect(find.text('Coming Soon'), findsOneWidget);
      expect(find.text('Connect LIFX Devices'), findsOneWidget);
      expect(find.text('Connect Philips Hue'), findsOneWidget);

      // Check for logout button
      expect(find.byIcon(Icons.logout), findsOneWidget);
    });

    testWidgets('Home screen has correct layout structure',
        (WidgetTester tester) async {
      await tester.pumpWidget(
        const ProviderScope(
          child: MaterialApp(
            home: HomeScreen(),
          ),
        ),
      );

      // Verify scaffold exists
      expect(find.byType(Scaffold), findsOneWidget);

      // Verify SafeArea exists
      expect(find.byType(SafeArea), findsOneWidget);

      // Verify scrollable content
      expect(find.byType(SingleChildScrollView), findsOneWidget);
    });
  });
}
