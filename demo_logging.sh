#!/bin/bash

echo "🎯 LOGGING & MONITORING WALKTHROUGH"
echo "=================================="

echo ""
echo "📋 What you'll see in this demo:"
echo "1. 🏗️  Structured JSON logging"
echo "2. 🔍 Request correlation IDs"
echo "3. 📊 Prometheus metrics"
echo "4. 🏥 Health check monitoring"
echo "5. ⚠️  Error handling and logging"

echo ""
echo "🚀 Starting Internal API with DEBUG logging..."
echo "Look for these log features:"
echo ""

# Start the API in background
cd "$(dirname "$0")"
export LOG_LEVEL=DEBUG
./internal-api.exe &
API_PID=$!

# Wait for startup
sleep 3

echo "📱 Making requests to demonstrate logging..."

# Test 1: Health check
echo ""
echo "🏥 TEST 1: Health Check (generates dependency health logs)"
curl -s http://localhost:8080/health | head -c 200
echo "..."
echo ""

# Test 2: Missing JWT
echo "🔐 TEST 2: Request without JWT (generates auth error logs)"
curl -s http://localhost:8080/albums | head -c 200
echo "..."
echo ""

# Test 3: Invalid JWT
echo "🎫 TEST 3: Request with invalid JWT (generates JWT validation logs)"
curl -s -H "Authorization: Bearer invalid-token" http://localhost:8080/albums | head -c 200
echo "..."
echo ""

# Test 4: Metrics
echo "📊 TEST 4: Metrics endpoint"
curl -s http://localhost:8080/metrics | grep "internal_api" | head -5
echo ""

echo "🏁 Demo complete!"
echo ""
echo "📝 In the logs, you should see:"
echo "✅ JSON structured logs with timestamps"
echo "✅ Request IDs for correlation"
echo "✅ Different log levels (DEBUG, INFO, WARN, ERROR)"
echo "✅ Service names and structured fields"
echo "✅ Performance metrics (duration, status codes)"
echo "✅ External service call tracking"

# Cleanup
sleep 2
kill $API_PID 2>/dev/null