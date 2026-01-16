package model

import "time"

type Subscription struct {
	ID        int64      `json:"id"`
	Service   string     `json:"service_name"`
	Price     int        `json:"price"`
	UserID    string     `json:"user_id"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
