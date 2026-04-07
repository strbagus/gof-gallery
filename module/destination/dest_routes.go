package destination

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router) {
	app := router.Group("/destinations")
	app.Get("/", ListDestination)
}
