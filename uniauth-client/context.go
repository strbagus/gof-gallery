package uniauthclient

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

// Private key type to prevent context/locals key collisions
type contextKey struct{}

var claimsKey = contextKey{}

// WithClaims returns a new context containing the UniAuth claims.
func WithClaims(ctx context.Context, claims *UniAuthClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// GetClaimsFromContext extracts UniAuthClaims from a generic context.Context.
func GetClaimsFromContext(ctx context.Context) (*UniAuthClaims, bool) {
	claims, ok := ctx.Value(claimsKey).(*UniAuthClaims)
	return claims, ok
}

// GetClaimsFromRequest extracts UniAuthClaims from a standard http.Request context.
func GetClaimsFromRequest(r *http.Request) (*UniAuthClaims, bool) {
	return GetClaimsFromContext(r.Context())
}

// GetClaimsFromFiber extracts UniAuthClaims from a Fiber context's locals.
func GetClaimsFromFiber(c fiber.Ctx) (*UniAuthClaims, bool) {
	claims, ok := c.Locals(claimsKey).(*UniAuthClaims)
	return claims, ok
}
