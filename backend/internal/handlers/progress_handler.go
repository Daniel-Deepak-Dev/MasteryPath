package handlers

import (
	"masterypath/internal/apperror"
	"masterypath/internal/models"
	"masterypath/internal/services"

	"github.com/gofiber/fiber/v2"
)

// ProgressHandler handles HTTP requests for ProgressItem resources.
type ProgressHandler struct {
	service *services.ProgressService
}

// NewProgressHandler creates a new ProgressHandler with the given ProgressService.
func NewProgressHandler(svc *services.ProgressService) *ProgressHandler {
	return &ProgressHandler{service: svc}
}

// CreateProgressItem godoc
// @Summary Create a progress item
// @Description Create a new progress item for a skill
// @Tags progress
// @Accept json
// @Produce json
// @Param item body models.ProgressItem true "Progress Item Data"
// @Success 201 {object} models.ProgressItem
// @Failure 400 {object} map[string]string
// @Router /progress [post]
func (h *ProgressHandler) CreateProgressItem(c *fiber.Ctx) error {
	item := new(models.ProgressItem)
	if err := c.BodyParser(item); err != nil {
		return SendError(c, apperror.ErrInvalidBody)
	}

	if err := h.service.Create(c.UserContext(), item); err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.Status(201).JSON(item)
}

// GetItemsBySkill godoc
// @Summary Get progress items by skill
// @Description Retrieve all progress items for a skill with calculated percentages
// @Tags progress
// @Accept json
// @Produce json
// @Param skillId path string true "Skill ID"
// @Success 200 {array} models.ProgressItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /progress/skill/{skillId} [get]
func (h *ProgressHandler) GetItemsBySkill(c *fiber.Ctx) error {
	parentID, ok := ParseObjectID(c, "skillId")
	if !ok {
		return nil
	}

	items, err := h.service.GetBySkill(c.UserContext(), parentID)
	if err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(items)
}

// UpdateProgressItem godoc
// @Summary Update a progress item
// @Description Update a progress item by ID
// @Tags progress
// @Accept json
// @Produce json
// @Param id path string true "Progress Item ID"
// @Param item body models.ProgressItem true "Updated Progress Item Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /progress/{id} [put]
func (h *ProgressHandler) UpdateProgressItem(c *fiber.Ctx) error {
	objID, ok := ParseObjectID(c, "id")
	if !ok {
		return nil
	}

	updateData := new(models.ProgressItem)
	if err := c.BodyParser(updateData); err != nil {
		return SendError(c, apperror.ErrInvalidBody)
	}

	if err := h.service.Update(c.UserContext(), objID, updateData); err != nil {
		return SendError(c, err.(*apperror.AppError))
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
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /progress/{id} [delete]
func (h *ProgressHandler) DeleteProgressItem(c *fiber.Ctx) error {
	objID, ok := ParseObjectID(c, "id")
	if !ok {
		return nil
	}

	if err := h.service.Delete(c.UserContext(), objID); err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(fiber.Map{"message": "Item deleted"})
}
