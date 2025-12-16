# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package.json pnpm-lock.yaml* package-lock.json* yarn.lock* ./

# Install pnpm and dependencies
RUN npm install -g pnpm@8.15.0
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

# Copy built application
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json

# Install pnpm and production dependencies only
RUN npm install -g pnpm@8.15.0

# Install adapter-node if not already in package.json
RUN pnpm add @sveltejs/adapter-node

# Expose port
EXPOSE 3001

# Environment variables
ENV PORT=3001
ENV HOST=0.0.0.0

# Run the application
CMD ["node", "build"]
