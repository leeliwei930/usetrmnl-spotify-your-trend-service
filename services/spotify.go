package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// SpotifyAuthParams contains parameters needed to construct a Spotify authorization URL
type SpotifyAuthParams struct {
	ClientID    string
	RedirectURI string
	Scopes      []string
	State       string
}

// CreateSpotifyAuthUrl constructs and returns a Spotify authorization URL
// according to the OAuth 2.0 specification
func CreateSpotifyAuthUrl(params SpotifyAuthParams) string {
	baseURL := "https://accounts.spotify.com/authorize"

	// Prepare query parameters
	queryParams := url.Values{}
	queryParams.Set("client_id", params.ClientID)
	queryParams.Set("response_type", "code")
	queryParams.Set("redirect_uri", params.RedirectURI)

	// Join scopes with space separator
	if len(params.Scopes) > 0 {
		queryParams.Set("scope", strings.Join(params.Scopes, " "))
	}

	// Add state for CSRF protection
	if params.State != "" {
		queryParams.Set("state", params.State)
	}

	// Construct the full URL
	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
}

// GenerateRandomState generates a cryptographically secure random state string
// for CSRF protection in the OAuth flow
func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
