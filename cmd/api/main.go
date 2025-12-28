package main

import (
	"log"
	"project-kelas-santai/internal/config"
	"project-kelas-santai/internal/database"
	"project-kelas-santai/internal/middleware"
	"project-kelas-santai/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 2. Connect Database
	dsn := cfg.Database.DSN()
	database.ConnectDB(dsn)

	// 3. Init Fiber App
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: middleware.ErrorHandler,
		//Prefork:      cfg.Web.Prefork,
	})

	// 4. Setup Routes
	routes.SetupRoutes(app, cfg)

	// 5. Start Server
	log.Fatal(app.Listen(cfg.Web.Port))
}
