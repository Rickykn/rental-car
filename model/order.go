package model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderID         uuid.UUID
	CarID           uuid.UUID
	OrderDate       time.Time
	PickupDate      time.Time
	DropoffDate     time.Time
	PickupLocation  string
	DropoffLocation string
	OrderCode       string
}

type BookReq struct {
	CarID           uuid.UUID `json:"car_id" `
	PickupDate      string    `json:"pickup_date" `
	DropoffDate     string    `json:"dropoff_date" `
	PickupLocation  string    `json:"pickup_location" `
	DropoffLocation string    `json:"dropoff_location" `
}

type CheckinReq struct {
	OrderCode string `json:"order_code" `
}
