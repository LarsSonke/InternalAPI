package routes

import (
	"InternalAPI/internal/config"
	"InternalAPI/internal/handlers"
	"InternalAPI/internal/middleware"
	
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Setup configures all routes for the application
func Setup(router *gin.Engine, config *config.Config) {
	// Create handler instances
	authHandlers := handlers.NewAuthHandlers(config)
	albumHandlers := handlers.NewAlbumHandlers(config)
	adminHandlers := handlers.NewAdminHandlers(config)

	// Public routes
	router.GET("/health", handlers.HealthHandler)
	router.GET("/health/circuit-breakers", handlers.GetCircuitBreakerStatusHandler)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	
	// Authentication routes
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandlers.Login)
		auth.POST("/refresh", authHandlers.RefreshToken)
	}

	// Protected routes (requires JWT authentication)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// Auth user info routes
		protected.POST("/auth/logout", authHandlers.Logout)
		protected.GET("/auth/me", authHandlers.GetUserInfo)
		protected.PUT("/auth/change-password", authHandlers.ChangePassword)

		// Album/Hotel management routes
		protected.GET("/albums", albumHandlers.GetAlbums)
		protected.GET("/albums/:id", albumHandlers.GetAlbumByID)
		protected.POST("/albums", albumHandlers.CreateAlbum)
		protected.PUT("/albums/:id", albumHandlers.UpdateAlbum)
		protected.DELETE("/albums/:id", albumHandlers.DeleteAlbum)
	}

	// Admin routes (requires JWT + admin role)
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.RequireRoles("admin", "super_admin"))
	{
		// User management
		admin.GET("/users", adminHandlers.GetUsers)
		admin.GET("/users/:id", adminHandlers.GetUserByID)
		admin.POST("/users", adminHandlers.CreateUser)
		admin.PUT("/users/:id", adminHandlers.UpdateUser)
		admin.DELETE("/users/:id", adminHandlers.DeleteUser)

		// Role management
		admin.GET("/roles", adminHandlers.GetRoles)
		admin.POST("/users/:id/roles", adminHandlers.AssignRole)
		admin.DELETE("/users/:id/roles/:role", adminHandlers.RemoveRole)

		// System management
		admin.GET("/system/stats", adminHandlers.GetSystemStats)
		admin.GET("/audit-logs", adminHandlers.GetAuditLogs)
		admin.POST("/circuit-breakers/:service/reset", handlers.ResetCircuitBreakerHandler)
	}
}