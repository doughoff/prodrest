package middleware

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hoffax/prodrest/constants"
)

func FiberCustomErrorHandler(c *fiber.Ctx, err error) error {
	var invalidBodyError *constants.InvalidBodyError
	if errors.As(err, &invalidBodyError) {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
			"code":    "bad_formatted_body",
			"message": err.Error(),
		})
	}

	var requiredFieldErr *constants.RequiredFieldError
	if errors.As(err, &requiredFieldErr) {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
			"code":    "required_field",
			"message": requiredFieldErr.Error(),
		})
	}

	var constraintError *constants.UniqueConstraintError
	if errors.As(err, &constraintError) {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
			"code":    "unique_violation",
			"message": constraintError.Error(),
		})
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorMessages := make([]map[string]string, 0)
		for _, err := range validationErrors {
			errorMessages = append(errorMessages, map[string]string{
				"field":   err.Field(),
				"message": fmt.Sprintf("failed on %v %v %v validation", err.Kind(), err.Tag(), err.Param()),
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{
			"code":    "validation_error",
			"message": "validation errors",
			"errors":  errorMessages,
		})
	}

	var e *fiber.Error
	if errors.As(err, &e) {
		fmt.Printf("fiber error: %v\n", e)
		return c.Status(e.Code).JSON(map[string]string{"message": e.Error()})
	}

	fmt.Printf("Not intercepted error: %v\n", err)
	// default error response
	return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"message": "internal server error"})
}
