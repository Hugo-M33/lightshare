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

      // Verify email and password form fields exist
      expect(find.widgetWithText(TextFormField, 'Email'), findsOneWidget);
      expect(find.widgetWithText(TextFormField, 'Password'), findsOneWidget);

      // Verify email icon is present
      expect(find.byIcon(Icons.email_outlined), findsOneWidget);

      // Verify password lock icon is present
      expect(find.byIcon(Icons.lock_outline), findsOneWidget);
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

      // Check for header icon
      expect(find.byIcon(Icons.person_add_outlined), findsOneWidget);

      // Check for signup form elements
      expect(find.text('Create Account'), findsAtLeastNWidgets(1));
      expect(find.text('Join LightShare today'), findsOneWidget);
      expect(find.byType(TextFormField), findsNWidgets(3));

      // Check for form fields
      expect(find.widgetWithText(TextFormField, 'Email'), findsOneWidget);
      expect(find.widgetWithText(TextFormField, 'Password'), findsOneWidget);
      expect(
          find.widgetWithText(TextFormField, 'Confirm Password'), findsOneWidget);

      // Check for terms checkbox
      expect(find.text('I agree to the Terms of Service and Privacy Policy'),
          findsOneWidget);

      // Check for sign in link
      expect(find.text('Already have an account? '), findsOneWidget);
      expect(find.text('Sign In'), findsOneWidget);
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

      // Check for smart lights section
      expect(find.text('Smart Lights'), findsOneWidget);
      expect(find.text('Manage Accounts'), findsOneWidget);
      expect(find.text('Share with Friends'), findsOneWidget);

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
