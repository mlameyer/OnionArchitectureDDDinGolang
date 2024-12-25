package handlers

import (
	"order-service/internal/application/dto"
	"order-service/internal/domain/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// OrderHandler handles order-related API requests
type OrderHandler struct {
	service *services.OrderService
}

// NewOrderHandler initializes the order handler with routes
func NewOrderHandler(app *fiber.App, service *services.OrderService) {
	handler := &OrderHandler{service}
	app.Post("/orders", handler.CreateOrder)
	app.Get("/orders/:id", handler.GetOrderByID)
	app.Get("/orders", handler.GetAllOrders)
	app.Post("/orders/:id/items", handler.AddItemToOrder)
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with items
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Order"
// @Success 201 {object} dto.OrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var order dto.OrderCreateDto
	var orderResponse dto.OrderResponse
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: err.Error()})
	}
	orderResponse, err := h.service.CreateOrder(order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(orderResponse)
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Get order details by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Success 404 {object} dto.OrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: err.Error()})
	}
	order, err := h.service.GetOrderByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: err.Error()})
	}
	if order == nil {
		return c.Status(fiber.StatusNotFound).JSON(order)
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// GetAllOrders godoc
// @Summary Get all orders
// @Description Get a list of all orders
// @Tags orders
// @Produce json
// @Success 200 {array} dto.OrderResponse
// @Success 404 {array} dto.OrderResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: err.Error()})
	}
	if len(orders) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(orders)
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

// AddItemToOrder godoc
// @Summary Add item to order
// @Description Add a new item to an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param item body models.OrderItem true "Order Item"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/{id}/items [post]
func (h *OrderHandler) AddItemToOrder(c *fiber.Ctx) error {
	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	id := uint(orderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: err.Error()})
	}

	var item dto.OrderItemDto
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: err.Error()})
	}

	response, err := h.service.AddItemToOrder(id, item)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
