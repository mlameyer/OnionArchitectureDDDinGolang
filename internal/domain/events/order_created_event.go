package events

type OrderCreatedEvent struct {
	OrderID     uint
	CustomerID  uint
	TotalAmount float64
}
