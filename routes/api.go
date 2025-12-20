package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leeliwei930/usetrmnl_spotify_service/handlers"
)

func RegisterApiRoutes(e *echo.Group) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Spotify API endpoints for NuxtJS frontend
	e.GET("/usetrmnl/spotify/authorize", handlers.SpotifyAuthorizeAPI)
	e.POST("/usetrmnl/spotify/callback", handlers.SpotifyCallbackAPI)

	// Spotify trends endpoint
	e.POST("/usetrmnl/spotify/trends", handlers.GetTrendsHandler)

	// Spotify recently played endpoint
	e.POST("/usetrmnl/spotify/recent-played", handlers.GetRecentPlayedHandler)
}
