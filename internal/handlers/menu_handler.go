package handlers

import (
	"project-kelas-santai/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	service services.MenuService
}

func NewMenuHandler(service services.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) GetAllMenus(c *fiber.Ctx) error {
	menus, err := h.service.GetAllMenus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"data": menus,
	})
}
