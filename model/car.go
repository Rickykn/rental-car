package model

import "github.com/google/uuid"

type Car struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CarName   string    `json:"car_name" db:"car_name"`
	DayRate   float64   `json:"day_rate" db:"day_rate"`
	MonthRate float64   `json:"month_rate" db:"month_rate"`
	Image     string    `json:"image" db:"image"`
}
