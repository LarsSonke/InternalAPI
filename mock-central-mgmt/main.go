package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Mock Central Management System
// Handles: permissions, business rules, audit logging, user filters, etc.

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Middleware to validate service key from your Internal API
	router.Use(func(c *gin.Context) {
		serviceKey := c.GetHeader("X-Service-Key")
		if serviceKey != "central-mgmt-service-key" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid service key"})
			c.Abort()
			return
		}
		c.Next()
	})

	// Permission checking endpoints
	router.POST("/check-permission", checkPermission)

	// Business rules endpoints
	router.GET("/business-rules/albums", getAlbumBusinessRules)
	router.GET("/user-filters/albums", getUserFilters)

	// Audit and logging endpoints
	router.POST("/audit-log", logAuditEvent)
	router.POST("/access-log", logAccessEvent)

	// Configuration endpoints
	router.GET("/config/:service", getServiceConfig)

	fmt.Println("🎛️  Mock Central Management System starting on :8082")
	fmt.Println("   🔐 Handles permissions, business rules, and audit logging")
	fmt.Println("   🔗 Ready to receive requests from your Internal API")
	router.Run(":8082")
}

// Check if user has permission to perform an action
func checkPermission(c *gin.Context) {
	var permissionRequest map[string]interface{}

	if err := c.ShouldBindJSON(&permissionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	userID := permissionRequest["userID"].(string)
	action := permissionRequest["action"].(string)
	resource := permissionRequest["resource"].(string)

	fmt.Printf("Permission check: User %s wants to %s on %s\n", userID, action, resource)

	// Mock permission logic (in real system, this would check database/RBAC)
	allowed := true
	reason := ""

	// Example business rules
	if action == "create_album" {
		if data, exists := permissionRequest["data"].(map[string]interface{}); exists {
			if price, hasPrice := data["price"].(float64); hasPrice && price > 100 {
				allowed = false
				reason = "Users cannot create albums with price over $100"
			}
		}
	}

	// Example: Different user roles
	if userID == "limited-user" {
		if action == "create_album" || action == "delete_album" {
			allowed = false
			reason = "Limited users can only read albums"
		}
	}

	response := gin.H{
		"allowed":   allowed,
		"userID":    userID,
		"action":    action,
		"resource":  resource,
		"timestamp": time.Now().Unix(),
	}

	if !allowed {
		response["reason"] = reason
	}

	c.JSON(http.StatusOK, response)
}

// Get business rules for albums
func getAlbumBusinessRules(c *gin.Context) {
	rules := gin.H{
		"maxPrice":       99.99,
		"minPrice":       0.01,
		"maxTitleLen":    100,
		"requiredFields": []string{"id", "title", "artist", "price"},
		"allowedGenres":  []string{"Jazz", "Rock", "Classical", "Blues", "Pop"},
		"priceValidation": gin.H{
			"currency":  "USD",
			"precision": 2,
		},
		"auditRequired": true,
		"version":       "1.2",
		"lastUpdated":   time.Now().Unix(),
	}

	c.JSON(http.StatusOK, rules)
}

// Get user-specific filters
func getUserFilters(c *gin.Context) {
	userID := c.Query("userID")

	fmt.Printf("Fetching filters for user: %s\n", userID)

	// Mock user-specific filters
	filters := gin.H{
		"userID": userID,
		"filters": gin.H{
			"maxPrice": 50.0,                      // This user can only see albums under $50
			"genres":   []string{"Jazz", "Blues"}, // Only these genres
			"region":   "US",                      // Regional restrictions
		},
		"permissions": gin.H{
			"canCreate": true,
			"canUpdate": false,
			"canDelete": false,
		},
	}

	// Different filters for different users
	if userID == "admin-user" {
		filters["filters"] = gin.H{} // No filters for admin
		filters["permissions"] = gin.H{
			"canCreate": true,
			"canUpdate": true,
			"canDelete": true,
		}
	}

	c.JSON(http.StatusOK, filters)
}

// Log audit events
func logAuditEvent(c *gin.Context) {
	var auditLog map[string]interface{}

	if err := c.ShouldBindJSON(&auditLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// In real system, this would store in audit database
	fmt.Printf("🔍 AUDIT LOG: %v\n", auditLog)

	c.JSON(http.StatusOK, gin.H{
		"status": "logged",
		"id":     fmt.Sprintf("audit-%d", time.Now().Unix()),
	})
}

// Log access events
func logAccessEvent(c *gin.Context) {
	var accessLog map[string]interface{}

	if err := c.ShouldBindJSON(&accessLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// In real system, this would store in analytics database
	fmt.Printf("📊 ACCESS LOG: %v\n", accessLog)

	c.JSON(http.StatusOK, gin.H{
		"status": "logged",
		"id":     fmt.Sprintf("access-%d", time.Now().Unix()),
	})
}

// Get service-specific configuration
func getServiceConfig(c *gin.Context) {
	service := c.Param("service")

	config := gin.H{
		"service": service,
		"version": "1.0",
		"features": gin.H{
			"auditEnabled":     true,
			"cacheEnabled":     true,
			"rateLimitEnabled": true,
		},
		"limits": gin.H{
			"maxRequestsPerMinute": 1000,
			"maxDataSize":          "10MB",
		},
		"lastUpdated": time.Now().Unix(),
	}

	c.JSON(http.StatusOK, config)
}
