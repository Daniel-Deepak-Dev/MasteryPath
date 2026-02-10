package handlers

import (
	"masterypath/internal/apperror"
	"masterypath/internal/models"
	"masterypath/internal/services"

	"github.com/gofiber/fiber/v2"
)

// SkillHandler handles HTTP requests for Skill resources.
type SkillHandler struct {
	service *services.SkillService
}

// NewSkillHandler creates a new SkillHandler with the given SkillService.
func NewSkillHandler(svc *services.SkillService) *SkillHandler {
	return &SkillHandler{service: svc}
}

// GetSkills godoc
// @Summary Get all skills
// @Description Retrieve all skills from the database
// @Tags skills
// @Accept json
// @Produce json
// @Success 200 {array} models.Skill
// @Failure 500 {object} map[string]string
// @Router /skills [get]
func (h *SkillHandler) GetSkills(c *fiber.Ctx) error {
	skills, err := h.service.GetAll(c.UserContext())
	if err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(skills)
}

// GetSkillsPopulated godoc
// @Summary Get all skills with parent info
// @Description Retrieve all skills with parent name populated
// @Tags skills
// @Accept json
// @Produce json
// @Success 200 {array} models.SkillPopulated
// @Failure 500 {object} map[string]string
// @Router /skills/populated [get]
func (h *SkillHandler) GetSkillsPopulated(c *fiber.Ctx) error {
	skills, err := h.service.GetAllPopulated(c.UserContext())
	if err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(skills)
}

// GetSkillTree godoc
// @Summary Get skill tree
// @Description Retrieve a skill and all its descendants
// @Tags skills
// @Accept json
// @Produce json
// @Param id path string true "Root Skill ID"
// @Success 200 {array} models.Skill
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills/{id}/tree [get]
func (h *SkillHandler) GetSkillTree(c *fiber.Ctx) error {
	objID, ok := ParseObjectID(c, "id")
	if !ok {
		return nil
	}

	skills, err := h.service.GetTree(c.UserContext(), objID)
	if err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(skills)
}

// GetSkill godoc
// @Summary Get a skill by ID
// @Description Retrieve a single skill by ID with parent info
// @Tags skills
// @Accept json
// @Produce json
// @Param id path string true "Skill ID"
// @Success 200 {object} models.SkillPopulated
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /skills/{id} [get]
func (h *SkillHandler) GetSkill(c *fiber.Ctx) error {
	objID, ok := ParseObjectID(c, "id")
	if !ok {
		return nil
	}

	skill, err := h.service.GetByID(c.UserContext(), objID)
	if err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(skill)
}

// CreateSkill godoc
// @Summary Create a new skill
// @Description Create a new skill with optional parent
// @Tags skills
// @Accept json
// @Produce json
// @Param skill body models.Skill true "Skill Data"
// @Success 201 {object} models.Skill
// @Failure 400 {object} map[string]string
// @Router /skills [post]
func (h *SkillHandler) CreateSkill(c *fiber.Ctx) error {
	skill := new(models.Skill)
	if err := c.BodyParser(skill); err != nil {
		return SendError(c, apperror.ErrInvalidBody)
	}

	if err := h.service.Create(c.UserContext(), skill); err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.Status(201).JSON(skill)
}

// UpdateSkill godoc
// @Summary Update a skill
// @Description Update a skill by ID, cascading ancestor changes to descendants
// @Tags skills
// @Accept json
// @Produce json
// @Param id path string true "Skill ID"
// @Param skill body models.Skill true "Updated Skill Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /skills/{id} [put]
func (h *SkillHandler) UpdateSkill(c *fiber.Ctx) error {
	objID, ok := ParseObjectID(c, "id")
	if !ok {
		return nil
	}

	updateData := new(models.Skill)
	if err := c.BodyParser(updateData); err != nil {
		return SendError(c, apperror.ErrInvalidBody)
	}

	if err := h.service.Update(c.UserContext(), objID, updateData); err != nil {
		return SendError(c, err.(*apperror.AppError))
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
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills/{id} [delete]
func (h *SkillHandler) DeleteSkill(c *fiber.Ctx) error {
	objID, ok := ParseObjectID(c, "id")
	if !ok {
		return nil
	}

	if err := h.service.Delete(c.UserContext(), objID); err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(fiber.Map{"message": "Skill deleted"})
}
