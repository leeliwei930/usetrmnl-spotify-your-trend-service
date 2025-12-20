package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// TokenResponse represents the response from Spotify's token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// SpotifyTokenParams contains parameters needed for token exchange
type SpotifyTokenParams struct {
	Code         string
	RedirectURI  string
	ClientID     string
	ClientSecret string
}

// ExchangeCodeForToken exchanges an authorization code for access and refresh tokens
func ExchangeCodeForToken(params SpotifyTokenParams) (*TokenResponse, error) {
	tokenURL := "https://accounts.spotify.com/api/token"

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", params.Code)
	data.Set("redirect_uri", params.RedirectURI)

	// Create request
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add Basic Auth (client_id:client_secret encoded in base64)
	auth := base64.StdEncoding.EncodeToString([]byte(params.ClientID + ":" + params.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// RefreshTokenParams contains parameters needed to refresh an access token
type RefreshTokenParams struct {
	RefreshToken string
	ClientID     string
	ClientSecret string
}

// RefreshAccessToken exchanges a refresh token for a new access token
func RefreshAccessToken(params RefreshTokenParams) (*TokenResponse, error) {
	tokenURL := "https://accounts.spotify.com/api/token"

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", params.RefreshToken)

	// Create request
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add Basic Auth (client_id:client_secret encoded in base64)
	auth := base64.StdEncoding.EncodeToString([]byte(params.ClientID + ":" + params.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// SpotifyImage represents an image from Spotify API
type SpotifyImage struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// SpotifyArtist represents an artist from Spotify API
type SpotifyArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URI  string `json:"uri"`
}

// SpotifyAlbum represents an album from Spotify API
type SpotifyAlbum struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Images []SpotifyImage `json:"images"`
}

// SpotifyTrack represents a track from Spotify API
type SpotifyTrack struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Artists []SpotifyArtist `json:"artists"`
	Album   SpotifyAlbum    `json:"album"`
	URI     string          `json:"uri"`
}

// SpotifyTopTracksResponse represents the response from Spotify's top tracks endpoint
type SpotifyTopTracksResponse struct {
	Items []SpotifyTrack `json:"items"`
	Total int            `json:"total"`
	Limit int            `json:"limit"`
}

// GetUserTopTracks fetches the user's top tracks from Spotify
func GetUserTopTracks(accessToken string, limit int, timeRange string) (*SpotifyTopTracksResponse, error) {
	baseURL := "https://api.spotify.com/v1/me/top/tracks"

	// Set default values
	if limit <= 0 {
		limit = 20
	}
	if timeRange == "" {
		timeRange = "medium_term" // Options: short_term, medium_term, long_term
	}

	// Build URL with query parameters
	queryParams := url.Values{}
	queryParams.Set("limit", fmt.Sprintf("%d", limit))
	queryParams.Set("time_range", timeRange)
	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	// Create request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var topTracksResp SpotifyTopTracksResponse
	if err := json.Unmarshal(body, &topTracksResp); err != nil {
		return nil, fmt.Errorf("failed to parse top tracks response: %w", err)
	}

	return &topTracksResp, nil
}

// SpotifyPlayHistory represents a play history item from Spotify API
type SpotifyPlayHistory struct {
	Track    SpotifyTrack `json:"track"`
	PlayedAt string       `json:"played_at"`
}

// SpotifyRecentlyPlayedResponse represents the response from Spotify's recently played endpoint
type SpotifyRecentlyPlayedResponse struct {
	Items []SpotifyPlayHistory `json:"items"`
	Limit int                  `json:"limit"`
}

// GetRecentlyPlayed fetches the user's recently played tracks from Spotify
func GetRecentlyPlayed(accessToken string, limit int) (*SpotifyRecentlyPlayedResponse, error) {
	baseURL := "https://api.spotify.com/v1/me/player/recently-played"

	// Set default values
	if limit <= 0 {
		limit = 20
	}

	// Build URL with query parameters
	queryParams := url.Values{}
	queryParams.Set("limit", fmt.Sprintf("%d", limit))
	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	// Create request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var recentlyPlayedResp SpotifyRecentlyPlayedResponse
	if err := json.Unmarshal(body, &recentlyPlayedResp); err != nil {
		return nil, fmt.Errorf("failed to parse recently played response: %w", err)
	}

	return &recentlyPlayedResp, nil
}
