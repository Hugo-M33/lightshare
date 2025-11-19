# Security Guidelines

## Overview

LightShare handles sensitive data including user credentials, payment information, and third-party API tokens. This document outlines security requirements and best practices for development.

## Threat Model

### Assets to Protect
1. User credentials (passwords, session tokens)
2. Provider tokens (LIFX, Hue access tokens)
3. Payment data (receipts, subscription status)
4. User data (email, usage patterns)

### Threat Actors
- External attackers (network-based attacks)
- Malicious users (abuse, token theft)
- Compromised dependencies (supply chain)

## Authentication Security

### Password Storage
- Use bcrypt with cost factor >= 12
- Never store plaintext passwords
- Implement password strength requirements

```go
// Example
hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
```

### Session Management
- Access tokens: Short-lived (1 hour)
- Refresh tokens: Longer-lived (30 days), stored in database
- Implement token rotation on refresh
- Store refresh tokens hashed in database

### Token Revocation
- Maintain revocation list in Redis for access tokens
- Delete refresh tokens from database on logout
- Revoke all tokens on password change

## Provider Token Security

### Token Flow
```
User -> Mobile App -> Backend -> Provider API
                |
                └── Provider tokens NEVER leave backend
```

### Storage Requirements

1. **Never store provider tokens in plaintext**

2. **Use envelope encryption**:
   - Master key in KMS (AWS KMS/GCP KMS/HashiCorp Vault)
   - Data Encryption Key (DEK) encrypted by master key
   - Tokens encrypted by DEK using AES-256-GCM

3. **Key rotation**:
   - Rotate DEK periodically
   - Re-encrypt all tokens with new DEK
   - Master key rotation handled by KMS

### Implementation

```go
// Encrypt token before storage
func EncryptToken(plaintext []byte, dek []byte) ([]byte, error) {
    block, err := aes.NewCipher(dek)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    return gcm.Seal(nonce, nonce, plaintext, nil), nil
}
```

### Token Validation
- Validate tokens on initial connection with test API call
- Check token validity periodically
- Implement automatic refresh for OAuth tokens
- Alert on repeated validation failures

## API Security

### Transport Security
- TLS 1.2+ required
- HSTS enabled with long max-age
- Strong cipher suites only
- Certificate pinning in mobile app (optional)

### Input Validation
- Validate all input on server side
- Use parameterized queries (prevent SQL injection)
- Sanitize output (prevent XSS in any web interfaces)
- Limit request body size

### Authorization
- Verify ownership/access for every resource request
- Check role permissions (viewer vs controller)
- Implement resource-level access control

```go
// Example middleware
func RequireAccountAccess(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    accountID := c.Params("id")

    hasAccess, err := checkAccess(userID, accountID)
    if err != nil || !hasAccess {
        return c.Status(403).JSON(fiber.Map{
            "error": "forbidden",
        })
    }

    return c.Next()
}
```

### Rate Limiting
- Per-user rate limits
- Per-IP rate limits for auth endpoints
- Provider-specific rate limits (respect LIFX/Hue quotas)
- Use Redis for distributed rate limiting

## Mobile Security

### Secure Storage
- Use `flutter_secure_storage` for session tokens
- **Never** store provider tokens on device
- Clear sensitive data on logout

### Network Security
- Use HTTPS exclusively
- Validate SSL certificates
- Consider certificate pinning for production

### Deep Links
- Validate deep link parameters
- Use HTTPS scheme for universal links
- Implement PKCE for OAuth flows

### App Security
- Enable code obfuscation for release builds
- Use ProGuard/R8 for Android
- Disable debugging in production

## Payment Security

### In-App Purchases
- **Always** validate receipts server-side
- Never trust client-side purchase status
- Maintain server-side subscription state

### Receipt Validation

**Apple:**
```go
func ValidateAppleReceipt(receipt string) (*AppleResponse, error) {
    // Always try production first
    resp, err := validateWithApple(receipt, productionURL)
    if err != nil {
        return nil, err
    }

    // If status 21007, retry with sandbox
    if resp.Status == 21007 {
        return validateWithApple(receipt, sandboxURL)
    }

    return resp, nil
}
```

**Google:**
```go
func ValidateGooglePurchase(token, productID string) (*GoogleResponse, error) {
    // Use service account authentication
    // Call Google Play Developer API
    // Verify purchase state and subscription status
}
```

### Stripe Security (Web)
- Use Stripe.js for payment collection
- Never handle raw card data
- Implement webhook signature verification
- Use idempotency keys for charges

## Data Protection

### Personal Data
- Minimize data collection
- Encrypt sensitive fields
- Implement data retention policies
- Support GDPR rights (export, deletion)

### Logging
- **Never** log secrets, tokens, or passwords
- Sanitize logs for PII
- Use structured logging
- Retain logs for appropriate period

```go
// Bad
log.Info("User login", "email", email, "password", password)

// Good
log.Info("User login", "email", email, "user_id", userID)
```

### Backups
- Encrypt database backups
- Test backup restoration
- Store backups in separate location
- Implement backup retention policy

## Sharing Security

### Invitation Flow
- Generate cryptographically random tokens
- Expire invitations (7 days default)
- One-time use tokens
- Rate limit invitation creation

### Access Control
- Enforce share limits server-side
- Log all access grant changes
- Allow owners to revoke access immediately
- Audit trail for shared account actions

## Security Monitoring

### Alerting
- Failed login attempts (brute force detection)
- Token validation failures
- Unusual API patterns
- Error rate spikes

### Audit Logging
Log security-relevant events:
- Login success/failure
- Password changes
- Provider connections/disconnections
- Share grants/revocations
- Subscription changes

```go
type AuditEvent struct {
    Timestamp time.Time
    UserID    string
    Action    string
    Resource  string
    IP        string
    UserAgent string
    Success   bool
    Details   map[string]interface{}
}
```

## Incident Response

### Preparation
- Document incident response procedures
- Maintain contact list
- Have rollback procedures ready
- Test recovery procedures

### If Provider Token Compromised
1. Revoke affected tokens immediately
2. Notify affected users
3. Investigate breach vector
4. Force re-authentication with providers

### If User Data Compromised
1. Assess scope of breach
2. Notify affected users (within GDPR timeline if applicable)
3. Force password resets
4. Review and fix vulnerability

## Development Practices

### Code Review
- Security-focused code review for auth/payment code
- Use static analysis tools (gosec, semgrep)
- Dependency vulnerability scanning

### Testing
- Unit tests for security functions
- Integration tests for auth flows
- Penetration testing before launch

### Dependencies
- Keep dependencies updated
- Monitor for security advisories
- Use lockfiles for reproducible builds
- Scan with Dependabot/Snyk

## Compliance Checklist

### Before Launch
- [ ] All passwords hashed with bcrypt
- [ ] Provider tokens encrypted at rest
- [ ] TLS enabled everywhere
- [ ] Receipt validation implemented
- [ ] Rate limiting in place
- [ ] Input validation on all endpoints
- [ ] Logging sanitized for secrets
- [ ] GDPR data export/deletion working
- [ ] Privacy policy published
- [ ] Penetration test completed

### Ongoing
- [ ] Dependencies updated monthly
- [ ] Security logs reviewed weekly
- [ ] Backup restoration tested quarterly
- [ ] Access reviews conducted quarterly
- [ ] Incident response plan reviewed annually
