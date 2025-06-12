package model

import "time"

type SignUpForm struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	ID   string
	Role string
}
