package dto

import "time"

// Banner represents a banner entity
type Banner struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Click represents a click entity
type Click struct {
	ID        int       `json:"id" db:"id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	BannerID  int       `json:"banner_id" db:"bannerid"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
