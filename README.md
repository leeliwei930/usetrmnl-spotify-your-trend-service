# UseTRMNL Spotify Service - Docker Setup

This directory contains Docker configuration files for running the UseTRMNL Spotify Service with all three components containerized.

## Architecture

The application consists of three services:

- **Backend** (Go/Echo) - Port 8080
- **Agents** (Python/FastAPI) - Port 8000  
- **Frontend** (NuxtJS) - Port 3000

All services run in separate containers and communicate via an internal Docker network.

## Quick Start

### 1. Configure Environment Variables

Before running the services, you need to populate the environment files with your actual credentials:

#### `.env.backend`
```bash
# Required
SPOTIFY_CLIENT_ID=<your_spotify_client_id>
SPOTIFY_CLIENT_SECRET=<your_spotify_client_secret>

# IMPORTANT: AGENT_BASE_URL uses container name for internal Docker networking
AGENT_BASE_URL=http://agents:8000

# Optional: Add additional CORS origins if needed
CORS_ALLOW_ORIGINS=http://localhost:3000,http://frontend:3000
```

#### `.env.agents`
```bash
# Required
OPENROUTER_API_KEY=<your_openrouter_api_key>
```

#### `.env.frontend`
```bash
# For local development (accessing from host machine)
NUXT_PUBLIC_API_BASE=http://localhost:8080/api

# For production, change to your domain
# NUXT_PUBLIC_API_BASE=https://your-domain.com/api
```

### 2. Build and Run

```bash
# Build all images
docker-compose build

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### 3. Access the Services

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Agents API**: http://localhost:8000

## Development Tips

### Rebuilding After Code Changes

```bash
# Rebuild specific service
docker-compose build backend
docker-compose build agents
docker-compose build frontend

# Rebuild and restart
docker-compose up -d --build
```

### Viewing Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f agents
docker-compose logs -f frontend
```

### Running Individual Services

```bash
# Start only agents
docker-compose up agents

# Start backend (automatically starts agents due to dependency)
docker-compose up backend

# Start all
docker-compose up
```

## Image Sizes

The multi-stage Alpine-based builds keep images small:

- **Backend**: ~20-30 MB (Go binary + Alpine)
- **Agents**: ~100-150 MB (Python + dependencies + Alpine)
- **Frontend**: ~200-300 MB (Node + built Nuxt app + Alpine)

## Internal Docker Networking

Services communicate using container names as hostnames:

- Backend → Agents: `http://agents:8000`
- Frontend → Backend: Accessible on host as `http://localhost:8080`

**Important**: The backend `.env` file must use `AGENT_BASE_URL=http://agents:8000` (container name) for internal Docker network communication, not `localhost`.

## Health Checks

The agents service has health checks configured:
- Endpoint: `/health`
- Interval: 10 seconds
- Timeout: 5 seconds
- Retries: 3

The backend service depends on the agents service being healthy before starting.

## Troubleshooting

### Service won't start

Check logs for the specific service:
```bash
docker-compose logs backend
```

### Can't connect to APIs

1. Ensure all environment variables are set correctly
2. Check that ports aren't already in use:
   ```bash
   lsof -i :3000
   lsof -i :8000
   lsof -i :8080
   ```

### CORS errors

Add your frontend URL to `CORS_ALLOW_ORIGINS` in `.env.backend`:
```bash
CORS_ALLOW_ORIGINS=http://localhost:3000,http://frontend:3000,https://your-domain.com
```

### Backend can't reach Agents

Verify `AGENT_BASE_URL` uses the container name:
```bash
# Correct (container name for internal Docker network)
AGENT_BASE_URL=http://agents:8000

# Incorrect (localhost won't work between containers)
AGENT_BASE_URL=http://localhost:8000
```

### Clean rebuild

Remove all containers, images, and volumes:
```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
```

## Production Deployment

For production:

1. Update environment variables:
   - Set proper `SPOTIFY_REDIRECT_URI` to your domain
   - Update `CORS_ALLOW_ORIGINS` to include your production domain
   - Set `NUXT_PUBLIC_API_BASE` to your production API URL

2. Consider using Docker secrets for sensitive data instead of .env files

3. Add volume mounts for persistent data if needed

4. Use a reverse proxy (nginx) in front of the services

5. Enable HTTPS with SSL certificates
