package models

type User struct {
	ID          int           `json:"id"`
	FullName    string        `json:"fullName" gorm:"type: varchar(255)" form:"fullName"`
	Email       string        `json:"email" gorm:"type: varchar(255)" form:"email"`
	Password    string        `json:"password" gorm:"type: varchar(255)" form:"password"`
	Phone       string        `json:"phone" gorm:"type: varchar(15)" form:"phone"`
	Address     string        `json:"address" gorm:"type: varchar(255)" form:"address"`
	Role        string        `json:"role" gorm:"type: varchar(255)"`
	Image       string        `json:"image" gorm:"type: varchar(255)" form:"phone"`
	Transaction []Transaction `json:"transaction" gorm:"foreignKey:UserID"`
}

type UserResponse struct {
	ID       int    `json:"-"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	// Role     string `json:"role"`
}

func (UserResponse) TableName() string {
	return "users"
}
