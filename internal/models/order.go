package models

import "time"

type Order struct {
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	Accrual   int       `json:"accural,omitempty"`
	Upload_at time.Time `json:"upload_at"`
}
