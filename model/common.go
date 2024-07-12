package model

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
	Message    interface{} `json:"message"`
	StatusCode int         `json:"status_code"`
	Error      error       `json:"error"`
	Data       interface{} `json:"data"`
}
