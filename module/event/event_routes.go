package event

import (
	"github.com/gofiber/fiber/v3"
	"github.com/strbagus/gof-gallery/internal/middleware"
)

func RegisterRoutes(router fiber.Router) {
	app := router.Group("/events")
	app.Get("/", ListEvent)
	app.Get("/:slug", GetEventDetail)

	protected := app.Group("/", middleware.AdminMiddleware)
	protected.Post("/", AddEvent)
	protected.Put("/:slug", UpdateEvent)
	protected.Delete("/:slug", DeleteEvent)

	token := protected.Group("/:slug/tokens")
	token.Post("/", GenerateToken)
	token.Get("/", ListTokens)
	token.Delete("/:token", RevokeToken)
}
