# Security Best Practices

## Password Policy

All user passwords must meet the following requirements:

- **Minimum length:** 12 characters
- **Maximum length:** 128 characters
- **Complexity requirements:**
  - At least one uppercase letter (A-Z)
  -At least one lowercase letter (a-z)
  - At least one number (0-9)
  - At least one special character (!@#$%^&*()_+-=[]{}; :'"\|,.<>/?)
- **Prohibited patterns:** Common passwords like "password", "123456", "qwerty", "admin", etc.

### Implementation

Password validation is enforced in `internal/core/services/validation_helpers.go` via the `validatePasswordStrength()` function.

## Input Validation

All user inputs must be validated before processing:

### Weight Logging

- **Range:** 20 kg - 500 kg (44 lbs - 1100 lbs)
- **Units:** Only `kg`, `lbs`, or `lb` accepted
- **Type:** Must be positive number

### Hydration Logging

- **Range:** 0 - 100 glasses
- **Type:** Must be non-negative number

### Implementation

Input validation is enforced in `internal/core/services/progress_service.go`.

## Authentication Security

### JWTJWT Token Security

- **Algorithm:** Strictly HS256 only
- **Secret:** Minimum 32 characters, cryptographically random
- **Expiration:** 24 hours
- **Claims validation:** Explicit expiration checking

### Rate Limiting

**Global Endpoints:**

- 100 requests per minute per IP

**Authentication Endpoints (`/auth/*`):**

- 5 requests per minute per IP
- Significantly reduces brute-force attack surface

### Login Attempt Tracking

- **Failed attempts:** Tracked per email address
- **Lockout threshold:** 5 failed attempts
- **Lockout duration:** 15 minutes
- **Auto-cleanup:** Attempts older than 24 hours are removed

### Implementation

- Rate limiting: `internal/adapters/middleware/rate_limit_auth.go`
- Login tracking: `internal/adapters/middleware/login_attempts.go`

## LLM Security (Prompt Injection Prevention)

### Input Sanitization

All prompts sent to the LLM are sanitized to remove injection patterns:

- "ignore previous instructions"
- "reveal your prompt"
- "forget everything"
- And similar patterns

### Output Validation

LLM responses are validated before being sent to the frontend:

- Check for suspicious patterns
- Length limits (max 5000 characters)
- XSS pattern detection

### Implementation

- Backend: `internal/adapters/secondary/llm/deepseek.go`
- Frontend: `frontend/src/components/dashboard/FastingTimer.tsx` (DOMPurify sanitization)

## CSRF Protection

### Double-Submit Cookie Pattern

- CSRF tokens generated for each session
- Tokens validated on all state-changing requests (POST, PUT, DELETE, PATCH)
- Safe methods (GET, HEAD, OPTIONS) exempt from CSRF checks

### Implementation

CSRF middleware in `internal/adapters/middleware/csrf.go`

## CORS Configuration

### Development

- Default: `http://localhost:5173`

### Production

- Set via `ALLOWED_ORIGINS` environment variable
- Supports multiple comma-separated origins
- Never use `*` (wildcard) in production

## Security Headers

The following security headers are automatically applied:

```
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Content-Security-Policy: [see middleware/security_headers.go]
Strict-Transport-Security: max-age=31536000; includeSubDomains (HTTPS only)
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

### Implementation

Security headers middleware in `internal/adapters/middleware/security_headers.go`

## Error Handling

### Information Disclosure Prevention

- Never expose internal error details to clients
- Log detailed errors server-side only
- Return generic error messages:
  - ❌ `"json: cannot unmarshal string into Go struct field..."`
  - ✅ `"invalid request format"`

### Panic Prevention

- All type assertions must use comma-ok idiom
- Example: `userID, ok := userIDVal.(uuid.UUID)`

## Database Security

### Connection Security

- Use SSL/TLS in production (`sslmode=require`)
- Never hardcode credentials
- Use environment variables

### SQL Injection Prevention

- Use parameterized queries
- Never concatenate user input into SQL
- Use ORM or prepared statements

## API Key Management

### Storage

- Store in environment variables only
- Never commit to version control
- Use secret management service in production (AWS Secrets Manager, HashiCorp Vault)

### Validation

- DeepSeek API keys must start with `sk-` and be ≥20 characters
- Reject placeholder values automatically

### Rotation

- Rotate API keys regularly
- Have revocation process in place

## Secret Scanning

### Git Pre-commit Hooks

```bash
# Install detect-secrets
pip install detect-secrets

# Scan repository
detect-secrets scan --baseline .secrets.baseline

# Audit findings
detect-secrets audit .secrets.baseline
```

### GitHub Secret Scanning

- Enable secret scanning in repository settings
- Configure custom patterns if needed

### CI/CD Integration

Add to CI pipeline:

```bash
# Fail build if secrets detected
detect-secrets scan --baseline .secrets.baseline
```

## Incident Response

### If a Secret is Leaked

1. **Immediate Actions:**
   - Rotate the compromised secret immediately
   - Revoke API keys
   - Review access logs for unauthorized usage

2. **Investigation:**
   - Determine scope of exposure
   - Check for unauthorized API calls
   - Review recent user account activity

3. **Remediation:**
   - Update all deployment environments
   - Force password reset if user credentials affected
   - Notify affected users if required by regulations

4. **Prevention:**
   - Add compromised pattern to secret scanning
   - Review how the leak occurred
   - Update documentation and training

## Security Monitoring

### Logging

- Log all authentication failures
- Log all type assertion failures
- Log suspicious LLM responses
- Log rate limit violations

### Alerts

Set up alerts for:

- Multiple failed login attempts
- Unusual API usage patterns
- Suspicious LLM prompt patterns
- Rate limit threshold exceeded

## Compliance

This security implementation addresses:

| Requirement | Implementation |
|-------------|----------------|
| OWASP A01: Broken Access Control | JWT validation, CSRF protection |
| OWASP A02: Cryptographic Failures | Strong JWT secrets, bcrypt password hashing |
| OWASP A03: Injection | Prompt sanitization, SQL parameterization, input validation |
| OWASP A04: Insecure Design | Rate limiting, account lockout |
| OWASP A07: Authentication Failures | Password policy, MFA-ready architecture |

## Testing

### Security Testing Commands

```bash
# Static analysis
gosec ./...

# Dependency vulnerability checking
govulncheck ./...

# Frontend dependency audit
cd frontend && npm audit

# Run all tests
go test ./... -v
cd frontend && npm test
```

### Manual Testing Checklist

- [ ] Try weak password → should be rejected
- [ ] Try invalid weight values → should be rejected
- [ ] Make 6 login attempts → 6th should be rate-limited
- [ ] Send malformed JSON → should get generic error
- [ ] Test CORS from disallowed origin → should fail
- [ ] Test CSRF token validation → should block without token

## Regular Maintenance

### Monthly

- Review and rotate API keys
- Audit user account security
- Review security logs for anomalies

### Quarterly

- Update dependencies
- Run penetration tests
- Review and update security policies

### Annually

- Third-party security audit
- Update compliance documentation
- Review incident response procedures

## Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [CWE Top 25](https://cwe.mitre.org/top25/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [Go Security Best Practices](https://golang.org/doc/security/best-practices)
