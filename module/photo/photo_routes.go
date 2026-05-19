package photo

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router) {
	app := router.Group("/photos")
	app.Get("/:slug", ListPhoto)
	app.Post("/", AddPhoto)
}
