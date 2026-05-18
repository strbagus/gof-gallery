package router

import (
	"github.com/strbagus/gof-gallery/module/event"
	"github.com/strbagus/gof-gallery/module/photo"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App) {
	event.RegisterRoutes(app)
	photo.RegisterRoutes(app)
}
