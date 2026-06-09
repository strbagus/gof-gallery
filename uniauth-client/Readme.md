# UniAuth Client

Centralized Single Sign-On (SSO) client package and HTTP middlewares for Go services consuming UniAuth authentication.

This package automatically manages JWKS (JSON Web Key Set) caching, signature verification using asymmetric RSA (RS256), background key rotation, and rate-limiting updates to mitigate cache-busting DDOS attacks.

It is designed to seamlessly consume **cookie-based authentication** (e.g. `access_token` and `refresh_token` cookies issued by UniAuth on login) and perform **automatic transparent token refreshing** under the hood.

---

## 📦 Installation

To use this client in your Go application, add the package to your dependencies:

```sh
go get github.com/strbagus/uniauth/uniauth-client
```

---

## 🚀 Usage

### 1. Initialize Client
Initialize the client during your service startup. Pass your UniAuth server's base URL:

```go
package main

import (
	"log"
	
	uniauthclient "github.com/strbagus/uniauth/uniauth-client"
)

func main() {
	// Initializes cache and starts background key-rotator workers
	client, err := uniauthclient.NewClient("https://auth.yourdomain.com")
	if err != nil {
		log.Fatalf("Failed to initialize UniAuth Client: %v", err)
	}
	defer client.Close()
}
```

---

### 2. Standard `net/http` Integration

Use `NewHTTPMiddleware` to protect standard Go endpoints. You can extract claims in downstream handlers using `GetClaimsFromRequest` or `GetClaimsFromContext`:

```go
package main

import (
	"fmt"
	"net/http"

	uniauthclient "github.com/strbagus/uniauth/uniauth-client"
)

func main() {
	client, _ := uniauthclient.NewClient("https://auth.yourdomain.com")
	defer client.Close()

	// 1. Create standard middleware
	authMiddleware := uniauthclient.NewHTTPMiddleware(client)

	// 2. Protected Route
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 3. Extract claims safely from context
		claims, ok := uniauthclient.GetClaimsFromRequest(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		fmt.Fprintf(w, "Hello, User Subject: %s!", claims.Subject)
	})

	// Bind middleware to endpoint
	http.Handle("/dashboard", authMiddleware(protectedHandler))
	http.ListenAndServe(":8080", nil)
}
```

---

### 3. Go Fiber v3 Integration

Use `NewFiberMiddleware` to protect Fiber routes. Extract claims inside handlers using `GetClaimsFromFiber`:

```go
package main

import (
	"github.com/gofiber/fiber/v3"
	uniauthclient "github.com/strbagus/uniauth/uniauth-client"
)

func main() {
	client, _ := uniauthclient.NewClient("https://auth.yourdomain.com")
	defer client.Close()

	app := fiber.New()

	// 1. Create Fiber v3 Middleware
	authMiddleware := uniauthclient.NewFiberMiddleware(client)

	// 2. Define protected route group
	api := app.Group("/api", authMiddleware)

	api.Get("/dashboard", func(c fiber.Ctx) error {
		// 3. Extract claims safely from fiber context locals
		claims, ok := uniauthclient.GetClaimsFromFiber(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Welcome!",
			"sub":     claims.Subject,
		})
	})

	app.Listen(":8080")
}
```

---

## 🔄 Automatic Token Refresh Flow

The HTTP and Fiber middlewares operate with the following cookie resolution priority:

1. **Check Access Token Cookie**: Reads the `access_token` cookie. If it is present and valid, the request proceeds immediately.
2. **Fallback to Refresh Token Cookie**: If the access token is missing or expired, the middleware reads the `refresh_token` cookie.
   - If a refresh token is found, the middleware automatically performs an HTTP call to the UniAuth `/refresh` endpoint.
   - On a successful refresh, a **new access token cookie** is written to the response writer so the client's browser session is automatically updated.
   - The request then validates the new access token and proceeds.
3. **Fallback to Authorization Header**: If no cookie-based auth succeeds, the middleware checks `Authorization: Bearer <token>` as a final fallback (useful for API/testing clients).
4. **Validation Failure**: If all checks fail, the middleware terminates the request with `401 Unauthorized` (JSON response).

---

## 🛡️ Security Features

1. **Anti-DDOS Cache Protection**: Key refreshes on unknown `kid` headers are rate-limited to once every 5 minutes to prevent cache-busting DDOS requests using forged headers.
2. **Access Token Type Lock**: Validates that `"token_type"` is strictly `"access"` to prevent refresh token hijack reuse on resource endpoints.
3. **Collision-free Context**: Uses private unexported struct keys (`contextKey{}`) for context injection, preventing third-party library keys from overlapping.
