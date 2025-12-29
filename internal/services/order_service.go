package services

import (
	"fmt"
	"project-kelas-santai/internal/config"
	"project-kelas-santai/internal/models"
	"project-kelas-santai/internal/repository"
	"project-kelas-santai/pkg/tools"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// DTOs for the incoming payload
type CreateOrderRequest struct {
	Customer           CustomerDTO           `json:"customer"`
	Items              []ItemDTO             `json:"items"`
	TransactionSummary TransactionSummaryDTO `json:"transaction_summary"`
	CreatedAt          string                `json:"created_at"`
}

type CustomerDTO struct {
	Name          string `json:"name"`
	TableNumber   string `json:"table_number"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
}

type ItemDTO struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"` // Reverted to float64
}

type TransactionSummaryDTO struct {
	Subtotal   float64 `json:"subtotal"`
	TaxAmount  float64 `json:"tax_amount"` // 11%
	GrandTotal float64 `json:"grand_total"`
}

type OrderService interface {
	CreateOrder(req *CreateOrderRequest) (string, string, string, error)
	HandleNotification(payload map[string]interface{}) error
}

type orderService struct {
	repo       repository.OrderRepository
	cfg        *config.Config
	snapClient snap.Client
}

func NewOrderService(repo repository.OrderRepository, cfg *config.Config) OrderService {
	// Initialize Midtrans
	var env midtrans.EnvironmentType
	if cfg.Midtrans.Environment == "production" {
		env = midtrans.Production
	} else {
		env = midtrans.Sandbox
	}

	var s snap.Client
	s.New(cfg.Midtrans.ServerKey, env)

	return &orderService{
		repo:       repo,
		cfg:        cfg,
		snapClient: s,
	}
}

func (s *orderService) CreateOrder(req *CreateOrderRequest) (string, string, string, error) {
	// Parse CreatedAt
	createdAt, err := time.Parse(time.RFC3339, req.CreatedAt)
	if err != nil {
		createdAt = time.Now()
	}

	// Map DTO to Models
	var orderItems []models.OrderItem
	for _, item := range req.Items {
		orderItems = append(orderItems, models.OrderItem{
			MenuID:     item.ID, // Assuming ID in payload is MenuID
			Name:       item.Name,
			Quantity:   item.Quantity,
			Price:      item.Price,
			TotalPrice: item.TotalPrice,
		})
	}

	transaction := &models.Transaction{
		Subtotal:   req.TransactionSummary.Subtotal,
		TaxAmount:  req.TransactionSummary.TaxAmount,
		GrandTotal: req.TransactionSummary.GrandTotal,
	}

	order := &models.Order{
		CustomerName:  req.Customer.Name,
		TableNumber:   req.Customer.TableNumber,
		PaymentMethod: req.Customer.PaymentMethod,
		OrderItems:    orderItems,
		Transaction:   transaction,
		Email:         req.Customer.Email,
		CreatedAt:     createdAt,
	}

	order.Status = "pending"

	// 1. Save to DB
	if err := s.repo.Create(order); err != nil {
		return "", "", "", err
	}

	// 2. Create Snap Request
	if req.Customer.PaymentMethod == "Cash" {
		return order.ID, "", "", nil
	}
	if req.Customer.PaymentMethod == "QRIS" {
		snapReq := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  order.ID,
				GrossAmt: int64(transaction.GrandTotal),
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: req.Customer.Name,
			},
		}
		snapResp, snapErr := s.snapClient.CreateTransaction(snapReq)
		if snapErr != nil {
			return "", "", "", snapErr
		}
		//Send Email
		tools.SendOrderSuccessEmail(order.Email, order.CustomerName, order.ID, "Terima kasih telah memesan di Cafe Santai, pesanan Anda sedang di proses")
		return order.ID, snapResp.RedirectURL, snapResp.Token, nil
	}

	return order.ID, "", "", nil
}

func (s *orderService) HandleNotification(payload map[string]interface{}) error {
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return nil // or error
	}

	transactionStatus, ok := payload["transaction_status"].(string)
	if !ok {
		return nil // or error
	}

	fraudStatus, _ := payload["fraud_status"].(string)

	var newStatus string

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			newStatus = "challenge"
		} else if fraudStatus == "accept" {
			newStatus = "paid"
		}
	} else if transactionStatus == "settlement" {
		newStatus = "paid"
	} else if transactionStatus == "deny" {
		newStatus = "cancelled"
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		newStatus = "cancelled"
	} else if transactionStatus == "pending" {
		newStatus = "pending"
	}
	if newStatus != "" {
		fmt.Printf("Updating order %s status to %s\n", orderID, newStatus)
		return s.repo.UpdateFromNotification(orderID, newStatus)
	}

	return nil
}
