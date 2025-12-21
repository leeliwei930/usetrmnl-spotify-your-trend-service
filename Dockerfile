# ========================================
# Stage 1: Builder
# ========================================
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o usetrmnl-spotify-service .

# ========================================
# Stage 2: Runtime
# ========================================
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/usetrmnl-spotify-service .

# Copy entrypoint script
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

# Copy views and public directories (needed for HTML templates and static files)
COPY --from=builder /build/views ./views
COPY --from=builder /build/public ./public

# Expose port
EXPOSE 8009

# Use entrypoint script that generates service token and starts server
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["start" , "-H", "0.0.0.0", "-p", "8009"]
