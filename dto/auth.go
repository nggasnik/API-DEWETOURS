package dto

type RegisterRequest struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Role    string `json:"role"`
}

type RegisterResponse struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Role    string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role 	string `json:"role"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
	Role string `json:"role"`
}
