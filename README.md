# 🏨 Hotel Internal API

> **Production-Ready Hotel Management Gateway with Enterprise-Grade Architecture**

A robust, scalable middleware service that serves as the secure gateway between hotel management frontends and backend services. Built with Go and designed for enterprise hotel operations with comprehensive authentication, circuit breaker resilience, and modular architecture.

## 🌟 Overview

The Hotel Internal API is the central orchestration layer in a distributed hotel management ecosystem, providing secure JWT authentication, intelligent request routing, circuit breaker protection, and seamless integration between user interfaces and backend data services.

### 🏗️ Enterprise Architecture

```
┌─────────────────┐    ┌──────────────────────┐    ┌─────────────────┐
│   User Portal   │───▶│  Hotel Internal API  │───▶│  API Beheerder  │
│                 │    │                      │    │                 │
│ • React/Vue     │    │ • JWT Auth           │    │ • Data Layer    │
│ • Admin Panel   │    │ • Circuit Breakers   │    │ • Database      │
│ • Mobile App    │    │ • Health Monitoring  │    │ • CRUD Ops      │
└─────────────────┘    │ • CORS & Middleware  │    └─────────────────┘
                       │ • Request Routing    │    
                       └──────────────────────┘    
                                │                  
                                ▼                  
                       ┌─────────────────┐         
                       │ Central Mgmt    │         
                       │                 │         
                       │ • Business Rules│         
                       │ • Permissions   │         
                       │ • Audit Logs    │         
                       └─────────────────┘         
```

## ✨ Enterprise Features

### 🔒 **Security & Authentication**
- **JWT Authentication**: Enterprise-grade token validation and user context extraction
- **Role-Based Access Control**: Admin, user, and guest permission levels with granular permissions
- **CORS Protection**: Configurable cross-origin policies for web security
- **Request Correlation**: Unique request IDs for distributed tracing and debugging
- **Service-to-Service Auth**: Secure API key authentication for backend services

### 🛡️ **Resilience & Reliability**
- **Circuit Breaker Pattern**: Automatic failure detection and recovery for external services
- **Health Monitoring**: Real-time dependency health checks and status reporting
- **Graceful Degradation**: Intelligent fallback mechanisms when services are unavailable
- **Retry Logic**: Configurable retry strategies with exponential backoff
- **Timeout Management**: Comprehensive timeout handling for all external calls

### 📊 **Observability & Monitoring**
- **Structured Logging**: JSON-formatted logs with correlation IDs and context
- **Prometheus Metrics**: Built-in performance, business, and system metrics
- **Health Endpoints**: Comprehensive health checks for all dependencies
- **Request Tracing**: End-to-end request tracking and performance monitoring
- **Error Tracking**: Detailed error reporting with context and stack traces

### 🏗️ **Modern Architecture**
- **Modular Design**: Clean separation of concerns with `internal/` package structure
- **Handler Pattern**: Dedicated handlers for different functional areas
- **Middleware Chain**: Composable middleware for cross-cutting concerns
- **Configuration Management**: Environment-based configuration with sensible defaults
- **Dependency Injection**: Clean dependency management and enhanced testability

## 🚀 Quick Start

### Prerequisites
- **Go 1.21+**: Modern Go version with generics and improved performance
- **Git**: For version control and repository management
- **API Beheerder**: Data layer service (mock available for development)
- **Central Management**: Business rules service (mock available for development)

### Installation & Setup

1. **Clone and Initialize**
   ```bash
   git clone https://github.com/LarsSonke/InternalAPI.git
   cd InternalAPI
   go mod tidy
   ```

2. **Configure Environment** (copy and customize)
   ```bash
   # Server Configuration
   export HOST=localhost
   export PORT=8080
   export GIN_MODE=release

   # JWT Authentication
   export JWT_SECRET=your-super-secure-jwt-secret-key-here
   
   # External Services
   export API_BEHEERDER_URL=http://localhost:8081
   export API_BEHEERDER_KEY=beheerder-service-key
   export CENTRAL_MGMT_URL=http://localhost:8082
   export CENTRAL_MGMT_KEY=central-mgmt-service-key
   
   # CORS Configuration
   export USER_PORTAL_URL=http://localhost:3000
   export ALLOWED_ORIGINS=http://localhost:3000,https://hotel-portal.example.com
   
   # Monitoring
   export LOG_LEVEL=INFO
   ```

3. **Build and Run**
   ```bash
   # Development
   go run main.go
   
   # Production Build
   go build -ldflags="-s -w" -o hotel-internal-api .
   ./hotel-internal-api
   ```

4. **Verify Installation**
   ```bash
   curl http://localhost:8080/health
   curl http://localhost:8080/metrics
   curl http://localhost:8080/health/circuit-breakers
   ```

## 📚 API Documentation

### 🌐 **Public Endpoints**

| Method | Endpoint | Description | Auth | Response |
|--------|----------|-------------|------|----------|
| `GET` | `/health` | System health check with dependencies | ❌ | Health status |
| `GET` | `/metrics` | Prometheus metrics for monitoring | ❌ | Metrics data |
| `GET` | `/health/circuit-breakers` | Circuit breaker status | ❌ | Breaker states |

### 🔐 **Authentication Endpoints**

| Method | Endpoint | Description | Auth | Response |
|--------|----------|-------------|------|----------|
| `POST` | `/api/auth/login` | User authentication | ❌ | JWT Token |
| `POST` | `/api/auth/refresh` | Token refresh | ✅ JWT | New JWT |
| `POST` | `/api/auth/logout` | User logout | ✅ JWT | Success |

### 🏨 **Hotel Management Endpoints**

| Method | Endpoint | Description | Auth | Response |
|--------|----------|-------------|------|----------|
| `GET` | `/api/albums` | Get hotel bookings/rooms | ✅ JWT | Album list |
| `GET` | `/api/albums/:id` | Get specific booking/room | ✅ JWT | Album details |
| `POST` | `/api/albums` | Create new booking/room | ✅ JWT | Created album |
| `PUT` | `/api/albums/:id` | Update booking/room | ✅ JWT | Updated album |
| `DELETE` | `/api/albums/:id` | Cancel booking/delete room | ✅ JWT | Deletion status |

### 👑 **Admin Endpoints**

| Method | Endpoint | Description | Auth | Response |
|--------|----------|-------------|------|----------|
| `GET` | `/admin/system-status` | System overview | ✅ Admin JWT | System info |
| `GET` | `/admin/audit-logs` | Audit trail | ✅ Admin JWT | Audit data |
| `GET` | `/admin/users` | User management | ✅ Admin JWT | User list |

## 🏗️ Project Structure

```
Hotel Internal API/
├── 📁 internal/                    # Private application code
│   ├── 📁 circuitbreaker/         # Circuit breaker implementation
│   │   └── circuitbreaker.go      # Breaker logic and state management
│   ├── 📁 config/                 # Configuration management
│   │   └── config.go              # Environment and configuration loading
│   ├── 📁 handlers/               # HTTP request handlers
│   │   ├── admin.go               # Admin functionality handlers
│   │   ├── albums.go              # Hotel/booking handlers
│   │   ├── auth.go                # Authentication handlers
│   │   └── health.go              # Health check handlers
│   ├── 📁 middleware/             # HTTP middleware
│   │   └── auth.go                # JWT authentication middleware
│   ├── 📁 models/                 # Data models and types
│   │   └── models.go              # Shared data structures
│   ├── 📁 routes/                 # Route configuration
│   │   └── routes.go              # Route setup and middleware chaining
│   └── 📁 services/               # External service clients
│       └── external.go            # External API communication
├── 📁 mock-beheerder/             # Mock API Beheerder for development
├── 📁 mock-central-mgmt/          # Mock Central Management for testing
├── 📄 main.go                     # Application entry point
├── 📄 go.mod                      # Go module definition
├── 📄 go.sum                      # Dependency checksums
├── 📄 README.md                   # This documentation
└── 📄 MERGE_VERIFICATION.md       # Architecture merge documentation
```

## 🔍 Request Flow Architecture

### 🔄 **Authenticated Request Lifecycle**

```
1. User Portal Request
   ↓
2. CORS Validation (middleware)
   ↓
3. Request ID Generation (middleware)
   ↓
4. JWT Authentication (middleware)
   ↓
5. Permission Validation (Central Management)
   ↓
6. Business Logic Processing (handlers)
   ↓
7. External Service Calls (API Beheerder)
   ↓ (with circuit breaker protection)
8. Response Formatting
   ↓
9. Audit Logging (async)
   ↓
10. Response to User Portal
```

### 📊 **Example Request with Monitoring**

```json
{
  "request_id": "req_1696248615_abc123",
  "user_id": "user_456",
  "action": "create_booking",
  "timestamp": "2025-10-09T14:30:15Z",
  "validation": {
    "jwt": "✅ valid",
    "permissions": "✅ allowed",
    "rate_limit": "✅ within_limits"
  },
  "external_calls": {
    "central_mgmt": {
      "status": "✅ success",
      "duration_ms": 45,
      "circuit_breaker": "closed"
    },
    "api_beheerder": {
      "status": "✅ success", 
      "duration_ms": 125,
      "circuit_breaker": "closed"
    }
  },
  "response_time_ms": 180,
  "status": "success"
}
```

## 📊 Monitoring & Observability

### 🏥 **Health Check Response**
```json
{
  "status": "healthy",
  "service": "hotel-internal-api",
  "version": "2.0.0",
  "timestamp": "2025-10-09T14:30:15Z",
  "dependencies": {
    "api_beheerder": {
      "status": "healthy",
      "response_time_ms": 45,
      "circuit_breaker": "closed"
    },
    "central_management": {
      "status": "healthy", 
      "response_time_ms": 32,
      "circuit_breaker": "closed"
    }
  },
  "system": {
    "memory_usage": "45MB",
    "goroutines": 23,
    "uptime": "2h30m15s"
  }
}
```

### 📈 **Prometheus Metrics**

#### **HTTP Metrics**
- `hotel_api_requests_total{method,endpoint,status}` - Total HTTP requests
- `hotel_api_request_duration_seconds{method,endpoint}` - Request duration histogram
- `hotel_api_active_requests` - Currently active requests

#### **Business Metrics**
- `hotel_bookings_created_total` - Total bookings created
- `hotel_auth_attempts_total{result}` - Authentication attempts
- `hotel_admin_actions_total{action}` - Admin actions performed

#### **External Service Metrics**
- `hotel_external_calls_total{service,method,status}` - External API calls
- `hotel_external_duration_seconds{service}` - External service response times
- `hotel_circuit_breaker_state{service}` - Circuit breaker states

#### **System Metrics**
- `hotel_api_uptime_seconds` - Service uptime
- `hotel_api_memory_usage_bytes` - Memory consumption
- `hotel_api_goroutines_active` - Active goroutines

## 🔧 Configuration Reference

### 🌍 **Environment Variables**

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `HOST` | `localhost` | Server bind address | `0.0.0.0` |
| `PORT` | `8080` | Server port | `8080` |
| `GIN_MODE` | `debug` | Gin framework mode | `release` |
| `JWT_SECRET` | `your-jwt-secret-key` | JWT signing secret | `super-secure-key-2025` |
| `API_BEHEERDER_URL` | `http://localhost:8081` | Data service URL | `https://api.hotel.com` |
| `API_BEHEERDER_KEY` | `beheerder-service-key` | Data service auth key | `bhr_sk_live_xxx` |
| `CENTRAL_MGMT_URL` | `http://localhost:8082` | Management service URL | `https://mgmt.hotel.com` |
| `CENTRAL_MGMT_KEY` | `central-mgmt-service-key` | Management auth key | `cmg_sk_live_xxx` |
| `USER_PORTAL_URL` | `http://localhost:3000` | Frontend URL for CORS | `https://portal.hotel.com` |
| `ALLOWED_ORIGINS` | `*` | CORS allowed origins | `https://portal.hotel.com,https://admin.hotel.com` |
| `LOG_LEVEL` | `INFO` | Logging level | `DEBUG,INFO,WARN,ERROR` |

### ⚙️ **Circuit Breaker Configuration**

```go
// Circuit breaker settings (internal/circuitbreaker/circuitbreaker.go)
MaxRequests:    3,     // Max requests in half-open state
Interval:       30s,   // Reset interval for failure counting  
Timeout:        60s,   // Time to wait before attempting reset
ReadyToTrip:    5,     // Failures needed to trip breaker
OnStateChange:  func() // Callback for state changes
```

## 🧪 Testing & Development

### 🚀 **Development Workflow**

1. **Start Mock Services**
   ```bash
   # Terminal 1: Mock API Beheerder
   cd mock-beheerder
   go run main.go
   
   # Terminal 2: Mock Central Management
   cd mock-central-mgmt  
   go run main.go
   
   # Terminal 3: Hotel Internal API
   go run main.go
   ```

2. **Run Test Suite**
   ```bash
   # Unit tests
   go test ./internal/...
   
   # Integration tests  
   go test -tags=integration ./...
   
   # Load tests
   go test -tags=load ./...
   ```

3. **Manual API Testing**
   ```bash
   # Health checks
   curl http://localhost:8080/health
   
   # Authentication flow
   curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"admin123"}'
   
   # Protected endpoint
   curl http://localhost:8080/api/albums \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

### 🏗️ **Build & Deployment**

#### **Local Development Build**
```bash
go build -o hotel-api-dev .
```

#### **Production Build**
```bash
go build -ldflags="-s -w -X main.version=2.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o hotel-api .
```

#### **Docker Deployment**
```dockerfile
# Multi-stage Docker build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o hotel-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/hotel-api .
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
CMD ["./hotel-api"]
```

#### **Kubernetes Deployment**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hotel-internal-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hotel-internal-api
  template:
    metadata:
      labels:
        app: hotel-internal-api
    spec:
      containers:
      - name: api
        image: hotel/internal-api:2.0.0
        ports:
        - containerPort: 8080
        env:
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: hotel-secrets
              key: jwt-secret
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

## 🔒 Security Best Practices

### 🛡️ **Production Security Checklist**

- [ ] **JWT Secret**: Use cryptographically secure, randomly generated JWT secrets (256-bit minimum)
- [ ] **HTTPS Only**: Deploy behind TLS termination, never expose HTTP in production
- [ ] **CORS Policy**: Configure strict CORS policies, avoid wildcards in production
- [ ] **Rate Limiting**: Implement request rate limiting to prevent abuse
- [ ] **Input Validation**: Validate and sanitize all inputs, use structured logging
- [ ] **Error Handling**: Never expose internal errors or stack traces to clients
- [ ] **Audit Logging**: Log all authentication events and sensitive operations
- [ ] **Dependency Updates**: Keep dependencies updated, regularly scan for vulnerabilities
- [ ] **Access Control**: Implement principle of least privilege for service accounts
- [ ] **Secret Management**: Use secure secret management (Vault, K8s secrets, etc.)

### 🔐 **JWT Token Configuration**

```go
// Example JWT claims structure
{
  "sub": "user_12345",           // Subject (user ID)
  "name": "John Manager",        // User display name
  "roles": ["user", "manager"],  // User roles for RBAC
  "hotel_id": "hotel_456",       // Associated hotel
  "permissions": [               // Granular permissions
    "bookings:read",
    "bookings:write", 
    "reports:read"
  ],
  "iat": 1696248615,            // Issued at
  "exp": 1696252215,            // Expires at (1 hour)
  "iss": "hotel-auth-service",   // Issuer
  "aud": "hotel-internal-api"    // Audience
}
```

## 🏨 Hotel Management Ecosystem

### 🌐 **System Components**

| Component | Repository | Description | Status |
|-----------|------------|-------------|---------|
| **Hotel Internal API** | `LarsSonke/InternalAPI` | This service - Gateway & orchestration | ✅ Active |
| **User Portal** | `hotel/user-portal` | React frontend for staff & management | 🔄 Development |
| **API Beheerder** | `hotel/api-beheerder` | Data layer service & database operations | 🔄 Development |
| **Central Management** | `hotel/central-mgmt` | Business rules, permissions & audit | 🔄 Development |
| **Plugin System** | `hotel/plugins` | Extensible third-party integrations | 📋 Planned |

### 🔄 **Data Flow Overview**

```
Guest Booking Request
        ↓
┌─────────────────┐
│   User Portal   │ ← Staff manages bookings, checks availability
└─────────────────┘
        ↓ HTTPS + JWT
┌─────────────────┐
│ Internal API    │ ← Authentication, authorization, orchestration  
│ (This Service)  │
└─────────────────┘
    ↓           ↓
┌─────────────┐ ┌─────────────────┐
│Central Mgmt │ │  API Beheerder  │ ← Database operations
│Business     │ │  Data Layer     │
│Rules        │ │                 │
└─────────────┘ └─────────────────┘
```

### 🔌 **Integration Points**

#### **Frontend Integration**
```javascript
// User Portal API client example
const apiClient = {
  baseURL: 'https://api.hotel.com',
  auth: {
    login: async (credentials) => {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(credentials)
      });
      return response.json();
    }
  },
  bookings: {
    list: async (token) => {
      const response = await fetch('/api/albums', {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      return response.json();
    }
  }
};
```

#### **Backend Service Integration**
```go
// External service client example
type APIBeheerderClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

func (c *APIBeheerderClient) CreateBooking(ctx context.Context, booking *Booking) error {
    req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/bookings", nil)
    req.Header.Set("X-Service-Key", c.apiKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return fmt.Errorf("api beheerder call failed: %w", err)
    }
    defer resp.Body.Close()
    
    return nil
}
```

## 🤝 Contributing

### 📋 **Development Guidelines**

1. **Fork & Branch**
   ```bash
   git fork https://github.com/LarsSonke/InternalAPI.git
   git checkout -b feature/amazing-hotel-feature
   ```

2. **Code Standards**
   - Follow Go effective practices and idioms
   - Use `gofmt` and `golint` for consistent formatting
   - Write comprehensive tests (aim for >80% coverage)
   - Add JSDoc-style comments for public functions
   - Use structured logging with context

3. **Testing Requirements**
   ```bash
   # Run all tests before submitting
   go test ./... -v
   go test -race ./...  # Race condition detection
   go test -bench=. ./... # Benchmark tests
   ```

4. **Commit & Submit**
   ```bash
   git commit -m "feat(bookings): add real-time availability checking
   
   - Implement WebSocket connection for live updates
   - Add booking conflict detection
   - Update API documentation
   
   Closes #123"
   
   git push origin feature/amazing-hotel-feature
   # Create Pull Request on GitHub
   ```

### 🔍 **Code Review Checklist**

- [ ] **Security**: No secrets in code, proper input validation
- [ ] **Performance**: Efficient algorithms, proper resource management
- [ ] **Testing**: Unit tests, integration tests, edge cases covered
- [ ] **Documentation**: Updated README, code comments, API docs
- [ ] **Backwards Compatibility**: No breaking changes without version bump
- [ ] **Error Handling**: Comprehensive error handling and logging

## 📝 License & Legal

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for complete details.

```
MIT License

Copyright (c) 2025 Lars Sonke

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

## 📧 Contact & Support

### 👨‍💻 **Maintainer**
- **Name**: Lars Sonke
- **GitHub**: [@LarsSonke](https://github.com/LarsSonke)
- **Project**: [InternalAPI](https://github.com/LarsSonke/InternalAPI)

### 🆘 **Support Channels**
- **GitHub Issues**: [Report bugs or request features](https://github.com/LarsSonke/InternalAPI/issues)
- **GitHub Discussions**: [Community discussions and Q&A](https://github.com/LarsSonke/InternalAPI/discussions)
- **Documentation**: [Wiki and detailed guides](https://github.com/LarsSonke/InternalAPI/wiki)

### 🏷️ **Project Status**

![Version](https://img.shields.io/badge/version-2.0.0-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)
![Coverage](https://img.shields.io/badge/coverage-85%25-yellowgreen.svg)

---

⭐ **If you find this project useful, please consider giving it a star on GitHub!**

🏨 **Built for the future of hotel management - scalable, secure, and developer-friendly.**