package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/leeliwei930/usetrmnl_spotify_service/config"
	"github.com/leeliwei930/usetrmnl_spotify_service/routes"
	"github.com/spf13/viper"
)

type StartParams struct {
	Port int
}

func Start(params StartParams) {
	// Configure server timeouts
	cfg := config.Get()
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Gzip())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(cfg.ServerTimeout) * time.Second,
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			e.Logger.Infof("%s %s [%d] %v", v.Method, v.URI, v.Status, v.Latency)
			return nil
		},
	}))
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Recover())

	// CORS middleware for NuxtJS frontend
	// Read allowed origins from environment variable (comma-separated)
	allowOriginsStr := viper.GetString("CORS_ALLOW_ORIGINS")
	if allowOriginsStr == "" {
		allowOriginsStr = "http://localhost:3000,http://localhost:3001"
	}
	allowOrigins := strings.Split(allowOriginsStr, ",")
	// Trim whitespace from each origin
	for i := range allowOrigins {
		allowOrigins[i] = strings.TrimSpace(allowOrigins[i])
	}
	e.Static("/", "./public")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-SPOTIFY-CLIENT-ID",
			"X-SPOTIFY-CLIENT-SECRET",
			"X-SPOTIFY-REFRESH-TOKEN",
		},
		AllowCredentials: true,
	}))

	apiGroup := e.Group("/api")
	routes.RegisterApiRoutes(apiGroup)

	webGroup := e.Group("")
	routes.RegisterWebRoutes(webGroup)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	e.Logger.Infof("Starting server on port %d", params.Port)
	go func() {
		if err := e.Start(":" + strconv.Itoa(params.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}
