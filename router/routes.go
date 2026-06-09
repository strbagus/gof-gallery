package router

import (
	"os"

	"github.com/strbagus/gof-gallery/module/event"
	"github.com/strbagus/gof-gallery/module/photo"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterRoutes(app *fiber.App) {

	base := app.Group(os.Getenv("APP_PATH"))

	if os.Getenv("APP_ENV") != "production" {
		base.Get("/swagger/*", swaggo.New(swaggo.Config{
			WithCredentials: true,
		}))
	}

	base.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Service up and running.")
	})
	base.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Register Modules
	event.RegisterRoutes(base)
	photo.RegisterRoutes(base)
}
