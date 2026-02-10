package main

//go:generate go run github.com/swaggo/swag/cmd/swag init -g main.go -o ../../docs

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// Import your internal packages
	"masterypath/internal/config"
	"masterypath/internal/database"
	"masterypath/internal/handlers"
	"masterypath/internal/routes"
	"masterypath/internal/services"

	_ "masterypath/docs"

	"github.com/gofiber/swagger"
)

// @title MasteryPath API
// @version 1.0
// @description This is the API for the MasteryPath task manager.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Connect to Database
	client := database.ConnectDB(cfg.MongoURI)
	defer client.Disconnect(context.Background())

	db := client.Database("masterypath")

	// 3. Create Services (business logic layer)
	goalService := services.NewGoalService(db.Collection("goals"))
	skillService := services.NewSkillService(db.Collection("skills"), db.Collection("metadata"))
	progressService := services.NewProgressService(db.Collection("progress"))

	// 4. Create Handlers (thin HTTP layer)
	goalHandler := handlers.NewGoalHandler(goalService)
	skillHandler := handlers.NewSkillHandler(skillService)
	progressHandler := handlers.NewProgressHandler(progressService)

	// 5. Setup Fiber App
	app := fiber.New()
	app.Use(cors.New())

	// 6. Setup Routes
	routes.SetupRoutes(app, goalHandler, skillHandler, progressHandler)

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
			log.Fatalf("Server forced shutdown: %v", err)
		}
	}()

	// 7. Start Server
	log.Println("Server running on port", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
