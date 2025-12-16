# Docker Troubleshooting Guide

Quick reference for common Docker issues with Suipic.

## Quick Diagnostics

### Check Service Status
```bash
docker-compose ps
```

### View All Logs
```bash
docker-compose logs -f
```

### Check Service Health
```bash
make health
# or
docker-compose ps
curl http://localhost:9200/_cluster/health?pretty
docker-compose exec postgres pg_isready -U suipic
```

## Common Issues

### 1. Services Won't Start

**Symptoms:**
- Container exits immediately
- "Container exited with code 1"

**Solutions:**

1. Check logs:
   ```bash
   docker-compose logs <service-name>
   ```

2. Verify environment variables:
   ```bash
   cat .env
   ```

3. Check port conflicts:
   ```bash
   # Windows
   netstat -ano | findstr ":<port>"
   
   # Linux/Mac
   netstat -tlnp | grep :<port>
   ```

4. Restart with clean state:
   ```bash
   docker-compose down
   docker-compose up -d
   ```

### 2. Database Connection Failed

**Symptoms:**
- Backend can't connect to PostgreSQL
- "connection refused" errors

**Solutions:**

1. Check if PostgreSQL is healthy:
   ```bash
   docker-compose ps postgres
   docker-compose exec postgres pg_isready -U suipic
   ```

2. Verify database credentials in `.env`:
   ```env
   DB_HOST=postgres
   DB_USER=suipic
   DB_PASSWORD=suipic_password
   ```

3. Check if database is initialized:
   ```bash
   docker-compose exec postgres psql -U suipic -d suipic -c "SELECT 1"
   ```

4. Reset database (WARNING: destroys data):
   ```bash
   docker-compose down -v
   docker-compose up -d
   ```

### 3. Migration Errors

**Symptoms:**
- Backend fails to start with migration errors
- "migration failed" in logs

**Solutions:**

1. Check migration logs:
   ```bash
   docker-compose logs backend | grep -i migration
   ```

2. Run migrations manually:
   ```bash
   docker-compose exec backend ./migrate -action=up
   ```

3. Check migration files exist:
   ```bash
   docker-compose exec backend ls -la db/migrations/
   ```

4. Rollback and retry:
   ```bash
   docker-compose exec backend ./migrate -action=down
   docker-compose exec backend ./migrate -action=up
   ```

### 4. Elasticsearch Not Ready

**Symptoms:**
- Backend can't connect to Elasticsearch
- Search functionality not working

**Solutions:**

1. Check Elasticsearch health:
   ```bash
   curl http://localhost:9200/_cluster/health?pretty
   ```

2. Check Elasticsearch logs:
   ```bash
   docker-compose logs elasticsearch
   ```

3. Increase memory if needed (edit `docker-compose.yml`):
   ```yaml
   elasticsearch:
     environment:
       - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
   ```

4. Restart Elasticsearch:
   ```bash
   docker-compose restart elasticsearch
   ```

### 5. MinIO Connection Issues

**Symptoms:**
- Photo upload fails
- "connection refused" to MinIO

**Solutions:**

1. Check MinIO health:
   ```bash
   curl http://localhost:9000/minio/health/live
   ```

2. Access MinIO Console:
   ```
   http://localhost:9001
   Login: minioadmin / minioadmin
   ```

3. Verify bucket exists:
   ```bash
   docker-compose exec backend ./main check-minio
   ```

4. Check MinIO logs:
   ```bash
   docker-compose logs minio
   ```

### 6. Frontend Build Fails

**Symptoms:**
- Frontend container exits
- Build errors in logs

**Solutions:**

1. Check frontend logs:
   ```bash
   docker-compose logs frontend
   ```

2. Verify Node.js version in Dockerfile:
   ```dockerfile
   FROM node:20-alpine
   ```

3. Clean build and retry:
   ```bash
   docker-compose down
   docker-compose build --no-cache frontend
   docker-compose up -d frontend
   ```

4. Check if pnpm is installed:
   ```bash
   docker-compose exec frontend pnpm --version
   ```

### 7. Nginx 502 Bad Gateway

**Symptoms:**
- "502 Bad Gateway" error
- Can't access application

**Solutions:**

1. Check if backend/frontend are running:
   ```bash
   docker-compose ps backend frontend
   ```

2. Test backend directly:
   ```bash
   curl http://localhost:3000/api/health
   ```

3. Test frontend directly:
   ```bash
   curl http://localhost:3001
   ```

4. Check nginx configuration:
   ```bash
   docker-compose exec nginx nginx -t
   ```

5. View nginx logs:
   ```bash
   docker-compose logs nginx
   ```

### 8. SSL/HTTPS Issues

**Symptoms:**
- "Certificate not valid" errors
- Can't access via HTTPS

**Solutions:**

1. Verify certificate files exist:
   ```bash
   ls -la docker/nginx/ssl/
   ```

2. Test certificate:
   ```bash
   openssl x509 -in docker/nginx/ssl/cert.pem -text -noout
   ```

3. Check nginx SSL configuration:
   ```bash
   docker-compose exec nginx cat /etc/nginx/conf.d/default.conf
   ```

4. Regenerate self-signed certificate:
   ```bash
   make ssl-cert
   docker-compose restart nginx
   ```

### 9. Out of Disk Space

**Symptoms:**
- "no space left on device"
- Services fail to start

**Solutions:**

1. Check Docker disk usage:
   ```bash
   docker system df
   ```

2. Clean up unused resources:
   ```bash
   docker system prune -a
   ```

3. Remove unused volumes (WARNING: may delete data):
   ```bash
   docker volume prune
   ```

4. Clean up old images:
   ```bash
   docker image prune -a
   ```

### 10. Port Already in Use

**Symptoms:**
- "port is already allocated"
- Can't bind to port

**Solutions:**

1. Check what's using the port:
   ```bash
   # Windows
   netstat -ano | findstr ":80"
   
   # Linux/Mac
   lsof -i :80
   ```

2. Stop conflicting service or change port in `docker-compose.yml`:
   ```yaml
   services:
     nginx:
       ports:
         - "8080:80"  # Use 8080 instead of 80
   ```

3. Kill process using the port (if safe):
   ```bash
   # Linux/Mac
   kill -9 <PID>
   
   # Windows
   taskkill /PID <PID> /F
   ```

### 11. Container Keeps Restarting

**Symptoms:**
- Container shows "Restarting" status
- Service cycles constantly

**Solutions:**

1. Check logs for crash reason:
   ```bash
   docker-compose logs <service-name>
   ```

2. Remove restart policy temporarily:
   ```bash
   docker update --restart=no <container-name>
   ```

3. Check health check configuration
4. Verify resource limits aren't too low

### 12. Slow Performance

**Symptoms:**
- Application is slow
- High CPU/memory usage

**Solutions:**

1. Check resource usage:
   ```bash
   docker stats
   ```

2. Increase resource limits in `docker-compose.prod.yml`:
   ```yaml
   deploy:
     resources:
       limits:
         cpus: '4'
         memory: 2G
   ```

3. Check for memory leaks in logs
4. Optimize Elasticsearch memory:
   ```yaml
   ES_JAVA_OPTS=-Xms2g -Xmx2g
   ```

## Complete Reset

If all else fails, perform a complete reset:

⚠️ **WARNING: This will delete ALL data!**

```bash
# Stop everything
docker-compose down -v

# Remove all containers, images, volumes
docker system prune -a --volumes -f

# Remove project-specific volumes
docker volume rm suipic_postgres-data suipic_elasticsearch-data suipic_minio-data

# Start fresh
docker-compose up -d
```

## Getting More Information

### Enable Debug Logging

Backend (`docker-compose.yml`):
```yaml
backend:
  environment:
    ENV: development
    LOG_LEVEL: debug
```

Nginx (`docker/nginx/nginx.conf`):
```nginx
error_log /var/log/nginx/error.log debug;
```

### Check Container Details

```bash
docker inspect <container-name>
```

### Access Container Shell

```bash
docker-compose exec backend sh
docker-compose exec postgres sh
```

### Network Issues

Check network connectivity:
```bash
docker-compose exec backend ping postgres
docker-compose exec backend nc -zv postgres 5432
```

List networks:
```bash
docker network ls
docker network inspect suipic_suipic-network
```

## Prevention Tips

1. **Always check logs first**: `docker-compose logs -f`
2. **Keep Docker updated**: Update Docker Desktop regularly
3. **Monitor resources**: Use `docker stats` to watch usage
4. **Regular backups**: Backup volumes regularly
5. **Use health checks**: Ensure all services have health checks
6. **Test changes locally**: Test configuration changes in development first
7. **Document customizations**: Keep track of any custom changes

## Still Having Issues?

1. Check the main [DOCKER.md](../DOCKER.md) for detailed setup
2. Review service-specific documentation
3. Check Docker and Docker Compose versions
4. Try with a clean slate (reset everything)
5. Search Docker logs for specific error messages
6. Check system resources (CPU, memory, disk)

## Useful Commands Reference

```bash
# Service Management
docker-compose up -d              # Start all services
docker-compose down               # Stop all services
docker-compose restart <service>  # Restart specific service
docker-compose ps                 # List services

# Logs
docker-compose logs -f            # Follow all logs
docker-compose logs backend       # Specific service logs
docker-compose logs --tail=50     # Last 50 lines

# Database
make db-shell                     # PostgreSQL shell
make db-backup                    # Backup database
make db-restore                   # Restore database

# Cleanup
docker-compose down -v            # Remove volumes
docker system prune               # Clean up unused resources

# Health Checks
make health                       # Check all services
curl http://localhost:9200/_cluster/health?pretty

# Rebuild
docker-compose build --no-cache   # Rebuild without cache
docker-compose up -d --force-recreate
```
