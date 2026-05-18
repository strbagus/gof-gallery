package event

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router) {
	app := router.Group("/events")
	app.Get("/", ListEvent)
}
