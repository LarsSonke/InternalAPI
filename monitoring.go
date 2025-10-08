package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Prometheus metrics
var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "internal_api_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "internal_api_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	externalServiceCalls = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "internal_api_external_calls_total",
			Help: "Total number of external service calls",
		},
		[]string{"service", "endpoint", "status"},
	)

	externalServiceDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "internal_api_external_duration_seconds",
			Help:    "Duration of external service calls",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "endpoint"},
	)
)

// Setup structured logging
func setupLogging() {
	log = logrus.New()
	
	// Set log level based on environment
	logLevel := getEnvOrDefault("LOG_LEVEL", "INFO")
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	case "WARN":
		log.SetLevel(logrus.WarnLevel)
	case "ERROR":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// Use JSON formatter for structured logs
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	log.WithField("service", "internal-api").Info("Logging initialized")
}

// Setup Prometheus metrics
func setupMetrics() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(externalServiceCalls)
	prometheus.MustRegister(externalServiceDuration)
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	
	log.WithField("request_id", requestID).Info("Health check requested")
	
	health := gin.H{
		"status":     "healthy",
		"service":    "internal-api",
		"version":    "1.0.0",
		"timestamp":  time.Now().Unix(),
		"request_id": requestID,
	}

	dependencies := gin.H{}
	
	// Check API Beheerder connectivity
	apiBeheerderStatus := checkAPIBeheerderHealth()
	dependencies["api_beheerder"] = apiBeheerderStatus
	
	// Check Central Management connectivity
	centralMgmtStatus := checkCentralManagementHealth()
	dependencies["central_management"] = centralMgmtStatus
	
	health["dependencies"] = dependencies
	
	// Determine overall health status
	allHealthy := apiBeheerderStatus["status"] == "healthy" && centralMgmtStatus["status"] == "healthy"
	
	if allHealthy {
		health["status"] = "healthy"
		log.WithField("request_id", requestID).Info("Health check passed - all dependencies healthy")
		c.JSON(http.StatusOK, health)
	} else {
		health["status"] = "degraded"
		log.WithField("request_id", requestID).Warn("Health check degraded - some dependencies unhealthy")
		c.JSON(http.StatusServiceUnavailable, health)
	}
}

// Check API Beheerder health
func checkAPIBeheerderHealth() gin.H {
	start := time.Now()
	config := getConfig()
	beheerderURL := config["api_beheerder_url"]
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", beheerderURL+"/health", nil)
	if err != nil {
		return gin.H{
			"status": "unhealthy",
			"error":  "failed to create health check request",
		}
	}
	
	req.Header.Set("X-Service-Key", config["api_beheerder_key"])
	
	resp, err := apiBeheerderClient.Do(req)
	duration := time.Since(start)
	
	if err != nil {
		log.WithError(err).Error("API Beheerder health check failed")
		return gin.H{
			"status":   "unhealthy",
			"error":    err.Error(),
			"duration": duration.Milliseconds(),
		}
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return gin.H{
			"status":      "unhealthy",
			"status_code": resp.StatusCode,
			"duration":    duration.Milliseconds(),
		}
	}
	
	return gin.H{
		"status":   "healthy",
		"duration": duration.Milliseconds(),
	}
}

// Check Central Management health
func checkCentralManagementHealth() gin.H {
	start := time.Now()
	config := getConfig()
	centralURL := config["central_mgmt_url"]
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, "GET", centralURL+"/health", nil)
	if err != nil {
		return gin.H{
			"status": "unhealthy",
			"error":  "failed to create health check request",
		}
	}
	
	req.Header.Set("X-Service-Key", config["central_mgmt_key"])
	
	resp, err := centralMgmtClient.Do(req)
	duration := time.Since(start)
	
	if err != nil {
		log.WithError(err).Error("Central Management health check failed")
		return gin.H{
			"status":   "unhealthy",
			"error":    err.Error(),
			"duration": duration.Milliseconds(),
		}
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return gin.H{
			"status":      "unhealthy",
			"status_code": resp.StatusCode,
			"duration":    duration.Milliseconds(),
		}
	}
	
	return gin.H{
		"status":   "healthy",
		"duration": duration.Milliseconds(),
	}
}