#!/bin/bash

echo "🏨 HOTEL PLUGIN HUB API DEMO"
echo "============================"

echo ""
echo "🎯 Testing Hotel Internal API with Plugin Authentication"
echo ""

API_URL="http://localhost:8080"
PLUGIN_KEY="hotel-plugin-key-12345"
JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTk4MjY1NjMsImlhdCI6MTc1OTc0MDE2MywibmFtZSI6IkpvaG4gRG9lIiwicm9sZSI6ImFkbWluIiwic3ViIjoiam9obi5kb2VAZXhhbXBsZS5jb20iLCJ1c2VyX2lkIjoidXNlcjEyMzQ1In0.5AVfEPUD5b9iWVkT-oA8k2Kk-S0B8X3bmE8DdY0n9Gs"

echo "📊 1. Health Check (Public endpoint)"
echo "curl $API_URL/health"
curl -s "$API_URL/health" | head -c 200
echo "..."
echo ""

echo "🔌 2. Plugin Authentication Test (should work)"
echo "curl -H 'X-Plugin-API-Key: $PLUGIN_KEY' $API_URL/api/v1/albums"
curl -s -H "X-Plugin-API-Key: $PLUGIN_KEY" -H "X-Plugin-Name: Booking-System" "$API_URL/api/v1/albums" | head -c 200
echo "..."
echo ""

echo "❌ 3. No Authentication Test (should fail)"
echo "curl $API_URL/api/v1/albums"
curl -s "$API_URL/api/v1/albums" | head -c 200
echo "..."
echo ""

echo "🎫 4. JWT Authentication Test"
echo "curl -H 'Authorization: Bearer <token>' $API_URL/api/v1/albums"
curl -s -H "Authorization: Bearer $JWT_TOKEN" "$API_URL/api/v1/albums" | head -c 200
echo "..."
echo ""

echo "🔐 5. Admin Endpoint Test (JWT only)"
echo "curl -H 'Authorization: Bearer <token>' $API_URL/admin/system-status"
curl -s -H "Authorization: Bearer $JWT_TOKEN" "$API_URL/admin/system-status" | head -c 200
echo "..."
echo ""

echo "❌ 6. Admin with Plugin Key (should fail)"
echo "curl -H 'X-Plugin-API-Key: $PLUGIN_KEY' $API_URL/admin/system-status"
curl -s -H "X-Plugin-API-Key: $PLUGIN_KEY" "$API_URL/admin/system-status" | head -c 200
echo "..."
echo ""

echo "📊 7. CORS Headers Test"
echo "curl -H 'Origin: http://localhost:3000' $API_URL/health"
curl -s -H "Origin: http://localhost:3000" -I "$API_URL/health" | grep -i "access-control"
echo ""

echo "📈 8. Metrics Endpoint"
echo "curl $API_URL/metrics | grep internal_api"
curl -s "$API_URL/metrics" | grep "internal_api" | head -3
echo ""

echo "🎉 Demo Complete!"
echo ""
echo "✅ Plugin Authentication: Use X-Plugin-API-Key header"
echo "✅ User Authentication: Use Authorization: Bearer token"
echo "✅ CORS Enabled: For web frontends"
echo "✅ Dual Auth Support: Both plugins and users can access /api/v1/*"
echo "✅ Admin Protection: Only JWT users can access /admin/*"