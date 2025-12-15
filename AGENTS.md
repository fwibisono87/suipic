# AGENTS.md

## Commands

### Initial Setup
```bash
# Backend setup
cd backend
go mod download
cp .env.example .env
# Edit .env with your configuration

# Run database migrations
go run cmd/migrate/main.go -action=up
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
- Modular structure with separate packages for config, handlers, db, models, services
- Environment-based configuration using .env files
- PostgreSQL database with migration support
- ElasticSearch integration for photo search
- Graceful shutdown support

## Database
- **Migrations**: Located in `backend/db/migrations/`
- **Models**: Go structs in `backend/models/`
- **Schema**: Users, Albums, Photos, AlbumUsers (junction), Comments
- **Run migrations**: `go run cmd/migrate/main.go -action=up`
- **Rollback**: `go run cmd/migrate/main.go -action=down`

## Code Style
- Follow standard Go conventions
- Use gofmt for formatting
- Organize code into logical packages
- No comments unless necessary for complex logic

## ElasticSearch Integration
- **Documentation**: See `backend/ELASTICSEARCH.md` for detailed documentation
- **Search Endpoint**: `GET /api/search` with query parameters for filtering
- **Bulk Indexing**: `POST /api/albums/:albumId/index` to index all photos in an album
- **Automatic Indexing**: Photos are indexed on create, update, and when comments are added
- **Docker Setup**: Use `backend/docker-compose.elasticsearch.yml` for local development
