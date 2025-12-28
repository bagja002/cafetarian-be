package handlers

import (
	"fmt"
	"project-kelas-santai/internal/services"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service services.OrderService
}

func NewOrderHandler(service services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req services.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	orderID, redirectURL, snapToken, err := h.service.CreateOrder(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fmt.Println("Redirect URL:", redirectURL)
	fmt.Println("Snap Token:", snapToken)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Order created successfully",
		"order_id":     orderID,
		"redirect_url": redirectURL,
		"snap_token":   snapToken,
	})
}

func (h *OrderHandler) CallBackNotification(c *fiber.Ctx) error {

	var notificationPayload map[string]interface{}
	if err := c.BodyParser(&notificationPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification payload",
		})
	}

	fmt.Println("Notification payload:", notificationPayload)

	if err := h.service.HandleNotification(notificationPayload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process notification",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Callback notification received successfully",
	})
}
