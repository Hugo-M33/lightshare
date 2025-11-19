# LightShare

A mobile application for connecting, controlling, and sharing access to smart lighting systems (LIFX, Philips Hue).

## Overview

LightShare allows users to:
- Connect their smart lighting accounts (LIFX, Philips Hue)
- Control lights from their mobile device
- Share light control access with family and friends
- Manage sharing permissions and access

## Features

### Free Tier
- Connect unlimited lighting accounts
- Full light control (on/off, brightness, color)
- Share with up to 2 users
- Ad-supported

### Pro Tier
- Everything in Free
- Share with up to 10+ users (or unlimited)
- Ad-free experience
- Priority support

## Architecture

LightShare uses a client-server architecture:

- **Mobile App** (Flutter): User interface for light control and account management
- **Backend API** (Go/Fiber): Handles authentication, proxies provider APIs, manages subscriptions
- **Database** (PostgreSQL): Stores users, accounts, permissions
- **Cache** (Redis): Sessions, rate limiting, invitation tokens

### Security Model

Provider tokens (LIFX/Hue) are **never** stored on client devices. The backend:
- Stores tokens encrypted at rest (KMS + AES-GCM)
- Proxies all API calls to providers
- Validates and refreshes tokens as needed

Clients only store session tokens for backend authentication.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Mobile | Flutter, Riverpod, Dio |
| Backend | Go, Fiber, sqlx |
| Database | PostgreSQL |
| Cache | Redis |
| Payments | Apple IAP, Google Play Billing, Stripe (web) |
| Monitoring | Sentry, Prometheus, Grafana |

## Getting Started

### Prerequisites

- Flutter SDK 3.x
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)

### Development Setup

```bash
# Clone the repository
git clone https://github.com/your-org/lightshare.git
cd lightshare

# Start dependencies
docker-compose up -d postgres redis

# Run backend
cd backend
cp .env.example .env  # Configure environment
go run cmd/server/main.go

# Run mobile app
cd mobile
flutter pub get
flutter run
```

### Environment Configuration

See `backend/.env.example` for required environment variables including:
- Database connection
- Redis connection
- KMS key configuration
- Provider OAuth credentials
- Receipt validation secrets

## Documentation

- [Architecture](docs/architecture.md) - System design and data flow
- [API Reference](docs/api.md) - Backend API endpoints
- [Security](docs/security.md) - Security guidelines and practices

## Project Structure

```
lightshare/
├── mobile/                 # Flutter mobile application
├── backend/                # Go backend service
├── docs/                   # Documentation
├── docker-compose.yml      # Local development setup
├── CLAUDE.md               # AI assistant context
└── README.md
```

## Contributing

1. Create a feature branch from `main`
2. Make your changes with appropriate tests
3. Ensure all tests pass and linting is clean
4. Submit a pull request

## License

[License Type] - See LICENSE file for details

## Privacy & Legal

- [Privacy Policy](docs/privacy-policy.md)
- GDPR compliant with data export and deletion
- See App Store / Play Store listings for data usage disclosure
