package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	uniauthclient "github.com/strbagus/gof-gallery/uniauth-client"
)

var UniAuthClient *uniauthclient.Client

// InitUniAuth initializes the UniAuth client safely.
func InitUniAuth() {
	authService := os.Getenv("AUTH_SERVICE")
	if authService == "" {
		log.Println("Peringatan: AUTH_SERVICE tidak disetel. AdminMiddleware tidak akan dapat memverifikasi token.")
		return
	}

	// Ensure the URL has a protocol scheme so url.Parse doesn't fail.
	if !strings.HasPrefix(authService, "http://") && !strings.HasPrefix(authService, "https://") {
		authService = "http://" + authService
	}

	client, err := uniauthclient.NewClient(authService)
	if err != nil {
		log.Printf("Peringatan: Gagal menginisialisasi UniAuth Client ke %s: %v\n", authService, err)
		return
	}
	UniAuthClient = client
	log.Printf("UniAuth Client berhasil diinisialisasi ke %s\n", authService)
}

// AdminMiddleware acts as the authentication middleware block for admin endpoints.
func AdminMiddleware(c fiber.Ctx) error {
	if UniAuthClient == nil {
		// In development mode, if the auth service is not running, bypass authentication
		// to allow developers to work without being blocked.
		if os.Getenv("APP_ENV") == "development" {
			log.Println("Peringatan: Memperbolehkan akses admin bypass di lingkungan development (UniAuth Client tidak aktif).")
			return c.Next()
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authentication service is unavailable",
		})
	}

	// Delegate verification to the UniAuth fiber middleware.
	handler := uniauthclient.NewFiberMiddleware(UniAuthClient)
	return handler(c)
}
