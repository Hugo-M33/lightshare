# Phase 3 Testing Guide

## Overview
This guide will help you test the Phase 3 implementation: Provider Connection with LIFX.

## Prerequisites

1. ✅ Backend running (should be running in Docker)
2. ✅ Database migrated (accounts table created)
3. ✅ Encryption key configured
4. 📱 Mobile app (Flutter)
5. 🔑 LIFX account with token

## Step 1: Get a LIFX Token

1. Visit https://cloud.lifx.com/settings
2. Log in to your LIFX account
3. Scroll to **"Personal Access Tokens"**
4. Click **"Generate New Token"**
5. Give it a label (e.g., "LightShare Development")
6. Click **"Generate"**
7. **Copy the token** - you'll need this for testing

## Step 2: Start the Backend

The backend should already be running in Docker. Verify it's running:

```bash
docker-compose ps
```

You should see:
- `lightshare-backend-1` - Status: Running (healthy)
- `lightshare-postgres-1` - Status: Running (healthy)
- `lightshare-redis-1` - Status: Running (healthy)

If not running, start it:

```bash
cd /Users/hugomartin/projets-perso/lightshare
export ENCRYPTION_KEY="7a0a53da98a0964f9a99438b7b9ae7aad8756e1571d2c8770977a71f4467bc2c"
docker-compose up -d
```

Check the logs to ensure no errors:

```bash
docker-compose logs -f backend
```

Look for: `"Services initialized successfully"`

## Step 3: Run the Mobile App

```bash
cd /Users/hugomartin/projets-perso/lightshare/mobile
flutter run
```

Or run from your IDE (VS Code/Android Studio).

## Step 4: Test the Flow

### 4.1 Login to the App
1. Open the app
2. Login with your test account (or signup if needed)
3. Verify email if required

### 4.2 Navigate to Accounts
1. On the home screen, tap **"Manage Accounts"** under "Smart Lights"
2. You should see the Accounts screen (empty initially)

### 4.3 Connect LIFX Account
1. Tap **"Connect Account"** button
2. You'll see Provider Selection screen
3. Tap on **"LIFX"** card
4. You'll see the Token Entry screen with instructions
5. Paste your LIFX token from Step 1
6. Tap **"Connect Account"**

### 4.4 Verify Success
1. You should see a success message
2. You'll be redirected to the Accounts screen
3. Your LIFX account should now be listed
4. You should see the account name (location name from LIFX)

### 4.5 Test Disconnect
1. On the Accounts screen, tap the **delete icon** (🗑️) on your account
2. The account should be removed
3. The list should be empty again

## Step 5: API Testing (Optional)

You can also test the API directly using curl:

### Login and Get Token
```bash
# Signup (if needed)
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq .

# Save the access_token from the response
export ACCESS_TOKEN="your-access-token-here"
```

### Connect Provider
```bash
curl -X POST http://localhost:8080/api/v1/providers/connect \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{
    "provider": "lifx",
    "token": "YOUR_LIFX_TOKEN_HERE"
  }' | jq .
```

Expected response:
```json
{
  "id": "uuid",
  "provider": "lifx",
  "provider_account_id": "location-id",
  "metadata": {
    "lights_count": 5
  },
  "created_at": "2025-01-23T..."
}
```

### List Accounts
```bash
curl -X GET http://localhost:8080/api/v1/accounts \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
```

Expected response:
```json
{
  "accounts": [
    {
      "id": "uuid",
      "provider": "lifx",
      "provider_account_id": "location-id",
      "metadata": {
        "lights_count": 5
      },
      "created_at": "2025-01-23T..."
    }
  ]
}
```

### Disconnect Account
```bash
curl -X DELETE http://localhost:8080/api/v1/accounts/ACCOUNT_ID_HERE \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
```

Expected response:
```json
{
  "message": "account disconnected successfully"
}
```

## Troubleshooting

### Backend Issues

**Error: "Failed to load encryption key"**
```bash
# Make sure ENCRYPTION_KEY is exported
export ENCRYPTION_KEY="7a0a53da98a0964f9a99438b7b9ae7aad8756e1571d2c8770977a71f4467bc2c"
docker-compose restart backend
```

**Error: "invalid provider token"**
- Your LIFX token may be invalid or expired
- Generate a new token at https://cloud.lifx.com/settings
- Make sure you copied the entire token

**Database connection errors**
```bash
# Check database is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres
```

### Mobile Issues

**"Unable to connect to server"**
- Make sure backend is running on port 8080
- Check the API_BASE_URL in mobile app (should be `http://localhost:8080` for simulator/emulator)
- For physical device, use your computer's IP address

**"Account already connected"**
- Each LIFX account can only be connected once per user
- Disconnect the existing account first
- Or use a different LIFX account

**Compilation errors**
```bash
cd mobile
flutter clean
flutter pub get
flutter run
```

### Common Issues

**Port 8080 already in use**
```bash
# Find what's using port 8080
lsof -i :8080

# Kill the process or change the port in docker-compose.yml
```

**Docker containers not healthy**
```bash
# Check all container health
docker-compose ps

# Restart all services
docker-compose down
docker-compose up -d
```

## What's Working

✅ User authentication (Phase 2)
✅ Database with accounts table
✅ Token encryption using AES-256-GCM
✅ LIFX token validation
✅ Account management (connect, list, disconnect)
✅ Mobile UI for provider connection
✅ End-to-end flow from mobile to backend

## What's NOT Implemented Yet

❌ Philips Hue support (coming in future phase)
❌ OAuth flows (coming in Phase 5)
❌ Actual light control (coming in Phase 4)
❌ Token refresh/expiration handling
❌ Sharing functionality (coming in Phase 6)

## Next Steps

After successful testing:
1. **Phase 4**: Implement light control (turn on/off, brightness, color)
2. Add device listing endpoints
3. Implement backend caching for device state
4. Add real-time updates

## Support

If you encounter issues not covered in this guide:
1. Check backend logs: `docker-compose logs backend`
2. Check mobile app logs in your IDE console
3. Verify all environment variables are set correctly
4. Ensure database migration ran successfully
