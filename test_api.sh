#!/bin/bash

# Test script to simulate API beheerder calling your internal API
API_KEY="dev-internal-key-12345"
BASE_URL="http://localhost:8080"

echo "🧪 Testing Internal API..."
echo "================================"

# Test 1: Health check
echo "1. Health Check:"
curl -s -H "X-Internal-API-Key: $API_KEY" "$BASE_URL/health" | jq '.' || echo "Failed"
echo ""

# Test 2: Get all albums
echo "2. Get All Albums:"
curl -s -H "X-Internal-API-Key: $API_KEY" "$BASE_URL/albums" | jq '.' || echo "Failed"
echo ""

# Test 3: Get specific album
echo "3. Get Album by ID (ID=1):"
curl -s -H "X-Internal-API-Key: $API_KEY" "$BASE_URL/albums/1" | jq '.' || echo "Failed"
echo ""

# Test 4: Create new album
echo "4. Create New Album:"
curl -s -X POST \
  -H "X-Internal-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"id":"4","title":"Test Album","artist":"Test Artist","price":29.99}' \
  "$BASE_URL/albums" | jq '.' || echo "Failed"
echo ""

# Test 5: Test authentication error
echo "5. Test Invalid API Key (should fail):"
curl -s -H "X-Internal-API-Key: wrong-key" "$BASE_URL/health" | jq '.' || echo "Failed as expected"
echo ""

# Test 6: Test missing API key
echo "6. Test Missing API Key (should fail):"
curl -s "$BASE_URL/health" | jq '.' || echo "Failed as expected"
echo ""

echo "✅ Test completed!"