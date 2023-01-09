package dto

type UserResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Image    string `json:"image"`
	Role     string `json:"role"`
}

type UserDeleteResponse struct {
	ID int `json:"id"`
}
