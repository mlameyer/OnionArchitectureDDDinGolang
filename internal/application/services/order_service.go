package services

import (
	"order-service/internal/application/dto"
)

type OrderService interface {
	CreateOrder(order dto.OrderCreateDto) (dto.OrderResponse, error)
	GetOrderByID(id uint) (*dto.OrderResponse, error)
	GetAllOrders() ([]dto.OrderResponse, error)
	AddItemToOrder(id uint, item dto.OrderItemDto) (*dto.OrderResponse, error)
}
