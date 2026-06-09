package uniauthclient

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func TestClientAndCookieRefresh(t *testing.T) {
	// Generate RS256 key pair for test signing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Build mock JWKS JSON
	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())

	kid := os.Getenv("JWKS_KID")
	if kid == "" {
		kid = "test-key-id"
	}
	jwksMap := map[string][]map[string]string{
		"keys": {
			{
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"kid": kid,
				"n":   n,
				"e":   e,
			},
		},
	}
	jwksJSON, _ := json.Marshal(jwksMap)

	// Helper to sign tokens for test cases
	signToken := func(claims jwt.Claims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token.Header["kid"] = "test-key-id"
		tokenStr, err := token.SignedString(privateKey)
		if err != nil {
			t.Fatalf("Failed to sign token: %v", err)
		}
		return tokenStr
	}

	// Spin up mock JWKS and Refresh server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/jwks.json" {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jwksJSON)
			return
		}

		if r.URL.Path == "/refresh" && r.Method == "POST" {
			var refreshReq RefreshTokenRequest
			if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Validate incoming refresh token
			claims := &UniAuthClaims{}
			token, err := jwt.ParseWithClaims(refreshReq.RefreshToken, claims, func(token *jwt.Token) (any, error) {
				return publicKey, nil
			})

			if err != nil || !token.Valid || claims.TokenType != "refresh" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Issue new access token
			newAccessClaims := &UniAuthClaims{
				TokenType: "access",
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   claims.Subject,
					Issuer:    "UniAuth",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
				},
			}
			newAccessToken := signToken(newAccessClaims)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(RefreshTokenResponse{
				AccessToken:  newAccessToken,
				AccessExpiry: time.Now().Add(15 * time.Minute).Format(time.RFC3339),
			})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	// Initialize the UniAuth client pointing to mock server with trailing slash
	client, err := NewClient(mockServer.URL + "/")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	defer client.Close()

	// Generate a valid access token and refresh token
	accessClaims := &UniAuthClaims{
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    "UniAuth",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	accessTokenStr := signToken(accessClaims)

	refreshClaims := &UniAuthClaims{
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    "UniAuth",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshTokenStr := signToken(refreshClaims)

	// 1. Verify standard net/http middleware - Success Case (Direct Access Token Cookie)
	middlewareHTTP := NewHTTPMiddleware(client)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claimsExt, ok := GetClaimsFromRequest(r)
		if !ok {
			t.Errorf("GetClaimsFromRequest claims not found")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if claimsExt.Subject != "user-123" {
			t.Errorf("Expected subject 'user-123', got '%s'", claimsExt.Subject)
		}
		w.WriteHeader(http.StatusOK)
	})

	reqCookie := httptest.NewRequest("GET", "/dashboard", nil)
	reqCookie.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: accessTokenStr,
	})
	respCookie := httptest.NewRecorder()
	middlewareHTTP(testHandler).ServeHTTP(respCookie, reqCookie)

	if respCookie.Code != http.StatusOK {
		t.Errorf("Expected access token cookie auth HTTP status 200, got %d", respCookie.Code)
	}

	// 2. Verify standard net/http middleware - Auto Refresh Case (Missing Access Token Cookie, Valid Refresh Cookie)
	reqRefreshCookie := httptest.NewRequest("GET", "/dashboard", nil)
	reqRefreshCookie.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshTokenStr,
	})
	respRefreshCookie := httptest.NewRecorder()
	middlewareHTTP(testHandler).ServeHTTP(respRefreshCookie, reqRefreshCookie)

	if respRefreshCookie.Code != http.StatusOK {
		t.Errorf("Expected auto refresh HTTP status 200, got %d", respRefreshCookie.Code)
	}

	// Verify that the new access token cookie was set in the response
	var newAccessCookie *http.Cookie
	for _, cookie := range respRefreshCookie.Result().Cookies() {
		if cookie.Name == "access_token" {
			newAccessCookie = cookie
		}
	}
	if newAccessCookie == nil || newAccessCookie.Value == "" {
		t.Errorf("Expected a new access_token cookie to be set in the response after refresh")
	}

	// 3. Verify Fiber v3 middleware - Auto Refresh Case
	app := fiber.New()
	middlewareFiber := NewFiberMiddleware(client)
	app.Get("/dashboard", middlewareFiber, func(c fiber.Ctx) error {
		claimsExt, ok := GetClaimsFromFiber(c)
		if !ok {
			t.Errorf("GetClaimsFromFiber claims not found")
			return c.SendStatus(http.StatusInternalServerError)
		}
		if claimsExt.Subject != "user-123" {
			t.Errorf("Expected Fiber subject 'user-123', got '%s'", claimsExt.Subject)
		}
		return c.SendStatus(http.StatusOK)
	})

	reqFiber := httptest.NewRequest("GET", "/dashboard", nil)
	reqFiber.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshTokenStr,
	})
	respFiber, err := app.Test(reqFiber)
	if err != nil {
		t.Fatalf("Failed to test Fiber app: %v", err)
	}

	if respFiber.StatusCode != http.StatusOK {
		t.Errorf("Expected Fiber refresh status 200, got %d", respFiber.StatusCode)
	}

	// Check if access_token cookie is returned in Fiber response
	var newAccessFiberCookie string
	for _, cookieStr := range respFiber.Header["Set-Cookie"] {
		if strings.HasPrefix(cookieStr, "access_token=") {
			newAccessFiberCookie = cookieStr
		}
	}
	if newAccessFiberCookie == "" {
		t.Errorf("Expected access_token cookie set in Fiber response")
	}
}
