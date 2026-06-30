package main

import (
	"log"
	"os"

	"github.com/strbagus/gof-gallery/docs"
	"github.com/strbagus/gof-gallery/internal/database"
	"github.com/strbagus/gof-gallery/internal/middleware"
	"github.com/strbagus/gof-gallery/internal/s3"
	"github.com/strbagus/gof-gallery/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env. Menggunakan variabel lingkungan sistem.")
	}
	docs.SwaggerInfo.Host = os.Getenv("APP_HOST")
	docs.SwaggerInfo.BasePath = os.Getenv("APP_PATH")
	docs.SwaggerInfo.Schemes = append(docs.SwaggerInfo.Schemes, os.Getenv("APP_HTTP_SCHEMA"))
}

// @title Go Gallery API
// @version 1.0
// @description API for managing events and photo galleries.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_token

func main() {

	database.InitPostgres()

	defer database.ClosePostgres()

	s3.InitS3()

	middleware.InitUniAuth()

	app := fiber.New(fiber.Config{})

	app.Use(middleware.Prometheus)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost", "http://127.0.0.1", "http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
	}))

	router.RegisterRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen("0.0.0.0:" + port)

}
