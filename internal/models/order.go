package models

import "time"

type Order struct {
	Number   string    `json:"number"`
	Status   string    `json:"status"`
	Accrual  int       `json:"accural,omitempty"`
	UploadAt time.Time `json:"upload_at"`
}
