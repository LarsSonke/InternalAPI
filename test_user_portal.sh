#!/bin/bash

# Test script to simulate User Portal calling your internal API
# Architecture: User Portal → Your API → [API Beheerder + Central Management]

TOKEN="valid-jwt-token-12345"  # Simulating JWT token from user login
BASE_URL="http://localhost:8080"

echo "🧪 Testing Internal API with Central Management Integration..."
echo "================================================================"
echo "Architecture: User Portal → Internal API → [API Beheerder + Central Management]"
echo ""

# Test 1: Health check (no authentication required)
echo "1. Health Check (Public endpoint):"
curl -s "$BASE_URL/health" | jq '.' 2>/dev/null || echo "Failed or no jq installed"
echo ""

# Test 2: Get all albums with authentication (requires both Central Management and API Beheerder)
echo "2. Get All Albums (Protected endpoint - needs permissions + data):"
curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/albums" | jq '.' 2>/dev/null || echo "Note: This will fail until both Central Management and API Beheerder are running"
echo ""

# Test 3: Get specific album
echo "3. Get Album by ID (ID=1):"
curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/albums/1" | jq '.' 2>/dev/null || echo "Note: This will fail until both services are running"
echo ""

# Test 4: Create new album (needs permission check + business rules + data storage)
echo "4. Create New Album (needs all services):"
curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":"4","title":"Test Album","artist":"Test Artist","price":29.99}' \
  "$BASE_URL/albums" | jq '.' 2>/dev/null || echo "Note: This will fail until all services are running"
echo ""

# Test 5: Try to create expensive album (should be blocked by Central Management)
echo "5. Create Expensive Album (should be blocked by business rules):"
curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":"5","title":"Expensive Album","artist":"Rich Artist","price":150.00}' \
  "$BASE_URL/albums" | jq '.' 2>/dev/null || echo "Note: Should be blocked by Central Management if running"
echo ""

# Test 6: Update album
echo "6. Update Album (ID=1):"
curl -s -X PUT \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":"1","title":"Updated Album","artist":"Updated Artist","price":35.99}' \
  "$BASE_URL/albums/1" | jq '.' 2>/dev/null || echo "Note: This will fail until all services are running"
echo ""

# Test 7: Delete album
echo "7. Delete Album (ID=4):"
curl -s -X DELETE \
  -H "Authorization: Bearer $TOKEN" \
  "$BASE_URL/albums/4" | jq '.' 2>/dev/null || echo "Note: This will fail until all services are running"
echo ""

# Test 8: Test limited user (different permissions)
echo "8. Test Limited User (different token):"
curl -s -X POST \
  -H "Authorization: Bearer limited-user-token" \
  -H "Content-Type: application/json" \
  -d '{"id":"6","title":"Limited User Album","artist":"Limited Artist","price":25.00}' \
  "$BASE_URL/albums" | jq '.' 2>/dev/null || echo "Note: Should be blocked by Central Management if running"
echo ""

# Test 9: Test invalid token (should fail)
echo "9. Test Invalid Token (should fail):"
curl -s -H "Authorization: Bearer invalid-token" "$BASE_URL/albums" | jq '.' 2>/dev/null || echo "Failed as expected"
echo ""

# Test 10: Test missing token (should fail)
echo "10. Test Missing Token (should fail):"
curl -s "$BASE_URL/albums" | jq '.' 2>/dev/null || echo "Failed as expected"
echo ""

echo "✅ Test completed!"
echo ""
echo "📝 Notes:"
echo "- Health check should work immediately"
echo "- Protected endpoints need BOTH Central Management (port 8082) AND API Beheerder (port 8081)"
echo "- Central Management handles: permissions, business rules, audit logging"
echo "- API Beheerder handles: data storage and retrieval"
echo "- Your Internal API orchestrates between both systems"
echo ""
echo "🚀 To run full test:"
echo "1. Terminal 1: cd mock-central-mgmt && go run main.go"
echo "2. Terminal 2: cd mock-beheerder && go run main.go" 
echo "3. Terminal 3: go run main.go"
echo "4. Terminal 4: ./test_user_portal.sh"