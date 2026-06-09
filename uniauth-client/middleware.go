package uniauthclient

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// NewHTTPMiddleware creates a standard net/http middleware that authenticates requests using UniAuth.
// It prioritizes access_token and refresh_token cookies to automatically refresh expired access tokens.
func NewHTTPMiddleware(client *Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenStr string

			// 1. Try reading access_token cookie first
			cookie, err := r.Cookie("access_token")
			if err == nil && cookie.Value != "" {
				tokenStr = cookie.Value
			}

			// Parse and validate the access token if found in cookie
			var claims *UniAuthClaims
			if tokenStr != "" {
				claims = &UniAuthClaims{}
				token, err := jwt.ParseWithClaims(tokenStr, claims, client.Keyfunc)
				if err != nil || !token.Valid || claims.TokenType != "access" {
					// Reset tokenStr to trigger refresh check below
					tokenStr = ""
					claims = nil
				}
			}

			// 2. If access_token cookie is missing/expired, check refresh_token cookie
			if tokenStr == "" {
				refreshCookie, err := r.Cookie("refresh_token")
				if err == nil && refreshCookie.Value != "" {
					newAccessToken, expiry, refreshErr := client.RefreshToken(refreshCookie.Value)
					if refreshErr == nil && newAccessToken != "" {
						tokenStr = newAccessToken

						// Write updated access_token cookie back to client response
						http.SetCookie(w, &http.Cookie{
							Name:     "access_token",
							Value:    newAccessToken,
							Path:     "/",
							Expires:  expiry,
							HttpOnly: true,
							Secure:   true,
							SameSite: http.SameSiteLaxMode,
						})
					}
				}
			}

			// 3. Fallback: Parse Authorization Bearer header
			if tokenStr == "" {
				authHeader := r.Header.Get("Authorization")
				if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
					tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			// Final parsing if refreshed or header fallback was resolved
			if claims == nil && tokenStr != "" {
				claims = &UniAuthClaims{}
				token, err := jwt.ParseWithClaims(tokenStr, claims, client.Keyfunc)
				if err != nil || !token.Valid || claims.TokenType != "access" {
					respondJSONError(w, http.StatusUnauthorized, "Invalid or expired token")
					return
				}
			}

			if claims == nil {
				respondJSONError(w, http.StatusUnauthorized, "Missing or expired authentication token")
				return
			}

			// Inject validated claims into standard request context
			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// NewFiberMiddleware creates a Go Fiber v3 middleware that authenticates requests using UniAuth.
// It prioritizes access_token and refresh_token cookies to automatically refresh expired access tokens.
func NewFiberMiddleware(client *Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		var tokenStr string

		// 1. Try reading access_token cookie first
		cookieVal := c.Cookies("access_token")
		if cookieVal != "" {
			tokenStr = cookieVal
		}

		// Parse and validate the access token if found in cookie
		var claims *UniAuthClaims
		if tokenStr != "" {
			claims = &UniAuthClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, client.Keyfunc)
			if err != nil || !token.Valid || claims.TokenType != "access" {
				// Reset tokenStr to trigger refresh check below
				tokenStr = ""
				claims = nil
			}
		}

		// 2. If access_token cookie is missing/expired, check refresh_token cookie
		if tokenStr == "" {
			refreshVal := c.Cookies("refresh_token")
			if refreshVal != "" {
				newAccessToken, expiry, refreshErr := client.RefreshToken(refreshVal)
				if refreshErr == nil && newAccessToken != "" {
					tokenStr = newAccessToken

					// Write updated access_token cookie back to client response
					c.Cookie(&fiber.Cookie{
						Name:     "access_token",
						Value:    newAccessToken,
						Path:     "/",
						Expires:  expiry,
						HTTPOnly: true,
						Secure:   true,
						SameSite: "Lax",
					})
				}
			}
		}

		// 3. Fallback: Parse Authorization Bearer header
		if tokenStr == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// Final parsing if refreshed or header fallback was resolved
		if claims == nil && tokenStr != "" {
			claims = &UniAuthClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, client.Keyfunc)
			if err != nil || !token.Valid || claims.TokenType != "access" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Invalid or expired token",
				})
			}
		}

		if claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing or expired authentication token",
			})
		}

		// Inject validated claims into Fiber locals for easy downstream extraction
		c.Locals(claimsKey, claims)

		// Inject into standard context.Context as well for standard handler compatibility
		ctx := WithClaims(c.Context(), claims)
		c.SetContext(ctx)

		return c.Next()
	}
}

// Helper function to respond with standard JSON error messages
func respondJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}
