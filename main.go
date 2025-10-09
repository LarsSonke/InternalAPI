package main

import (
	"fmt"

	"InternalAPI/internal/circuitbreaker"
	"InternalAPI/internal/config"
	"InternalAPI/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Global logger
var log *logrus.Logger

// Initialize logging and circuit breakers
func init() {
	setupLogging()
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize circuit breakers for external services
	circuitbreaker.Init("api-beheerder", cfg.CircuitBreakerFailureThreshold, cfg.CircuitBreakerTimeout, cfg.CircuitBreakerMaxRetries, cfg.CircuitBreakerRetryDelay)
	circuitbreaker.Init("central-mgmt", cfg.CircuitBreakerFailureThreshold, cfg.CircuitBreakerTimeout, cfg.CircuitBreakerMaxRetries, cfg.CircuitBreakerRetryDelay)

	log.WithFields(logrus.Fields{
		"failure_threshold": cfg.CircuitBreakerFailureThreshold,
		"timeout":          cfg.CircuitBreakerTimeout,
		"max_retries":      cfg.CircuitBreakerMaxRetries,
		"retry_delay":      cfg.CircuitBreakerRetryDelay,
	}).Info("Circuit breakers initialized")

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router with middleware
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Add CORS middleware for User Portal access
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:3001", 
		"https://hotel-portal.local",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Internal-API-Key"}
	router.Use(cors.New(corsConfig))

	log.WithFields(logrus.Fields{
		"valid_origins": corsConfig.AllowOrigins,
	}).Info("Configured CORS origins for User Portal access")

	// Setup routes with handlers
	routes.Setup(router, cfg)

	// Start server
	address := cfg.Host + ":" + cfg.Port
	
	log.WithFields(logrus.Fields{
		"address":            address,
		"api_beheerder_url":  cfg.APIBeheerderURL,
		"central_mgmt_url":   cfg.CentralMgmtURL,
		"cors_origins":       cfg.AllowedOrigins,
		"user_portal_url":    cfg.UserPortalURL,
		"api_endpoint":       "http://" + address + "/api/v1",
		"health_endpoint":    "http://" + address + "/health", 
		"metrics_endpoint":   "http://" + address + "/metrics",
	}).Info("Hotel Internal API started successfully")

	// Pretty startup messages
	fmt.Printf("🚀 Internal API starting on %s\n", address)
	fmt.Printf("   🔗 API Beheerder: %s\n", cfg.APIBeheerderURL)
	fmt.Printf("   🎛️  Central Management: %s\n", cfg.CentralMgmtURL)
	fmt.Printf("   👤 User Portal: %s\n", cfg.UserPortalURL)
	fmt.Printf("   📊 Metrics: http://%s/metrics\n", address)
	fmt.Printf("   💚 Health: http://%s/health\n", address)

	if err := router.Run(address); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}

// setupLogging configures structured logging
func setupLogging() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	
	log.WithFields(logrus.Fields{
		"service": "internal-api",
	}).Info("Logging initialized")
}