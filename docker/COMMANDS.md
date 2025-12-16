# Docker Commands Reference

Quick reference for all Docker-related commands for Suipic.

## Quick Start

### First Time Setup
```bash
# Linux/Mac
./docker-start.sh

# Windows
docker-start.bat

# Or manually
cp .env.docker .env
docker-compose up -d
```

## Using Makefile

The Makefile provides convenient shortcuts for common operations.

### Service Management
```bash
make help          # Show all available commands
make build         # Build all Docker images
make up            # Start all services
make down          # Stop all services
make restart       # Restart all services
make ps            # Show running containers
make stats         # Show container resource usage
make health        # Check health of all services
```

### Viewing Logs
```bash
make logs              # Show logs from all services
make backend-logs      # Show backend logs
make frontend-logs     # Show frontend logs
make nginx-logs        # Show nginx logs
make postgres-logs     # Show PostgreSQL logs
make es-logs           # Show Elasticsearch logs
make minio-logs        # Show MinIO logs
```

### Database Operations
```bash
make db-shell      # Open PostgreSQL shell
make db-backup     # Backup database to backup.sql
make db-restore    # Restore database from backup.sql
```

### Cleanup
```bash
make clean         # Stop and remove containers, networks
make prune         # Clean and remove volumes (WARNING: destroys data)
```

### SSL/Security
```bash
make ssl-cert      # Generate self-signed SSL certificate
```

### Deployment
```bash
make infra-up      # Start infrastructure services only
make prod-up       # Start with production configuration
```

## Docker Compose Commands

### Basic Operations
```bash
# Start services
docker-compose up -d

# Start specific service
docker-compose up -d backend

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Restart services
docker-compose restart

# Restart specific service
docker-compose restart backend

# View status
docker-compose ps

# View logs
docker-compose logs -f
docker-compose logs -f backend
docker-compose logs --tail=100 backend
```

### Building Images
```bash
# Build all images
docker-compose build

# Build specific service
docker-compose build backend

# Build without cache
docker-compose build --no-cache

# Pull latest base images
docker-compose pull
```

### Service Control
```bash
# Start only specific services
docker-compose up -d postgres elasticsearch minio

# Stop specific service
docker-compose stop backend

# Remove specific service container
docker-compose rm backend

# Recreate containers
docker-compose up -d --force-recreate
```

## Infrastructure Only (For Local Development)

```bash
# Start only database, search, and storage
docker-compose -f docker-compose.infra.yml up -d

# Stop infrastructure
docker-compose -f docker-compose.infra.yml down
```

## Production Deployment

```bash
# Start with production configuration
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Build for production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

# Stop production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml down
```

## Accessing Containers

### Execute Commands in Containers
```bash
# Backend shell
docker-compose exec backend sh

# PostgreSQL shell
docker-compose exec postgres sh
docker-compose exec postgres psql -U suipic -d suipic

# Nginx shell
docker-compose exec nginx sh

# Frontend shell
docker-compose exec frontend sh
```

### Run One-off Commands
```bash
# Run migration
docker-compose exec backend ./migrate -action=up

# Check database
docker-compose exec postgres pg_isready -U suipic

# Test nginx config
docker-compose exec nginx nginx -t

# Reload nginx
docker-compose exec nginx nginx -s reload
```

## Database Commands

### PostgreSQL Operations
```bash
# Connect to database
docker-compose exec postgres psql -U suipic -d suipic

# Backup database
docker-compose exec postgres pg_dump -U suipic suipic > backup.sql

# Restore database
docker-compose exec -T postgres psql -U suipic -d suipic < backup.sql

# Check database status
docker-compose exec postgres pg_isready -U suipic

# View database size
docker-compose exec postgres psql -U suipic -d suipic -c "SELECT pg_size_pretty(pg_database_size('suipic'))"
```

### SQL Queries
```bash
# List tables
docker-compose exec postgres psql -U suipic -d suipic -c "\dt"

# Count users
docker-compose exec postgres psql -U suipic -d suipic -c "SELECT COUNT(*) FROM users"

# List albums
docker-compose exec postgres psql -U suipic -d suipic -c "SELECT * FROM albums LIMIT 10"
```

## Elasticsearch Commands

### Health and Status
```bash
# Check cluster health
curl http://localhost:9200/_cluster/health?pretty

# List indices
curl http://localhost:9200/_cat/indices?v

# Get cluster stats
curl http://localhost:9200/_cluster/stats?pretty

# Check node info
curl http://localhost:9200/_nodes?pretty
```

### Index Operations
```bash
# Get index mapping
curl http://localhost:9200/photos/_mapping?pretty

# Count documents
curl http://localhost:9200/photos/_count?pretty

# Delete index
curl -X DELETE http://localhost:9200/photos

# Refresh index
curl -X POST http://localhost:9200/photos/_refresh
```

### Search
```bash
# Search all documents
curl http://localhost:9200/photos/_search?pretty

# Search with query
curl -X GET "http://localhost:9200/photos/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match_all": {}
  }
}
'
```

## MinIO Commands

### Access MinIO Console
```
URL: http://localhost:9001
Username: minioadmin
Password: minioadmin
```

### MinIO Health
```bash
# Check MinIO health
curl http://localhost:9000/minio/health/live

# Check MinIO readiness
curl http://localhost:9000/minio/health/ready
```

## Nginx Commands

### Configuration
```bash
# Test configuration
docker-compose exec nginx nginx -t

# Reload configuration
docker-compose exec nginx nginx -s reload

# View configuration
docker-compose exec nginx cat /etc/nginx/nginx.conf
docker-compose exec nginx cat /etc/nginx/conf.d/default.conf
```

### Logs
```bash
# View access logs
docker-compose exec nginx tail -f /var/log/nginx/access.log

# View error logs
docker-compose exec nginx tail -f /var/log/nginx/error.log
```

## Monitoring and Debugging

### Resource Usage
```bash
# Show resource stats
docker stats

# Show specific service stats
docker stats suipic-backend

# Show disk usage
docker system df

# Show volume usage
docker volume ls
```

### Logs
```bash
# Follow all logs
docker-compose logs -f

# Logs for specific service
docker-compose logs -f backend

# Last 100 lines
docker-compose logs --tail=100 backend

# Logs since timestamp
docker-compose logs --since 2023-01-01T00:00:00 backend

# Logs with timestamps
docker-compose logs -t backend
```

### Network
```bash
# List networks
docker network ls

# Inspect network
docker network inspect suipic_suipic-network

# Test connectivity
docker-compose exec backend ping postgres
docker-compose exec backend nc -zv postgres 5432
```

### Container Details
```bash
# Inspect container
docker inspect suipic-backend

# Show container processes
docker top suipic-backend

# Show port mappings
docker port suipic-backend
```

## Cleanup Commands

### Remove Containers
```bash
# Stop and remove all containers
docker-compose down

# Remove specific container
docker-compose rm backend

# Force remove running container
docker rm -f suipic-backend
```

### Remove Images
```bash
# Remove all project images
docker-compose down --rmi all

# Remove specific image
docker rmi suipic-backend:latest

# Remove unused images
docker image prune

# Remove all unused images
docker image prune -a
```

### Remove Volumes
```bash
# Remove all project volumes
docker-compose down -v

# Remove specific volume
docker volume rm suipic_postgres-data

# Remove unused volumes
docker volume prune
```

### Complete Cleanup
```bash
# Remove everything (containers, images, volumes, networks)
docker-compose down -v --rmi all

# Clean up all Docker resources
docker system prune -a --volumes -f

# Remove specific project resources
docker-compose down -v
docker volume rm suipic_postgres-data suipic_elasticsearch-data suipic_minio-data
```

## SSL Certificate Commands

### Generate Self-Signed Certificate
```bash
# Using Makefile
make ssl-cert

# Manually
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout docker/nginx/ssl/key.pem \
  -out docker/nginx/ssl/cert.pem \
  -subj "/C=US/ST=State/L=City/O=Suipic/CN=localhost"
```

### Let's Encrypt Certificate
```bash
# Stop nginx
docker-compose stop nginx

# Generate certificate
sudo certbot certonly --standalone -d yourdomain.com

# Copy certificates
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem docker/nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem docker/nginx/ssl/key.pem
sudo chown $USER:$USER docker/nginx/ssl/*.pem

# Start nginx
docker-compose start nginx
```

### Verify Certificate
```bash
# Check certificate info
openssl x509 -in docker/nginx/ssl/cert.pem -text -noout

# Check certificate expiry
openssl x509 -in docker/nginx/ssl/cert.pem -noout -dates

# Test SSL connection
openssl s_client -connect localhost:443
```

## Backup and Restore

### Database Backup
```bash
# Backup to file
docker-compose exec -T postgres pg_dump -U suipic suipic > backup.sql

# Backup with compression
docker-compose exec -T postgres pg_dump -U suipic suipic | gzip > backup.sql.gz

# Backup all databases
docker-compose exec -T postgres pg_dumpall -U suipic > backup-all.sql
```

### Database Restore
```bash
# Restore from file
docker-compose exec -T postgres psql -U suipic -d suipic < backup.sql

# Restore from compressed file
gunzip < backup.sql.gz | docker-compose exec -T postgres psql -U suipic -d suipic
```

### Volume Backup
```bash
# Backup PostgreSQL volume
docker run --rm \
  -v suipic_postgres-data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/postgres-backup.tar.gz -C /data .

# Backup all volumes
docker run --rm \
  -v suipic_postgres-data:/postgres \
  -v suipic_elasticsearch-data:/elasticsearch \
  -v suipic_minio-data:/minio \
  -v $(pwd):/backup \
  alpine tar czf /backup/full-backup.tar.gz /postgres /elasticsearch /minio
```

### Volume Restore
```bash
# Restore PostgreSQL volume
docker run --rm \
  -v suipic_postgres-data:/data \
  -v $(pwd):/backup \
  alpine tar xzf /backup/postgres-backup.tar.gz -C /data
```

## Health Checks

### Check All Services
```bash
# Using Makefile
make health

# Manually
docker-compose ps
curl http://localhost:9200/_cluster/health?pretty
docker-compose exec postgres pg_isready -U suipic
curl http://localhost:9000/minio/health/live
curl http://localhost:3000/api/health
```

### Individual Service Checks
```bash
# PostgreSQL
docker-compose exec postgres pg_isready -U suipic

# Elasticsearch
curl http://localhost:9200/_cluster/health?pretty

# MinIO
curl http://localhost:9000/minio/health/live

# Backend
curl http://localhost:3000/api/health

# Frontend
curl http://localhost:3001

# Nginx
curl -I http://localhost
```

## Troubleshooting Commands

### View Errors
```bash
# All errors in logs
docker-compose logs | grep -i error

# Backend errors
docker-compose logs backend | grep -i error

# Recent errors
docker-compose logs --since 5m | grep -i error
```

### Check Port Usage
```bash
# Windows
netstat -ano | findstr ":80"

# Linux/Mac
lsof -i :80
netstat -tlnp | grep :80
```

### Restart Specific Service
```bash
# Restart backend
docker-compose restart backend

# Rebuild and restart backend
docker-compose build backend
docker-compose up -d --force-recreate backend
```

### Reset Service
```bash
# Reset backend (remove container and recreate)
docker-compose stop backend
docker-compose rm -f backend
docker-compose up -d backend
```

## Advanced Commands

### Scale Services
```bash
# Run multiple backend instances
docker-compose up -d --scale backend=3
```

### Resource Limits
```bash
# Set memory limit
docker-compose up -d --memory="1g" backend

# Set CPU limit
docker-compose up -d --cpus="1.5" backend
```

### Copy Files
```bash
# Copy file from container
docker cp suipic-backend:/app/logs/app.log ./app.log

# Copy file to container
docker cp config.json suipic-backend:/app/config.json
```

## CI/CD Commands

### Build for CI
```bash
# Build with build args
docker-compose build --build-arg VERSION=1.0.0

# Build with no cache
docker-compose build --no-cache --parallel

# Export image
docker save suipic-backend:latest | gzip > backend-image.tar.gz

# Import image
docker load < backend-image.tar.gz
```

### Testing
```bash
# Start for testing
docker-compose -f docker-compose.yml -f docker-compose.test.yml up -d

# Run tests
docker-compose exec backend go test ./...
docker-compose exec frontend pnpm test
```

## Environment-Specific Commands

### Development
```bash
# Start with dev settings (automatic with docker-compose.override.yml)
docker-compose up -d

# With explicit override
docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d
```

### Production
```bash
# Start with production settings
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Build for production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build
```

## Reference

For more detailed information, see:
- Main Docker Guide: [../DOCKER.md](../DOCKER.md)
- Troubleshooting: [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
- Nginx Config: [nginx/README.md](nginx/README.md)
- Quick Start: [README.md](README.md)
