package dto

// OrderResponse represents an order response
type OrderResponse struct {
	OrderID     uint                `json:"order_id"`
	CustomerID  uint                `json:"customer_id"`
	Items       []OrderItemResponse `json:"items"`
	TotalAmount float64             `json:"total_amount"`
}

// OrderItemResponse represents an order item response
type OrderItemResponse struct {
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
