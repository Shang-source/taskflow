package models

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

// LoginInput Login request body
// swagger:model
type LoginInput struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"123456"`
}

// LoginResponse Login success response
// swagger:model
type LoginResponse struct {
	Message string `json:"message" example:"Login successful"`
	Token   string `json:"token"   example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..."` // JWT example
}

// ErrorResponse Standard error response
// swagger:model
type ErrorResponse struct {
	Error string `json:"error" example:"Username or Password is incorrect"`
}
