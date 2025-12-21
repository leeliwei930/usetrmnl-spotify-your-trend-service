package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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
	Host string
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

	// Configure IP extraction for proxy support
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

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

	// Proxy middleware - must come before CORS when running behind a proxy
	// This ensures X-Forwarded-* headers are properly handled

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.PATCH},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-SPOTIFY-CLIENT-ID",
			"X-SPOTIFY-CLIENT-SECRET",
			"X-SPOTIFY-REFRESH-TOKEN",
			// Headers that might be added by proxy
			"X-Forwarded-For",
			"X-Forwarded-Proto",
			"X-Forwarded-Host",
			"X-Real-IP",
		},
		ExposeHeaders: []string{
			echo.HeaderContentType,
			echo.HeaderContentLength,
			echo.HeaderAcceptEncoding,
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           3600, // Cache preflight requests for 1 hour
	}))

	// Security headers middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Content Security Policy
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'")
			// Prevent MIME sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			// Prevent clickjacking
			c.Response().Header().Set("X-Frame-Options", "DENY")
			// Enable XSS protection
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			// Enforce HTTPS (in production, adjust max-age as needed)
			c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			// Control referrer information
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			return next(c)
		}
	})

	apiGroup := e.Group("/api")
	routes.RegisterApiRoutes(apiGroup)

	webGroup := e.Group("")
	routes.RegisterWebRoutes(webGroup)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		port := viper.GetInt("SERVER_PORT")
		host := viper.GetString("SERVER_HOST")
		// Start server
		if len(host) == 0 {
			host = params.Host
		}

		if port == 0 {
			port = params.Port
		}

		address := fmt.Sprintf("%s:%d", host, port)
		e.Logger.Infof("Starting server on %s", address)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
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
