# Mobile App Configuration

This document describes how to configure the LightShare mobile application for different environments.

## Environment Variables

The mobile app uses compile-time environment variables (via `--dart-define`) to configure environment-specific settings.

### Available Variables

#### `API_BASE_URL`

The base URL for the backend API.

- **Default**: `http://localhost:8080`
- **Type**: String
- **Example**: `http://192.168.1.100:8080`, `https://api.lightshare.app`

## Usage

### Running in Development

#### Default (localhost)

```bash
cd mobile
flutter run
```

This will use the default `http://localhost:8080`.

#### Custom Development Server

```bash
cd mobile
flutter run --dart-define=API_BASE_URL=http://192.168.1.100:8080
```

Replace `192.168.1.100:8080` with your local development server's IP and port.

### Building for Production

#### Android

```bash
cd mobile
flutter build apk --dart-define=API_BASE_URL=https://api.lightshare.app
```

Or for App Bundle:

```bash
flutter build appbundle --dart-define=API_BASE_URL=https://api.lightshare.app
```

#### iOS

```bash
cd mobile
flutter build ios --dart-define=API_BASE_URL=https://api.lightshare.app
```

### Multiple Environment Variables

You can define multiple environment variables by repeating the `--dart-define` flag:

```bash
flutter run \
  --dart-define=API_BASE_URL=https://staging.lightshare.app \
  --dart-define=ENVIRONMENT=staging
```

## Environment Configuration Files

For convenience, you can create shell scripts or configuration files for different environments:

### `.env.development` (bash script example)

```bash
#!/bin/bash
flutter run --dart-define=API_BASE_URL=http://192.168.1.100:8080
```

### `.env.staging` (bash script example)

```bash
#!/bin/bash
flutter run --dart-define=API_BASE_URL=https://staging.lightshare.app
```

### `.env.production` (bash script example)

```bash
#!/bin/bash
flutter build apk --dart-define=API_BASE_URL=https://api.lightshare.app
```

Make the scripts executable:

```bash
chmod +x .env.development .env.staging .env.production
```

Then run them:

```bash
./.env.development
```

## IDE Configuration

### VS Code

Create launch configurations in `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Development",
      "request": "launch",
      "type": "dart",
      "program": "lib/main.dart",
      "args": [
        "--dart-define=API_BASE_URL=http://192.168.1.100:8080"
      ]
    },
    {
      "name": "Staging",
      "request": "launch",
      "type": "dart",
      "program": "lib/main.dart",
      "args": [
        "--dart-define=API_BASE_URL=https://staging.lightshare.app"
      ]
    },
    {
      "name": "Production",
      "request": "launch",
      "type": "dart",
      "program": "lib/main.dart",
      "args": [
        "--dart-define=API_BASE_URL=https://api.lightshare.app"
      ]
    }
  ]
}
```

### Android Studio / IntelliJ IDEA

1. Go to **Run > Edit Configurations**
2. Select your Flutter run configuration
3. Add to **Additional run args**: `--dart-define=API_BASE_URL=http://your-api-url`

## Testing with Different APIs

When running on a physical device, make sure:

1. Your device and development machine are on the same network
2. Use your machine's local IP address (not `localhost` or `127.0.0.1`)
3. The backend server is accessible from your device

### Finding Your Local IP

**macOS/Linux:**
```bash
ifconfig | grep "inet " | grep -v 127.0.0.1
```

**Windows:**
```bash
ipconfig
```

Look for your IPv4 address (usually starts with `192.168.x.x` or `10.0.x.x`).

## CI/CD Integration

### GitHub Actions

```yaml
- name: Build Android APK
  run: |
    cd mobile
    flutter build apk \
      --dart-define=API_BASE_URL=${{ secrets.API_BASE_URL }}
```

### Environment Secrets

Store your API URLs as secrets in your CI/CD platform:

- **Development**: `DEV_API_BASE_URL`
- **Staging**: `STAGING_API_BASE_URL`
- **Production**: `PROD_API_BASE_URL`

## Security Considerations

- Never commit actual production URLs to version control if they contain sensitive information
- Use CI/CD secrets for production builds
- For development, use `.env.*` files that are in `.gitignore`
- The `--dart-define` values are embedded in the compiled app, so ensure they point to public endpoints only
- Sensitive configuration (API keys, tokens) should be retrieved from the backend, not hardcoded in the app

## Troubleshooting

### Changes not taking effect

1. Stop the app completely
2. Run `flutter clean`
3. Rebuild with the `--dart-define` flag

### Can't connect to local backend

1. Verify backend is running: `curl http://localhost:8080/health`
2. Check your local IP address
3. Ensure firewall allows connections on the port
4. For Android emulator, use `10.0.2.2` instead of `localhost`
5. For iOS simulator, `localhost` should work

### Environment variable not found

Ensure you're using the exact variable name: `API_BASE_URL` (case-sensitive).
