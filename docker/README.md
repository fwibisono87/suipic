# Docker Configuration for Suipic

This directory contains Docker configuration files for running Suipic in containers.

## Services

The Docker Compose setup includes the following services:

1. **PostgreSQL** - Database server with persistent storage
2. **Elasticsearch** - Search engine for photo indexing
3. **MinIO** - S3-compatible object storage for photos
4. **Backend** - Go Fiber API server
5. **Frontend** - SvelteKit web application
6. **Nginx** - Reverse proxy with SSL support

## Quick Start

### Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+

### Running the Application

1. **Configure environment variables:**
   ```bash
   cp .env.docker .env
   # Edit .env with your configuration
   ```

2. **Start all services:**
   ```bash
   docker-compose up -d
   ```

3. **Check service status:**
   ```bash
   docker-compose ps
   ```

4. **View logs:**
   ```bash
   docker-compose logs -f
   ```

5. **Access the application:**
   - Frontend: http://localhost
   - Backend API: http://localhost/api
   - MinIO Console: http://localhost:9001
   - Elasticsearch: http://localhost:9200

### Stopping the Application

```bash
docker-compose down
```

To also remove volumes (database data will be lost):
```bash
docker-compose down -v
```

## Service Details

### PostgreSQL
- **Port**: 5432
- **Volume**: `postgres-data`
- **Init Scripts**: `docker/postgres/init-scripts/`
- **Credentials**: Configured via environment variables

### Elasticsearch
- **Port**: 9200 (HTTP), 9300 (Transport)
- **Volume**: `elasticsearch-data`
- **Mode**: Single-node (for development)
- **Memory**: 512MB heap size

### MinIO
- **Port**: 9000 (API), 9001 (Console)
- **Volume**: `minio-data`
- **Default Credentials**: minioadmin/minioadmin
- **Bucket**: suipic (auto-created by backend)

### Backend
- **Port**: 3000
- **Image**: Built from `backend/Dockerfile`
- **Dependencies**: PostgreSQL, Elasticsearch, MinIO
- **Migrations**: Run automatically on startup

### Frontend
- **Port**: 3001
- **Image**: Built from `frontend/Dockerfile`
- **Adapter**: Node.js adapter for SvelteKit
- **Dependencies**: Backend API

### Nginx
- **Port**: 80 (HTTP), 443 (HTTPS)
- **Configuration**: `docker/nginx/nginx.conf`
- **Site Config**: `docker/nginx/conf.d/default.conf`
- **SSL Certificates**: `docker/nginx/ssl/`

## Configuration

### Environment Variables

Key environment variables can be customized in `.env` file:

- `POSTGRES_PASSWORD` - Database password
- `JWT_SECRET` - Secret key for JWT tokens (CHANGE IN PRODUCTION!)
- `ADMIN_PASSWORD` - Initial admin user password
- `MINIO_ROOT_PASSWORD` - MinIO admin password
- `PUBLIC_API_URL` - Frontend API URL

### SSL Configuration

To enable HTTPS:

1. **Generate or obtain SSL certificates:**
   ```bash
   # Self-signed for development
   cd docker/nginx/ssl
   openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
     -keyout key.pem -out cert.pem \
     -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
   ```

2. **Update nginx configuration:**
   - Edit `docker/nginx/conf.d/default.conf`
   - Uncomment the HTTPS server block
   - Uncomment the HTTP to HTTPS redirect

3. **Restart nginx:**
   ```bash
   docker-compose restart nginx
   ```

### Database Migrations

Migrations are automatically run when the backend starts. To run migrations manually:

```bash
docker-compose exec backend ./main migrate
```

Or access the container:
```bash
docker-compose exec backend sh
```

## Development Workflow

### Building Images

To rebuild images after code changes:

```bash
docker-compose build
docker-compose up -d
```

Build specific service:
```bash
docker-compose build backend
docker-compose up -d backend
```

### Viewing Logs

All services:
```bash
docker-compose logs -f
```

Specific service:
```bash
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Database Access

Connect to PostgreSQL:
```bash
docker-compose exec postgres psql -U suipic -d suipic
```

### MinIO Access

Access MinIO Console at http://localhost:9001 with credentials from environment variables.

### Elasticsearch Access

Check cluster health:
```bash
curl http://localhost:9200/_cluster/health?pretty
```

## Production Considerations

1. **Change default passwords** in `.env`
2. **Use strong JWT_SECRET**
3. **Configure SSL certificates** for HTTPS
4. **Set appropriate CORS_ORIGINS**
5. **Configure backup strategy** for volumes
6. **Enable Elasticsearch security** if needed
7. **Use Docker secrets** for sensitive data
8. **Set resource limits** in docker-compose.yml
9. **Monitor logs and metrics**
10. **Regular security updates** for base images

## Troubleshooting

### Service won't start

Check logs:
```bash
docker-compose logs <service-name>
```

### Database connection errors

Ensure PostgreSQL is healthy:
```bash
docker-compose ps postgres
docker-compose logs postgres
```

### Port conflicts

Check if ports are already in use:
```bash
netstat -tlnp | grep -E ':(80|443|3000|3001|5432|9000|9001|9200)'
```

### Reset everything

Remove all containers, volumes, and data:
```bash
docker-compose down -v
docker-compose up -d
```

## Network Architecture

All services are connected via the `suipic-network` bridge network:

```
┌─────────────────────────────────────────┐
│             Nginx (Reverse Proxy)       │
│         HTTP: 80, HTTPS: 443           │
└─────────────┬───────────────────────────┘
              │
     ┌────────┴─────────┐
     │                  │
┌────▼─────┐      ┌────▼─────┐
│ Frontend │      │ Backend  │
│  :3001   │      │  :3000   │
└──────────┘      └────┬─────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
   ┌────▼────┐   ┌────▼────┐   ┌────▼────┐
   │Postgres │   │  Elastic│   │  MinIO  │
   │  :5432  │   │  :9200  │   │ :9000/1 │
   └─────────┘   └─────────┘   └─────────┘
```

## Volumes

Persistent data is stored in Docker volumes:

- `postgres-data` - Database files
- `elasticsearch-data` - Search indices
- `minio-data` - Uploaded photos

To backup volumes:
```bash
docker run --rm -v postgres-data:/data -v $(pwd):/backup alpine tar czf /backup/postgres-backup.tar.gz -C /data .
```

To restore volumes:
```bash
docker run --rm -v postgres-data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres-backup.tar.gz -C /data
```
