package services

import (
	"order-service/internal/application/dto"
	"order-service/internal/domain/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository is a mock implementation of the OrderRepository interface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(order models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByID(id uint) (*models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) FindAll() ([]models.Order, error) {
	args := m.Called()
	return args.Get(0).([]models.Order), args.Error(1)
}

// MockEventPublisher is a mock implementation of the EventPublisher interface
type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) Publish(event interface{}) error {
	args := m.Called(event)
	return args.Error(0)
}

// TestCreateOrder tests the CreateOrder method for a successful case
func TestCreateOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockPublisher := new(MockEventPublisher)

	service := NewOrderService(mockRepo, mockPublisher)

	mockTime := time.Date(2025, time.January, 8, 19, 38, 18, 365284100, time.Local)

	orderDto := dto.OrderCreateDto{
		OrderID:    "test-123",
		CustomerID: 123,
		OrderItems: []dto.OrderItemDto{
			{ProductID: 1, Quantity: 2, Price: 9.99},
		},
	}

	orderDto.NewOrderCreateDto()

	expectedOrder := models.Order{
		OrderID:    orderDto.OrderID,
		CustomerID: orderDto.CustomerID,
		OrderDate:  mockTime,
		CreatedAt:  mockTime,
		UpdatedAt:  mockTime,
	}
	expectedOrder.AddItem(models.OrderItem{
		ProductID: 1,
		Quantity:  2,
		Price:     9.99,
	})

	// Set up mock expectations
	mockRepo.On("Save", expectedOrder).Return(0)
	mockPublisher.On("Publish", mock.AnythingOfType("events.OrderCreatedEvent")).Return(nil)

	orderResponse, err := service.CreateOrder(orderDto)
	assert.NoError(t, err)
	assert.Equal(t, "test-123", orderResponse.OrderID)
	assert.Equal(t, uint(123), orderResponse.CustomerID)
	assert.Equal(t, 19.98, orderResponse.TotalAmount)
	assert.Len(t, orderResponse.Items, 1)

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

// TestGetOrderByID tests the GetOrderByID method for a successful case
func TestGetOrderByID(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockPublisher := new(MockEventPublisher)

	service := NewOrderService(mockRepo, mockPublisher)

	sampleOrder := models.Order{
		ID:         1,
		CustomerID: 123,
		OrderItems: []models.OrderItem{
			{ProductID: 1, Quantity: 2, Price: 9.99},
		},
		TotalAmount: 19.98,
	}

	// Set up mock expectations
	mockRepo.On("FindByID", uint(1)).Return(&sampleOrder, nil)

	orderResponse, err := service.GetOrderByID(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), orderResponse.OrderID)
	assert.Equal(t, uint(123), orderResponse.CustomerID)
	assert.Equal(t, 19.98, orderResponse.TotalAmount)
	assert.Len(t, orderResponse.Items, 1)

	mockRepo.AssertExpectations(t)
}

// TestGetAllOrders tests the GetAllOrders method for a successful case
func TestGetAllOrders(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockPublisher := new(MockEventPublisher)

	service := NewOrderService(mockRepo, mockPublisher)

	sampleOrders := []models.Order{
		{
			ID:         1,
			CustomerID: 123,
			OrderItems: []models.OrderItem{
				{ProductID: 1, Quantity: 2, Price: 9.99},
			},
			TotalAmount: 19.98,
		},
	}

	// Set up mock expectations
	mockRepo.On("FindAll").Return(sampleOrders, nil)

	ordersResponse, err := service.GetAllOrders()
	assert.NoError(t, err)
	assert.Len(t, ordersResponse, 1)
	assert.Equal(t, uint(1), ordersResponse[0].OrderID)
	assert.Equal(t, uint(123), ordersResponse[0].CustomerID)
	assert.Equal(t, 19.98, ordersResponse[0].TotalAmount)
	assert.Len(t, ordersResponse[0].Items, 1)

	mockRepo.AssertExpectations(t)
}

// TestAddItemToOrder tests the AddItemToOrder method for a successful case
func TestAddItemToOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockPublisher := new(MockEventPublisher)

	service := NewOrderService(mockRepo, mockPublisher)

	sampleOrder := models.Order{
		ID:         1,
		CustomerID: 123,
		OrderItems: []models.OrderItem{
			{ProductID: 1, Quantity: 2, Price: 9.99},
		},
		TotalAmount: 19.98,
	}

	newItem := dto.OrderItemDto{
		ProductID: 2,
		Quantity:  1,
		Price:     5.99,
	}

	updatedOrder := sampleOrder
	updatedOrder.AddItem(models.OrderItem{
		ProductID: newItem.ProductID,
		Quantity:  newItem.Quantity,
		Price:     newItem.Price,
	})

	// Set up mock expectations
	mockRepo.On("FindByID", uint(1)).Return(&sampleOrder, nil)
	mockRepo.On("Save", updatedOrder).Return(nil)

	orderResponse, err := service.AddItemToOrder(1, newItem)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), orderResponse.OrderID)
	assert.Equal(t, uint(123), orderResponse.CustomerID)
	assert.Equal(t, 25.97, orderResponse.TotalAmount)
	assert.Len(t, orderResponse.Items, 2)

	mockRepo.AssertExpectations(t)
}
