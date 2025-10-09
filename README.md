# 🏨 Hotel Internal API

A secure, scalable middleware service for hotel management systems, providing authentication, authorization, and request routing between the User Portal and backend services.

## � Overview

The Hotel Internal API serves as the central gateway in a distributed hotel management architecture, handling JWT authentication, permission validation, and seamless communication with data and business logic services.

### Architecture

```
User Portal → Internal API → [API Beheerder + Central Management]
     ↑              ↑                    ↑              ↑
   Frontend     Middleware           Data Layer    Business Rules
```

## ✨ Features

### 🔒 Security First
- **JWT Authentication**: Secure token validation for User Portal access
- **CORS Protection**: Configurable cross-origin resource sharing
- **Request Correlation**: Unique request IDs for tracing and debugging
- **Permission Validation**: Integration with Central Management for authorization

### 📊 Production Ready
- **Structured Logging**: JSON logging with Logrus for observability
- **Prometheus Metrics**: Built-in monitoring and performance tracking
- **Health Checks**: Dependency monitoring for API Beheerder and Central Management
- **Error Handling**: Standardized error responses with codes and timestamps

### 🏗️ Modular Architecture
- **Clean Code Structure**: Separated into focused modules for maintainability
- **Microservice Pattern**: Designed for distributed hotel management systems
- **Scalable Design**: Easy to extend and modify for growing requirements

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Access to API Beheerder service
- Access to Central Management service

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/LarsSonke/InternalAPI.git
   cd InternalAPI
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables** (optional)
   ```bash
   export HOST=localhost
   export PORT=8080
   export JWT_SECRET=your-secret-key
   export API_BEHEERDER_URL=http://localhost:8081
   export CENTRAL_MGMT_URL=http://localhost:8082
   ```

4. **Build and run**
   ```bash
   go build -o internal-api .
   ./internal-api
   ```

## 📚 API Documentation

### Public Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Health check with dependency status | ❌ |
| GET | `/metrics` | Prometheus metrics | ❌ |

### User Portal Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/albums` | Get hotel bookings/rooms | ✅ JWT |
| GET | `/api/v1/albums/:id` | Get specific booking/room | ✅ JWT |
| POST | `/api/v1/albums` | Create new booking/room | ✅ JWT |
| PUT | `/api/v1/albums/:id` | Update booking/room | ✅ JWT |
| DELETE | `/api/v1/albums/:id` | Cancel booking/delete room | ✅ JWT |

### Admin Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/admin/system-status` | System status overview | ✅ JWT |
| GET | `/admin/audit-logs` | Audit trail for compliance | ✅ JWT |

## 🔧 Configuration

The API can be configured through environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | `localhost` | Server host |
| `PORT` | `8080` | Server port |
| `JWT_SECRET` | `your-jwt-secret-key` | JWT signing secret |
| `API_BEHEERDER_URL` | `http://localhost:8081` | API Beheerder service URL |
| `API_BEHEERDER_KEY` | `beheerder-service-key` | API Beheerder authentication key |
| `CENTRAL_MGMT_URL` | `http://localhost:8082` | Central Management service URL |
| `CENTRAL_MGMT_KEY` | `central-mgmt-service-key` | Central Management authentication key |
| `USER_PORTAL_URL` | `http://localhost:3000` | User Portal URL for CORS |
| `ALLOWED_ORIGINS` | `http://localhost:3000,http://localhost:3001,https://hotel-portal.local` | CORS allowed origins |
| `LOG_LEVEL` | `INFO` | Logging level (DEBUG, INFO, WARN, ERROR) |

## 🏗️ Project Structure

```
.
├── main.go           # Server setup, routing, and external service calls
├── auth.go           # JWT validation and user extraction
├── config.go         # Configuration management
├── handlers.go       # API endpoint handlers
├── middleware.go     # Request ID, metrics, and authentication middleware
├── monitoring.go     # Health checks, metrics, and logging setup
├── go.mod           # Go module dependencies
└── go.sum           # Dependency checksums
```

## 🔍 Request Flow

### Authenticated Request Example

1. **User Portal** sends request with JWT token
2. **Internal API** validates JWT and extracts user information
3. **Permission Check** with Central Management System
4. **Data Request** forwarded to API Beheerder
5. **Response** processed and returned to User Portal
6. **Audit Log** recorded for compliance

```json
{
  "request_id": "12345-67890-abcdef",
  "user_id": "user_123",
  "action": "get_bookings",
  "permission_check": "✅ allowed",
  "data_service": "✅ success",
  "response_time": "89ms"
}
```

## 📊 Monitoring

### Health Check Response
```json
{
  "status": "healthy",
  "service": "internal-api",
  "version": "1.0.0",
  "dependencies": {
    "api_beheerder": {
      "status": "healthy",
      "duration": 45
    },
    "central_management": {
      "status": "healthy",
      "duration": 32
    }
  }
}
```

### Prometheus Metrics
- `internal_api_requests_total` - Total HTTP requests
- `internal_api_request_duration_seconds` - Request duration
- `internal_api_external_calls_total` - External service calls
- `internal_api_external_duration_seconds` - External service duration

## 🛠️ Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -ldflags="-s -w" -o internal-api .
```

### Docker Support
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o internal-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/internal-api .
EXPOSE 8080
CMD ["./internal-api"]
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🏨 Hotel Management Ecosystem

This Internal API is part of a larger hotel management system:

- **User Portal**: React frontend for hotel staff and management
- **Internal API**: This service - secure middleware and authentication
- **API Beheerder**: Data layer service handling database operations
- **Central Management**: Business rules, permissions, and audit logging
- **Plugin System**: Extensible architecture for third-party integrations

## 📧 Contact

**Author**: Lars Sonke  
**GitHub**: [@LarsSonke](https://github.com/LarsSonke)  
**Project**: [InternalAPI](https://github.com/LarsSonke/InternalAPI)

---

⭐ **If you find this project useful, please consider giving it a star!**
- **Permission Checking**: Before any operation, check if user is allowed
- **Business Rules**: Dynamic rules fetched from Central Management
- **Audit Logging**: All actions logged for compliance and security
- **User Filtering**: Role-based data filtering and access control
- **Configuration**: Dynamic system configuration

### ✅ **API Beheerder Integration**
- HTTP client to call API Beheerder for all data operations
- Service-to-service authentication with API keys
- Complete error handling and response formatting

### ✅ **Orchestrated Business Logic**
Your Internal API now follows this pattern for every operation:

1. **Authenticate User** (JWT validation)
2. **Check Permissions** (Central Management)
3. **Get Business Rules** (Central Management)
4. **Validate Input** (Your business logic)
5. **Perform Data Operation** (API Beheerder)
6. **Log Action** (Central Management - async)
7. **Return Response** (User Portal)

## 📋 **Configuration**

### **Environment Variables:**
```bash
# Your API settings
PORT=8080
HOST=localhost

# JWT authentication
JWT_SECRET=your-jwt-secret-key

# API Beheerder connection (data operations)
API_BEHEERDER_URL=http://localhost:8081
API_BEHEERDER_KEY=beheerder-service-key

# Central Management System connection (business rules, permissions, audit)
CENTRAL_MGMT_URL=http://localhost:8082
CENTRAL_MGMT_KEY=central-mgmt-service-key
```

## 🧪 **Testing the Complete Implementation**

### **1. Start All Services:**

**Terminal 1 - Central Management System:**
```bash
cd mock-central-mgmt
go run main.go
```
Output: `🎛️  Mock Central Management System starting on :8082`

**Terminal 2 - API Beheerder:**
```bash
cd mock-beheerder
go run main.go
```
Output: `🔧 Mock API Beheerder starting on :8081`

**Terminal 3 - Your Internal API:**
```bash
go run main.go
```
Output: 
```
🚀 Internal API starting on localhost:8080
   📱 Accepts requests from User Portal
   🔗 Connects to API Beheerder at http://localhost:8081
   🎛️  Connects to Central Management at http://localhost:8082
   🔄 Architecture: User Portal → Internal API → [API Beheerder + Central Management]
```

### **2. Test Complete Flow:**
```bash
./test_user_portal.sh
```

This will test:
- Permission checking with Central Management
- Business rules enforcement  
- Data operations via API Beheerder
- Audit logging to Central Management
- Error handling when services are down

## 📡 **Request Flow Example**

### **User Creates Album:**

1. **User Portal → Your API:**
   ```http
   POST /albums HTTP/1.1
   Host: localhost:8080
   Authorization: Bearer valid-jwt-token-12345
   Content-Type: application/json
   
   {
     "id": "4",
     "title": "Kind of Blue",
     "artist": "Miles Davis",
     "price": 45.99
   }
   ```

2. **Your API processes:**
   - ✅ Validates JWT token
   - ✅ Validates album data
   - ✅ Adds user context (who created it)

3. **Your API → API Beheerder:**
   ```http
   POST /albums HTTP/1.1
   Host: localhost:8081
   X-Service-Key: beheerder-service-key
   Content-Type: application/json
   
   {
     "id": "4",
     "title": "Kind of Blue",
     "artist": "Miles Davis", 
     "price": 45.99,
     "createdBy": "user123",
     "createdAt": 1696248615
   }
   ```

4. **API Beheerder → Database:** Stores album

5. **API Beheerder → Your API:** Returns success

6. **Your API → User Portal:** Returns formatted response:
   ```json
   {
     "message": "Album created successfully",
     "data": {
       "id": "4",
       "title": "Kind of Blue",
       "artist": "Miles Davis",
       "price": 45.99,
       "createdBy": "user123",
       "createdAt": 1696248615
     }
   }
   ```

## 🔒 **Security Features**

### **User Authentication (from User Portal):**
- JWT token validation in `Authorization` header
- User context extraction and audit trails
- Protected endpoints require valid tokens

### **Service Authentication (to API Beheerder):**
- Service key authentication via `X-Service-Key` header
- Timeout handling for service calls
- Error propagation and logging

## 🎯 **Your API's Responsibilities**

### ✅ **What Your API Handles:**
1. **User-facing concerns**: Authentication, validation, formatting
2. **Business logic**: Rules, workflows, user context
3. **API orchestration**: Calling API Beheerder, error handling
4. **Response formatting**: Consistent responses to User Portal

### ❌ **What Your API Doesn't Handle:**
- Direct database access (API Beheerder does this)
- Complex data queries (API Beheerder does this)
- Data persistence logic (API Beheerder does this)
- Cross-service data management (API Beheerder does this)

## 🔄 **Migration to Real Business Logic**

When you're ready to replace albums with your real business logic:

1. **Keep all infrastructure** (auth, logging, HTTP client, error handling)
2. **Replace data structures** (`album` → your real types)
3. **Update endpoints** (`/albums` → your real endpoints)
4. **Update validation** (album rules → your business rules)
5. **Update API Beheerder calls** (album endpoints → your data endpoints)

The foundation supports any business domain! 🚀

## 📝 **Next Steps**

1. **Test the current setup** with mock API Beheerder
2. **Replace album logic** with your real business domain
3. **Connect to real API Beheerder** when it's ready
4. **Add more endpoints** as your business grows
5. **Enhance authentication** with real JWT validation
6. **Add monitoring** and metrics as needed

This architecture gives you a clean separation of concerns and a solid foundation for any business logic!