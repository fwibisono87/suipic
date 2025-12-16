# Suipic Setup Guide

This guide provides detailed instructions for setting up and deploying Suipic in various environments.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Prerequisites](#prerequisites)
3. [Environment Configuration](#environment-configuration)
4. [Docker Deployment](#docker-deployment)
5. [Local Development Setup](#local-development-setup)
6. [Production Deployment](#production-deployment)
7. [Post-Installation](#post-installation)
8. [Troubleshooting](#troubleshooting)

## Quick Start

### Using Docker (Recommended)

The fastest way to get Suipic running:

```bash
# Clone the repository
git clone <repository-url>
cd suipic

# Run setup script
chmod +x setup.sh
./setup.sh --docker

# Start all services
./docker-start.sh
```

On Windows:
```powershell
# Run setup script
.\setup.ps1 -Docker

# Start all services
.\docker-start.bat
```

Access the application at:
- **Frontend**: http://localhost
- **Backend API**: http://localhost/api
- **MinIO Console**: http://localhost:9001

Default credentials:
- **Admin**: admin@suipic.local / admin123
- **MinIO**: minioadmin / minioadmin

## Prerequisites

### Docker Deployment

- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **RAM**: 4GB minimum (8GB recommended)
- **Disk**: 10GB minimum for images and data

### Local Development

- **Go**: 1.21 or higher
- **Node.js**: 18 or higher
- **PostgreSQL**: 15 or higher
- **Elasticsearch**: 8.x
- **MinIO**: Latest
- **pnpm**: Latest (for frontend)

## Environment Configuration

### Backend Environment Variables

The backend requires configuration via environment variables. Copy the example file and customize:

```bash
cd backend
cp .env.example .env
# Edit .env with your settings
```

#### Required Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PORT` | Backend server port | `3000` | `3000` |
| `ENV` | Environment (development/production) | `development` | `production` |
| `DB_HOST` | PostgreSQL hostname | `localhost` | `postgres` |
| `DB_PORT` | PostgreSQL port | `5432` | `5432` |
| `DB_USER` | PostgreSQL username | `suipic` | `suipic` |
| `DB_PASSWORD` | PostgreSQL password | `password` | `strong_password` |
| `DB_NAME` | Database name | `suipic` | `suipic` |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` | `require` |
| `ES_ADDRESSES` | Elasticsearch URLs (comma-separated) | `http://localhost:9200` | `http://elasticsearch:9200` |
| `MINIO_ENDPOINT` | MinIO endpoint | `localhost:9000` | `minio:9000` |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` | `your_access_key` |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin` | `your_secret_key` |
| `MINIO_USE_SSL` | Use SSL for MinIO | `false` | `true` |
| `MINIO_BUCKET` | MinIO bucket name | `suipic` | `suipic-photos` |
| `JWT_SECRET` | JWT signing secret | - | `your-random-secret` |
| `JWT_EXPIRY` | JWT token expiry | `24h` | `168h` |
| `CORS_ORIGINS` | Allowed CORS origins (comma-separated) | `http://localhost:5173,http://localhost:3001` | `https://yourdomain.com` |
| `ADMIN_EMAIL` | Initial admin email | `admin@suipic.local` | `admin@company.com` |
| `ADMIN_PASSWORD` | Initial admin password | `admin123` | `strong_password` |
| `ADMIN_USERNAME` | Initial admin username | `admin` | `administrator` |

### Docker Environment Variables

For Docker deployments, copy and customize `.env.docker`:

```bash
cp .env.docker .env
# Edit .env with your settings
```

The Docker environment file includes additional variables for infrastructure services:

- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` - PostgreSQL configuration
- `MINIO_ROOT_USER`, `MINIO_ROOT_PASSWORD` - MinIO root credentials
- `NGINX_HTTP_PORT`, `NGINX_HTTPS_PORT` - NGINX port configuration
- `PUBLIC_API_URL`, `ORIGIN` - Frontend configuration

## Docker Deployment

### Standard Deployment

1. **Run the setup script**:
   ```bash
   ./setup.sh --docker
   ```

2. **Review and update `.env`**:
   ```bash
   nano .env  # or use your preferred editor
   ```
   
   At minimum, change:
   - `POSTGRES_PASSWORD`
   - `JWT_SECRET`
   - `ADMIN_PASSWORD`
   - `MINIO_ROOT_PASSWORD`

3. **Start all services**:
   ```bash
   docker-compose up -d
   ```

4. **Check service health**:
   ```bash
   docker-compose ps
   docker-compose logs -f backend
   ```

### Infrastructure Only

To run only the infrastructure services (PostgreSQL, Elasticsearch, MinIO) for local development:

```bash
docker-compose -f docker-compose.infra.yml up -d
```

### Production Deployment

For production with additional configuration:

```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## Local Development Setup

### 1. Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Create and configure .env
cp .env.example .env
# Edit .env with your configuration

# Start infrastructure services (if using Docker)
docker-compose -f ../docker-compose.infra.yml up -d

# Run database migrations
go run cmd/migrate/main.go -action=up

# Start the backend server
go run main.go
```

The backend will be available at `http://localhost:3000`.

### 2. Frontend Setup

```bash
cd frontend

# Install pnpm if not already installed
npm install -g pnpm

# Install dependencies
pnpm install

# Create and configure .env (if needed)
cp .env.example .env

# Start the development server
pnpm run dev
```

The frontend will be available at `http://localhost:5173`.

## Production Deployment

### Security Checklist

Before deploying to production, ensure you:

- [ ] Change all default passwords
- [ ] Generate strong JWT_SECRET: `openssl rand -base64 32`
- [ ] Use strong database passwords
- [ ] Configure SSL/HTTPS (see [SSL Configuration](#ssl-configuration))
- [ ] Set `ENV=production` for backend
- [ ] Enable database SSL (`DB_SSLMODE=require`)
- [ ] Update `CORS_ORIGINS` to your domain
- [ ] Review and restrict network access
- [ ] Set up regular backups
- [ ] Configure log aggregation/monitoring
- [ ] Use strong MinIO credentials
- [ ] Never commit `.env` files to version control

### SSL Configuration

For HTTPS support:

1. Place SSL certificates in `docker/nginx/ssl/`:
   ```bash
   docker/nginx/ssl/
   ├── cert.pem
   └── key.pem
   ```

2. Update `docker/nginx/conf.d/default.conf` to enable SSL

3. Update `.env`:
   ```env
   NGINX_HTTPS_PORT=443
   CORS_ORIGINS=https://yourdomain.com
   PUBLIC_API_URL=https://yourdomain.com/api
   ```

4. Restart services:
   ```bash
   docker-compose restart nginx frontend
   ```

### Database Backups

#### Manual Backup

```bash
docker exec suipic-postgres pg_dump -U suipic suipic > backup_$(date +%Y%m%d_%H%M%S).sql
```

#### Automated Backups

Add to crontab:
```cron
0 2 * * * docker exec suipic-postgres pg_dump -U suipic suipic > /path/to/backups/suipic_$(date +\%Y\%m\%d).sql
```

#### Restore from Backup

```bash
docker exec -i suipic-postgres psql -U suipic suipic < backup.sql
```

### Monitoring and Logs

View logs for all services:
```bash
docker-compose logs -f
```

View logs for specific service:
```bash
docker-compose logs -f backend
docker-compose logs -f postgres
docker-compose logs -f elasticsearch
```

## Post-Installation

### Initial Setup Steps

1. **Access the application**:
   - Navigate to http://localhost (or your domain)
   - Login with admin credentials

2. **Change admin password**:
   - Go to user settings
   - Update password immediately

3. **Create photographer accounts**:
   - As admin, create photographer accounts
   - Photographers can then create client accounts

4. **Create first album**:
   - Create an album
   - Upload photos
   - Share with clients

### Verify Services

Check that all services are running properly:

```bash
# Backend health check
curl http://localhost/api/health

# Elasticsearch
curl http://localhost:9200/_cluster/health

# MinIO
curl http://localhost:9000/minio/health/live

# PostgreSQL
docker exec suipic-postgres pg_isready -U suipic
```

### Testing Photo Upload

1. Login as photographer
2. Create an album
3. Upload a test photo
4. Verify thumbnail generation
5. Test search functionality

## Troubleshooting

### Common Issues

#### Database Connection Failed

**Problem**: Backend cannot connect to PostgreSQL

**Solution**:
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check PostgreSQL logs
docker-compose logs postgres

# Verify credentials in .env match POSTGRES_* variables
```

#### MinIO Bucket Not Created

**Problem**: Photos fail to upload

**Solution**:
```bash
# Check MinIO logs
docker-compose logs minio

# Verify MinIO is accessible
curl http://localhost:9000/minio/health/live

# The bucket is created automatically on first backend start
# Check backend logs for bucket initialization
docker-compose logs backend | grep -i bucket
```

#### Elasticsearch Not Responding

**Problem**: Search functionality not working

**Solution**:
```bash
# Check Elasticsearch health
curl http://localhost:9200/_cluster/health

# Check Elasticsearch logs
docker-compose logs elasticsearch

# Increase Elasticsearch memory if needed (in docker-compose.yml)
ES_JAVA_OPTS=-Xms1g -Xmx1g
```

#### Admin User Not Created

**Problem**: Cannot login with admin credentials

**Solution**:
```bash
# Verify ADMIN_* variables are set in .env
grep ADMIN .env

# Check backend logs for admin creation
docker-compose logs backend | grep -i admin

# Admin user is only created if it doesn't exist
# Check database directly
docker exec -it suipic-postgres psql -U suipic suipic -c "SELECT * FROM users WHERE role='admin';"
```

#### Frontend Cannot Reach Backend

**Problem**: API calls fail from frontend

**Solution**:
```bash
# Verify PUBLIC_API_URL in frontend environment
# Should be http://localhost/api for Docker setup

# Check CORS_ORIGINS in backend .env
# Must include frontend origin

# Check nginx logs
docker-compose logs nginx
```

#### Port Already in Use

**Problem**: Docker fails to start due to port conflict

**Solution**:
```bash
# Check what's using the port
sudo lsof -i :80  # or :5432, :9200, etc.

# Change port in docker-compose.yml or stop conflicting service
```

### Resetting the Application

To completely reset and start fresh:

```bash
# Stop all services
docker-compose down

# Remove all volumes (WARNING: destroys all data)
docker-compose down -v

# Remove images (optional)
docker-compose down --rmi all

# Start fresh
docker-compose up -d
```

### Getting Help

- Check logs: `docker-compose logs -f`
- Review [AGENTS.md](AGENTS.md) for development commands
- Review [DOCKER.md](DOCKER.md) for Docker-specific documentation
- Check [backend/ELASTICSEARCH.md](backend/ELASTICSEARCH.md) for search issues

### Performance Tuning

#### For Production Workloads

1. **Increase Elasticsearch memory**:
   ```yaml
   # In docker-compose.yml
   ES_JAVA_OPTS: "-Xms2g -Xmx2g"
   ```

2. **Optimize PostgreSQL**:
   ```yaml
   # Add to postgres service in docker-compose.yml
   command: postgres -c shared_buffers=256MB -c max_connections=200
   ```

3. **Scale services**:
   ```bash
   docker-compose up -d --scale backend=3
   ```

4. **Use external managed services**:
   - Use managed PostgreSQL (RDS, Cloud SQL, etc.)
   - Use managed Elasticsearch
   - Use S3 instead of MinIO

## Next Steps

- Review [README.md](README.md) for feature overview
- Check [backend/README.md](backend/README.md) for API documentation
- Review [backend/AUTH.md](backend/AUTH.md) for authentication details
- See [backend/SEARCH_EXAMPLES.md](backend/SEARCH_EXAMPLES.md) for search usage

## Support

For issues and questions:
- Create an issue in the repository
- Check existing documentation
- Review logs for error messages
