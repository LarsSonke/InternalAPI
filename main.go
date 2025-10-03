package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Note: We no longer store albums locally, they're managed by API Beheerder

// ============================================
// INFRASTRUCTURE LAYER (Keep for real system)
// ============================================

// Standard error response format
type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// Send standardized error response
func sendError(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, ErrorResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}

// Configuration
func getConfig() map[string]string {
	return map[string]string{
		// Your API settings
		"port": getEnvOrDefault("PORT", "8080"),
		"host": getEnvOrDefault("HOST", "localhost"),

		// JWT settings for User Portal authentication
		"jwt_secret": getEnvOrDefault("JWT_SECRET", "your-jwt-secret-key"),

		// API Beheerder settings (for data operations)
		"api_beheerder_url": getEnvOrDefault("API_BEHEERDER_URL", "http://localhost:8081"),
		"api_beheerder_key": getEnvOrDefault("API_BEHEERDER_KEY", "beheerder-service-key"),

		// Central Management System settings (for business rules, permissions, audit, etc.)
		"central_mgmt_url": getEnvOrDefault("CENTRAL_MGMT_URL", "http://localhost:8082"),
		"central_mgmt_key": getEnvOrDefault("CENTRAL_MGMT_KEY", "central-mgmt-service-key"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Authentication middleware - validates requests from User Portal
func userPortalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for Authorization header (JWT token from user login)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sendError(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization token required")
			c.Abort()
			return
		}

		// Extract Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			sendError(c, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Token must be in format: Bearer <token>")
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate JWT token (simple validation for now)
		if !isValidJWTToken(token) {
			sendError(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user info in context for use in handlers
		userID := extractUserIDFromToken(token)
		c.Set("userID", userID)
		c.Set("token", token)

		c.Next()
	}
}

// Simple JWT token validation (replace with proper JWT library in production)
func isValidJWTToken(token string) bool {
	// For development: accept any token that's not empty and has minimum length
	// In production: use proper JWT validation with signature verification
	return len(token) > 10 && !strings.Contains(token, "invalid")
}

// Extract user ID from token (placeholder implementation)
func extractUserIDFromToken(token string) string {
	// For development: return a dummy user ID
	// In production: decode JWT and extract user ID from claims
	return "user123"
}

// HTTP client for calling API Beheerder
var apiBeheerderClient = &http.Client{
	Timeout: 30 * time.Second,
}

// Call API Beheerder for data operations
func callAPIBeheerder(method, endpoint string, data interface{}) (map[string]interface{}, error) {
	config := getConfig()
	beheerderURL := config["api_beheerder_url"]
	beheerderKey := config["api_beheerder_key"]

	var body []byte
	var err error

	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %v", err)
		}
	}

	req, err := http.NewRequest(method, beheerderURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Key", beheerderKey) // Authenticate to API Beheerder

	resp, err := apiBeheerderClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to API Beheerder failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Handle non-200 responses
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API Beheerder returned error %d: %s", resp.StatusCode, string(responseBody))
	}

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %v", err)
	}

	return result, nil
}

// HTTP client for calling Central Management System
var centralMgmtClient = &http.Client{
	Timeout: 30 * time.Second,
}

// Call Central Management System for business rules, permissions, audit, etc.
func callCentralManagement(method, endpoint string, data interface{}) (map[string]interface{}, error) {
	config := getConfig()
	centralURL := config["central_mgmt_url"]
	centralKey := config["central_mgmt_key"]

	var body []byte
	var err error

	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data for Central Management: %v", err)
		}
	}

	req, err := http.NewRequest(method, centralURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request to Central Management: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Key", centralKey) // Authenticate to Central Management

	resp, err := centralMgmtClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to Central Management failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Central Management response: %v", err)
	}

	// Handle non-200 responses
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Central Management returned error %d: %s", resp.StatusCode, string(responseBody))
	}

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Central Management response JSON: %v", err)
	}

	return result, nil
}

// Request logging middleware
func requestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[INTERNAL-API] %s | %s | %s %s | %d | %s | %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.ClientIP,
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
		)
	})
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "internal-api",
		"version":   "1.0.0",
		"timestamp": time.Now().Unix(),
		"uptime":    "system_uptime_here", // You can add real uptime tracking
	})
}

// ============================================
// BUSINESS LAYER (Replace with real logic later)
// ============================================

func main() {
	config := getConfig()

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New() // Use gin.New() instead of Default() for more control

	// Add infrastructure middleware
	router.Use(requestLogger())
	router.Use(gin.Recovery()) // Recover from panics

	// Public endpoints (no authentication required)
	router.GET("/health", healthCheck)

	// Protected endpoints (require user authentication from User Portal)
	protected := router.Group("/")
	protected.Use(userPortalAuthMiddleware())
	{
		// Business endpoints - these now call API Beheerder
		protected.GET("/albums", getAlbums)
		protected.GET("/albums/:id", getAlbumByID)
		protected.POST("/albums", postAlbums)
		protected.PUT("/albums/:id", updateAlbum)
		protected.DELETE("/albums/:id", deleteAlbum)
	}

	// Start server
	address := config["host"] + ":" + config["port"]
	fmt.Printf("🚀 Internal API starting on %s\n", address)
	fmt.Printf("   📱 Accepts requests from User Portal\n")
	fmt.Printf("   🔗 Connects to API Beheerder at %s\n", config["api_beheerder_url"])
	fmt.Printf("   🎛️  Connects to Central Management at %s\n", config["central_mgmt_url"])
	fmt.Printf("   🔄 Architecture: User Portal → Internal API → [API Beheerder + Central Management]\n")
	router.Run(address)
}

// ============================================
// BUSINESS LAYER (Now calls API Beheerder)
// ============================================

// getAlbums responds with the list of all albums as JSON.
// Checks permissions with Central Management and gets data from API Beheerder
func getAlbums(c *gin.Context) {
	// STEP 1: Check permissions with Central Management System
	userID, _ := c.Get("userID")
	permissionCheck := map[string]interface{}{
		"userID":   userID,
		"action":   "read_albums",
		"resource": "albums",
	}

	permission, err := callCentralManagement("POST", "/check-permission", permissionCheck)
	if err != nil {
		fmt.Printf("Error checking permissions with Central Management: %v\n", err)
		sendError(c, http.StatusInternalServerError, "PERMISSION_SERVICE_ERROR", "Failed to check permissions")
		return
	}

	if !permission["allowed"].(bool) {
		reason := "Permission denied"
		if msg, exists := permission["reason"]; exists {
			reason = msg.(string)
		}
		sendError(c, http.StatusForbidden, "PERMISSION_DENIED", reason)
		return
	}

	// STEP 2: Get user-specific filters from Central Management (if any)
	userFilters, err := callCentralManagement("GET", "/user-filters/albums?userID="+userID.(string), nil)
	if err != nil {
		fmt.Printf("Warning: Could not fetch user filters: %v\n", err)
		// Continue without filters
	}

	// STEP 3: Call API Beheerder to get albums data
	result, err := callAPIBeheerder("GET", "/albums", nil)
	if err != nil {
		fmt.Printf("Error calling API Beheerder: %v\n", err)
		sendError(c, http.StatusInternalServerError, "DATA_SERVICE_ERROR", "Failed to fetch albums from data service")
		return
	}

	// STEP 4: Apply user-specific filters (if any) from Central Management
	albums := result["albums"]
	if userFilters != nil && userFilters["filters"] != nil {
		// Apply filters based on user role, department, etc.
		// This is where you'd implement filtering logic based on Central Management rules
		fmt.Printf("Applying user filters: %v\n", userFilters["filters"])
	}

	// STEP 5: Log access to Central Management (async)
	go func() {
		accessLog := map[string]interface{}{
			"userID":    userID,
			"action":    "albums_accessed",
			"resource":  "albums",
			"timestamp": time.Now().Unix(),
			"count":     result["count"],
		}

		_, err := callCentralManagement("POST", "/access-log", accessLog)
		if err != nil {
			fmt.Printf("Warning: Failed to send access log to Central Management: %v\n", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"data":   albums,
		"count":  result["count"],
		"source": "api-beheerder", // Indicate data source
	})
}

// postAlbums adds an album by calling both Central Management and API Beheerder
func postAlbums(c *gin.Context) {
	var newAlbum album

	// STEP 1: Validate JSON binding from User Portal
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format: "+err.Error())
		return
	}

	// STEP 2: Basic validation (your business rules)
	if err := validateAlbum(&newAlbum); err != nil {
		sendError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	// STEP 3: Check permissions with Central Management System
	userID, _ := c.Get("userID")
	permissionCheck := map[string]interface{}{
		"userID":   userID,
		"action":   "create_album",
		"resource": "albums",
		"data": map[string]interface{}{
			"albumID": newAlbum.ID,
			"price":   newAlbum.Price,
		},
	}

	permission, err := callCentralManagement("POST", "/check-permission", permissionCheck)
	if err != nil {
		fmt.Printf("Error checking permissions with Central Management: %v\n", err)
		sendError(c, http.StatusInternalServerError, "PERMISSION_SERVICE_ERROR", "Failed to check permissions")
		return
	}

	if !permission["allowed"].(bool) {
		reason := "Permission denied"
		if msg, exists := permission["reason"]; exists {
			reason = msg.(string)
		}
		sendError(c, http.StatusForbidden, "PERMISSION_DENIED", reason)
		return
	}

	// STEP 4: Get business rules/configuration from Central Management
	businessRules, err := callCentralManagement("GET", "/business-rules/albums", nil)
	if err != nil {
		fmt.Printf("Warning: Could not fetch business rules: %v\n", err)
		// Continue without business rules (fallback behavior)
	} else {
		// Apply dynamic business rules from Central Management
		if maxPrice, exists := businessRules["maxPrice"]; exists && newAlbum.Price > maxPrice.(float64) {
			sendError(c, http.StatusBadRequest, "PRICE_TOO_HIGH", fmt.Sprintf("Price cannot exceed %.2f", maxPrice.(float64)))
			return
		}
	}

	// STEP 5: Prepare data with user context for API Beheerder
	albumData := map[string]interface{}{
		"id":        newAlbum.ID,
		"title":     newAlbum.Title,
		"artist":    newAlbum.Artist,
		"price":     newAlbum.Price,
		"createdBy": userID,
		"createdAt": time.Now().Unix(),
	}

	// STEP 6: Store data via API Beheerder
	result, err := callAPIBeheerder("POST", "/albums", albumData)
	if err != nil {
		fmt.Printf("Error calling API Beheerder: %v\n", err)
		sendError(c, http.StatusInternalServerError, "DATA_SERVICE_ERROR", "Failed to create album in data service")
		return
	}

	// STEP 7: Send audit log to Central Management System (async - don't block response)
	go func() {
		auditLog := map[string]interface{}{
			"userID":    userID,
			"action":    "album_created",
			"resource":  "albums",
			"albumID":   newAlbum.ID,
			"timestamp": time.Now().Unix(),
			"metadata": map[string]interface{}{
				"title":  newAlbum.Title,
				"artist": newAlbum.Artist,
				"price":  newAlbum.Price,
			},
		}

		_, err := callCentralManagement("POST", "/audit-log", auditLog)
		if err != nil {
			fmt.Printf("Warning: Failed to send audit log to Central Management: %v\n", err)
		}
	}()

	// STEP 8: Return success response to User Portal
	c.JSON(http.StatusCreated, gin.H{
		"message": "Album created successfully",
		"data":    result["album"], // Data from API Beheerder
	})
}

// Basic validation function (replace with your real validation)
func validateAlbum(album *album) error {
	if strings.TrimSpace(album.ID) == "" {
		return fmt.Errorf("ID is required")
	}
	if strings.TrimSpace(album.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if strings.TrimSpace(album.Artist) == "" {
		return fmt.Errorf("artist is required")
	}
	if album.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	if album.Price > 99999 {
		return fmt.Errorf("price cannot exceed 99999")
	}
	return nil
}

// getAlbumByID gets a specific album by calling API Beheerder
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		sendError(c, http.StatusBadRequest, "INVALID_ID", "ID parameter is required")
		return
	}

	// Call API Beheerder to get specific album
	result, err := callAPIBeheerder("GET", "/albums/"+id, nil)
	if err != nil {
		fmt.Printf("Error calling API Beheerder: %v\n", err)

		// Check if it's a 404 from API Beheerder
		if strings.Contains(err.Error(), "404") {
			sendError(c, http.StatusNotFound, "NOT_FOUND", "Album not found")
			return
		}

		sendError(c, http.StatusInternalServerError, "DATA_SERVICE_ERROR", "Failed to fetch album from data service")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result["album"], // Assuming API Beheerder returns {"album": {...}}
	})
}

// updateAlbum updates an album by calling API Beheerder
func updateAlbum(c *gin.Context) {
	id := c.Param("id")
	var updateData album

	if strings.TrimSpace(id) == "" {
		sendError(c, http.StatusBadRequest, "INVALID_ID", "ID parameter is required")
		return
	}

	// Validate JSON binding from User Portal
	if err := c.ShouldBindJSON(&updateData); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format: "+err.Error())
		return
	}

	// Basic validation
	if err := validateAlbum(&updateData); err != nil {
		sendError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	// Add user context
	userID, _ := c.Get("userID")
	albumData := map[string]interface{}{
		"id":        updateData.ID,
		"title":     updateData.Title,
		"artist":    updateData.Artist,
		"price":     updateData.Price,
		"updatedBy": userID,
		"updatedAt": time.Now().Unix(),
	}

	// Send to API Beheerder for update
	result, err := callAPIBeheerder("PUT", "/albums/"+id, albumData)
	if err != nil {
		fmt.Printf("Error calling API Beheerder: %v\n", err)

		if strings.Contains(err.Error(), "404") {
			sendError(c, http.StatusNotFound, "NOT_FOUND", "Album not found")
			return
		}

		sendError(c, http.StatusInternalServerError, "DATA_SERVICE_ERROR", "Failed to update album in data service")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Album updated successfully",
		"data":    result["album"],
	})
}

// deleteAlbum deletes an album by calling API Beheerder
func deleteAlbum(c *gin.Context) {
	id := c.Param("id")

	if strings.TrimSpace(id) == "" {
		sendError(c, http.StatusBadRequest, "INVALID_ID", "ID parameter is required")
		return
	}

	// Add user context for audit
	userID, _ := c.Get("userID")
	deleteData := map[string]interface{}{
		"deletedBy": userID,
		"deletedAt": time.Now().Unix(),
	}

	// Send to API Beheerder for deletion
	_, err := callAPIBeheerder("DELETE", "/albums/"+id, deleteData)
	if err != nil {
		fmt.Printf("Error calling API Beheerder: %v\n", err)

		if strings.Contains(err.Error(), "404") {
			sendError(c, http.StatusNotFound, "NOT_FOUND", "Album not found")
			return
		}

		sendError(c, http.StatusInternalServerError, "DATA_SERVICE_ERROR", "Failed to delete album from data service")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Album deleted successfully",
	})
}
