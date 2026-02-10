package handlers

import (
	"masterypath/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	service services.AnalyticsService
}

func NewAnalyticsHandler(service services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

// GetContributionData godoc
// @Summary Get contribution graph data
// @Description Returns the daily contribution counts for the last 365 days
// @Tags analytics
// @Produce json
// @Success 200 {array} services.ContributionDay
// @Router /analytics/contributions [get]
func (h *AnalyticsHandler) GetContributionData(c *fiber.Ctx) error {
	data, err := h.service.GetContributionData(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch contribution data",
		})
	}

	return c.JSON(data)
}
