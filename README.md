# Suipic

Yet another photo manager - A modern, full-featured photo management application.

## Features

- 📸 Photo upload and management
- 🔍 Advanced search with Elasticsearch
- 📁 Album organization
- 💬 Photo comments
- 👥 Multi-user support with role-based access
- 🔐 JWT-based authentication
- 📦 S3-compatible storage with MinIO
- 🐳 Full Docker support

## Tech Stack

- **Backend**: Go 1.21+ with Fiber v2
- **Frontend**: SvelteKit with TypeScript
- **Database**: PostgreSQL 15
- **Search**: Elasticsearch 8
- **Storage**: MinIO (S3-compatible)
- **Reverse Proxy**: Nginx
- **Containerization**: Docker & Docker Compose

## Quick Start with Docker

The easiest way to get started is using Docker:

```bash
# Clone the repository
git clone <repository-url>
cd suipic

# Linux/Mac
./docker-start.sh

# Windows
docker-start.bat
```

For detailed Docker documentation, see [DOCKER.md](DOCKER.md).

## Manual Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Elasticsearch 8+
- MinIO

### Backend Setup

```bash
cd backend
go mod download
cp .env.example .env
# Edit .env with your configuration
go run cmd/migrate/main.go -action=up
go run main.go
```

### Frontend Setup

```bash
cd frontend
npm install -g pnpm
pnpm install
pnpm run dev
```

## Development

### Using Docker for Infrastructure Only

Run only the infrastructure services (PostgreSQL, Elasticsearch, MinIO) while developing the application locally:

```bash
docker-compose -f docker-compose.infra.yml up -d
cd backend && go run main.go
cd frontend && pnpm run dev
```

### Using Makefile

```bash
make help          # Show all available commands
make up            # Start all services
make logs          # View logs
make down          # Stop all services
```

## Documentation

- [Docker Deployment Guide](DOCKER.md) - Complete Docker setup and deployment guide
- [Backend Documentation](backend/README.md) - Backend API documentation
- [Frontend Documentation](frontend/README.md) - Frontend development guide
- [Elasticsearch Integration](backend/ELASTICSEARCH.md) - Search functionality
- [Authentication](backend/AUTH.md) - Authentication and authorization
- [AGENTS.md](AGENTS.md) - Development commands and architecture

## Project Structure

```
suipic/
├── backend/              # Go backend application
│   ├── cmd/             # Command line tools
│   ├── config/          # Configuration management
│   ├── db/              # Database migrations
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   ├── repository/      # Data access layer
│   ├── services/        # Business logic
│   └── main.go          # Application entry point
├── frontend/            # SvelteKit frontend
│   ├── src/            # Source code
│   │   ├── lib/        # Shared components
│   │   └── routes/     # Page routes
│   └── static/         # Static assets
├── docker/              # Docker configuration
│   ├── nginx/          # Nginx config and SSL
│   └── postgres/       # PostgreSQL init scripts
├── docker-compose.yml   # Main Docker Compose file
├── DOCKER.md           # Docker documentation
└── README.md           # This file
```

## Deployment

### Development

```bash
docker-compose up -d
```

### Production

```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

See [DOCKER.md](DOCKER.md) for comprehensive deployment instructions including:
- SSL/HTTPS setup
- Environment configuration
- Backup and restore
- Monitoring and logs
- Security hardening

## Environment Variables

Key configuration options (see `.env.docker` for full list):

```env
# Database
POSTGRES_PASSWORD=changeme

# Security
JWT_SECRET=your-secret-key-here

# Admin User
ADMIN_EMAIL=admin@suipic.local
ADMIN_PASSWORD=admin123

# Storage
MINIO_ROOT_PASSWORD=changeme
```

## API Endpoints

Once running, the API is available at:

- Development: http://localhost:3000/api
- With Nginx: http://localhost/api

Key endpoints:
- `POST /api/auth/login` - User authentication
- `GET /api/albums` - List albums
- `POST /api/photos` - Upload photo
- `GET /api/search` - Search photos

## Access Points

After starting with Docker:

- **Frontend**: http://localhost
- **Backend API**: http://localhost/api
- **MinIO Console**: http://localhost:9001
- **Elasticsearch**: http://localhost:9200

Default credentials:
- **Admin**: admin@suipic.local / admin123
- **MinIO**: minioadmin / minioadmin
- **PostgreSQL**: suipic / suipic_password

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
