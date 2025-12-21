package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leeliwei930/usetrmnl_spotify_service/handlers"
	"github.com/leeliwei930/usetrmnl_spotify_service/middleware"
)

func RegisterApiRoutes(e *echo.Group) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Spotify API endpoints for NuxtJS frontend (public)
	e.GET("/usetrmnl/spotify/authorize", handlers.SpotifyAuthorizeAPI)
	e.POST("/usetrmnl/spotify/callback", handlers.SpotifyCallbackAPI)

	// Protected Spotify API endpoints (require JWT authentication)
	protected := e.Group("")
	protected.Use(middleware.JWTMiddleware())

	// Spotify trends endpoint
	protected.POST("/usetrmnl/spotify/trends", handlers.GetTrendsHandler)

	// Spotify recently played endpoint
	protected.POST("/usetrmnl/spotify/recent-played", handlers.GetRecentPlayedHandler)
}
