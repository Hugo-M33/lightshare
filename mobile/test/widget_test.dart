import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lightshare/main.dart';

void main() {
  testWidgets('App renders home screen correctly', (WidgetTester tester) async {
    await tester.pumpWidget(
      const ProviderScope(
        child: LightShareApp(),
      ),
    );

    expect(find.text('LightShare'), findsOneWidget);
    expect(find.text('Welcome to LightShare'), findsOneWidget);
    expect(find.text('Control your smart lights'), findsOneWidget);
    expect(find.byIcon(Icons.lightbulb_outline), findsOneWidget);
  });

  testWidgets('Home screen has correct layout', (WidgetTester tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: HomeScreen(),
      ),
    );

    final scaffold = tester.widget<Scaffold>(find.byType(Scaffold));
    expect(scaffold.appBar, isNotNull);

    final icon = tester.widget<Icon>(find.byIcon(Icons.lightbulb_outline));
    expect(icon.size, equals(80));
  });
}
