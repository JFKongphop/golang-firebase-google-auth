package errs

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// business logic error use in customer service

type AppError struct {
	Code int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) error {
	return AppError{	
		Code: fiber.StatusNotFound,
		Message: message,
	}
}

func NewUnexpectedError() error {
	return AppError{	
		Code: fiber.StatusInternalServerError,
		Message: "unexpected error",
	}
}

func NewValidationError(message string) error {
	return AppError{
		Code: http.StatusUnprocessableEntity,
		Message: message,
	}
}

func NewUnAuthorization() error {
	return AppError{
		Code: fiber.StatusUnauthorized,
		Message: "unauthorized error",
	}
}