package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	// Spotify API Configuration
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string

	// Server Configuration
	ServerPort int
	ServerHost string

	// Agent Service Configuration
	AgentBaseURL string

	// Timeout Configuration (in seconds)
	AgentClientTimeout  int
	AgentRequestTimeout int
	ServerTimeout       int
}

var instance *Config

// Init initializes the configuration by loading .env file and setting up Viper
func Init() error {
	// Load .env file - don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or unable to load: %v", err)
		log.Println("Continuing with environment variables and defaults...")
	}

	// Initialize Viper
	viper.AutomaticEnv() // Read environment variables

	// Set default values
	viper.SetDefault("SPOTIFY_REDIRECT_URI", "http://localhost:8080/usetrmnl/spotify/callback")
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("AGENT_BASE_URL", "http://localhost:8000")
	// Timeout defaults
	viper.SetDefault("AGENT_CLIENT_TIMEOUT", 30)
	viper.SetDefault("AGENT_REQUEST_TIMEOUT", 25)
	viper.SetDefault("SERVER_READ_TIMEOUT", 30)
	viper.SetDefault("SERVER_WRITE_TIMEOUT", 30)

	// Load configuration into struct
	instance = &Config{
		SpotifyClientID:     viper.GetString("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: viper.GetString("SPOTIFY_CLIENT_SECRET"),
		SpotifyRedirectURI:  viper.GetString("SPOTIFY_REDIRECT_URI"),
		ServerPort:          viper.GetInt("SERVER_PORT"),
		ServerHost:          viper.GetString("SERVER_HOST"),
		AgentBaseURL:        viper.GetString("AGENT_BASE_URL"),
		// Timeout configuration
		AgentClientTimeout:  viper.GetInt("AGENT_CLIENT_TIMEOUT"),
		AgentRequestTimeout: viper.GetInt("AGENT_REQUEST_TIMEOUT"),
		ServerTimeout:       viper.GetInt("SERVER_TIMEOUT"),
	}

	// Validate required configuration
	if err := instance.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

// Get returns the singleton config instance
// Must call Init() before using this function
func Get() *Config {
	if instance == nil {
		log.Fatal("Config not initialized. Call config.Init() first.")
	}
	return instance
}

// Validate checks that all required configuration is present
func (c *Config) Validate() error {
	if c.SpotifyClientID == "" {
		return fmt.Errorf("SPOTIFY_CLIENT_ID is required but not set")
	}
	if c.SpotifyClientSecret == "" {
		return fmt.Errorf("SPOTIFY_CLIENT_SECRET is required but not set")
	}
	return nil
}

// MustInit initializes the configuration and exits if there's an error
func MustInit() {
	if err := Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize configuration: %v\n", err)
		fmt.Fprintln(os.Stderr, "\nPlease ensure you have:")
		fmt.Fprintln(os.Stderr, "  1. Created a .env file (copy from .env.example)")
		fmt.Fprintln(os.Stderr, "  2. Set SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET")
		fmt.Fprintln(os.Stderr, "  3. Or set these as environment variables")
		os.Exit(1)
	}
}
