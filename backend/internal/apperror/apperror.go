package apperror

import "fmt"

// AppError represents a structured application error with an HTTP status code.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// New creates a new AppError with a custom message.
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap creates a new AppError from a base error with additional context.
func Wrap(base *AppError, detail string) *AppError {
	return &AppError{
		Code:    base.Code,
		Message: fmt.Sprintf("%s: %s", base.Message, detail),
	}
}

// Common application errors.
var (
	ErrNotFound         = &AppError{Code: 404, Message: "Resource not found"}
	ErrInvalidBody      = &AppError{Code: 400, Message: "Invalid request body"}
	ErrInvalidID        = &AppError{Code: 400, Message: "Invalid ID format"}
	ErrValidationFailed = &AppError{Code: 400, Message: "Validation failed"}
	ErrInternal         = &AppError{Code: 500, Message: "Internal server error"}
	ErrDuplicateKey     = &AppError{Code: 409, Message: "Resource already exists"}
)
