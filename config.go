package main

import "os"

// Configuration management
func getConfig() map[string]string {
	return map[string]string{
		// Server settings
		"host": getEnvOrDefault("HOST", "localhost"),
		"port": getEnvOrDefault("PORT", "8080"),

		// JWT settings for User Portal authentication
		"jwt_secret": getEnvOrDefault("JWT_SECRET", "your-jwt-secret-key"),

		// API Beheerder settings (for data operations)
		"api_beheerder_url": getEnvOrDefault("API_BEHEERDER_URL", "http://localhost:8081"),
		"api_beheerder_key": getEnvOrDefault("API_BEHEERDER_KEY", "beheerder-service-key"),

		// Central Management System settings (for business rules, permissions, audit, etc.)
		"central_mgmt_url": getEnvOrDefault("CENTRAL_MGMT_URL", "http://localhost:8082"),
		"central_mgmt_key": getEnvOrDefault("CENTRAL_MGMT_KEY", "central-mgmt-service-key"),

		// User Portal CORS settings (plugins communicate with API Beheerder directly)
		"user_portal_url": getEnvOrDefault("USER_PORTAL_URL", "http://localhost:3000"),
		"allowed_origins": getEnvOrDefault("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001,https://hotel-portal.local"),
	}
}

// Get environment variable with fallback default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}