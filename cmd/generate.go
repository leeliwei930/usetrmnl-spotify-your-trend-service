/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/leeliwei930/usetrmnl_spotify_service/middleware"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate various resources for the service",
	Long: `Generate various resources for the service such as JWT tokens.

This command provides subcommands for generating different types of resources
needed for the service operation.`,
}

// serviceKeyCmd represents the service:key subcommand
var serviceKeyCmd = &cobra.Command{
	Use:   "service:key",
	Short: "Generate a JWT service token",
	Long: `Generate a JWT service token for authenticating API requests.

This token should be included in the Authorization header as:
  Authorization: Bearer <token>

The token is used to authenticate requests to protected endpoints such as:
  - /api/usetrmnl/spotify/trends
  - /api/usetrmnl/spotify/recent-played

Example:
  usetrmnl_spotify_service generate service:key
  usetrmnl_spotify_service generate service:key --expires 30`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get expiration days from flag
		expires, _ := cmd.Flags().GetInt("expires")

		// Generate token
		token, err := middleware.GenerateServiceToken(expires)
		if err != nil {
			fmt.Printf("Error generating token: %v\n", err)
			fmt.Println("\nPlease ensure JWT_SECRET is set in your .env file.")
			return
		}

		// Display token and usage instructions
		fmt.Println("✓ Service token generated successfully!")
		fmt.Println()
		fmt.Println("Token:")
		fmt.Println(token)
		fmt.Println()
		fmt.Printf("This token will expire in %d days.\n", expires)
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("Include this token in the Authorization header of your API requests:")
		fmt.Println()
		fmt.Println("  curl -X POST http://localhost:1323/api/usetrmnl/spotify/trends \\")
		fmt.Println("    -H \"Authorization: Bearer " + token + "\" \\")
		fmt.Println("    -H \"X-SPOTIFY-CLIENT-ID: your-client-id\" \\")
		fmt.Println("    -H \"X-SPOTIFY-CLIENT-SECRET: your-client-secret\" \\")
		fmt.Println("    -H \"X-SPOTIFY-REFRESH-TOKEN: your-refresh-token\"")
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(serviceKeyCmd)

	// Add flags for the service:key command
	serviceKeyCmd.Flags().IntP("expires", "e", 365, "Number of days until token expiration")
}
