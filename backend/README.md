# Suipic Backend

Backend API for Suipic photo manager built with Go and Fiber.

## Tech Stack

- **Framework**: Fiber v2
- **Database**: PostgreSQL (lib/pq driver)
- **Search**: Elasticsearch v8
- **Storage**: MinIO
- **Authentication**: JWT (golang-jwt/jwt)
- **Configuration**: godotenv

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Elasticsearch
- MinIO

### Installation

1. Install dependencies:
```bash
go mod download
```

2. Copy environment configuration:
```bash
cp .env.example .env
```

3. Update `.env` with your configuration

### Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:3000` (or the port specified in `.env`)

## API Endpoints

### Health Check
```
GET /api/v1/health
```

Returns service health status.

## Project Structure

```
backend/
├── config/          # Configuration management
├── handlers/        # HTTP request handlers
├── main.go         # Application entry point
├── go.mod          # Go module dependencies
├── .env.example    # Example environment configuration
└── README.md       # This file
```
