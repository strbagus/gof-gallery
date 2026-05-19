package main

import (
	"log"
	"os"

	"github.com/strbagus/gof-gallery/internal/database"
	"github.com/strbagus/gof-gallery/internal/middleware"
	"github.com/strbagus/gof-gallery/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Peringatan: Tidak dapat memuat file .env. Menggunakan variabel lingkungan sistem.")
	}

	database.InitPostgres()

	defer database.ClosePostgres()

	app := fiber.New(fiber.Config{})

	app.Use(middleware.Prometheus)

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost", "http://127.0.0.1", "http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
	}))

	router.RegisterRoutes(app)

	base := app.Group(os.Getenv("BASE_URL"))


	base.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Service up and running.")
	})
	base.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen("0.0.0.0:" + port)

}
