package models

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // 返回时可以省略密码
}

// LoginInput 登录请求体
// swagger:model
type LoginInput struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"123456"`
}

// LoginResponse 登录成功返回
// swagger:model
type LoginResponse struct {
	Message string `json:"message" example:"Login successful"`
	Token   string `json:"token"   example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..."` // JWT 示例
}

// ErrorResponse 通用错误返回
// swagger:model
type ErrorResponse struct {
	Error string `json:"error" example:"Username or Password is incorrect"`
}

type Message struct {
	Message string `json:"message"`
}
