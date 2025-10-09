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