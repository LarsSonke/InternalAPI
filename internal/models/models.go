package models

// Album represents an album in the system
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// UserInfo represents user information from JWT or external service
type UserInfo struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Exp      int64    `json:"exp"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest represents a change password request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// User represents a user in the system
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Active   bool     `json:"active"`
	Created  int64    `json:"created"`
	Modified int64    `json:"modified"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	Roles    []string `json:"roles"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email  string   `json:"email,omitempty" binding:"omitempty,email"`
	Roles  []string `json:"roles,omitempty"`
	Active *bool    `json:"active,omitempty"`
}

// Role represents a role in the system
type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// AssignRoleRequest represents a request to assign a role to a user
type AssignRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// SystemStats represents system statistics
type SystemStats struct {
	Timestamp      int64                  `json:"timestamp"`
	Uptime         float64                `json:"uptime_seconds"`
	TotalRequests  int64                  `json:"total_requests"`
	ActiveRequests int                    `json:"active_requests"`
	TotalUsers     int                    `json:"total_users"`
	ActiveUsers    int                    `json:"active_users"`
	TotalAlbums    int                    `json:"total_albums"`
	TotalRoles     int                    `json:"total_roles"`
	Services       map[string]interface{} `json:"services"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Timestamp string `json:"timestamp"`
	Details   string `json:"details,omitempty"`
}
