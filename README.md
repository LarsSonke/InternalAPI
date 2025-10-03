# Internal API - Central Management System Integration

## 🏗️ **Complete Architecture Flow**

```
User Portal → YOUR INTERNAL API → API Beheerder → Database/Storage
                    ↓
              Central Management System
```

Your Internal API now acts as the orchestration layer that:
1. **Receives requests from User Portal** (frontend/web app)
2. **Authenticates users** via JWT tokens
3. **Checks permissions and business rules** with Central Management System
4. **Validates and processes business logic** 
5. **Calls API Beheerder** for data operations
6. **Sends audit logs** to Central Management System
7. **Returns formatted responses** to User Portal

## 🎛️ **System Responsibilities**

### **Central Management System (Port 8082)**
- **Permission Management**: Who can do what
- **Business Rules Engine**: Dynamic rules and policies  
- **Audit Logging**: Track all user actions
- **User Filters**: Role-based data filtering
- **Configuration Management**: System settings
- **Analytics**: Usage tracking and reporting

### **API Beheerder (Port 8081)**  
- **Data Storage**: CRUD operations
- **Data Integrity**: Constraints and validation
- **Database Management**: Connections and queries
- **Data Consistency**: Transactions and locking

### **Your Internal API (Port 8080)**
- **User Authentication**: JWT validation
- **Request Orchestration**: Coordinate between systems
- **Business Logic**: Application-specific rules
- **Response Formatting**: Consistent API responses
- **Error Handling**: Graceful error management

## 🚀 **What's Implemented**

### ✅ **User Portal Authentication**
- JWT token validation from `Authorization: Bearer <token>` headers
- User context extraction and tracking
- Protected endpoints that require authentication

### ✅ **Central Management System Integration**
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