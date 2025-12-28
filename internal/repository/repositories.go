package repository

import (
	"project-kelas-santai/internal/database"
	"project-kelas-santai/internal/models"

	"gorm.io/gorm"
)

// Menu Repository
type MenuRepository interface {
	FindAll() ([]models.Menu, error)
	FindByID(id string) (*models.Menu, error)
	// Add Create/Update/Delete if needed
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository() MenuRepository {
	return &menuRepository{db: database.DB}
}

func (r *menuRepository) FindAll() ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Preload("Category").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindByID(id string) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.Preload("Category").First(&menu, "id = ?", id).Error
	return &menu, err
}

// Order Repository
type OrderRepository interface {
	Create(order *models.Order) error
	UpdateStatus(orderID string, status string) error
	UpdateFromNotification(orderID string, status string) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() OrderRepository {
	return &orderRepository{db: database.DB}
}

func (r *orderRepository) Create(order *models.Order) error {
	// Create order with associations (Items, Transaction)
	return r.db.Create(order).Error
}

func (r *orderRepository) UpdateStatus(orderID string, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *orderRepository) UpdateFromNotification(orderID string, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}
