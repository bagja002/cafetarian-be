package routes

import (
	"project-kelas-santai/internal/config"
	"project-kelas-santai/internal/database"
	"project-kelas-santai/internal/handlers"

	"project-kelas-santai/internal/models"
	"project-kelas-santai/internal/repository"
	"project-kelas-santai/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Health Check
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Dependency Injection & Auto Migration
	database.DB.AutoMigrate(
		&models.Category{},
		&models.Menu{},
		&models.Order{},
		&models.OrderItem{},
		&models.Transaction{},
	)

	// Repositories
	menuRepo := repository.NewMenuRepository()
	orderRepo := repository.NewOrderRepository()

	// Services
	menuService := services.NewMenuService(menuRepo)
	orderService := services.NewOrderService(orderRepo, cfg)

	// Handlers
	menuHandler := handlers.NewMenuHandler(menuService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Routes
	v1.Get("/menus", menuHandler.GetAllMenus)
	v1.Post("/orders", orderHandler.CreateOrder)
	v1.Post("/callback-notification", orderHandler.CallBackNotification)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend Cafe Reporting for Duty!")
	})
}
