package models

import (
	"errors"
	"time"
)

type Order struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	OrderID     string `gorm:"index"`
	CustomerID  uint
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
	TotalAmount float64
	OrderDate   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// OrderItem struct definition
type OrderItem struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	OrderID   string `gorm:"index"`
	ProductID uint
	Quantity  int
	Price     float64
}

// Validate method for the Order struct
func (o *Order) Validate() error {
	if o.CustomerID == 0 {
		return errors.New("customer ID is required")
	}
	if len(o.OrderItems) == 0 {
		return errors.New("order must contain at least one item")
	}
	for _, item := range o.OrderItems {
		if item.Quantity <= 0 {
			return errors.New("order item quantity must be greater than zero")
		}
		if item.Price <= 0 {
			return errors.New("order item price must be greater than zero")
		}
	}
	if o.TotalAmount <= 0 {
		return errors.New("total amount must be greater than zero")
	}
	if o.OrderDate.IsZero() {
		return errors.New("order date is required")
	}
	return nil
}

func (o *Order) AddItem(item OrderItem) {
	o.OrderItems = append(o.OrderItems, item)
	o.TotalAmount += item.Price * float64(item.Quantity)
}
