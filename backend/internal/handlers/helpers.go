package handlers

import (
	"masterypath/internal/apperror"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParseObjectID extracts and validates an ObjectID from the URL params.
// Returns the parsed ObjectID or sends a 400 error response and returns false.
func ParseObjectID(c *fiber.Ctx, param string) (primitive.ObjectID, bool) {
	id := c.Params(param)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		SendError(c, apperror.ErrInvalidID)
		return primitive.NilObjectID, false
	}
	return objID, true
}

// SendError sends a standardized JSON error response from an AppError.
func SendError(c *fiber.Ctx, err *apperror.AppError) error {
	return c.Status(err.Code).JSON(fiber.Map{"error": err.Message})
}

// SendInternalError sends a 500 response, logging the original error context.
func SendInternalError(c *fiber.Ctx, detail string) error {
	return c.Status(500).JSON(fiber.Map{"error": detail})
}
