package router

import (
	"os"

	"github.com/strbagus/gof-gallery/module/event"
	"github.com/strbagus/gof-gallery/module/photo"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App) {

	base := app.Group(os.Getenv("BASE_URL"))
	event.RegisterRoutes(base)
	photo.RegisterRoutes(base)
}
