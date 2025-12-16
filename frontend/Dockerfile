# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm@8.15.0

# Copy package files
COPY package.json pnpm-lock.yaml* ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code
COPY . .

# Build argument for API URL
ARG PUBLIC_API_URL=http://localhost/api
ENV PUBLIC_API_URL=$PUBLIC_API_URL

# Build the application
RUN pnpm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Create non-root user
RUN getent group 1000 || addgroup -g 1000 appuser && \
    getent passwd 1000 || adduser -D -u 1000 -G appuser appuser

# Install pnpm
RUN npm install -g pnpm@8.15.0

# Copy package files and install production dependencies
COPY package.json ./
RUN pnpm add @sveltejs/adapter-node

# Copy built application from builder
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json

# Change ownership to appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3001

# Environment variables
ENV PORT=3001
ENV HOST=0.0.0.0
ENV NODE_ENV=production

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3001/ || exit 1

# Run the application
CMD ["node", "build"]
