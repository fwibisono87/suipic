# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Build migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate cmd/migrate/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl netcat-openbsd

# Copy binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

# Copy migration files and entrypoint
COPY --from=builder /app/db/migrations ./db/migrations
COPY --from=builder /app/cmd ./cmd
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:3000/api/health || exit 1

# Run entrypoint
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./main"]
