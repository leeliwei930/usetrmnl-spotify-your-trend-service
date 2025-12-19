package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leeliwei930/usetrmnl_spotify_service/handlers"
)

func RegisterWebRoutes(e *echo.Group) {
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "pages/index.html", nil)
	})

	// Spotify authorization route
	e.GET("/usetrmnl/spotify/authorize", handlers.SpotifyAuthorizePage)

}
