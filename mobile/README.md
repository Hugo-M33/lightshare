# LightShare Mobile App

Flutter mobile application for LightShare - a smart lighting control and sharing platform.

## Prerequisites

- Flutter SDK 3.24.0 or higher
- Dart SDK 3.0.0 or higher
- Android Studio / Xcode (for platform-specific builds)

## Getting Started

### 1. Install Dependencies

```bash
flutter pub get
```

### 2. Configure Environment

The app requires the backend API URL to be configured. By default, it uses `http://localhost:8080`.

**For development with a custom backend URL:**

```bash
flutter run --dart-define=API_BASE_URL=http://192.168.1.100:8080
```

**For production builds:**

```bash
flutter build apk --dart-define=API_BASE_URL=https://api.lightshare.app
```

See [docs/mobile-configuration.md](../docs/mobile-configuration.md) for detailed environment configuration instructions.

### 3. Run the App

**Development mode:**
```bash
flutter run
```

**With custom API URL:**
```bash
flutter run --dart-define=API_BASE_URL=http://YOUR_LOCAL_IP:8080
```

## Project Structure

```
lib/
├── core/
│   ├── models/         # Data models
│   ├── providers/      # Riverpod providers
│   └── services/       # API client, auth service
├── features/           # Feature modules
│   └── auth/          # Authentication feature
│       ├── models/
│       ├── providers/
│       ├── screens/
│       └── widgets/
└── main.dart          # App entry point
```

## Testing

**Run all tests:**
```bash
flutter test
```

**Run with coverage:**
```bash
flutter test --coverage
```

## Building

### Android

**Debug APK:**
```bash
flutter build apk --debug --dart-define=API_BASE_URL=http://your-api-url
```

**Release APK:**
```bash
flutter build apk --release --dart-define=API_BASE_URL=https://api.lightshare.app
```

### iOS

**Debug build:**
```bash
flutter build ios --debug --dart-define=API_BASE_URL=http://your-api-url
```

**Release build:**
```bash
flutter build ios --release --dart-define=API_BASE_URL=https://api.lightshare.app
```

## Key Features

- User authentication (signup, login, email verification)
- Magic link authentication
- Secure token storage (flutter_secure_storage)
- State management with Riverpod
- HTTP client with Dio

## Documentation

- [Mobile Configuration Guide](../docs/mobile-configuration.md)
- [API Documentation](../docs/api.md)
- [Architecture Overview](../docs/architecture.md)

## Troubleshooting

### Can't connect to backend

1. Verify backend is running
2. Check your local IP address (use `ipconfig` or `ifconfig`)
3. For Android emulator, use `10.0.2.2` instead of `localhost`
4. For iOS simulator, `localhost` should work
5. Ensure firewall allows connections

### Environment changes not taking effect

1. Stop the app completely
2. Run `flutter clean`
3. Rebuild with the `--dart-define` flag
