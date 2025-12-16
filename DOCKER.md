# Docker Deployment Guide for Suipic

Complete guide for deploying Suipic using Docker and Docker Compose.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Configuration](#configuration)
- [Development Setup](#development-setup)
- [Production Deployment](#production-deployment)
- [SSL/HTTPS Setup](#sslhttps-setup)
- [Maintenance](#maintenance)
- [Troubleshooting](#troubleshooting)

## Prerequisites

- Docker Engine 20.10 or higher
- Docker Compose 2.0 or higher
- At least 4GB of available RAM
- At least 10GB of available disk space

### Install Docker

**Linux:**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

**macOS:**
Download and install Docker Desktop from https://www.docker.com/products/docker-desktop

**Windows:**
Download and install Docker Desktop from https://www.docker.com/products/docker-desktop

## Quick Start

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd suipic
   ```

2. **Configure environment:**
   ```bash
   cp .env.docker .env
   # Edit .env and change default passwords
   ```

3. **Start services:**
   ```bash
   docker-compose up -d
   ```

4. **Check status:**
   ```bash
   docker-compose ps
   docker-compose logs -f
   ```

5. **Access the application:**
   - Frontend: http://localhost
   - Backend API: http://localhost/api
   - MinIO Console: http://localhost:9001 (minioadmin/minioadmin)
   - Elasticsearch: http://localhost:9200

## Architecture

```
┌─────────────────────────────────────────┐
│         Nginx (Reverse Proxy)          │
│         HTTP: 80, HTTPS: 443           │
└─────────────┬───────────────────────────┘
              │
     ┌────────┴─────────┐
     │                  │
┌────▼─────┐      ┌────▼─────┐
│ Frontend │      │ Backend  │
│ SvelteKit│      │ Go Fiber │
│  :3001   │      │  :3000   │
└──────────┘      └────┬─────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   ┌────▼────┐   ┌────▼────┐   ┌────▼────┐
   │Postgres │   │Elastic  │   │  MinIO  │
   │  :5432  │   │Search   │   │:9000/1  │
   └─────────┘   │ :9200   │   └─────────┘
                 └─────────┘
```

### Services

- **Nginx**: Reverse proxy handling routing, SSL termination, rate limiting, and gzip compression
- **Frontend**: SvelteKit app with Node.js adapter (multi-stage pnpm build)
- **Backend**: Go Fiber REST API (multi-stage build with optimized binaries)
- **PostgreSQL**: Primary database
- **Elasticsearch**: Photo search and indexing
- **MinIO**: S3-compatible object storage

## Docker Images

### Backend Dockerfile

The backend uses a multi-stage Docker build for optimal production images:

**Builder Stage:**
- Base: `golang:1.24-alpine`
- Compiles Go binaries with CGO disabled
- Build optimizations: `-ldflags='-w -s -extldflags "-static"'`
- Strips debug symbols for smaller binary size
- Builds both main application and migration tool

**Final Stage:**
- Base: `alpine:latest`
- Minimal runtime with only required dependencies
- Runs as non-root user (`appuser:1000`)
- Includes health checks on `/api/health`
- Total image size: ~20-30MB

**Security Features:**
- Non-root user execution
- Minimal attack surface
- No build tools in final image
- Read-only file system compatible

### Frontend Dockerfile

The frontend uses a multi-stage Docker build with pnpm:

**Builder Stage:**
- Base: `node:20-alpine`
- Uses pnpm 8.15.0 for efficient package management
- Frozen lockfile for reproducible builds
- Builds production-optimized SvelteKit application
- Supports build-time API URL configuration

**Final Stage:**
- Base: `node:20-alpine`
- Runs as non-root user (`appuser:1000`)
- Only includes built application and minimal dependencies
- Uses SvelteKit's Node adapter for SSR
- Includes health checks

**Build Arguments:**
- `PUBLIC_API_URL`: Configure API endpoint at build time

### Nginx Configuration

The nginx reverse proxy includes advanced features:

**Rate Limiting:**
- Login endpoint: 5 requests/minute (burst: 2)
- API endpoints: 10 requests/second (burst: 20)
- General traffic: 30 requests/second (burst: 50)
- Connection limit: 10 concurrent per IP

**Gzip Compression:**
- Enabled for all text-based content
- Compression level: 6
- Minimum file size: 256 bytes
- Supports 30+ MIME types including JSON, JS, CSS, SVG, fonts

**SSL/TLS Configuration:**
- TLS 1.2 and 1.3 support
- Modern cipher suites
- OCSP stapling enabled
- HSTS header with preload support
- Session caching for performance

**Security Headers:**
- X-Frame-Options: SAMEORIGIN
- X-Content-Type-Options: nosniff
- X-XSS-Protection: 1; mode=block
- Referrer-Policy: no-referrer-when-downgrade
- Strict-Transport-Security (HTTPS only)

**Performance Features:**
- HTTP/2 support
- Keepalive connections to upstreams
- Static asset caching (1 year for images, fonts, etc.)
- Connection pooling
- Optimized buffer sizes

## Configuration

### Environment Variables

Copy `.env.docker` to `.env` and customize:

```bash
# Critical settings to change for production
POSTGRES_PASSWORD=strong_password_here
JWT_SECRET=long_random_string_here
ADMIN_PASSWORD=secure_admin_password
MINIO_ROOT_PASSWORD=strong_minio_password
```

### Service Ports

Default port mapping:

| Service       | Internal | External | Description          |
|--------------|----------|----------|----------------------|
| Nginx        | 80       | 80       | HTTP                 |
| Nginx        | 443      | 443      | HTTPS                |
| Backend      | 3000     | 3000     | API (direct access)  |
| Frontend     | 3001     | 3001     | App (direct access)  |
| PostgreSQL   | 5432     | 5432     | Database             |
| Elasticsearch| 9200     | 9200     | Search API           |
| MinIO        | 9000     | 9000     | S3 API               |
| MinIO        | 9001     | 9001     | Web Console          |

## Development Setup

For development with hot-reload:

1. **Use development override:**
   ```bash
   docker-compose up -d
   ```
   The `docker-compose.override.yml` is automatically applied.

2. **Watch logs:**
   ```bash
   docker-compose logs -f backend frontend
   ```

3. **Rebuild after changes:**
   ```bash
   docker-compose build backend frontend
   docker-compose up -d backend frontend
   ```

### Development Commands

Using the Makefile:
```bash
make help          # Show all available commands
make build         # Build all images
make up            # Start services
make down          # Stop services
make logs          # View all logs
make backend-logs  # View backend logs only
make frontend-logs # View frontend logs only
```

## Production Deployment

### 1. Prepare Environment

```bash
# Copy and edit production config
cp .env.docker .env

# Generate strong passwords
openssl rand -base64 32  # For JWT_SECRET
openssl rand -base64 24  # For POSTGRES_PASSWORD
```

### 2. Configure SSL Certificates

See [SSL/HTTPS Setup](#sslhttps-setup) section below.

### 3. Deploy with Production Settings

```bash
# Build images
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

# Start services
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Check status
docker-compose ps
```

### 4. Security Hardening

**Update .env:**
- Change all default passwords
- Use strong JWT_SECRET (32+ characters)
- Set appropriate CORS_ORIGINS
- Review exposed ports

**File Permissions:**
```bash
chmod 600 .env
chmod 600 docker/nginx/ssl/*.pem
```

**Firewall:**
```bash
# Only expose HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

## SSL/HTTPS Setup

### Self-Signed Certificate (Development)

```bash
# Generate certificate
make ssl-cert

# Or manually:
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout docker/nginx/ssl/key.pem \
  -out docker/nginx/ssl/cert.pem \
  -subj "/C=US/ST=State/L=City/O=Suipic/CN=localhost"
```

### Let's Encrypt (Production)

**Option 1: Using Certbot**

```bash
# Install certbot
sudo apt-get install certbot

# Stop nginx temporarily
docker-compose stop nginx

# Generate certificate
sudo certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com

# Copy certificates
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem docker/nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem docker/nginx/ssl/key.pem
sudo chown $USER:$USER docker/nginx/ssl/*.pem
```

**Option 2: Using Docker Certbot**

```bash
docker run -it --rm \
  -v /etc/letsencrypt:/etc/letsencrypt \
  -v /var/lib/letsencrypt:/var/lib/letsencrypt \
  -p 80:80 \
  certbot/certbot certonly --standalone -d yourdomain.com
```

### Enable HTTPS in Nginx

The nginx configuration is pre-configured with HTTPS support. Once certificates are in place:

1. **Verify certificates are mounted:**
   ```bash
   docker-compose exec nginx ls -la /etc/nginx/ssl/
   ```

2. **The HTTPS server block is active** - nginx will automatically serve on port 443 when certificates are present

3. **Optional: Enable HTTP to HTTPS redirect**
   Edit `docker/nginx/conf.d/default.conf` and uncomment these lines:
   ```nginx
   # location / {
   #     return 301 https://$host$request_uri;
   # }
   ```

4. **Update server_name** for your domain (optional):
   Replace `server_name _;` with `server_name yourdomain.com;`

5. **Restart nginx:**
   ```bash
   docker-compose restart nginx
   ```

### Certificate Mounting

SSL certificates are mounted via Docker volumes in `docker-compose.yml`:

```yaml
nginx:
  volumes:
    - ./docker/nginx/ssl:/etc/nginx/ssl:ro
```

The `:ro` flag mounts certificates as read-only for security. Place your certificate files in `docker/nginx/ssl/`:
- `cert.pem` - SSL certificate (public key)
- `key.pem` - Private key
- `chain.pem` - Certificate chain (optional)

**Important:** Certificate files are excluded from git via `.gitignore` for security.

### Auto-Renewal Setup

Add to crontab:
```bash
0 0 * * 0 certbot renew --quiet && docker-compose restart nginx
```

## Maintenance

### Backup

**Database Backup:**
```bash
# Create backup
docker-compose exec postgres pg_dump -U suipic suipic > backup.sql

# With Docker volume
docker run --rm \
  -v suipic_postgres-data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/postgres-backup.tar.gz -C /data .
```

**MinIO Backup:**
```bash
docker run --rm \
  -v suipic_minio-data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/minio-backup.tar.gz -C /data .
```

**Full Backup:**
```bash
# Backup all volumes
docker-compose down
docker run --rm \
  -v suipic_postgres-data:/postgres \
  -v suipic_elasticsearch-data:/elasticsearch \
  -v suipic_minio-data:/minio \
  -v $(pwd):/backup \
  alpine sh -c "tar czf /backup/full-backup-$(date +%Y%m%d).tar.gz /postgres /elasticsearch /minio"
docker-compose up -d
```

### Restore

```bash
# Stop services
docker-compose down

# Restore volume
docker run --rm \
  -v suipic_postgres-data:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/postgres-backup.tar.gz -C /data

# Start services
docker-compose up -d
```

### Updates

**Update application:**
```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose build
docker-compose up -d
```

**Update base images:**
```bash
docker-compose pull
docker-compose up -d
```

### Logs

**View logs:**
```bash
docker-compose logs -f                    # All services
docker-compose logs -f backend            # Backend only
docker-compose logs --tail=100 backend    # Last 100 lines
```

**Log location:**
- Container logs: `/var/log/nginx/` (nginx)
- JSON logs: `/var/lib/docker/containers/`

### Monitoring

**Resource usage:**
```bash
docker stats
```

**Health checks:**
```bash
docker-compose ps
curl http://localhost:9200/_cluster/health?pretty
curl http://localhost/api/health
```

## Troubleshooting

### Service Won't Start

**Check logs:**
```bash
docker-compose logs <service-name>
```

**Check health:**
```bash
docker-compose ps
docker inspect <container-name>
```

### Database Connection Issues

**Verify PostgreSQL:**
```bash
docker-compose exec postgres pg_isready -U suipic
docker-compose exec postgres psql -U suipic -d suipic -c "SELECT 1"
```

**Check connectivity:**
```bash
docker-compose exec backend nc -zv postgres 5432
```

### Migration Errors

**Run migrations manually:**
```bash
docker-compose exec backend ./migrate -action=up
```

**Reset database (WARNING: destroys data):**
```bash
docker-compose down -v
docker-compose up -d
```

### Port Conflicts

**Check port usage:**
```bash
sudo netstat -tlnp | grep -E ':(80|443|3000|3001|5432|9000|9001|9200)'
```

**Change ports in docker-compose.yml:**
```yaml
services:
  nginx:
    ports:
      - "8080:80"  # Use port 8080 instead of 80
```

### Elasticsearch Issues

**Check health:**
```bash
curl http://localhost:9200/_cluster/health?pretty
```

**Increase memory:**
Edit `docker-compose.yml`:
```yaml
elasticsearch:
  environment:
    - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
```

### Out of Disk Space

**Clean up:**
```bash
docker system prune -a --volumes
docker volume prune
```

**Check usage:**
```bash
docker system df
du -sh /var/lib/docker
```

### Reset Everything

**Complete reset (destroys all data):**
```bash
docker-compose down -v
docker system prune -a --volumes -f
docker-compose up -d
```

## Advanced Configuration

### Custom Network

```yaml
networks:
  suipic-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### Resource Limits

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

### Multiple Environments

```bash
# Staging
docker-compose -f docker-compose.yml -f docker-compose.staging.yml up -d

# Production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Elasticsearch Documentation](https://www.elastic.co/guide/)
- [MinIO Documentation](https://min.io/docs/)

## Support

For issues and questions:
- Check the logs: `docker-compose logs -f`
- Review this guide's troubleshooting section
- Check Docker service status: `docker-compose ps`
