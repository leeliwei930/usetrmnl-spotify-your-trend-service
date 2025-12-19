package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leeliwei930/usetrmnl_spotify_service/config"
	"github.com/leeliwei930/usetrmnl_spotify_service/services"
)

// AuthorizePageData contains data to pass to the authorization page template
type AuthorizePageData struct {
	AuthURL string
}

// SpotifyAuthorizePage handles the GET request to display the Spotify authorization page
func SpotifyAuthorizePage(c echo.Context) error {
	// Get configuration from the centralized config package
	cfg := config.Get()

	// Define required scopes
	scopes := []string{
		"user-read-private",
		"user-read-email",
	}

	// Generate random state for CSRF protection
	state, err := services.GenerateRandomState()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate state")
	}

	// Create authorization URL using config values
	authURL := services.CreateSpotifyAuthUrl(services.SpotifyAuthParams{
		ClientID:    cfg.SpotifyClientID,
		RedirectURI: cfg.SpotifyRedirectURI,
		Scopes:      scopes,
		State:       state,
	})

	// Render the authorization page
	data := AuthorizePageData{
		AuthURL: authURL,
	}

	return c.Render(http.StatusOK, "pages/spotify/authorize.html", data)
}
