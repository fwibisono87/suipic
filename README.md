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

## Quick Start

**🚀 Want to get started fast? See [QUICKSTART.md](QUICKSTART.md) for a 5-minute setup guide!**

### Automated Setup (Recommended)

Use the setup script for quick installation:

```bash
# Clone the repository
git clone <repository-url>
cd suipic

# Linux/Mac
chmod +x setup.sh
./setup.sh --docker

# Windows
.\setup.ps1 -Docker

# Start all services
./docker-start.sh    # Linux/Mac
docker-start.bat     # Windows
```

Access the application at http://localhost with default credentials:
- **Admin**: admin@suipic.local / admin123
- **MinIO Console**: http://localhost:9001 (minioadmin / minioadmin)

### Manual Docker Setup

```bash
# Copy environment configuration
cp .env.docker .env

# Edit .env with your settings (change passwords!)
nano .env

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

**📖 For detailed setup instructions, see [SETUP.md](SETUP.md)**

## Prerequisites

### Docker Deployment (Recommended)
- Docker 20.10+
- Docker Compose 2.0+
- 4GB RAM minimum (8GB recommended)

### Local Development
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Elasticsearch 8+
- MinIO (latest)

## Local Development

### Infrastructure Only (Docker)

Run infrastructure services while developing locally:

```bash
# Start PostgreSQL, Elasticsearch, and MinIO
docker-compose -f docker-compose.infra.yml up -d

# Terminal 1: Backend
cd backend
cp .env.example .env  # Edit with your config
go run cmd/migrate/main.go -action=up
go run main.go

# Terminal 2: Frontend
cd frontend
pnpm install
pnpm run dev
```

### Using Makefile

```bash
make help          # Show all available commands
make up            # Start all services
make logs          # View logs
make down          # Stop all services
```

### Development Commands

See [AGENTS.md](AGENTS.md) for:
- Build commands
- Test commands
- Lint commands
- Migration commands

## Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - 5-minute quick start guide
- **[SETUP.md](SETUP.md)** - Complete setup and deployment guide
- [Docker Deployment Guide](DOCKER.md) - Docker-specific documentation
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
./setup.sh --docker          # Linux/Mac
.\setup.ps1 -Docker          # Windows
docker-compose up -d
```

### Production

```bash
# Review security checklist in SETUP.md
# Update .env with production settings
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

**Important**: Before production deployment:
- Change all default passwords
- Generate strong JWT secret: `openssl rand -base64 32`
- Configure SSL/HTTPS certificates
- Set `ENV=production`
- Update `CORS_ORIGINS` to your domain
- Enable database SSL
- Set up regular backups

See [SETUP.md](SETUP.md) for comprehensive deployment instructions including:
- Security checklist
- SSL/HTTPS configuration
- Database backup and restore
- Monitoring and troubleshooting
- Production optimization

## Environment Variables

All configuration is done via environment variables. Key settings:

```env
# Database
DB_HOST=postgres
DB_PASSWORD=strong_password_here

# Security  
JWT_SECRET=generate-with-openssl-rand-base64-32

# Admin User (created automatically on first start)
ADMIN_EMAIL=admin@suipic.local
ADMIN_PASSWORD=change_me_immediately
ADMIN_USERNAME=admin

# Storage
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=strong_secret_here
MINIO_BUCKET=suipic

# Search
ES_ADDRESSES=http://elasticsearch:9200

# CORS
CORS_ORIGINS=http://localhost,https://yourdomain.com
```

See [SETUP.md](SETUP.md) for complete environment variable documentation.

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

## License

[Add your license here]

## Contributing

[Add contribution guidelines here]
