package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// Replace 'backend' with your actual module name from go.mod
	"masterypath/internal/models"
)

// GoalHandler struct holds the database collection
type GoalHandler struct {
	collection *mongo.Collection
}

// NewGoalHandler creates a new instance of GoalHandler
func NewGoalHandler(col *mongo.Collection) *GoalHandler {
	return &GoalHandler{collection: col}
}

// GetGoals godoc
// @Summary List goals
// @Description Get all goals
// @Tags goals
// @Accept json
// @Produce json
// @Success 200 {array} models.Goal
// @Failure 500 {object} map[string]string
// @Router /goals [get]
func (h *GoalHandler) GetGoals(c *fiber.Ctx) error {
	var goals []models.Goal
	cursor, err := h.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &goals); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Return empty list if nil
	if goals == nil {
		goals = []models.Goal{}
	}
	return c.JSON(goals)
}

// CreateGoal godoc
// @Summary Create a goal
// @Description Create a new goal
// @Tags goals
// @Accept json
// @Produce json
// @Param goal body models.Goal true "Goal object"
// @Success 201 {object} models.Goal
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goals [post]
func (h *GoalHandler) CreateGoal(c *fiber.Ctx) error {
	goal := new(models.Goal)
	if err := c.BodyParser(goal); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	goal.ID = primitive.NewObjectID()
	goal.CreatedAt = time.Now()

	_, err := h.collection.InsertOne(context.Background(), goal)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create goal"})
	}

	return c.Status(201).JSON(goal)
}

// UpdateGoal handles updating an existing goal
// UpdateGoal godoc
// @Summary Update a goal
// @Description Update an existing goal by ID
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "Goal ID"
// @Param goal body models.Goal true "Goal object"
// @Success 200 {object} models.Goal
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goals/{id} [put]
func (h *GoalHandler) UpdateGoal(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Use a pointer to the Goal model so we handle optional fields correctly
	updateData := new(models.Goal)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Create the update query
	update := bson.M{
		"$set": bson.M{
			"title":       updateData.Title,
			"description": updateData.Description,
			"completed":   updateData.Completed,
		},
	}

	// Use h.collection here
	result, err := h.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update goal"})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Goal not found"})
	}

	// Return the updated data (or fetch the fresh document if strict accuracy is needed)
	updateData.ID = objectID
	return c.JSON(updateData)
}

// DeleteGoal handles removing a goal
// DeleteGoal godoc
// @Summary Delete a goal
// @Description Delete a goal by ID
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "Goal ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goals/{id} [delete]
func (h *GoalHandler) DeleteGoal(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Use h.collection here
	result, err := h.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete goal"})
	}

	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Goal not found"})
	}

	return c.JSON(fiber.Map{"message": "Goal deleted successfully"})
}

// GetGoal godoc
// @Summary Get a goal
// @Description Get a single goal by ID
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "Goal ID"
// @Success 200 {object} models.Goal
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goals/{id} [get]
func (h *GoalHandler) GetGoal(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var goal models.Goal

	// FindOne returns a SingleResult. We call Decode(&goal) to get the error.
	err = h.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&goal)

	if err != nil {
		// Check specifically if the error is "document not found"
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Goal not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get goal"})
	}

	return c.JSON(goal)
}
