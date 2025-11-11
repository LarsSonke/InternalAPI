# Security & Production Hardening Documentation

## Overview
This document describes the security and production hardening features implemented in the Internal API.

## Security Features Implemented

### 1. Rate Limiting
**Purpose**: Prevent abuse, brute force attacks, and DDoS.

**Implementation**:
- **Token Bucket Algorithm**: Efficient, memory-safe rate limiting
- **Multi-tier Limits**: Different limits for different endpoint types
  - Login endpoints: 5 requests per 5 minutes (configurable)
  - General API: 100 requests per minute (configurable)
  - Admin endpoints: 50 requests per minute (configurable)
- **Per-IP and Per-User**: Tracks both authenticated users and IPs
- **Automatic Cleanup**: Removes stale rate limit buckets

**Configuration**:
```bash
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_INTERVAL_SECONDS=60
LOGIN_RATE_LIMIT_REQUESTS=5
LOGIN_RATE_LIMIT_INTERVAL_SECONDS=300
```

**Response on Rate Limit**:
```json
{
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "Too many requests. Please try again later.",
  "retry_after": 60
}
```

### 2. Request Size Limits
**Purpose**: Prevent memory exhaustion and upload bombing attacks.

**Implementation**:
- Maximum request body size: 5MB (default, configurable)
- Applied to all incoming requests
- Returns 413 Request Entity Too Large when exceeded

**Configuration**:
```bash
MAX_REQUEST_BODY_SIZE=5242880  # 5MB in bytes
```

### 3. Input Validation
**Purpose**: Prevent injection attacks and ensure data integrity.

**Implementation**:
- Comprehensive validation tags on all models
- String length limits (min/max)
- Email format validation
- Alphanumeric constraints where appropriate
- Array element validation (dive)

**Examples**:
```go
Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
Email    string `json:"email" binding:"required,email,max=100"`
Password string `json:"password" binding:"required,min=8,max=100"`
```

### 4. JWT Token Security
**Purpose**: Secure authentication with proper token validation.

**Implementation**:
- **Signature Verification**: HMAC-SHA256 signing
- **Expiration Checking**: Automatic token expiry validation
- **Token Blacklisting**: Support for revoked tokens
- **Automatic Cleanup**: Removes expired blacklisted tokens

**Features**:
- ✅ Token signature validation
- ✅ Expiry time checking
- ✅ Blacklist support for logout
- ✅ Claims extraction and validation
- ✅ Role-based access control

**Configuration**:
```bash
JWT_SECRET=your-super-secret-jwt-key  # ⚠️ MUST change in production!
```

### 5. Security Headers
**Purpose**: Protect against common web vulnerabilities.

**Headers Implemented**:
- `X-Frame-Options: DENY` - Prevents clickjacking
- `X-Content-Type-Options: nosniff` - Prevents MIME sniffing
- `X-XSS-Protection: 1; mode=block` - XSS protection
- `Content-Security-Policy` - Restricts resource loading
- `Referrer-Policy: strict-origin-when-cross-origin` - Referrer control
- `Permissions-Policy` - Feature access control
- `Strict-Transport-Security` (when HTTPS enabled) - Forces HTTPS

**Configuration**:
```bash
ENABLE_SECURITY_HEADERS=true
```

### 6. Request ID Tracking
**Purpose**: Trace requests through the system for debugging and security auditing.

**Implementation**:
- Unique UUID generated for each request
- Included in response headers: `X-Request-ID`
- Logged with all audit entries
- Can be provided by client or auto-generated

### 7. Audit Logging
**Purpose**: Complete security audit trail for compliance and forensics.

**Logged Information**:
- Request ID
- Timestamp
- HTTP method and path
- Client IP address
- User agent
- Authenticated user ID
- Response status code
- Request/response duration
- Request/response sizes
- Request body (sanitized, excludes passwords)

**Log Levels**:
- INFO: Successful requests (2xx, 3xx)
- WARN: Client errors (4xx)
- ERROR: Server errors (5xx)

**Configuration**:
```bash
ENABLE_AUDIT_LOGGING=true
```

**Sample Audit Log**:
```json
{
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": 1699632000,
  "method": "POST",
  "path": "/api/v1/albums",
  "ip": "192.168.1.100",
  "user_id": "user123",
  "status": 201,
  "duration_ms": 45,
  "level": "info"
}
```

### 8. Timeout Configuration
**Purpose**: Prevent resource exhaustion and hanging connections.

**Timeouts Implemented**:
- **Read Timeout**: 15 seconds (default) - Max time to read request
- **Write Timeout**: 15 seconds (default) - Max time to write response
- **Idle Timeout**: 60 seconds (default) - Max idle time for keep-alive
- **Request Timeout**: 30 seconds (default) - Total request duration

**Configuration**:
```bash
READ_TIMEOUT_SECONDS=15
WRITE_TIMEOUT_SECONDS=15
IDLE_TIMEOUT_SECONDS=60
REQUEST_TIMEOUT_SECONDS=30
```

### 9. Graceful Shutdown
**Purpose**: Allow in-flight requests to complete before shutdown.

**Implementation**:
- Listens for SIGINT/SIGTERM signals
- 30-second grace period for active requests
- Prevents data loss during deployment
- Clean connection closure

**Behavior**:
1. Receive shutdown signal
2. Stop accepting new requests
3. Wait up to 30s for active requests
4. Force shutdown if timeout exceeded
5. Log shutdown status

### 10. Pagination
**Purpose**: Prevent resource exhaustion from large dataset queries.

**Implementation**:
- Default page size: 20 items
- Maximum page size: 100 items (hard limit)
- Query parameters: `page` and `page_size`
- Standardized response format

**Usage**:
```bash
GET /api/v1/albums?page=1&page_size=20
```

**Response**:
```json
{
  "data": [...],
  "page": 1,
  "page_size": 20,
  "total_pages": 5,
  "total_items": 87
}
```

## Production Deployment Checklist

### Critical Security Items
- [ ] **Change JWT_SECRET** - Use strong random string (32+ chars)
- [ ] **Enable HTTPS/TLS** - Configure reverse proxy or load balancer
- [ ] **Restrict CORS Origins** - Remove localhost, add production domains
- [ ] **Review Rate Limits** - Adjust based on expected traffic
- [ ] **Set Strong Service Keys** - Change all API_BEHEERDER_KEY and CENTRAL_MGMT_KEY
- [ ] **Enable All Security Features** - Verify ENABLE_SECURITY_HEADERS=true
- [ ] **Configure Audit Log Storage** - Set up log aggregation (ELK, Splunk, etc.)

### Infrastructure
- [ ] Configure load balancer with health checks (`/health`)
- [ ] Set up TLS termination at load balancer
- [ ] Configure firewall rules (restrict to known IPs if possible)
- [ ] Set up monitoring and alerting (Prometheus metrics at `/metrics`)
- [ ] Configure log rotation and retention
- [ ] Set up database backups (if applicable)
- [ ] Implement secrets management (Vault, AWS Secrets Manager, etc.)

### Monitoring
- [ ] Monitor rate limit violations
- [ ] Track authentication failures
- [ ] Alert on 5xx error rates
- [ ] Monitor circuit breaker states
- [ ] Track request latencies
- [ ] Monitor resource usage (CPU, memory)

### Testing
- [ ] Load testing with realistic traffic patterns
- [ ] Penetration testing
- [ ] Verify rate limits work as expected
- [ ] Test graceful shutdown under load
- [ ] Verify all timeouts trigger correctly
- [ ] Test CORS restrictions

## Security Best Practices

### 1. API Key Rotation
Regularly rotate service keys:
```bash
# Generate new keys periodically
API_BEHEERDER_KEY=$(openssl rand -base64 32)
CENTRAL_MGMT_KEY=$(openssl rand -base64 32)
```

### 2. JWT Secret Management
Never commit JWT secrets to version control:
```bash
# Generate strong secret
JWT_SECRET=$(openssl rand -base64 64)
```

### 3. HTTPS Only in Production
Always use HTTPS in production. Update CORS config:
```bash
CORS_ORIGINS=https://your-domain.com,https://www.your-domain.com
```

Uncomment HSTS header in `middleware/security.go`:
```go
c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
```

### 4. Regular Security Updates
- Keep Go version updated
- Update dependencies regularly: `go get -u ./...`
- Monitor security advisories
- Run `go audit` or similar tools

### 5. Defense in Depth
This API implements multiple security layers:
- Network (firewall, VPC)
- Transport (TLS)
- Application (rate limiting, validation)
- Authentication (JWT)
- Authorization (RBAC)
- Audit (logging)

## Rate Limit Recommendations

### Development
```bash
RATE_LIMIT_REQUESTS=1000
LOGIN_RATE_LIMIT_REQUESTS=20
```

### Production
```bash
RATE_LIMIT_REQUESTS=100        # Adjust based on legitimate traffic
LOGIN_RATE_LIMIT_REQUESTS=5    # Strict to prevent brute force
ADMIN_RATE_LIMIT_REQUESTS=50   # Lower for admin operations
```

### High Traffic Production
```bash
RATE_LIMIT_REQUESTS=500        # Higher for known high-traffic
LOGIN_RATE_LIMIT_REQUESTS=10   # Slightly relaxed but still protective
```

## Incident Response

### Rate Limit Violations
1. Check audit logs for patterns
2. Identify source IP
3. Investigate user behavior
4. Consider IP blocking if malicious

### Authentication Failures
1. Review audit logs for user_id
2. Check for credential stuffing patterns
3. Consider temporary account lockout
4. Alert user of suspicious activity

### 5xx Errors
1. Check circuit breaker status
2. Verify external service health
3. Review application logs
4. Scale if resource constrained

## Performance Impact

### Rate Limiting
- Memory: ~100 bytes per active IP/user
- CPU: Negligible (<0.1% overhead)
- Cleanup: Runs hourly, minimal impact

### Audit Logging
- Disk I/O: Moderate (asynchronous writes recommended)
- CPU: Low (<1% overhead)
- Storage: Plan for log rotation

### JWT Validation
- CPU: Low (HMAC validation is fast)
- Memory: Blacklist grows with active logouts (auto-cleanup)

### Security Headers
- Negligible performance impact

## Compliance Notes

These features support compliance with:
- **GDPR**: Audit logging, data access controls
- **PCI DSS**: Encryption, access controls, logging
- **SOC 2**: Security controls, monitoring, audit trails
- **ISO 27001**: Information security management

## Additional Resources

- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- [OWASP Cheat Sheet Series](https://cheatsheetseries.owasp.org/)
- [Go Security Policy](https://golang.org/security)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
