package models

// Album represents an album in the system
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title" binding:"required,min=1,max=200"`
	Artist string  `json:"artist" binding:"required,min=1,max=100"`
	Price  float64 `json:"price" binding:"required,min=0,max=999999"`
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
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=100"`
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
	CurrentPassword string `json:"current_password" binding:"required,min=8,max=100"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=100"`
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
	Username string   `json:"username" binding:"required,min=3,max=50,alphanum"`
	Email    string   `json:"email" binding:"required,email,max=100"`
	Password string   `json:"password" binding:"required,min=8,max=100"`
	Roles    []string `json:"roles" binding:"dive,min=1,max=50"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email  string   `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Roles  []string `json:"roles,omitempty" binding:"omitempty,dive,min=1,max=50"`
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
	Role string `json:"role" binding:"required,min=1,max=50"`
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

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	TotalItems int         `json:"total_items"`
}

// GetPage returns the page number (defaults to 1)
func (p *PaginationParams) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

// GetPageSize returns the page size (defaults to 20, max 100)
func (p *PaginationParams) GetPageSize() int {
	if p.PageSize < 1 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset calculates the offset for database queries
func (p *PaginationParams) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}
