package main

import (
	"InternalAPI/internal/broker"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"InternalAPI/internal/circuitbreaker"
	"InternalAPI/internal/config"
	"InternalAPI/internal/middleware"
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

	// Validate JWT secret
	if cfg.JWTSecret == "your-jwt-secret-key" {
		log.Warn("⚠️  WARNING: Using default JWT secret! Set JWT_SECRET environment variable in production!")
	}

	// Initialize JWT middleware with secret
	middleware.InitJWT(cfg.JWTSecret)

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

	// Add security middleware
	if cfg.EnableSecurityHeaders {
		router.Use(middleware.SecurityHeaders())
		log.Info("Security headers enabled")
	}

	// Add request ID tracking
	router.Use(middleware.RequestID())

	// Add audit logging
	if cfg.EnableAuditLogging {
		router.Use(middleware.AuditLogger())
		log.Info("Audit logging enabled")
	}

	// Add request size limit
	router.Use(middleware.RequestSizeLimit(cfg.MaxRequestBodySize))
	log.WithField("max_size_mb", cfg.MaxRequestBodySize/(1024*1024)).Info("Request size limit configured")

	// Add CORS middleware for User Portal access
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:3001", 
		"https://hotel-portal.local",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Internal-API-Key", "X-Request-ID"}
	router.Use(cors.New(corsConfig))

	log.WithFields(logrus.Fields{
		"valid_origins": corsConfig.AllowOrigins,
	}).Info("Configured CORS origins for User Portal access")

	// Setup routes with handlers
	routes.Setup(router, cfg)

	// Create HTTP server with timeouts
	address := cfg.Host + ":" + cfg.Port
	srv := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.WithFields(logrus.Fields{
		"address":              address,
		"api_beheerder_url":    cfg.APIBeheerderURL,
		"central_mgmt_url":     cfg.CentralMgmtURL,
		"cors_origins":         cfg.AllowedOrigins,
		"user_portal_url":      cfg.UserPortalURL,
		"api_endpoint":         "http://" + address + "/api/v1",
		"health_endpoint":      "http://" + address + "/health",
		"metrics_endpoint":     "http://" + address + "/metrics",
		"read_timeout":         cfg.ReadTimeout,
		"write_timeout":        cfg.WriteTimeout,
		"idle_timeout":         cfg.IdleTimeout,
		"max_request_body_mb":  cfg.MaxRequestBodySize / (1024 * 1024),
		"rate_limit_enabled":   cfg.RateLimitEnabled,
		"security_headers":     cfg.EnableSecurityHeaders,
		"audit_logging":        cfg.EnableAuditLogging,
	}).Info("Hotel Internal API started successfully")

	// Pretty startup messages
	fmt.Printf("🚀 Internal API starting on %s\n", address)
	fmt.Printf("   🔗 API Beheerder: %s\n", cfg.APIBeheerderURL)
	fmt.Printf("   🎛️  Central Management: %s\n", cfg.CentralMgmtURL)
	fmt.Printf("   👤 User Portal: %s\n", cfg.UserPortalURL)
	fmt.Printf("   📊 Metrics: http://%s/metrics\n", address)
	fmt.Printf("   💚 Health: http://%s/health\n", address)
	fmt.Printf("   🔒 Security: Headers=%v, Audit=%v, RateLimit=%v\n", 
		cfg.EnableSecurityHeaders, cfg.EnableAuditLogging, cfg.RateLimitEnabled)
	fmt.Printf("   ⏱️  Timeouts: Read=%v, Write=%v, Idle=%v\n", 
		cfg.ReadTimeout, cfg.WriteTimeout, cfg.IdleTimeout)

		// Register with broker (non-blocking)
	broker.RegisterWithBroker(cfg.Host, cfg.Port)

// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server gracefully...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
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