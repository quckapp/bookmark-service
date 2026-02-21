# Build stage
# NOTE: Build context is repo root (.) to resolve replace directive for go-auth
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go-auth package (referenced via replace directive)
COPY packages/go-auth /app/packages/go-auth

# Copy service go mod files
WORKDIR /app/services/bookmark-service
COPY services/bookmark-service/go.mod services/bookmark-service/go.sum ./
RUN go mod download

# Copy service source code
COPY services/bookmark-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy binary
COPY --from=builder /app/server .

# Set ownership
RUN chown -R appuser:appgroup /app

USER appuser

# Expose port
EXPOSE 5010

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:5010/health || exit 1

# Run the application
CMD ["./server"]
