package model

import "time"

type SignUpForm struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	// SignUpMethod string `json:"signup_method,omitempty"`
}

type LoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Response struct {
	Message    any   `json:"message"`
	StatusCode int   `json:"status_code"`
	Error      error `json:"error"`
	Data       any   `json:"data"` // alias for interface{} and equivalent in all ways.
}

type LogData struct {
	ClientIp          string
	UserAgent         string
	RequestedResource string
	At                time.Time
}

type UserCtxData struct {
	ID          string   `json:"id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}
