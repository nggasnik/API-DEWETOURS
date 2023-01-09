package models

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	CounterQty  int       `json:"counterQty" gorm:"type: int"`
	Total       int       `json:"total" gorm:"type: int"`
	Status      string    `json:"status" gorm:"type: varchar(255)"`
	Attachment  string    `json:"attachment" gorm:"type: varchar(255)"`
	BookingDate time.Time `json:"bookingDate"`
	TripID      int       `json:"tripId" gorm:"type: int"`
	Trip        TripResponse
	UserID      int `json:"id_user" gorm:"type: int"`
	User        UserResponse
}

type TransactionResponse struct {
	ID         int    `json:"id"`
	CounterQty int    `json:"counterQty" gorm:"type: int"`
	Total      int    `json:"total" gorm:"type: int"`
	Status     string `json:"status" gorm:"type: varchar(255)"`
	Attachment string `json:"attachment" gorm:"type: varchar(255)"`
	TripID     int    `json:"id_trip" gorm:"type: int"`
	UserID     int    `json:"id_user" gorm:"type: int"`
}

func (TransactionResponse) TableName() string {
	return "transactions"
}
