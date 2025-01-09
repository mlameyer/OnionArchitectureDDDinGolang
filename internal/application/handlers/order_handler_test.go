package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"order-service/internal/application/dto"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderService is a mock implementation of the OrderService
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) GetOrderByID(id uint) (*dto.OrderResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*dto.OrderResponse), args.Error(1)
}

func (m *MockOrderService) CreateOrder(order dto.OrderCreateDto) (dto.OrderResponse, error) {
	args := m.Called(order)
	return args.Get(0).(dto.OrderResponse), args.Error(1)
}

func (m *MockOrderService) GetAllOrders() ([]dto.OrderResponse, error) {
	args := m.Called()
	return args.Get(0).([]dto.OrderResponse), args.Error(1)
}

func (m *MockOrderService) AddItemToOrder(orderID uint, item dto.OrderItemDto) (*dto.OrderResponse, error) {
	args := m.Called(orderID, item)
	return args.Get(0).(*dto.OrderResponse), args.Error(1)
}

// TestGetOrderByID tests the GetOrderByID handler for a successful case
func TestGetOrderByID(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Create a mock order service
	mockService := new(MockOrderService)

	// Create a sample order response
	sampleOrder := &dto.OrderResponse{
		OrderID:    "Test-123",
		CustomerID: 123,
		Items: []dto.OrderItemResponse{
			{ProductID: 1, Quantity: 2, Price: 9.99},
		},
		TotalAmount: 19.98,
	}

	// Set up mock expectations
	mockService.On("GetOrderByID", uint(1)).Return(sampleOrder, nil)

	// Initialize the order handler with routes
	NewOrderHandler(app, mockService)

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/orders/1", nil)
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP response recorder
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the response body
	var orderResponse dto.OrderResponse
	err = json.NewDecoder(resp.Body).Decode(&orderResponse)
	assert.NoError(t, err)
	assert.Equal(t, sampleOrder, &orderResponse)

	// Assert that the expectations were met
	mockService.AssertExpectations(t)
}
