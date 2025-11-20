# Broker Integration Guide

## Overview

InternalAPI now automatically registers with the Broker service on startup. This enables the Broker to act as a reverse proxy gateway, routing all client requests through a single entry point.

## Architecture

```
Client Request → Broker (:8081) → InternalAPI (:8080)
                    ↓
              - JWT validation
              - Route matching
              - Request forwarding
```

## Setup

### 1. Start the Broker

```powershell
cd ..\modulair-achterkantje\broker
go run main.go
```

The broker will start on port `8081`.

### 2. Configure InternalAPI Environment

Add these environment variables to InternalAPI:

```powershell
# Broker connection settings
$env:BROKER_URL = "http://localhost:8081"
$env:BROKER_AUTH_TOKEN = "your-jwt-token-here"
```

**Getting a JWT token:**
You can generate a test token from the broker or use an existing valid JWT that the broker accepts.

### 3. Start InternalAPI

```powershell
cd ..\InternalAPI
go run main.go
```

You should see:
```
✓ Successfully registered with broker
```

## How It Works

### Registration Process

1. **InternalAPI starts** and initializes all services
2. **After 2 seconds**, it registers with the broker (non-blocking)
3. **Broker receives registration** with these details:
   - `slug`: "internal-api"
   - `host`: "http://localhost:8080"
   - `base-api-route`: "/api/v1"
   - List of available endpoints

4. **Broker stores registration** and begins routing requests

### Request Flow

**Before (Direct to InternalAPI):**
```
Client → http://localhost:8080/api/v1/albums
```

**After (Through Broker):**
```
Client → http://localhost:8081/api/v1/albums → Broker forwards to → http://localhost:8080/api/v1/albums
```

## Testing

### 1. Test Direct Access (InternalAPI)
```powershell
curl http://localhost:8080/health
```

### 2. Test Broker Proxying
```powershell
# Request goes through broker
curl http://localhost:8081/api/v1/albums

# Broker checks its registry
# Finds plugin with base-api-route="/api/v1"
# Forwards to http://localhost:8080/api/v1/albums
```

### 3. Check Plugin Registration
```powershell
curl http://localhost:8081/api/v1/routes
```

Should show:
```json
{
  "plugins": [
    {
      "slug": "internal-api",
      "name": "Hotel Internal API",
      "host": "http://localhost:8080",
      "base-api-route": "/api/v1",
      "enabled": true
    }
  ]
}
```

## Troubleshooting

### Registration Failed

**Error:** `Failed to register with broker`

**Causes:**
1. Broker not running
2. Invalid `BROKER_AUTH_TOKEN`
3. Network connectivity issue

**Solution:**
```powershell
# Check if broker is running
curl http://localhost:8081/api/v1/status

# Check InternalAPI logs for details
# The service continues running even if registration fails
```

### Authentication Required

**Error:** `registration failed with status 401`

**Solution:**
Set a valid JWT token:
```powershell
$env:BROKER_AUTH_TOKEN = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Route Conflict

**Error:** `route conflict`

**Cause:** Another plugin already registered with `/api/v1`

**Solution:**
```powershell
# Check existing registrations
curl http://localhost:8081/api/v1/routes

# Delete conflicting plugin
curl -X DELETE http://localhost:8081/api/v1/route/conflicting-slug \
  -H "Authorization: Bearer $env:BROKER_AUTH_TOKEN"
```

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `BROKER_URL` | Broker service URL | `http://localhost:8081` | No |
| `BROKER_AUTH_TOKEN` | JWT for broker authentication | None | Yes* |

*Required if broker has authentication enabled

## Production Considerations

### Security
- Use HTTPS for broker URL in production
- Rotate `BROKER_AUTH_TOKEN` regularly
- Use environment-specific tokens

### High Availability
- If broker is down, InternalAPI continues running
- Registration failure is logged but non-fatal
- Consider implementing retry logic for critical deployments

### Monitoring
- Check broker logs for registration events
- Monitor `/api/v1/health/plugins` endpoint
- Set up alerts for failed registrations

## Benefits of Using the Broker

1. **Single Entry Point**: All services accessible through one gateway
2. **Centralized Authentication**: JWT validation in one place
3. **Service Discovery**: Services register themselves automatically
4. **Traffic Management**: Future support for load balancing, rate limiting
5. **Observability**: Centralized logging and monitoring
