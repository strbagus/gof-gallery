package uniauthclient

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UniAuthClaims represents the custom JWT claims returned by the UniAuth SSO service.
type UniAuthClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// Validate validates the custom and registered claims of the UniAuth JWT token.
// Compatible with golang-jwt/jwt/v5 validation specs.
func (c *UniAuthClaims) Validate() error {

	// Verify expiration and activation time bounds
	now := time.Now()
	if c.ExpiresAt != nil && c.ExpiresAt.Time.Before(now) {
		return errors.New("token has expired")
	}
	if c.NotBefore != nil && c.NotBefore.Time.After(now) {
		return errors.New("token is not valid yet")
	}

	return nil
}
