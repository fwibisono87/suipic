# AGENTS.md

## Commands

### Initial Setup
```bash
# Backend setup
cd backend
go mod download
cp .env.example .env
# Edit .env with your configuration
```

### Build
```bash
# Backend build
cd backend
go build -o bin/suipic main.go
```

### Lint
```bash
# Backend lint (requires golangci-lint)
cd backend
golangci-lint run
```

### Tests
```bash
# Backend tests
cd backend
go test ./...
```

### Dev Server
```bash
# Backend dev server
cd backend
go run main.go
```

## Tech Stack
- **Project**: suipic - photo manager
- **Backend**: Go 1.21+ with Fiber v2
- **Database**: PostgreSQL
- **Search**: Elasticsearch v8
- **Storage**: MinIO
- **Authentication**: JWT

## Architecture
- RESTful API with Fiber framework
- Modular structure with separate packages for config, handlers
- Environment-based configuration using .env files
- Graceful shutdown support

## Code Style
- Follow standard Go conventions
- Use gofmt for formatting
- Organize code into logical packages
- No comments unless necessary for complex logic
