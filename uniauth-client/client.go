package uniauthclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Client encapsulates the JWKS memory-caching controller for UniAuth signature verification.
type Client struct {
	authBaseURL string
	jwks        *keyfunc.JWKS
	cancel      context.CancelFunc
}

// NewClient initializes the JWKS memory-caching controller.
// It will query the /.well-known/jwks.json endpoint of the provided authBaseURL.
func NewClient(authBaseURL string) (*Client, error) {
	// Parse base URL
	u, err := url.Parse(authBaseURL)
	if err != nil {
		return nil, err
	}

	// Build the JWKS URL path cleanly
	u.Path = strings.TrimSuffix(u.Path, "/") + "/.well-known/jwks.json"
	jwksURL := u.String()

	// Create a cancelable context to prevent memory leaks on server teardown
	ctx, cancel := context.WithCancel(context.Background())

	options := keyfunc.Options{
		Ctx:               ctx,
		RefreshInterval:   12 * time.Hour,
		RefreshRateLimit:  5 * time.Minute,
		RefreshTimeout:    5 * time.Second,
		RefreshUnknownKID: true, // Auto-refreshes key when unknown kid is encountered (with rate limits)
	}

	// Fetch and start background refresh controller
	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		cancel()
		return nil, err
	}

	return &Client{
		authBaseURL: authBaseURL,
		jwks:        jwks,
		cancel:      cancel,
	}, nil
}

// JWKS returns the underlying keyfunc.JWKS instance.
func (c *Client) JWKS() *keyfunc.JWKS {
	return c.jwks
}

// Keyfunc wraps keyfunc.JWKS.Keyfunc and falls back to using the first available key if kid is missing.
func (c *Client) Keyfunc(token *jwt.Token) (interface{}, error) {
	key, err := c.jwks.Keyfunc(token)
	if err == nil {
		return key, nil
	}

	// Fallback for missing kid
	kid, _ := token.Header["kid"].(string)
	if kid == "" {
		kids := c.jwks.KIDs()
		if len(kids) > 0 {
			keys := c.jwks.ReadOnlyKeys()
			if firstKey, ok := keys[kids[0]]; ok {
				return firstKey, nil
			}
		}
	}

	return nil, err
}

// Close stops the background refresh goroutine and cleans up resources.
func (c *Client) Close() {
	if c.cancel != nil {
		c.cancel()
	}
}

// RefreshTokenRequest represents the request body payload for token refresh calls.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse represents the success payload returned by the refresh endpoint.
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	AccessExpiry string `json:"access_expiry"`
}

// RefreshToken invokes the UniAuth server's `/refresh` endpoint to exchange a long-lived refresh token
// for a new short-lived access token and returns the new token value and its expiry time.
func (c *Client) RefreshToken(refreshToken string) (string, time.Time, error) {
	u, err := url.Parse(c.authBaseURL)
	if err != nil {
		return "", time.Time{}, err
	}

	u.Path = strings.TrimSuffix(u.Path, "/") + "/refresh"
	refreshURL := u.String()

	payload := RefreshTokenRequest{RefreshToken: refreshToken}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", time.Time{}, err
	}

	req, err := http.NewRequest("POST", refreshURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return "", time.Time{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 5 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("token refresh failed with status: %d", resp.StatusCode)
	}

	var res RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", time.Time{}, err
	}

	expiry, err := time.Parse(time.RFC3339, res.AccessExpiry)
	if err != nil {
		// Fallback default duration (15 minutes) if parsing fails
		expiry = time.Now().Add(15 * time.Minute)
	}

	return res.AccessToken, expiry, nil
}
