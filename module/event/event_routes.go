package event

import (
	"github.com/gofiber/fiber/v3"
	"github.com/strbagus/gof-gallery/internal/middleware"
)

func RegisterRoutes(router fiber.Router) {
	app := router.Group("/events")
	app.Get("/", ListEvent)
	app.Get("/:slug", GetEventDetail)
	app.Post("/", AddEvent)

	token := app.Group("/:slug/tokens", middleware.AdminMiddleware)
	token.Post("/", GenerateToken)
	token.Get("/", ListTokens)
	token.Delete("/:token", RevokeToken)
}
