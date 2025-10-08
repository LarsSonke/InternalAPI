package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Global logger
var log *logrus.Logger

// HTTP clients for external services
var (
	apiBeheerderClient  = &http.Client{Timeout: 30 * time.Second}
	centralMgmtClient   = &http.Client{Timeout: 30 * time.Second}
)

// Initialize logging and metrics
func init() {
	setupLogging()
	setupMetrics()
}

// Main server setup and routing
func main() {
	config := getConfig()

	log.WithFields(logrus.Fields{
		"service": "internal-api",
		"port":    config["port"],
		"host":    config["host"],
		"role":    "hotel-operations-api",
	}).Info("Starting Hotel Internal API")

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New() // Use gin.New() instead of Default() for more control

	// CORS configuration for User Portal access
	validOrigins := []string{
		"http://localhost:3000",      // User Portal development
		"http://localhost:3001",      // User Portal staging
		"https://hotel-portal.local", // User Portal production
	}
	
	// Add configured origins from environment if they're valid HTTP/HTTPS URLs
	if originsStr := config["allowed_origins"]; originsStr != "" {
		for _, origin := range strings.Split(originsStr, ",") {
			origin = strings.TrimSpace(origin)
			if origin != "" && (strings.HasPrefix(origin, "http://") || strings.HasPrefix(origin, "https://")) {
				validOrigins = append(validOrigins, origin)
			}
		}
	}
	
	log.WithField("valid_origins", validOrigins).Info("Configured CORS origins for User Portal access")
	
	corsConfig := cors.Config{
		AllowOrigins: validOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Origin", 
			"Content-Type", 
			"Authorization", 
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"X-Request-ID",
			"X-Rate-Limit-Remaining",
			"X-Rate-Limit-Reset",
		},
		AllowCredentials: true,
		MaxAge:          12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))

	// Add infrastructure middleware
	router.Use(requestIDMiddleware())
	router.Use(metricsMiddleware())
	router.Use(gin.Recovery()) // Recover from panics

	// Public endpoints (no authentication required)
	router.GET("/health", healthCheck)
	router.GET("/metrics", gin.WrapH(promhttp.Handler())) // Prometheus metrics endpoint

	// User Portal API endpoints (require JWT authentication from User Portal)
	api := router.Group("/api/v1")
	api.Use(userPortalAuthMiddleware())
	{
		// Hotel management endpoints for User Portal
		api.GET("/albums", getAlbums)          // Replace with hotel bookings/rooms
		api.GET("/albums/:id", getAlbumByID)   // Replace with specific booking/room
		api.POST("/albums", postAlbums)        // Replace with create booking/room  
		api.PUT("/albums/:id", updateAlbum)    // Replace with update booking/room
		api.DELETE("/albums/:id", deleteAlbum) // Replace with cancel booking/delete room
	}

	// Admin endpoints (require user authentication from User Portal only)
	admin := router.Group("/admin")
	admin.Use(userPortalAuthMiddleware())
	{
		// Admin endpoints - only for authenticated users from User Portal
		admin.GET("/system-status", getSystemStatus)
		admin.GET("/audit-logs", getAuditLogs)
	}

	// Start server
	address := config["host"] + ":" + config["port"]
	
	log.WithFields(logrus.Fields{
		"address":                  address,
		"api_beheerder_url":       config["api_beheerder_url"],
		"central_mgmt_url":        config["central_mgmt_url"],
		"user_portal_url":         config["user_portal_url"],
		"metrics_endpoint":        "http://" + address + "/metrics",
		"health_endpoint":         "http://" + address + "/health",
		"api_endpoint":            "http://" + address + "/api/v1",
		"cors_origins":            config["allowed_origins"],
	}).Info("Hotel Internal API started successfully")
	
	log.Info("🏨 Hotel Internal API starting")
	log.Info("   👤 Accepts requests from User Portal")
	log.Info("   🔗 Connects to API Beheerder for data operations")
	log.Info("   🎛️  Connects to Central Management for business rules")
	log.Info("   📊 Metrics available at /metrics")
	log.Info("   💚 Health check available at /health")
	log.Info("   🎯 Main API endpoints at /api/v1/*")
	log.Info("   🔄 Architecture: User Portal → Internal API → [API Beheerder + Central Management]")

	log.Fatal(router.Run(address))
}

// ============================================
// EXTERNAL SERVICE COMMUNICATION
// ============================================

// Call API Beheerder for data operations
func callAPIBeheerder(method, endpoint string, data interface{}) (map[string]interface{}, error) {
	start := time.Now()
	config := getConfig()
	beheerderURL := config["api_beheerder_url"]
	beheerderKey := config["api_beheerder_key"]

	log.WithFields(logrus.Fields{
		"service":  "api-beheerder",
		"method":   method,
		"endpoint": endpoint,
	}).Debug("Calling API Beheerder")

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, beheerderURL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Key", beheerderKey)

	resp, err := apiBeheerderClient.Do(req)
	duration := time.Since(start)

	// Record metrics
	status := "error"
	if err == nil {
		status = fmt.Sprintf("%d", resp.StatusCode)
	}
	externalServiceCalls.WithLabelValues("api-beheerder", endpoint, status).Inc()
	externalServiceDuration.WithLabelValues("api-beheerder", endpoint).Observe(duration.Seconds())

	if err != nil {
		log.WithError(err).Error("API Beheerder request failed")
		return nil, fmt.Errorf("API Beheerder request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API Beheerder returned error %d: %s", resp.StatusCode, string(responseBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.WithFields(logrus.Fields{
		"service":     "api-beheerder",
		"endpoint":    endpoint,
		"status_code": resp.StatusCode,
		"duration":    duration.Milliseconds(),
	}).Debug("API Beheerder call completed")

	return result, nil
}

// Call Central Management for business rules, permissions, etc.
func callCentralManagement(method, endpoint string, data interface{}) (map[string]interface{}, error) {
	start := time.Now()
	config := getConfig()
	centralURL := config["central_mgmt_url"]
	centralKey := config["central_mgmt_key"]

	log.WithFields(logrus.Fields{
		"service":  "central-management",
		"method":   method,
		"endpoint": endpoint,
	}).Debug("Calling Central Management")

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, centralURL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Key", centralKey)

	resp, err := centralMgmtClient.Do(req)
	duration := time.Since(start)

	// Record metrics
	status := "error"
	if err == nil {
		status = fmt.Sprintf("%d", resp.StatusCode)
	}
	externalServiceCalls.WithLabelValues("central-management", endpoint, status).Inc()
	externalServiceDuration.WithLabelValues("central-management", endpoint).Observe(duration.Seconds())

	if err != nil {
		log.WithError(err).Error("Central Management request failed")
		return nil, fmt.Errorf("Central Management request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Central Management returned error %d: %s", resp.StatusCode, string(responseBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.WithFields(logrus.Fields{
		"service":     "central-management",
		"endpoint":    endpoint,
		"status_code": resp.StatusCode,
		"duration":    duration.Milliseconds(),
	}).Debug("Central Management call completed")

	return result, nil
}