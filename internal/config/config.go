package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server settings
	Host string
	Port string

	// JWT settings for User Portal authentication
	JWTSecret string

	// External services
	APIBeheerderURL string
	APIBeheerderKey string
	CentralMgmtURL  string
	CentralMgmtKey  string

	// CORS settings
	UserPortalURL  string
	AllowedOrigins string

	// Circuit breaker configuration
	CircuitBreakerFailureThreshold int
	CircuitBreakerTimeout          time.Duration
	CircuitBreakerMaxRetries       int
	CircuitBreakerRetryDelay       time.Duration

	// Security settings
	MaxRequestBodySize     int64         // Maximum request body size in bytes
	RequestTimeout         time.Duration // Maximum time for a request
	ReadTimeout            time.Duration // Maximum time to read request
	WriteTimeout           time.Duration // Maximum time to write response
	IdleTimeout            time.Duration // Maximum time for idle connections
	EnableSecurityHeaders  bool          // Enable security headers
	EnableAuditLogging     bool          // Enable audit logging

	// Rate limiting settings
	RateLimitEnabled       bool          // Enable rate limiting
	RateLimitRequests      int           // Requests per interval for general API
	RateLimitInterval      time.Duration // Time window for rate limiting
	LoginRateLimitRequests int           // Requests per interval for login
	LoginRateLimitInterval time.Duration // Time window for login rate limiting
	AdminRateLimitRequests int           // Requests per interval for admin endpoints
	AdminRateLimitInterval time.Duration // Time window for admin rate limiting
}

// Load loads configuration from environment variables with sensible defaults
func Load() *Config {
	return &Config{
		// Server settings
		Host: getEnv("HOST", "localhost"),
		Port: getEnv("PORT", "8080"),

		// JWT settings
		JWTSecret: getEnv("JWT_SECRET", "your-jwt-secret-key"),

		// External services
		APIBeheerderURL: getEnv("API_BEHEERDER_URL", "http://localhost:8081"),
		APIBeheerderKey: getEnv("API_BEHEERDER_KEY", "beheerder-service-key"),
		CentralMgmtURL:  getEnv("CENTRAL_MGMT_URL", "http://localhost:8082"),
		CentralMgmtKey:  getEnv("CENTRAL_MGMT_KEY", "central-mgmt-service-key"),

		// CORS settings
		UserPortalURL:  getEnv("USER_PORTAL_URL", "http://localhost:3000"),
		AllowedOrigins: getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001,https://hotel-portal.local"),

		// Circuit breaker defaults
		CircuitBreakerFailureThreshold: getEnvInt("CB_FAILURE_THRESHOLD", 5),
		CircuitBreakerTimeout:          time.Duration(getEnvInt("CB_TIMEOUT_SECONDS", 60)) * time.Second,
		CircuitBreakerMaxRetries:       getEnvInt("CB_MAX_RETRIES", 3),
		CircuitBreakerRetryDelay:       time.Duration(getEnvInt("CB_RETRY_DELAY_MS", 1000)) * time.Millisecond,

		// Security settings
		MaxRequestBodySize:    int64(getEnvInt("MAX_REQUEST_BODY_SIZE", 5*1024*1024)), // 5MB default
		RequestTimeout:        time.Duration(getEnvInt("REQUEST_TIMEOUT_SECONDS", 30)) * time.Second,
		ReadTimeout:           time.Duration(getEnvInt("READ_TIMEOUT_SECONDS", 15)) * time.Second,
		WriteTimeout:          time.Duration(getEnvInt("WRITE_TIMEOUT_SECONDS", 15)) * time.Second,
		IdleTimeout:           time.Duration(getEnvInt("IDLE_TIMEOUT_SECONDS", 60)) * time.Second,
		EnableSecurityHeaders: getEnvBool("ENABLE_SECURITY_HEADERS", true),
		EnableAuditLogging:    getEnvBool("ENABLE_AUDIT_LOGGING", true),

		// Rate limiting settings
		RateLimitEnabled:       getEnvBool("RATE_LIMIT_ENABLED", true),
		RateLimitRequests:      getEnvInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitInterval:      time.Duration(getEnvInt("RATE_LIMIT_INTERVAL_SECONDS", 60)) * time.Second,
		LoginRateLimitRequests: getEnvInt("LOGIN_RATE_LIMIT_REQUESTS", 5),
		LoginRateLimitInterval: time.Duration(getEnvInt("LOGIN_RATE_LIMIT_INTERVAL_SECONDS", 300)) * time.Second, // 5 minutes
		AdminRateLimitRequests: getEnvInt("ADMIN_RATE_LIMIT_REQUESTS", 50),
		AdminRateLimitInterval: time.Duration(getEnvInt("ADMIN_RATE_LIMIT_INTERVAL_SECONDS", 60)) * time.Second,
	}
}

// GetAPIBeheerderURL returns the API Beheerder URL for services package
func (c *Config) GetAPIBeheerderURL() string {
	return c.APIBeheerderURL
}

// GetAPIBeheerderKey returns the API Beheerder key for services package
func (c *Config) GetAPIBeheerderKey() string {
	return c.APIBeheerderKey
}

// GetCentralMgmtURL returns the Central Management URL for services package
func (c *Config) GetCentralMgmtURL() string {
	return c.CentralMgmtURL
}

// GetCentralMgmtKey returns the Central Management key for services package
func (c *Config) GetCentralMgmtKey() string {
	return c.CentralMgmtKey
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as int or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool gets an environment variable as bool or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
