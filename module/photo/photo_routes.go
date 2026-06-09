package photo

import (
	"github.com/gofiber/fiber/v3"
	"github.com/strbagus/gof-gallery/internal/middleware"
)

func RegisterRoutes(router fiber.Router) {
	app := router.Group("/photos")
	app.Post("/", AddPhoto)

	protected := app.Group("/", middleware.AdminMiddleware)
	protected.Get("/:slug", ListPhoto)
}
