package event

import (
	"github.com/gofiber/fiber/v3"
	"github.com/strbagus/gof-gallery/internal/middleware"
)

func RegisterRoutes(router fiber.Router) {
	app := router.Group("/events")
	app.Get("/", ListEvent)

	protected := app.Group("/", middleware.AdminMiddleware)
	protected.Get("/:slug", GetEventDetail)
	protected.Post("/", AddEvent)
	protected.Put("/:slug", UpdateEvent)
	protected.Delete("/:slug", DeleteEvent)

	token := protected.Group("/:slug/tokens", middleware.AdminMiddleware)
	token.Post("/", GenerateToken)
	token.Get("/", ListTokens)
	token.Delete("/:token", RevokeToken)
}
