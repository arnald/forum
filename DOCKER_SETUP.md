# Docker Setup - Quick Start Guide

## What Was Implemented

### Files Created/Modified:
- ✅ **Dockerfile** - Multi-stage build for both server and client
- ✅ **docker-compose.yml** - Production orchestration
- ✅ **docker-compose.dev.yml** - Development environment
- ✅ **entrypoint.sh** - Startup script for both services
- ✅ **.dockerignore** - Exclude unnecessary files from builds
- ✅ **Makefile** - Added Docker commands
- ✅ **Client config** - Made backend URL configurable via BACKEND_URL env var

## Quick Start

### Production Mode:
```bash
# Build and start
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### Development Mode:
```bash
# Start with volume mounts and seed data
make docker-dev

# Or build first
make docker-dev-build
```

## Access Points:
- Frontend: http://localhost:3001
- Backend API: http://localhost:8080/api/v1
- Health Check: http://localhost:8080/api/v1/health

## Available Commands:
```bash
make docker-build        # Build image
make docker-up           # Start services
make docker-down         # Stop services
make docker-logs         # View logs
make docker-restart      # Restart services
make docker-ps           # Show containers
make docker-clean        # Remove everything
make docker-dev          # Start dev mode
make docker-dev-build    # Build and start dev
```

## Architecture:
- **Single Container**: Runs both backend (port 8080) and frontend (port 3001)
- **Persistent Volumes**: Database and uploads survive container restarts
- **Health Checks**: Automatic monitoring of backend health
- **Non-root User**: Runs as appuser (uid 1000) for security

## Environment Variables:
All configurable via docker-compose.yml:
- `BACKEND_URL` - How frontend reaches backend (default: http://localhost:8080/api/v1)
- `SERVER_PORT` - Backend port (default: 8080)
- `CLIENT_PORT` - Frontend port (default: 3001)
- `DB_SEED_ON_START` - Load seed data (true/false)
- OAuth credentials (optional)

## Development Features:
- Volume mounts for frontend/db files (see docker-compose.dev.yml)
- Seed data automatically loaded
- Local uploads directory mounted

## Production Features:
- Multi-stage build (optimized image size)
- No seed data
- Persistent volumes for data
- Health checks and restart policies
- Non-root user execution

## Next Steps:
1. Test locally: `make docker-up`
2. Configure OAuth credentials in docker-compose.yml if needed
3. For production deployment:
   - Set `SESSION_SECURE_COOKIE=true`
   - Update `BACKEND_URL` if using different domains
   - Consider HTTPS/TLS setup (reverse proxy recommended)
