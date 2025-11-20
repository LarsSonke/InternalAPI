# Test Broker Integration
# This script helps verify that InternalAPI successfully registers with the Broker

Write-Host "=== Broker Integration Test ===" -ForegroundColor Cyan
Write-Host ""

# Configuration
$brokerUrl = "http://localhost:8081"
$internalApiUrl = "http://localhost:8080"

# Test 1: Check if Broker is running
Write-Host "Test 1: Checking Broker status..." -ForegroundColor Yellow
try {
    $brokerStatus = Invoke-RestMethod -Uri "$brokerUrl/api/v1/status" -Method GET -ErrorAction Stop
    Write-Host "✓ Broker is running" -ForegroundColor Green
    Write-Host "  Version: $($brokerStatus.version)" -ForegroundColor Gray
    Write-Host "  Status: $($brokerStatus.status)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Broker is not running" -ForegroundColor Red
    Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "  Start the broker with: cd ..\modulair-achterkantje\broker; go run main.go" -ForegroundColor Yellow
    exit 1
}
Write-Host ""

# Test 2: Check if InternalAPI is running
Write-Host "Test 2: Checking InternalAPI status..." -ForegroundColor Yellow
try {
    $apiHealth = Invoke-RestMethod -Uri "$internalApiUrl/health" -Method GET -ErrorAction Stop
    Write-Host "✓ InternalAPI is running" -ForegroundColor Green
    Write-Host "  Status: $($apiHealth.status)" -ForegroundColor Gray
} catch {
    Write-Host "✗ InternalAPI is not running" -ForegroundColor Red
    Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "  Start InternalAPI with: go run main.go" -ForegroundColor Yellow
    exit 1
}
Write-Host ""

# Test 3: Check if InternalAPI registered with Broker
Write-Host "Test 3: Checking plugin registration..." -ForegroundColor Yellow
try {
    $plugins = Invoke-RestMethod -Uri "$brokerUrl/api/v1/routes" -Method GET -ErrorAction Stop
    $internalApi = $plugins.plugins | Where-Object { $_.slug -eq "internal-api" }
    
    if ($internalApi) {
        Write-Host "✓ InternalAPI is registered with Broker" -ForegroundColor Green
        Write-Host "  Slug: $($internalApi.slug)" -ForegroundColor Gray
        Write-Host "  Name: $($internalApi.name)" -ForegroundColor Gray
        Write-Host "  Host: $($internalApi.host)" -ForegroundColor Gray
        Write-Host "  Base Route: $($internalApi.'base-api-route')" -ForegroundColor Gray
        Write-Host "  Enabled: $($internalApi.enabled)" -ForegroundColor Gray
    } else {
        Write-Host "✗ InternalAPI not found in broker registry" -ForegroundColor Red
        Write-Host "  Registered plugins: $($plugins.plugins.Count)" -ForegroundColor Yellow
        foreach ($plugin in $plugins.plugins) {
            Write-Host "    - $($plugin.slug)" -ForegroundColor Yellow
        }
    }
} catch {
    Write-Host "✗ Failed to check registration" -ForegroundColor Red
    Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 4: Test proxying through Broker
Write-Host "Test 4: Testing request proxying..." -ForegroundColor Yellow
try {
    # Test health endpoint through broker
    $proxyResponse = Invoke-RestMethod -Uri "$brokerUrl/api/v1/health" -Method GET -ErrorAction Stop
    Write-Host "✓ Broker successfully proxied request to InternalAPI" -ForegroundColor Green
    Write-Host "  Endpoint: /api/v1/health" -ForegroundColor Gray
    Write-Host "  Response status: $($proxyResponse.status)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Proxy request failed" -ForegroundColor Red
    Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Red
    
    # Try to determine the issue
    if ($_.Exception.Response.StatusCode -eq 404) {
        Write-Host "  Possible cause: InternalAPI not registered or route mismatch" -ForegroundColor Yellow
    } elseif ($_.Exception.Response.StatusCode -eq 503) {
        Write-Host "  Possible cause: InternalAPI is down or unreachable" -ForegroundColor Yellow
    }
}
Write-Host ""

# Test 5: Compare direct vs proxied requests
Write-Host "Test 5: Comparing direct vs proxied requests..." -ForegroundColor Yellow
try {
    # Direct request
    $directStart = Get-Date
    $directResponse = Invoke-RestMethod -Uri "$internalApiUrl/health" -Method GET -ErrorAction Stop
    $directTime = (Get-Date) - $directStart
    
    # Proxied request
    $proxyStart = Get-Date
    $proxyResponse = Invoke-RestMethod -Uri "$brokerUrl/api/v1/health" -Method GET -ErrorAction Stop
    $proxyTime = (Get-Date) - $proxyStart
    
    Write-Host "✓ Both direct and proxied requests successful" -ForegroundColor Green
    Write-Host "  Direct request time: $($directTime.TotalMilliseconds)ms" -ForegroundColor Gray
    Write-Host "  Proxied request time: $($proxyTime.TotalMilliseconds)ms" -ForegroundColor Gray
    Write-Host "  Overhead: $([math]::Round($proxyTime.TotalMilliseconds - $directTime.TotalMilliseconds, 2))ms" -ForegroundColor Gray
} catch {
    Write-Host "✗ Comparison failed" -ForegroundColor Red
    Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

Write-Host "=== Test Complete ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Integration Status:" -ForegroundColor Cyan
Write-Host "  Broker URL: $brokerUrl" -ForegroundColor Gray
Write-Host "  InternalAPI URL: $internalApiUrl" -ForegroundColor Gray
Write-Host "  Proxied Endpoint: $brokerUrl/api/v1/* → $internalApiUrl/api/v1/*" -ForegroundColor Gray
Write-Host ""
Write-Host "Next Steps:" -ForegroundColor Cyan
Write-Host "  1. Use $brokerUrl as the main entry point for all client requests" -ForegroundColor Gray
Write-Host "  2. Configure clients to point to the Broker instead of InternalAPI directly" -ForegroundColor Gray
Write-Host "  3. Monitor logs in both services for registration and proxy events" -ForegroundColor Gray
