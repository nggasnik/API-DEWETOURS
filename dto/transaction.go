package dto

import (
	"erlangga/models"

	
)

type TransactionRequest struct {
	CounterQty int    `json:"counterQty"  form:"counterQty"`
	Total      int    `json:"total"  form:"total"`
	Status     string `json:"status"  form:"status"`
	TripID     int    `json:"tripId" form:"tripId"`
	Attachment string `json:"attachment" form:"attachment"`
	UserId    int    `json:"userId" form:"userId"`
}

type TransactionResponse struct {
	ID          int                 `json:"id"`
	CounterQty  int                 `json:"counterQty"`
	Total       int                 `json:"total"`
	Status      string              `json:"status"`
	Attachment  string              `json:"attachment"`
	BookingDate string              `json:"bookingDate"`
	Trip        TripResponse        `json:"trip"`
	User        models.UserResponse `json:"user"`
}
