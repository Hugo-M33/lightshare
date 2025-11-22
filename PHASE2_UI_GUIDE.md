# Phase 2 UI Implementation Guide

## 🎨 Design Overview

LightShare now features a stunning **dark purple glassmorphism** UI that combines modern aesthetics with excellent usability.

### Color Palette

```dart
Primary Purple:    #8B5CF6  // Main brand color
Deep Purple:       #6D28D9  // Darker accent
Accent Pink:       #EC4899  // Highlight color
Dark Background:   #0F0A1E  // Main background
Card Background:   #1A1333  // Glass containers
Text Primary:      #FFFFFF  // Main text
Text Secondary:    #B4B4C8  // Muted text
Glass Border:      #33FFFFFF // 20% white for borders
```

### Design Features

✨ **Glassmorphism Effects**
- Backdrop blur on all containers
- Subtle gradient borders
- Layered depth with transparency

🌈 **Gradient Accents**
- Purple-to-pink gradients on buttons
- Multi-stop gradients on backgrounds
- Glowing effects with box shadows

🎯 **Modern UI Elements**
- Rounded corners (16-20px radius)
- Floating action buttons
- Material 3 design principles
- Smooth transitions and animations

## 📱 Screens

### 1. Login Screen (`/auth/login`)

**Features:**
- Email and password fields with validation
- Password visibility toggle
- "Sign in with Magic Link" option
- Link to signup screen
- Glowing logo with shadow effects
- Gradient submit button

**UX Flow:**
1. User enters email and password
2. Form validation on submit
3. Loading state during authentication
4. Auto-redirect to home on success
5. Error snackbar on failure

### 2. Signup Screen (`/auth/signup`)

**Features:**
- Email, password, and confirm password fields
- Real-time password strength indicator
- Terms and conditions checkbox
- Back button to navigate to login
- Success dialog after signup

**Password Strength Indicator:**
- Weak (red): < 8 characters
- Medium (orange): 8+ chars + some requirements
- Strong (green): All requirements met
  - 8+ characters
  - Uppercase letter
  - Number
  - Special character

**UX Flow:**
1. User fills in email and password
2. Real-time strength feedback
3. Must accept terms to proceed
4. Email verification dialog on success
5. Redirect to login

### 3. Email Verification Screen (`/auth/verify-email`)

**Two Modes:**

**A. With Token (Deep Link):**
- Auto-verifies email
- Shows loading spinner
- Success icon and message
- Auto-redirects to login

**B. Without Token:**
- Shows "Check your email" message
- Back to login button
- Helpful instructions

**UX Flow:**
1. User clicks link in email
2. App opens with token in URL
3. Auto-verification happens
4. Success message shows
5. Redirect to login after 2 seconds

### 4. Magic Link Screen (`/auth/magic-link`)

**Two Modes:**

**A. Request Mode (Default):**
- Email input field
- "Send Magic Link" button
- Success message after sending
- Option to request another link
- Back to login link

**B. Verification Mode (With Token):**
- Auto-login with token
- Loading state
- Error handling
- Redirect to home on success

**UX Flow:**
```
Request:
1. User enters email
2. Click "Send Magic Link"
3. See success message
4. Check email
5. Click link in email

Verification:
1. App opens with token
2. Auto-login happens
3. Redirect to home
```

### 5. Home Screen (`/`)

**Features:**
- Welcome message with user email
- Email verification status badge
- Quick stats cards (Devices, Shared)
- "Coming Soon" feature list
- Logout button with confirmation dialog
- Gradient logo header

**Stats Displayed:**
- Number of devices (currently 0)
- Number of shares (currently 0)

**Coming Soon Features:**
- Connect LIFX Devices
- Connect Philips Hue
- Share with Friends
- Remote Control

## 🧩 Reusable Components

### GlassContainer

A versatile container with glassmorphism effects.

```dart
GlassContainer(
  padding: EdgeInsets.all(24),
  borderRadius: 20,
  blur: 10,
  child: YourWidget(),
)
```

**Properties:**
- `child`: Widget to display inside
- `width`, `height`: Optional dimensions
- `padding`, `margin`: Spacing
- `borderRadius`: Corner radius (default: 20)
- `blur`: Blur amount (default: 10)
- `color`: Background color override
- `border`: Custom border

### GradientButton

Beautiful gradient button with loading state.

```dart
GradientButton(
  text: 'Sign In',
  onPressed: () => handleLogin(),
  isLoading: isLoading,
  width: double.infinity,
  height: 56,
)
```

**Properties:**
- `text`: Button text
- `onPressed`: Callback function
- `isLoading`: Shows spinner when true
- `width`, `height`: Dimensions
- `gradientColors`: Custom gradient colors

**Default Gradient:**
- Deep Purple → Primary Purple → Accent Pink

## 🎭 Theme Configuration

### App Theme (`AppTheme.darkTheme`)

The app uses a custom dark theme with:
- Material 3 design system
- Purple-based color scheme
- Custom input decorations
- Glassmorphic card theme
- Consistent button styles

### Text Styles

```dart
displayLarge:  32px, bold      // Large headings
displayMedium: 28px, bold      // Medium headings
displaySmall:  24px, semibold  // Small headings
bodyLarge:     16px, regular   // Body text
bodyMedium:    14px, regular   // Secondary text
```

## 🚀 Running the App

### Prerequisites

```bash
cd mobile
flutter pub get
```

### Development

```bash
# Run on emulator/device
flutter run

# Run with hot reload
flutter run --hot

# Build for release (Android)
flutter build apk --release

# Build for release (iOS)
flutter build ios --release
```

### API Configuration

Update the API base URL in `lib/core/providers/app_providers.dart`:

```dart
final apiBaseUrlProvider = Provider<String>((ref) {
  // Development
  return 'http://10.0.2.2:8080';  // Android emulator
  // return 'http://localhost:8080';  // iOS simulator

  // Production
  // return 'https://api.lightshare.com';
});
```

## 🔗 Deep Linking

### Android Setup

Add to `android/app/src/main/AndroidManifest.xml`:

```xml
<intent-filter android:autoVerify="true">
    <action android:name="android.intent.action.VIEW" />
    <category android:name="android.intent.category.DEFAULT" />
    <category android:name="android.intent.category.BROWSABLE" />

    <!-- HTTPS deep links -->
    <data
        android:scheme="https"
        android:host="app.lightshare.com" />

    <!-- Custom scheme -->
    <data android:scheme="lightshare" />
</intent-filter>
```

### iOS Setup

Add to `ios/Runner/Info.plist`:

```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleTypeRole</key>
        <string>Editor</string>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>lightshare</string>
        </array>
    </dict>
</array>
```

### Supported Deep Links

```
# Email verification
lightshare://auth/verify-email?token=abc123
https://app.lightshare.com/auth/verify-email?token=abc123

# Magic link login
lightshare://auth/magic-link?token=xyz789
https://app.lightshare.com/auth/magic-link?token=xyz789
```

## 📐 Layout Guidelines

### Spacing

- Screen padding: 24px
- Card padding: 24-32px
- Element spacing: 8-16px
- Section spacing: 24-48px

### Sizing

- Button height: 56px
- Icon sizes: 16-80px
- Border radius: 12-20px
- Max form width: 600px (centered)

### Responsive Design

All screens are responsive and work on:
- Small phones (320px width)
- Regular phones (375px-428px)
- Tablets (768px+)
- Different aspect ratios

## 🎨 Customization

### Changing Colors

Edit `lib/core/theme/app_theme.dart`:

```dart
class AppTheme {
  static const Color primaryPurple = Color(0xFF8B5CF6);
  static const Color deepPurple = Color(0xFF6D28D9);
  // ... change colors here
}
```

### Adding New Screens

1. Create screen in `lib/features/[feature]/screens/`
2. Add route in `lib/core/router/app_router.dart`
3. Use existing components for consistency

### Modifying Gradients

```dart
GradientButton(
  gradientColors: [
    Color(0xFF...),  // Start color
    Color(0xFF...),  // Middle color
    Color(0xFF...),  // End color
  ],
)
```

## 🐛 Troubleshooting

### Common Issues

**1. White screen on startup**
- Check API URL configuration
- Ensure backend is running
- Check console for errors

**2. Forms not validating**
- Verify all fields have validators
- Check GlobalKey is properly set
- Ensure formKey.currentState!.validate() is called

**3. Navigation not working**
- Verify route paths match exactly
- Check auth state for redirect logic
- Ensure GoRouter is properly configured

**4. Glassmorphism not showing**
- Ensure backdrop_filter is supported on platform
- Check if blur values are reasonable (10-30)
- Verify parent containers allow transparency

## 📚 Best Practices

1. **State Management**
   - Use Riverpod providers
   - Keep state minimal
   - Update UI reactively

2. **Error Handling**
   - Show user-friendly messages
   - Log errors for debugging
   - Provide recovery actions

3. **Performance**
   - Use const constructors
   - Avoid rebuilds with ConsumerWidget
   - Optimize images and assets

4. **Accessibility**
   - Provide semantic labels
   - Ensure good contrast
   - Support screen readers

## 🎯 Next Steps

After Phase 2 UI is complete:

1. **Test all flows**
   - Signup → Email verification → Login
   - Magic link request → Verification
   - Logout → Login

2. **Backend integration**
   - Connect to real API
   - Test with actual SMTP server
   - Verify deep links work

3. **Phase 3 preparation**
   - Provider connection UI
   - Device list screens
   - Sharing screens

## 📸 Screenshots

*Screenshots will be added after running the app*

## 🎉 Congratulations!

You now have a beautiful, modern authentication UI for LightShare with:
- ✅ Dark purple glassmorphism design
- ✅ Complete auth flows (signup, login, verification, magic links)
- ✅ Responsive layouts
- ✅ Reusable components
- ✅ Professional UX

Ready for Phase 3: Provider Integration! 🚀
