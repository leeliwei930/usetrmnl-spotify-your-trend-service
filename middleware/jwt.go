package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/leeliwei930/usetrmnl_spotify_service/config"
)

// JWTClaims represents the custom claims for the service token
type JWTClaims struct {
	Service string `json:"service"`
	jwt.RegisteredClaims
}

// GenerateServiceToken generates a JWT token for the service
func GenerateServiceToken(expirationDays int) (string, error) {
	cfg := config.Get()

	if cfg.JWTSecret == "" {
		return "", fmt.Errorf("JWT_SECRET is not configured")
	}

	// Create claims with service name and expiration
	claims := JWTClaims{
		Service: "usetrmnl_spotify",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "usetrmnl_spotify_service",
			Subject:   "service_auth",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token string and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	cfg := config.Get()

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not configured")
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// JWTMiddleware returns an Echo middleware function that validates JWT tokens
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Check for Bearer token format
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format, expected 'Bearer <token>'")
			}

			tokenString := parts[1]

			// Validate token
			claims, err := ValidateToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("invalid token: %v", err))
			}

			// Store claims in context for potential use in handlers
			c.Set("jwt_claims", claims)

			return next(c)
		}
	}
}
