package services

import (
	"order-service/internal/application/dto"
	"order-service/internal/domain/events"
	"order-service/internal/domain/models"
	"order-service/internal/domain/repositories"
	"order-service/internal/infrastructure/logging"
)

type OrderService struct {
	repo           repositories.OrderRepository
	eventPublisher EventPublisher
}

func NewOrderService(repo repositories.OrderRepository, eventPublisher EventPublisher) *OrderService {
	return &OrderService{repo: repo, eventPublisher: eventPublisher}
}

func (s *OrderService) CreateOrder(orderDto dto.OrderCreateDto) (dto.OrderResponse, error) {
	newOrder := models.Order{
		OrderID:    orderDto.OrderID,
		CustomerID: orderDto.CustomerID,
		OrderDate:  orderDto.OrderDate,
		CreatedAt:  orderDto.OrderDate,
		UpdatedAt:  orderDto.OrderDate,
	}

	for _, item := range orderDto.OrderItems {
		newOrder.AddItem(
			models.OrderItem{
				OrderID:   newOrder.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			},
		)
	}

	if err := newOrder.Validate(); err != nil {
		return dto.OrderResponse{}, err
	}

	err := s.repo.Save(newOrder)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	event := events.OrderCreatedEvent{
		OrderID:     newOrder.ID,
		CustomerID:  newOrder.CustomerID,
		TotalAmount: newOrder.TotalAmount,
	}

	err = s.eventPublisher.Publish(event)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	return dto.OrderResponse{
		OrderID:     newOrder.OrderID,
		CustomerID:  newOrder.CustomerID,
		TotalAmount: newOrder.TotalAmount,
		Items:       convertToOrderItemResponse(newOrder.OrderItems),
	}, nil
}

func (s *OrderService) GetOrderByID(id uint) (*dto.OrderResponse, error) {
	order, err := s.repo.FindByID(id)
	logging.Logger.Info().Msgf("%v", order)
	if err != nil {
		return nil, err
	}

	return &dto.OrderResponse{
		OrderID:     order.OrderID,
		CustomerID:  order.CustomerID,
		TotalAmount: order.TotalAmount,
		Items:       convertToOrderItemResponse(order.OrderItems),
	}, nil
}

func (s *OrderService) GetAllOrders() ([]dto.OrderResponse, error) {
	orders, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	ordersResponse := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		ordersResponse[i] = dto.OrderResponse{
			OrderID:     order.OrderID,
			CustomerID:  order.CustomerID,
			Items:       convertToOrderItemResponse(order.OrderItems),
			TotalAmount: order.TotalAmount,
		}
	}

	return ordersResponse, nil
}

func (s *OrderService) AddItemToOrder(id uint, item dto.OrderItemDto) (*dto.OrderResponse, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	order.AddItem(models.OrderItem{
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Price:     item.Price,
	})
	if err := order.Validate(); err != nil {
		return nil, err
	}

	err = s.repo.Save(*order)
	if err != nil {
		return nil, err
	}

	response := dto.OrderResponse{
		OrderID:     order.OrderID,
		CustomerID:  order.CustomerID,
		Items:       convertToOrderItemResponse(order.OrderItems),
		TotalAmount: order.TotalAmount,
	}

	return &response, nil
}

func convertToOrderItemResponse(items []models.OrderItem) []dto.OrderItemResponse {
	response := make([]dto.OrderItemResponse, len(items))
	for i, item := range items {
		response[i] = dto.OrderItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return response
}
