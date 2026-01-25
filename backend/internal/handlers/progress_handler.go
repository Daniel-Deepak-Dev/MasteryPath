package handlers

import (
	"context"
	"masterypath/internal/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProgressHandler struct {
	collection *mongo.Collection
}

func NewProgressHandler(col *mongo.Collection) *ProgressHandler {
	return &ProgressHandler{collection: col}
}

// CreateProgressItem
// CreateProgressItem godoc
// @Summary Create a progress item
// @Description Create a new progress item for a skill
// @Tags progress
// @Accept json
// @Produce json
// @Param item body models.ProgressItem true "ProgressItem object"
// @Success 201 {object} models.ProgressItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /progress [post]
func (h *ProgressHandler) CreateProgressItem(c *fiber.Ctx) error {
	item := new(models.ProgressItem)
	if err := c.BodyParser(item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if err := item.ValidateWeightage(); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	item.ID = primitive.NewObjectID()
	_, err := h.collection.InsertOne(context.Background(), item)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save"})
	}

	return c.Status(201).JSON(item)
}

// GetItemsBySkill returns all progress items for a specific skill
// AND calculates the percentages on the fly.
// GetItemsBySkill godoc
// @Summary List progress items by skill
// @Description Get all progress items for a specific skill and calculate percentages
// @Tags progress
// @Accept json
// @Produce json
// @Param skillId path string true "Skill ID"
// @Success 200 {array} models.ProgressItem
// @Failure 500 {object} map[string]string
// @Router /progress/skill/{skillId} [get]
func (h *ProgressHandler) GetItemsBySkill(c *fiber.Ctx) error {
	skillIDHex := c.Params("skillId")
	parentID, _ := primitive.ObjectIDFromHex(skillIDHex)

	pipeline := mongo.Pipeline{
		// 1. Match items for this skill
		{{Key: "$match", Value: bson.D{{Key: "parent_skill_id", Value: parentID}}}},

		// 2. Calculate Total Weight
		{{Key: "$setWindowFields", Value: bson.D{
			{Key: "partitionBy", Value: nil},
			{Key: "output", Value: bson.D{
				{Key: "total_weight", Value: bson.D{{Key: "$sum", Value: "$weightage"}}},
			}},
		}}},

		// 3. Calculate Percentage
		{{Key: "$addFields", Value: bson.D{
			{Key: "weight_percent", Value: bson.D{
				{Key: "$multiply", Value: bson.A{
					bson.D{{Key: "$divide", Value: bson.A{"$weightage", "$total_weight"}}},
					100,
				}},
			}},
		}}},
	}

	cursor, err := h.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Calculation failed"})
	}

	var results []models.ProgressItem
	if err = cursor.All(context.Background(), &results); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Parsing failed"})
	}

	if results == nil {
		results = []models.ProgressItem{}
	}
	return c.JSON(results)
}

// UpdateProgressItem godoc
// @Summary Update a progress item
// @Description Update an existing progress item by ID
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Progress Item ID"
// @Param item body models.ProgressItem true "ProgressItem object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /progress/{id} [put]
func (h *ProgressHandler) UpdateProgressItem(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	updateData := new(models.ProgressItem)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if err := updateData.ValidateWeightage(); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	update := bson.M{
		"$set": bson.M{
			"name":      updateData.Name,
			"weightage": updateData.Weightage,
			"achieved":  updateData.Achieved,
			"comments":  updateData.Comments,
		},
	}

	_, err := h.collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Update failed"})
	}
	return c.JSON(fiber.Map{"message": "Item updated"})
}

// DeleteProgressItem godoc
// @Summary Delete a progress item
// @Description Delete a progress item by ID
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Progress Item ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /progress/{id} [delete]
func (h *ProgressHandler) DeleteProgressItem(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	_, err := h.collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Delete failed"})
	}
	return c.JSON(fiber.Map{"message": "Item deleted"})
}
