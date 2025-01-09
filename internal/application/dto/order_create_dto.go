package dto

import "time"

type OrderCreateDto struct {
	OrderID    string
	CustomerID uint
	OrderItems []OrderItemDto
	OrderDate  time.Time
}

type OrderItemDto struct {
	ProductID uint
	Quantity  int
	Price     float64
}

// NewOrderCreateDto is a constructor for OrderCreateDto
func (o *OrderCreateDto) NewOrderCreateDto() OrderCreateDto {
	return OrderCreateDto{
		OrderDate: time.Now(),
	}
}
