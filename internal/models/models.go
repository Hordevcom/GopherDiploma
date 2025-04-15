package models

import "time"

type Order struct {
	Number   string    `json:"number"`
	Status   string    `json:"status"`
	Accrual  int       `json:"accrual,omitempty"`
	UploadAt time.Time `json:"upload_at"`
}

type OrderFloat struct {
	Number   string    `json:"number"`
	Status   string    `json:"status"`
	Accrual  float64   `json:"accrual,omitempty"`
	UploadAt time.Time `json:"upload_at"`
}

type UserBalance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}
