package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID            string         `gorm:"type:char(36);primaryKey" json:"id"`
	CustomerName  string         `gorm:"type:varchar(100);not null" json:"customer_name"`
	TableNumber   string         `gorm:"type:varchar(10);not null" json:"table_number"`
	PaymentMethod string         `gorm:"type:varchar(50)" json:"payment_method"`
	OrderItems    []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`
	Transaction   *Transaction   `gorm:"foreignKey:OrderID" json:"transaction_summary"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Status        string         `json:"status"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	return
}

type OrderItem struct {
	ID         string  `gorm:"type:char(36);primaryKey" json:"id"`
	OrderID    string  `gorm:"type:char(36);not null" json:"order_id"`
	MenuID     string  `gorm:"type:char(36)" json:"menu_id"` // Optional link to actual menu
	Name       string  `gorm:"type:varchar(255);not null" json:"name"`
	Quantity   int     `gorm:"not null" json:"quantity"`
	Price      float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	TotalPrice float64 `gorm:"type:decimal(10,2);not null" json:"total_price"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if oi.ID == "" {
		oi.ID = uuid.NewString()
	}
	return
}

type Transaction struct {
	ID         string  `gorm:"type:char(36);primaryKey" json:"id"`
	OrderID    string  `gorm:"type:char(36);not null" json:"order_id"`
	Subtotal   float64 `gorm:"type:decimal(10,2)" json:"subtotal"`
	TaxAmount  float64 `gorm:"type:decimal(10,2)" json:"tax_amount"`
	GrandTotal float64 `gorm:"type:decimal(10,2)" json:"grand_total"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}
