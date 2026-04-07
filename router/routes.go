package router

import (
	"strbagus/travel-agent/module/destination"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App) {

	destination.RegisterRoutes(app)
}
