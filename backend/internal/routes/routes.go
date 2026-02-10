package routes

import (
	"masterypath/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, goalHandler *handlers.GoalHandler, skillHandler *handlers.SkillHandler, progressHandler *handlers.ProgressHandler) {
	api := app.Group("/api")

	api.Get("/goals", goalHandler.GetGoals)
	api.Post("/goals", goalHandler.CreateGoal)
	api.Get("/goals/:id", goalHandler.GetGoal)
	api.Put("/goals/:id", goalHandler.UpdateGoal)
	api.Delete("/goals/:id", goalHandler.DeleteGoal)

	// Skill Routes
	api.Get("/skills/dashboard", skillHandler.GetDashboard)
	api.Get("/skills", skillHandler.GetSkills)
	api.Post("/skills", skillHandler.CreateSkill)
	api.Get("/skills/:id", skillHandler.GetSkill)
	api.Get("/skills/:id/tree", skillHandler.GetSkillTree)
	api.Put("/skills/:id", skillHandler.UpdateSkill)
	api.Delete("/skills/:id", skillHandler.DeleteSkill)

	// Progress Routes
	api.Post("/progress", progressHandler.CreateProgressItem)
	api.Get("/progress/skill/:skillId", progressHandler.GetItemsBySkill)
	api.Put("/progress/:id", progressHandler.UpdateProgressItem)
	api.Delete("/progress/:id", progressHandler.DeleteProgressItem)
}
