#!/bin/bash

# JWT Flow Demonstration Script
echo "🧪 JWT Validation Flow Examples"
echo "================================"

BASE_URL="http://localhost:8080"

# Example 1: Missing token (should fail)
echo "1. Request WITHOUT Authorization header:"
echo "   Expected: 401 Unauthorized with MISSING_TOKEN error"
echo ""
curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"id":"test","title":"Test Album","artist":"Test Artist","price":29.99}' \
  "$BASE_URL/albums" | jq '.' 2>/dev/null || echo "Request failed (expected)"
echo ""

# Example 2: Wrong token format (should fail)  
echo "2. Request with WRONG token format:"
echo "   Expected: 401 Unauthorized with INVALID_TOKEN_FORMAT error"
echo ""
curl -s -X POST \
  -H "Authorization: JustAString" \
  -H "Content-Type: application/json" \
  -d '{"id":"test","title":"Test Album","artist":"Test Artist","price":29.99}' \
  "$BASE_URL/albums" | jq '.' 2>/dev/null || echo "Request failed (expected)"
echo ""

# Example 3: Invalid JWT token (should fail)
echo "3. Request with INVALID JWT token:"
echo "   Expected: 401 Unauthorized with INVALID_TOKEN error"
echo ""
curl -s -X POST \
  -H "Authorization: Bearer invalid.jwt.token" \
  -H "Content-Type: application/json" \
  -d '{"id":"test","title":"Test Album","artist":"Test Artist","price":29.99}' \
  "$BASE_URL/albums" | jq '.' 2>/dev/null || echo "Request failed (expected)"
echo ""

# Example 4: Valid JWT token format (will work if services are running)
echo "4. Request with VALID JWT token format:"
echo "   Expected: Success if Central Management & API Beheerder are running"
echo "   Token claims will show: user_id=user123, action logged to Central Management"
echo ""
curl -s -X POST \
  -H "Authorization: Bearer valid-jwt-token-12345" \
  -H "Content-Type: application/json" \
  -d '{"id":"test","title":"Test Album","artist":"Test Artist","price":29.99}' \
  "$BASE_URL/albums" | jq '.' 2>/dev/null || echo "Request failed (need Central Management & API Beheerder running)"
echo ""

echo "✅ JWT Flow Demonstration Complete!"
echo ""
echo "💡 To test with REAL JWT tokens:"
echo "   1. Set JWT_SECRET environment variable"
echo "   2. Create tokens signed with same secret"
echo "   3. Include user_id or sub claims in token payload"