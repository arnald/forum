#INITIAL PLANNING FOR PDOCKER FILE IN PROJECT

# Build the Go binary

FROM golang:1.24.4-alpine AS builder

# Install runtime dependencies (SQLite3 needs these)
RUN apk --no-cache add ca-certificates libc6-compat

# Install build dependencies and latest stable SQLite3
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory and copy source files
WORKDIR /app
COPY . .

# Download dependencies
RUN go mod download

# Build the backend server binary with CGO enabled for SQLite3 support
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-extldflags=-static" -o /bin/server ./cmd/server/main.go

# Build the client (frontend server) binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/client ./cmd/client/main.go

#Minimal runtime image
FROM alpine:latest AS runtime

# Install runtime dependencies (e.g., SSL certs)
RUN apk --no-cache add ca-certificates

# Copy binaries and assets from builder
COPY --from=builder /bin/server /app/server
COPY --from=builder /app/storage /app/storage 
COPY --from=builder /app/RUNNING-FORUM/docs /app/docs 

# Copy environment file template
COPY --from=builder /app/env.example /app/.env

# Set environment variables
ENV GIN_MODE=release \
    # Server Configuration
    SERVER_HOST=0.0.0.0 \
    SERVER_PORT=8080 \
    SERVER_ENVIRONMENT=production \
    SERVER_API_CONTEXT_V1=/api/v1 \
    SERVER_READ_TIMEOUT=5 \
    SERVER_WRITE_TIMEOUT=10 \
    SERVER_IDLE_TIMEOUT=15 \
    # Client Configuration
    CLIENT_HOST=0.0.0.0 \
    CLIENT_PORT=3000 \
    CLIENT_ENVIRONMENT=production \
    # Database Configuration
    SQLITE_DB_PATH=/app/storage/database.db \
    # Path Configuration
    CONFIG_PATH=/app/config

# Expose ports
EXPOSE ${SERVER_PORT}
EXPOSE ${CLIENT_PORT}

# Run both server and client
CMD ["/app/server"]