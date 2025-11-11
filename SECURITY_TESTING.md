# Security Features Testing Guide

## Testing Rate Limiting

### Test Login Rate Limit (5 requests per 5 minutes)

```powershell
# Try 6 login attempts rapidly - the 6th should be blocked
for ($i=1; $i -le 6; $i++) {
    Write-Host "Attempt $i"
    curl -X POST http://localhost:8080/auth/login `
        -H "Content-Type: application/json" `
        -d '{"username":"testuser","password":"password123"}'
    Start-Sleep -Milliseconds 500
}
```

**Expected Result**: First 5 succeed, 6th returns:
```json
{
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "Too many login attempts. Please try again later.",
  "retry_after": 300
}
```

### Test General API Rate Limit (100 requests per minute)

```powershell
# Test authenticated endpoint rate limiting
$token = "your-jwt-token-here"

for ($i=1; $i -le 105; $i++) {
    Write-Host "Request $i"
    curl http://localhost:8080/api/v1/albums `
        -H "Authorization: Bearer $token"
    Start-Sleep -Milliseconds 100
}
```

**Expected Result**: First 100 succeed, 101+ return rate limit error.

## Testing Security Headers

```powershell
# Check security headers in response
curl -I http://localhost:8080/health
```

**Expected Headers**:
```
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'; script-src 'self'; object-src 'none';
Referrer-Policy: strict-origin-when-cross-origin
X-Request-ID: <unique-uuid>
```

## Testing Request Size Limit

```powershell
# Test with large payload (should fail at 5MB+)
$largeData = "a" * 6000000  # 6MB of data

curl -X POST http://localhost:8080/api/v1/albums `
    -H "Content-Type: application/json" `
    -H "Authorization: Bearer $token" `
    -d "{\"title\":\"$largeData\",\"artist\":\"test\",\"price\":10}"
```

**Expected Result**: 413 Request Entity Too Large or connection reset

## Testing Input Validation

### Invalid Email
```powershell
curl -X POST http://localhost:8080/admin/users `
    -H "Content-Type: application/json" `
    -H "Authorization: Bearer $adminToken" `
    -d '{
        "username": "newuser",
        "email": "not-an-email",
        "password": "password123",
        "roles": ["user"]
    }'
```

**Expected Result**: 400 Bad Request with validation error

### Password Too Short
```powershell
curl -X POST http://localhost:8080/admin/users `
    -H "Content-Type: application/json" `
    -H "Authorization: Bearer $adminToken" `
    -d '{
        "username": "newuser",
        "email": "user@example.com",
        "password": "short",
        "roles": ["user"]
    }'
```

**Expected Result**: 400 Bad Request - password must be at least 8 characters

### Username Too Long
```powershell
$longUsername = "a" * 100  # 100 characters

curl -X POST http://localhost:8080/auth/login `
    -H "Content-Type: application/json" `
    -d "{\"username\":\"$longUsername\",\"password\":\"password123\"}"
```

**Expected Result**: 400 Bad Request - username exceeds maximum length

## Testing JWT Token Validation

### Test with Expired Token
```powershell
# Use an expired token (you'll need to generate one or wait for expiry)
$expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl http://localhost:8080/api/v1/albums `
    -H "Authorization: Bearer $expiredToken"
```

**Expected Result**: 401 Unauthorized - "Token has expired"

### Test with Invalid Signature
```powershell
# Tamper with a valid token (change last few characters)
$tamperedToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzIn0.TAMPERED"

curl http://localhost:8080/api/v1/albums `
    -H "Authorization: Bearer $tamperedToken"
```

**Expected Result**: 401 Unauthorized - "Token validation failed"

### Test Token Blacklisting (Logout)
```powershell
$validToken = "your-valid-token"

# 1. Use token successfully
curl http://localhost:8080/api/v1/auth/me `
    -H "Authorization: Bearer $validToken"

# 2. Logout (blacklist token)
curl -X POST http://localhost:8080/api/v1/auth/logout `
    -H "Authorization: Bearer $validToken"

# 3. Try to use same token again (should fail)
curl http://localhost:8080/api/v1/auth/me `
    -H "Authorization: Bearer $validToken"
```

**Expected Result**: Third request returns 401 - "token has been revoked"

## Testing Request ID Tracking

```powershell
# Make a request and capture the Request ID
$response = curl -i http://localhost:8080/api/v1/albums `
    -H "Authorization: Bearer $token"

# Look for X-Request-ID in response headers
# Then check audit logs for this request ID
```

## Testing Pagination

### Default Pagination
```powershell
curl "http://localhost:8080/api/v1/albums?page=1&page_size=20" `
    -H "Authorization: Bearer $token"
```

**Expected Response**:
```json
{
  "data": [...],
  "page": 1,
  "page_size": 20,
  "total_pages": 5,
  "total_items": 87
}
```

### Max Page Size Enforcement
```powershell
# Try to request more than max (100)
curl "http://localhost:8080/api/v1/albums?page=1&page_size=200" `
    -H "Authorization: Bearer $token"
```

**Expected Result**: page_size capped at 100

## Testing Graceful Shutdown

```powershell
# Terminal 1: Start the server
.\internal-api.exe

# Terminal 2: Make a long-running request (if you have one)
# Then in Terminal 1: Press Ctrl+C

# Server should:
# 1. Log "Shutting down server gracefully..."
# 2. Wait for active requests (up to 30s)
# 3. Log "Server exited"
```

## Testing Audit Logging

```powershell
# Make various requests and check the logs
curl -X POST http://localhost:8080/auth/login `
    -H "Content-Type: application/json" `
    -d '{"username":"testuser","password":"password123"}'

# Check console output or log file for JSON audit entries with:
# - request_id
# - timestamp
# - method, path
# - ip, user_agent
# - user_id (if authenticated)
# - status, duration_ms
```

## Load Testing (Optional)

### Using Apache Bench
```powershell
# Test rate limiting under load
ab -n 200 -c 10 http://localhost:8080/health
```

### Using PowerShell
```powershell
# Concurrent requests test
$jobs = 1..50 | ForEach-Object {
    Start-Job -ScriptBlock {
        curl http://localhost:8080/health
    }
}

$jobs | Wait-Job | Receive-Job
$jobs | Remove-Job
```

## Security Checklist Before Production

### Configuration
- [ ] Changed JWT_SECRET to strong random value
- [ ] Updated CORS_ORIGINS to production domains only
- [ ] Rotated all service API keys
- [ ] Reviewed and adjusted rate limits
- [ ] Configured appropriate timeouts for your use case

### Testing
- [ ] Verified rate limiting works correctly
- [ ] Tested with expired/invalid JWT tokens
- [ ] Confirmed input validation rejects bad data
- [ ] Verified request size limits
- [ ] Tested graceful shutdown
- [ ] Confirmed audit logs are being written

### Monitoring
- [ ] Set up log aggregation
- [ ] Configured alerts for rate limit violations
- [ ] Monitoring authentication failures
- [ ] Tracking 5xx error rates
- [ ] Monitoring circuit breaker states

### Infrastructure
- [ ] HTTPS/TLS configured
- [ ] Firewall rules in place
- [ ] Load balancer health checks configured
- [ ] Secrets stored securely (not in code)

## Common Issues & Solutions

### Issue: Rate limit too strict for legitimate users
**Solution**: Increase `RATE_LIMIT_REQUESTS` or `RATE_LIMIT_INTERVAL_SECONDS`

### Issue: Timeouts occurring frequently
**Solution**: Increase timeout values or optimize external service calls

### Issue: Audit logs taking too much disk space
**Solution**: Configure log rotation, reduce log retention period

### Issue: Valid tokens being rejected
**Solution**: Ensure JWT_SECRET is the same across all instances, check token expiry

### Issue: Request size limit too small
**Solution**: Increase `MAX_REQUEST_BODY_SIZE` (be careful with memory implications)

## Performance Benchmarks

Expected overhead from security features:
- Rate Limiting: <0.1% CPU overhead
- JWT Validation: ~1-2ms per request
- Audit Logging: ~0.5ms per request (async recommended)
- Security Headers: Negligible
- Input Validation: ~0.1-0.5ms per request

Total expected overhead: ~3-5ms per request on average hardware.
