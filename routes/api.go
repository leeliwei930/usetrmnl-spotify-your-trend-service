package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterApiRoutes(e *echo.Group) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
}
