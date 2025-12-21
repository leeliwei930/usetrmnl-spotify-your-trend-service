#!/bin/sh

# Generate and display service token on container startup
echo "========================================="
echo "Generating Service Token..."
echo "========================================="
./usetrmnl-spotify-service generate service:key
echo ""
echo "========================================="
echo "Starting Server..."
echo "========================================="

# Start the server with the provided arguments
exec ./usetrmnl-spotify-service "$@"
