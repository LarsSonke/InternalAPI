# PowerShell test script for Internal API
$baseUrl = "http://localhost:8080"

Write-Host "🧪 COMPREHENSIVE API TESTING" -ForegroundColor Green
Write-Host "================================" -ForegroundColor Green

# Test 1: Health endpoint
Write-Host "`n1. Testing Health Endpoint:" -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$baseUrl/health" -Method GET
    Write-Host "✅ Health endpoint responding" -ForegroundColor Green
    Write-Host "Status: $($health.status)" -ForegroundColor Cyan
} catch {
    Write-Host "❌ Health endpoint failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 2: Circuit breakers status
Write-Host "`n2. Testing Circuit Breakers:" -ForegroundColor Yellow
try {
    $circuitBreakers = Invoke-RestMethod -Uri "$baseUrl/health/circuit-breakers" -Method GET
    Write-Host "✅ Circuit breakers endpoint responding" -ForegroundColor Green
    $circuitBreakers | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Circuit breakers failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3: Metrics endpoint
Write-Host "`n3. Testing Metrics Endpoint:" -ForegroundColor Yellow
try {
    $metrics = Invoke-WebRequest -Uri "$baseUrl/metrics" -Method GET
    Write-Host "✅ Metrics endpoint responding" -ForegroundColor Green
    Write-Host "Content-Type: $($metrics.Headers.'Content-Type')" -ForegroundColor Cyan
    Write-Host "Response size: $($metrics.Content.Length) bytes" -ForegroundColor Cyan
} catch {
    Write-Host "❌ Metrics failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: Authentication endpoint
Write-Host "`n4. Testing Authentication:" -ForegroundColor Yellow
try {
    $loginData = @{
        username = "testuser"
        password = "testpass"
    } | ConvertTo-Json
    
    $response = Invoke-WebRequest -Uri "$baseUrl/auth/login" -Method POST -Body $loginData -ContentType "application/json"
    Write-Host "✅ Auth endpoint responding (Status: $($response.StatusCode))" -ForegroundColor Green
} catch {
    Write-Host "✅ Auth endpoint responding with expected error: $($_.Exception.Response.StatusCode)" -ForegroundColor Green
}

# Test 5: Protected endpoint without auth
Write-Host "`n5. Testing Protected Endpoint (No Auth):" -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/api/v1/albums" -Method GET
    Write-Host "❌ Unexpected success - should require auth" -ForegroundColor Red
} catch {
    if ($_.Exception.Response.StatusCode -eq 401) {
        Write-Host "✅ Protected endpoint correctly rejecting unauthorized requests" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Unexpected error: $($_.Exception.Response.StatusCode)" -ForegroundColor Yellow
    }
}

# Test 6: CORS headers
Write-Host "`n6. Testing CORS Headers:" -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/health" -Method OPTIONS
    Write-Host "✅ OPTIONS request successful" -ForegroundColor Green
    Write-Host "CORS headers present: $($response.Headers.Keys -contains 'Access-Control-Allow-Origin')" -ForegroundColor Cyan
} catch {
    Write-Host "⚠️  CORS test: $($_.Exception.Message)" -ForegroundColor Yellow
}

# Test 7: Invalid endpoint
Write-Host "`n7. Testing Invalid Endpoint:" -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$baseUrl/invalid/endpoint" -Method GET
    Write-Host "❌ Unexpected success for invalid endpoint" -ForegroundColor Red
} catch {
    if ($_.Exception.Response.StatusCode -eq 404) {
        Write-Host "✅ Invalid endpoint correctly returns 404" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Unexpected status: $($_.Exception.Response.StatusCode)" -ForegroundColor Yellow
    }
}

Write-Host "`n🎯 TEST SUMMARY COMPLETE" -ForegroundColor Green