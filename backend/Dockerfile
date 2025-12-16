# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y libwebp-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -a \
    -ldflags='-w -s' \
    -o main .

# Build migration tool
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -a \
    -ldflags='-w -s' \
    -o migrate cmd/migrate/main.go

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates curl netcat-openbsd libwebp7 && \
    (getent group 1000 || groupadd -g 1000 appuser) && \
    (getent passwd 1000 || useradd -u 1000 -g appuser appuser)

# Copy binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

# Copy migration files and entrypoint
COPY --from=builder /app/db/migrations ./db/migrations
COPY --from=builder /app/cmd ./cmd
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# Change ownership to appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:3000/api/health || exit 1

# Run entrypoint
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./main"]
