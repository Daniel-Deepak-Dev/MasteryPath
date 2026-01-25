package main

//go:generate go run github.com/swaggo/swag/cmd/swag init -g main.go -o ../../docs

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// Import your internal packages
	"masterypath/internal/config"
	"masterypath/internal/database"
	"masterypath/internal/handlers"
	"masterypath/internal/routes"

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
	// Disconnect when main exits
	defer client.Disconnect(context.Background())

	// 3. Initialize Handler
	// We inject the specific collection into the handler
	coll := client.Database("masterypath").Collection("goals")
	goalHandler := handlers.NewGoalHandler(coll)

	// Skill Handler
	skillCol := client.Database("masterypath").Collection("skills")
	metadataCol := client.Database("masterypath").Collection("metadata")
	skillHandler := handlers.NewSkillHandler(skillCol, metadataCol)

	// Progress Handler
	progressCol := client.Database("masterypath").Collection("progress")
	progressHandler := handlers.NewProgressHandler(progressCol)

	// 4. Setup Fiber App
	app := fiber.New()
	app.Use(cors.New())

	// 5. Setup Routes
	routes.SetupRoutes(app, goalHandler, skillHandler, progressHandler)

	// Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// 6. Start Server
	log.Println("Server running on port", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
