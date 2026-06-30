package photo

import (
	"github.com/gofiber/fiber/v3"
	"github.com/strbagus/gof-gallery/internal/middleware"
)

func RegisterRoutes(router fiber.Router) {
	app := router.Group("/photos")
	app.Get("/:slug", ListPhoto)


	protected := app.Group("/", middleware.AdminMiddleware)
	protected.Post("/", AddPhoto)

}
