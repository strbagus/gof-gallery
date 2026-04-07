package response

import (
	"github.com/gofiber/fiber/v3"
)

func Success(c fiber.Ctx, message string, data any, meta Metadata) error {
	return c.Status(fiber.StatusOK).JSON(JSONResponse{
		Status:   "success",
		Message:  message,
		Data:     data,
		Metadata: meta,
	})
}

func Error(c fiber.Ctx, code int, message string, errors any) error {
	return c.Status(code).JSON(JSONResponse{
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}
