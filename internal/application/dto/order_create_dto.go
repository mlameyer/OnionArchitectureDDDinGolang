package dto

type OrderCreateDto struct {
	OrderID    string
	CustomerID uint
	OrderItems []OrderItemDto
}

type OrderItemDto struct {
	ProductID uint
	Quantity  int
	Price     float64
}
