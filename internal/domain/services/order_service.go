package services

import (
	"order-service/internal/domain/events"
	"order-service/internal/domain/models"
	"order-service/internal/domain/repositories"
)

type OrderService struct {
	repo           repositories.OrderRepository
	eventPublisher EventPublisher
}

func NewOrderService(repo repositories.OrderRepository, eventPublisher EventPublisher) *OrderService {
	return &OrderService{repo: repo, eventPublisher: eventPublisher}
}

func (s *OrderService) CreateOrder(order models.Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	err := s.repo.Save(order)
	if err != nil {
		return err
	}

	event := events.OrderCreatedEvent{
		OrderID:     order.ID,
		CustomerID:  order.CustomerID,
		TotalAmount: order.TotalAmount,
	}
	return s.eventPublisher.Publish(event)
}

func (s *OrderService) GetOrderByID(id uint) (*models.Order, error) {
	return s.repo.FindByID(id)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.repo.FindAll()
}

func (s *OrderService) AddItemToOrder(order *models.Order, item models.OrderItem) error {
	order.AddItem(item)
	if err := order.Validate(); err != nil {
		return err
	}
	return s.repo.Save(*order)
}
