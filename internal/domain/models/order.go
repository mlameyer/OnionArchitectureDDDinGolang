package models

import "time"

type Order struct {
	ID        uint   `gorm:"primaryKey"`
	OrderId   string `gorm:"index"` // Create an index on OrderId
	CreatedAt time.Time
	UpdatedAt time.Time
}
