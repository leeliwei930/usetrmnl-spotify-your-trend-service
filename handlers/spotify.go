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
		"user-top-read",
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

// CallbackPageData contains data to pass to the callback page template
type CallbackPageData struct {
	Success      bool
	RefreshToken string
	AccessToken  string
	ExpiresIn    int
	ErrorMessage string
}

// SpotifyCallbackPage handles the OAuth callback from Spotify
func SpotifyCallbackPage(c echo.Context) error {
	// Get configuration
	cfg := config.Get()

	// Extract query parameters
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	errorParam := c.QueryParam("error")

	// Check if user denied authorization
	if errorParam != "" {
		return c.Render(http.StatusOK, "pages/spotify/callback.html", CallbackPageData{
			Success:      false,
			ErrorMessage: "Authorization was denied. Please try again.",
		})
	}

	// Validate code parameter
	if code == "" {
		return c.Render(http.StatusBadRequest, "pages/spotify/callback.html", CallbackPageData{
			Success:      false,
			ErrorMessage: "Missing authorization code. Please start the authorization process again.",
		})
	}

	// TODO: In production, validate the state parameter against stored session state
	// For now, we'll just log it
	if state != "" {
		c.Logger().Debug("Received state parameter: ", state)
	}

	// Exchange code for tokens
	tokenResp, err := services.ExchangeCodeForToken(services.SpotifyTokenParams{
		Code:         code,
		RedirectURI:  cfg.SpotifyRedirectURI,
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifyClientSecret,
	})

	if err != nil {
		c.Logger().Error("Token exchange failed: ", err)
		return c.Render(http.StatusInternalServerError, "pages/spotify/callback.html", CallbackPageData{
			Success:      false,
			ErrorMessage: "Failed to exchange authorization code for tokens. Please try again.",
		})
	}

	// Render success page with tokens
	return c.Render(http.StatusOK, "pages/spotify/callback.html", CallbackPageData{
		Success:      true,
		RefreshToken: tokenResp.RefreshToken,
		AccessToken:  tokenResp.AccessToken,
		ExpiresIn:    tokenResp.ExpiresIn,
	})
}

// --- API Handlers (for NuxtJS frontend) ---

// AuthorizeAPIResponse contains the authorization URL for the frontend
type AuthorizeAPIResponse struct {
	AuthURL string `json:"authUrl"`
}

// SpotifyAuthorizeAPI returns the Spotify authorization URL as JSON
func SpotifyAuthorizeAPI(c echo.Context) error {
	// Get configuration from the centralized config package
	cfg := config.Get()

	// Define required scopes
	scopes := []string{
		"user-read-private",
		"user-read-email",
		"user-top-read",
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

	// Return JSON response
	return c.JSON(http.StatusOK, AuthorizeAPIResponse{
		AuthURL: authURL,
	})
}

// CallbackAPIRequest contains the authorization code from the frontend
type CallbackAPIRequest struct {
	Code  string `json:"code"`
	State string `json:"state,omitempty"`
}

// CallbackAPIResponse contains the tokens or error message
type CallbackAPIResponse struct {
	Success      bool   `json:"success"`
	RefreshToken string `json:"refreshToken,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	ExpiresIn    int    `json:"expiresIn,omitempty"`
	Error        string `json:"error,omitempty"`
}

// SpotifyCallbackAPI handles the OAuth callback from the frontend
func SpotifyCallbackAPI(c echo.Context) error {
	// Get configuration
	cfg := config.Get()

	// Parse request body
	var req CallbackAPIRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, CallbackAPIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate code parameter
	if req.Code == "" {
		return c.JSON(http.StatusBadRequest, CallbackAPIResponse{
			Success: false,
			Error:   "Missing authorization code",
		})
	}

	// TODO: In production, validate the state parameter against stored session state
	if req.State != "" {
		c.Logger().Debug("Received state parameter: ", req.State)
	}

	// Exchange code for tokens
	tokenResp, err := services.ExchangeCodeForToken(services.SpotifyTokenParams{
		Code:         req.Code,
		RedirectURI:  cfg.SpotifyRedirectURI,
		ClientID:     cfg.SpotifyClientID,
		ClientSecret: cfg.SpotifyClientSecret,
	})

	if err != nil {
		c.Logger().Error("Token exchange failed: ", err)
		return c.JSON(http.StatusInternalServerError, CallbackAPIResponse{
			Success: false,
			Error:   "Failed to exchange authorization code for tokens",
		})
	}

	// Return success response with tokens
	return c.JSON(http.StatusOK, CallbackAPIResponse{
		Success:      true,
		RefreshToken: tokenResp.RefreshToken,
		AccessToken:  tokenResp.AccessToken,
		ExpiresIn:    tokenResp.ExpiresIn,
	})
}

// Custom response types for Trends API
type TrendsResponse struct {
	Trends TrendsData `json:"trends"`
}

type TrendsData struct {
	Tracks          []TrendTrack    `json:"tracks"`
	BehindTheLyrics BehindTheLyrics `json:"behindTheLyrics"`
}

type BehindTheLyrics struct {
	Summary string      `json:"summary"`
	Title   string      `json:"title"`
	Artist  string      `json:"artist"`
	Album   BehindAlbum `json:"album"`
}

type BehindAlbum struct {
	Name     string `json:"name"`
	CoverUrl string `json:"coverUrl"`
}

type TrendTrack struct {
	Rank   int        `json:"rank"`
	Title  string     `json:"title"`
	Artist string     `json:"artist"`
	Album  TrendAlbum `json:"album"`
}

type TrendAlbum struct {
	Name     string `json:"name"`
	CoverUrl string `json:"coverUrl"`
}

// GetTrendsHandler handles requests to fetch user's top Spotify tracks
func GetTrendsHandler(c echo.Context) error {
	// Extract custom headers
	clientID := c.Request().Header.Get("X-SPOTIFY-CLIENT-ID")
	clientSecret := c.Request().Header.Get("X-SPOTIFY-CLIENT-SECRET")
	refreshToken := c.Request().Header.Get("X-SPOTIFY-REFRESH-TOKEN")

	// Validate required headers
	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Missing required headers: X-SPOTIFY-CLIENT-ID, X-SPOTIFY-CLIENT-SECRET, X-SPOTIFY-REFRESH-TOKEN",
		)
	}

	// Refresh access token
	tokenResp, err := services.RefreshAccessToken(services.RefreshTokenParams{
		RefreshToken: refreshToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		c.Logger().Error("Token refresh failed: ", err)
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"Failed to refresh access token: "+err.Error(),
		)
	}

	// Fetch user's top tracks (default: 20 tracks, medium term)
	topTracks, err := services.GetUserTopTracks(tokenResp.AccessToken, 20, "medium_term")
	if err != nil {
		c.Logger().Error("Failed to fetch top tracks: ", err)
		return echo.NewHTTPError(
			http.StatusBadGateway,
			"Failed to fetch top tracks from Spotify: "+err.Error(),
		)
	}

	// Transform response to custom format
	tracks := make([]TrendTrack, 0, len(topTracks.Items))
	for i, track := range topTracks.Items {
		// Get the first artist name (tracks can have multiple artists)
		artistName := ""
		if len(track.Artists) > 0 {
			artistName = track.Artists[0].Name
		}

		// Get album cover URL (prefer the largest image)
		coverURL := ""
		if len(track.Album.Images) > 0 {
			coverURL = track.Album.Images[0].URL
		}

		tracks = append(tracks, TrendTrack{
			Rank:   i + 1, // Rank starts from 1
			Title:  track.Name,
			Artist: artistName,
			Album: TrendAlbum{
				Name:     track.Album.Name,
				CoverUrl: coverURL,
			},
		})
	}

	// Populate behindTheLyrics with data from the top-ranked track
	var behindLyrics BehindTheLyrics
	if len(tracks) > 0 {
		// Use the first (top-ranked) track
		topTrack := tracks[0]
		behindLyrics = BehindTheLyrics{
			Summary: "", // TODO: Replace with actual summary source (LLM, API, etc.)
			Title:   topTrack.Title,
			Artist:  topTrack.Artist,
			Album: BehindAlbum{
				Name:     topTrack.Album.Name,
				CoverUrl: topTrack.Album.CoverUrl,
			},
		}
	}

	// Create response with structured behindTheLyrics data
	response := TrendsResponse{
		Trends: TrendsData{
			Tracks:          tracks,
			BehindTheLyrics: behindLyrics,
		},
	}

	return c.JSON(http.StatusOK, response)
}
