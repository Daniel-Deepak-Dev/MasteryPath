package handlers

import (
	"masterypath/internal/apperror"
	"masterypath/internal/models"
	"masterypath/internal/services"

	"github.com/gofiber/fiber/v2"
)

// GoalHandler handles HTTP requests for Goal resources.
type GoalHandler struct {
	service *services.GoalService
}

// NewGoalHandler creates a new GoalHandler with the given GoalService.
func NewGoalHandler(svc *services.GoalService) *GoalHandler {
	return &GoalHandler{service: svc}
}

// GetGoals godoc
// @Summary Get all goals
// @Description Retrieve all goals from the database
// @Tags goals
// @Accept json
// @Produce json
// @Success 200 {array} models.Goal
// @Failure 500 {object} map[string]string
// @Router /goals [get]
func (h *GoalHandler) GetGoals(c *fiber.Ctx) error {
	goals, err := h.service.GetAll(c.UserContext())
	if err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(goals)
}

// GetGoal godoc
// @Summary Get a goal by ID
// @Description Retrieve a single goal by its ID
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "Goal ID"
// @Success 200 {object} models.Goal
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /goals/{id} [get]
func (h *GoalHandler) GetGoal(c *fiber.Ctx) error {
	objID, ok := ParseObjectID(c, "id")
	if !ok {
		return nil
	}

	goal, err := h.service.GetByID(c.UserContext(), objID)
	if err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(goal)
}

// CreateGoal godoc
// @Summary Create a new goal
// @Description Create a new goal in the database
// @Tags goals
// @Accept json
// @Produce json
// @Param goal body models.Goal true "Goal Data"
// @Success 201 {object} models.Goal
// @Failure 400 {object} map[string]string
// @Router /goals [post]
func (h *GoalHandler) CreateGoal(c *fiber.Ctx) error {
	goal := new(models.Goal)
	if err := c.BodyParser(goal); err != nil {
		return SendError(c, apperror.ErrInvalidBody)
	}

	if err := h.service.Create(c.UserContext(), goal); err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.Status(201).JSON(goal)
}

// UpdateGoal godoc
// @Summary Update a goal
// @Description Update a goal by ID
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "Goal ID"
// @Param goal body models.Goal true "Updated Goal Data"
// @Success 200 {object} models.Goal
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /goals/{id} [put]
func (h *GoalHandler) UpdateGoal(c *fiber.Ctx) error {
	objID, ok := ParseObjectID(c, "id")
	if !ok {
		return nil
	}

	goal := new(models.Goal)
	if err := c.BodyParser(goal); err != nil {
		return SendError(c, apperror.ErrInvalidBody)
	}

	if err := h.service.Update(c.UserContext(), objID, goal); err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(goal)
}

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
// @Router /goals/{id} [delete]
func (h *GoalHandler) DeleteGoal(c *fiber.Ctx) error {
	objID, ok := ParseObjectID(c, "id")
	if !ok {
		return nil
	}

	if err := h.service.Delete(c.UserContext(), objID); err != nil {
		return SendError(c, err.(*apperror.AppError))
	}
	return c.JSON(fiber.Map{"message": "Goal deleted successfully"})
}
