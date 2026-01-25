package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// Replace with your actual module name
	"masterypath/internal/models"
)

type SkillHandler struct {
	skillCollection    *mongo.Collection
	metadataCollection *mongo.Collection
}

func NewSkillHandler(skillCol *mongo.Collection, metaCol *mongo.Collection) *SkillHandler {
	return &SkillHandler{
		skillCollection:    skillCol,
		metadataCollection: metaCol,
	}
}

// Helper: Check if category is valid
func (h *SkillHandler) isValidCategory(category string) bool {
	if category == "" {
		return true
	} // Allow empty category if not strict
	filter := bson.M{"type": "SKILL_CATEGORY", "value": category, "is_active": true}
	count, _ := h.metadataCollection.CountDocuments(context.Background(), filter)
	return count > 0
}

// --- CRUD OPERATIONS ---

// CreateSkill handles both Parent and Child skills
// CreateSkill godoc
// @Summary Create a skill
// @Description Create a new skill (parent or child)
// @Tags skills
// @Accept json
// @Produce json
// @Param skill body models.Skill true "Skill object"
// @Success 201 {object} models.Skill
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills [post]
func (h *SkillHandler) CreateSkill(c *fiber.Ctx) error {
	skill := new(models.Skill)
	if err := c.BodyParser(skill); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if !h.isValidCategory(skill.Category) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Category"})
	}

	skill.ID = primitive.NewObjectID()
	skill.CreatedAt = time.Now()

	// Build Ancestors array (Materialized Path)
	if skill.ParentID != nil {
		var parent models.Skill
		err := h.skillCollection.FindOne(context.Background(), bson.M{"_id": *skill.ParentID}).Decode(&parent)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Parent skill not found"})
		}
		// New skill's ancestors = parent's ancestors + parent's ID
		skill.Ancestors = append(parent.Ancestors, parent.ID)
	} else {
		skill.Ancestors = []primitive.ObjectID{} // Root skill
	}

	_, err := h.skillCollection.InsertOne(context.Background(), skill)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create skill"})
	}

	return c.Status(201).JSON(skill)
}

// GetSkills godoc
// @Summary List skills
// @Description Get all skills
// @Tags skills
// @Accept json
// @Produce json
// @Success 200 {array} models.Skill
// @Failure 500 {object} map[string]string
// @Router /skills [get]
func (h *SkillHandler) GetSkills(c *fiber.Ctx) error {
	cursor, err := h.skillCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch skills"})
	}
	var skills []models.Skill
	if err = cursor.All(context.Background(), &skills); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse skills"})
	}

	// Return empty list if nil
	if skills == nil {
		skills = []models.Skill{}
	}
	return c.JSON(skills)
}

// GetSkillTree fetches a skill and all its descendants in one query
// GetSkillTree godoc
// @Summary Get skill tree
// @Description Get a skill and all its descendants using the ancestors array
// @Tags skills
// @Accept json
// @Produce json
// @Param id path string true "Root Skill ID"
// @Success 200 {array} models.Skill
// @Failure 500 {object} map[string]string
// @Router /skills/{id}/tree [get]
func (h *SkillHandler) GetSkillTree(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	// Find the root skill + all skills that have this ID in their ancestors
	filter := bson.M{
		"$or": bson.A{
			bson.M{"_id": objID},
			bson.M{"ancestors": objID},
		},
	}

	cursor, err := h.skillCollection.Find(context.Background(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch skill tree"})
	}

	var skills []models.Skill
	if err = cursor.All(context.Background(), &skills); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse skills"})
	}

	if skills == nil {
		skills = []models.Skill{}
	}
	return c.JSON(skills)
}

// GetSkill fetches a single skill and populates its Parent data
// GetSkill godoc
// @Summary Get a skill
// @Description Get a single skill by ID with parent data populated
// @Tags skills
// @Accept json
// @Produce json
// @Param id path string true "Skill ID"
// @Success 200 {object} models.SkillPopulated
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills/{id} [get]
func (h *SkillHandler) GetSkill(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: objID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "skills"},
			{Key: "localField", Value: "parent_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "parent_data"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$parent_data"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	cursor, err := h.skillCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Query failed"})
	}

	var results []models.SkillPopulated
	if err = cursor.All(context.Background(), &results); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Parsing failed"})
	}

	if len(results) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Skill not found"})
	}

	return c.JSON(results[0])
}

// UpdateSkill godoc
// @Summary Update a skill
// @Description Update an existing skill by ID
// @Tags skills
// @Accept json
// @Produce json
// @Param id path string true "Skill ID"
// @Param skill body models.Skill true "Skill object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills/{id} [put]
func (h *SkillHandler) UpdateSkill(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	updateData := new(models.Skill)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if !h.isValidCategory(updateData.Category) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Category"})
	}

	// Calculate new ancestors if ParentID is changing
	var newAncestors []primitive.ObjectID
	if updateData.ParentID != nil {
		var parent models.Skill
		err := h.skillCollection.FindOne(context.Background(), bson.M{"_id": *updateData.ParentID}).Decode(&parent)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Parent skill not found"})
		}
		newAncestors = append(parent.Ancestors, parent.ID)
	} else {
		newAncestors = []primitive.ObjectID{}
	}

	update := bson.M{
		"$set": bson.M{
			"name":        updateData.Name,
			"category":    updateData.Category,
			"description": updateData.Description,
			"parent_id":   updateData.ParentID,
			"ancestors":   newAncestors,
		},
	}

	_, err := h.skillCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Update failed"})
	}

	// Cascade: Update ancestors for all descendants of this skill
	// Find all skills where objID is in their ancestors array
	descendantFilter := bson.M{"ancestors": objID}
	cursor, err := h.skillCollection.Find(context.Background(), descendantFilter)
	if err == nil {
		var descendants []models.Skill
		cursor.All(context.Background(), &descendants)
		for _, desc := range descendants {
			// Find index of objID in desc.Ancestors
			idx := -1
			for i, a := range desc.Ancestors {
				if a == objID {
					idx = i
					break
				}
			}
			if idx != -1 {
				// Replace everything from index 0 to idx (inclusive) with newAncestors + objID
				updatedAncestors := append(newAncestors, objID)
				updatedAncestors = append(updatedAncestors, desc.Ancestors[idx+1:]...)
				h.skillCollection.UpdateOne(context.Background(), bson.M{"_id": desc.ID}, bson.M{"$set": bson.M{"ancestors": updatedAncestors}})
			}
		}
	}

	return c.JSON(fiber.Map{"message": "Skill updated"})
}

// DeleteSkill godoc
// @Summary Delete a skill
// @Description Delete a skill by ID
// @Tags skills
// @Accept json
// @Produce json
// @Param id path string true "Skill ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills/{id} [delete]
func (h *SkillHandler) DeleteSkill(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	// Pro Tip: Check if this skill is a parent to others before deleting!
	// For now, we just delete.
	_, err := h.skillCollection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Delete failed"})
	}

	return c.JSON(fiber.Map{"message": "Skill deleted"})
}
