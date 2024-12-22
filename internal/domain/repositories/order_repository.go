package repositories

import "order-service/internal/domain/models"

type OrderRepository interface {
	Save(order models.Order) error
	FindByID(id uint) (*models.Order, error)
	FindAll() ([]models.Order, error)
}
