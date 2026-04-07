package main

import (
	"log"
	"os"
	"strbagus/travel-agent/internal/database"
	"strbagus/travel-agent/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Peringatan: Tidak dapat memuat file .env. Menggunakan variabel lingkungan sistem.")
	}

	database.InitPostgres()

	defer database.ClosePostgres()

	app := fiber.New(fiber.Config{})

	app.Use(cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost", "http://127.0.0.1"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
	}))

	router.RegisterRoutes(app)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Service up and running.")
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}
